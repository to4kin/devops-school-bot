package handler

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (handler *Handler) handleMyReport(c telebot.Context) error {
	if !c.Message().Private() {
		return c.EditOrReply(fmt.Sprintf(helper.ErrWrongChatType, "PRIVATE"), &telebot.SendOptions{ParseMode: "HTML"})
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

	hlpr := helper.NewHelper(handler.store, handler.logger)

	reportMessage := fmt.Sprintf(
		"Hello, @%v!\n\n<b>Account info:</b>\nFirst name: %v\nLast name: %v\n\n",
		account.Username,
		account.FirstName,
		account.LastName,
	)

	message, err := hlpr.GetUserReport(account)
	if err != nil {
		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	reportMessage += message

	handler.logger.Info("report sent")
	return c.EditOrReply(reportMessage, &telebot.SendOptions{ParseMode: "HTML"})
}
