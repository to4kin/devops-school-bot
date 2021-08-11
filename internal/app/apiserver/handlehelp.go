package apiserver

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleHelp(c telebot.Context) error {
	message := msgHelpCommand
	if c.Message().Private() {
		message += msgUserPrivateCmd
	} else {
		message += msgUserGroupCmd
	}

	srv.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Debug("get account from database by telegram_id")
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		if err == store.ErrRecordNotFound {
			srv.logger.Debug(err)
		} else {
			srv.logger.Error(err)
		}
	} else {
		srv.logger.WithFields(account.LogrusFields()).Debug("account found")

		if account.Superuser {
			if c.Message().Private() {
				message += msgSuperuserPrivateCmd
				message += "\n\n"
				message += fmt.Sprintf(msgBotInfo, msgVersion, msgBuildDate, runtime.Version())
			} else {
				message += msgSuperuserGroupCmd
			}
		}
	}

	return c.Reply(message, &telebot.SendOptions{ParseMode: "HTML"})
}
