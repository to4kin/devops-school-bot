package handler

import (
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (handler *Handler) handleMyReport(c telebot.Context) error {
	if c.Message().Private() {
		return c.EditOrReply(helper.ErrWrongChatType, &telebot.SendOptions{ParseMode: "HTML"})
	}

	handler.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Info("get account from database by telegram_id")
	account, err := handler.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		handler.logger.Error(err)

		if err == store.ErrRecordNotFound {
			return c.EditOrReply(helper.ErrUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	handler.logger.WithFields(account.LogrusFields()).Info("account found")

	handler.logger.WithFields(logrus.Fields{
		"chat_id": c.Message().Chat.ID,
	}).Info("get school by chat_id")
	school, err := handler.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		handler.logger.Error(err)

		if err == store.ErrRecordNotFound {
			return c.EditOrReply(helper.ErrSchoolNotStarted, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	handler.logger.WithFields(school.LogrusFields()).Info("school found")

	hlpr := helper.NewHelper(handler.store, handler.logger)
	reportMessage, err := hlpr.GetUserReport(account, school)
	if err != nil {
		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	handler.logger.Info("report sent")
	return c.EditOrReply(reportMessage, &telebot.SendOptions{ParseMode: "HTML"})
}
