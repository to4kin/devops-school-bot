package handler

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	msgSchoolStarted string = "school <b>%v</b> started"
)

func (handler *Handler) handleStartSchool(c telebot.Context) error {
	if c.Message().Private() {
		return c.EditOrReply(fmt.Sprintf(helper.ErrWrongChatType, "SCHOOL"), &telebot.SendOptions{ParseMode: "HTML"})
	}

	handler.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Info("get account from database by telegram_id")
	account, err := handler.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	handler.logger.WithFields(account.LogrusFields()).Info("account found")

	if !account.Superuser {
		handler.logger.WithFields(account.LogrusFields()).Info("account has insufficient permissions")
		return c.EditOrReply(helper.ErrInsufficientPermissions, &telebot.SendOptions{ParseMode: "HTML"})
	}

	handler.logger.WithFields(logrus.Fields{
		"chat_id": c.Message().Chat.ID,
	}).Info("get school by chat_id")
	school, err := handler.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			handler.logger.Info("school not found, will create a new one")
			school = &model.School{
				Created: time.Now(),
				Title:   c.Message().Chat.Title,
				ChatID:  c.Message().Chat.ID,
				Active:  true,
			}

			if err := handler.store.School().Create(school); err != nil {
				handler.logger.Error(err)
				return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
			}

			handler.logger.WithFields(school.LogrusFields()).Info("school created")
			return c.EditOrReply(fmt.Sprintf(msgSchoolStarted, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
		}

		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	if !school.Active {
		handler.logger.WithFields(school.LogrusFields()).Info("school will be started")
		school.Active = true
		if err := handler.store.School().Update(school); err != nil {
			handler.logger.Error(err)
			return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
		}
	}

	handler.logger.WithFields(school.LogrusFields()).Info("school started")
	return c.EditOrReply(fmt.Sprintf(msgSchoolStarted, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
}
