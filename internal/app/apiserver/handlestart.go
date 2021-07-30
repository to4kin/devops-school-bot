package apiserver

import (
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleStart(c telebot.Context) error {
	logrus.Debug("get account from database by telegram_id: ", c.Sender().ID)
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		if err == store.ErrRecordNotFound {
			logrus.Debug("account not found, will create a new one")
			account = &model.Account{
				TelegramID: int64(c.Sender().ID),
				FirstName:  c.Sender().FirstName,
				LastName:   c.Sender().LastName,
				Username:   c.Sender().Username,
				Superuser:  false,
			}

			if err := srv.store.Account().Create(account); err != nil {
				logrus.Error(err)
				return nil
			}
		} else {
			logrus.Error(err)
			return nil
		}
	}
	logrus.Debug(account.ToString())
	return nil
}
