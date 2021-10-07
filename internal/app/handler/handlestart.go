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

func (srv *Handler) handleStart(c telebot.Context) error {
	srv.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Debug("get account from database by telegram_id")
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		if err == store.ErrRecordNotFound {
			srv.logger.Debug("account not found, will create a new one")
			account = &model.Account{
				Created:    time.Now(),
				TelegramID: int64(c.Sender().ID),
				FirstName:  c.Sender().FirstName,
				LastName:   c.Sender().LastName,
				Username:   c.Sender().Username,
				Superuser:  false,
			}

			if err := srv.store.Account().Create(account); err != nil {
				srv.logger.Error(err)
				return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
			}

			srv.logger.WithFields(account.LogrusFields()).Debug("account created")
			return c.EditOrReply(
				fmt.Sprintf(msgUserCreated, account.FirstName, account.FirstName, account.LastName, account.Username, account.Superuser),
				&telebot.SendOptions{ParseMode: "HTML"},
			)
		}

		srv.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	srv.logger.WithFields(account.LogrusFields()).Debug("account found")
	return c.EditOrReply(
		fmt.Sprintf(msgUserExist, account.FirstName, account.FirstName, account.LastName, account.Username, account.Superuser),
		&telebot.SendOptions{ParseMode: "HTML"},
	)
}
