package apiserver

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	msgVersion   string = "dev"
	msgBuildDate string = ""
	msgBotInfo   string = "<b>Bot information:</b>\nVersion: %v\nBuild date: %v\nBuilt with: %v"

	msgHelpCommand string = `I'll manage students homeworks

<b>User Commands</b>
/start - Add user to database
/joinstudent - Join school as student
/joinmodule - Join school as listener
/myreport - Your progress
/homeworks - Homeworks list
/help - Help message

<b>Superuser Commands</b>
/schools - Manage schools
/startschool - Start school
/stopschool - Finish school
/report - School progress
/fullreport - School progress with homework list

/users - Manage users
/setsuperuser - Set Superuser
/unsetsuperuser - Unset Superuser
`
)

func (srv *server) handleHelp(c telebot.Context) error {
	message := msgHelpCommand

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
				message += fmt.Sprintf(msgBotInfo, msgVersion, msgBuildDate, runtime.Version())
			}
		}
	}

	return c.EditOrReply(msgHelpCommand, &telebot.SendOptions{ParseMode: "HTML"})
}
