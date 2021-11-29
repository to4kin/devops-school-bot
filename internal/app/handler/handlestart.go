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
	msgUserCreated string = "Hello, <b>%v!</b>\nAccount created successfully!\n\n" + msgUserInfo
	msgUserExist   string = "Hello, <b>%v!</b>\nAccount already exist!\n\n" + msgUserInfo
	msgUserInfo    string = "Account info:\n\nFirst name: %v\nLast name: %v\nUsername: @%v\nSuperuser: %v"
)

func (handler *Handler) handleStart(c telebot.Context) error {
	if !c.Message().Private() {
		return c.EditOrReply(fmt.Sprintf(helper.ErrWrongChatType, "PRIVATE"), &telebot.SendOptions{ParseMode: "HTML"})
	}

	handler.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Info("get account from database by telegram_id")
	account, err := handler.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		if err == store.ErrRecordNotFound {
			handler.logger.Info("account not found, will create a new one")
			account = &model.Account{
				Created:    time.Now(),
				TelegramID: int64(c.Sender().ID),
				FirstName:  c.Sender().FirstName,
				LastName:   c.Sender().LastName,
				Username:   c.Sender().Username,
				Superuser:  false,
			}

			if err := handler.store.Account().Create(account); err != nil {
				handler.logger.Error(err)
				return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
			}

			handler.logger.WithFields(account.LogrusFields()).Info("account created")
			return c.EditOrReply(
				fmt.Sprintf(msgUserCreated, account.FirstName, account.FirstName, account.LastName, account.Username, account.Superuser),
				&telebot.SendOptions{ParseMode: "HTML"},
			)
		}

		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	handler.logger.WithFields(account.LogrusFields()).Info("account found")
	return c.EditOrReply(
		fmt.Sprintf(msgUserExist, account.FirstName, account.FirstName, account.LastName, account.Username, account.Superuser),
		&telebot.SendOptions{ParseMode: "HTML"},
	)
}
