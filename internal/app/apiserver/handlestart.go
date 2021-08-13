package apiserver

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleStart(c telebot.Context) error {
	if c.Message().Private() {
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
					return nil
				}

				srv.logger.WithFields(account.LogrusFields()).Debug("account created")
				return c.Reply(
					fmt.Sprintf(msgUserCreated, account.FirstName, account.FirstName, account.LastName, account.Username, account.Superuser),
					&telebot.SendOptions{ParseMode: "HTML"},
				)
			}

			srv.logger.Error(err)
			return nil
		}

		srv.logger.WithFields(account.LogrusFields()).Debug("account found")
		return c.Reply(
			fmt.Sprintf(msgUserExist, account.FirstName, account.FirstName, account.LastName, account.Username, account.Superuser),
			&telebot.SendOptions{ParseMode: "HTML"},
		)
	}

	srv.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Debug("get account from database by telegram_id")
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		srv.logger.Error(err)
		return nil
	}
	srv.logger.WithFields(account.LogrusFields()).Debug("account found")

	if !account.Superuser {
		srv.logger.WithFields(account.LogrusFields()).Debug("account has insufficient permissions")
		return c.Reply(msgUserInsufficientPermissions, &telebot.SendOptions{ParseMode: "HTML"})
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
				return nil
			}

			srv.logger.WithFields(school.LogrusFields()).Debug("school created")
			return c.Reply(fmt.Sprintf(msgSchoolStarted, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
		}

		srv.logger.Error(err)
		return nil
	}

	if !school.Active {
		srv.logger.WithFields(school.LogrusFields()).Debug("school already finished")
		return c.Reply(fmt.Sprintf(msgSchoolAlreadyFinished, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
	}

	srv.logger.WithFields(school.LogrusFields()).Debug("school already started")
	return c.Reply(fmt.Sprintf(msgSchoolAlreadyStarted, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
}
