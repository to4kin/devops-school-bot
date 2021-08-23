package apiserver

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

func (srv *server) handleStartSchool(c telebot.Context) error {
	srv.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Debug("get account from database by telegram_id")
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		srv.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	srv.logger.WithFields(account.LogrusFields()).Debug("account found")

	if !account.Superuser {
		srv.logger.WithFields(account.LogrusFields()).Debug("account has insufficient permissions")
		return c.EditOrReply(helper.ErrInsufficientPermissions, &telebot.SendOptions{ParseMode: "HTML"})
	}

	if c.Message().Private() {
		callback := &model.Callback{
			ID:          0,
			Type:        "school",
			Command:     "start",
			ListCommand: "start",
		}

		replyMessage, replyMarkup, err := helper.GetSchoolsList(srv.store, callback)
		if err != nil {
			srv.logger.Error(err)
			return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return c.EditOrReply(replyMessage, replyMarkup)
	}

	srv.logger.WithFields(logrus.Fields{
		"chat_id": c.Message().Chat.ID,
	}).Debug("get school by chat_id")
	school, err := srv.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			srv.logger.Debug("school not found, will create a new one")
			school = &model.School{
				Created: time.Now(),
				Title:   c.Message().Chat.Title,
				ChatID:  c.Message().Chat.ID,
				Active:  true,
			}

			if err := srv.store.School().Create(school); err != nil {
				srv.logger.Error(err)
				return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
			}

			srv.logger.WithFields(school.LogrusFields()).Debug("school created")
			return c.EditOrReply(fmt.Sprintf(msgSchoolStarted, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
		}

		srv.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	if !school.Active {
		srv.logger.WithFields(school.LogrusFields()).Debug("school will be started")
		school.Active = true
		if err := srv.store.School().Update(school); err != nil {
			srv.logger.Error(err)
			return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
		}
	}

	srv.logger.WithFields(school.LogrusFields()).Debug("school started")
	return c.EditOrReply(fmt.Sprintf(msgSchoolStarted, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
}
