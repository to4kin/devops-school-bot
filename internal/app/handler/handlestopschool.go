package handler

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	msgSchoolFinished string = `school <b>%v</b> finished`
)

func (handler *Handler) handleStopSchool(c telebot.Context) error {
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

	if c.Message().Private() {
		callback := &model.Callback{
			ID:          0,
			Type:        "school",
			Command:     "stop",
			ListCommand: "stop",
		}

		hlpr := helper.NewHelper(handler.store, handler.logger)
		replyMessage, replyMarkup, err := hlpr.GetSchoolsList(callback)
		if err != nil {
			handler.logger.Error(err)
			return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return c.EditOrReply(replyMessage, replyMarkup)
	}

	handler.logger.WithFields(logrus.Fields{
		"chat_id": c.Message().Chat.ID,
	}).Info("get school by chat_id")
	school, err := handler.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrSchoolNotStarted, &telebot.SendOptions{ParseMode: "HTML"})
	}
	handler.logger.WithFields(school.LogrusFields()).Info("school found")

	if !school.Active {
		handler.logger.WithFields(school.LogrusFields()).Info("school already finished")
		return c.EditOrReply(fmt.Sprintf(msgSchoolFinished, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
	}

	school.Active = false
	if err := handler.store.School().Update(school); err != nil {
		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	handler.logger.WithFields(school.LogrusFields()).Info("school finished")
	return c.EditOrReply(fmt.Sprintf(msgSchoolFinished, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
}
