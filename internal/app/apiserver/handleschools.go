package apiserver

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleSchools(c telebot.Context) error {
	if !c.Message().Private() {
		return nil
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

	return srv.schoolNaviButtons(c, 0)
}
