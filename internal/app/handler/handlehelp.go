package handler

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	msgBotInfo string = "\n\n<b>Bot information:</b>\nVersion: %v\nBuild date: %v\nBuilt with: %v"

	msgGroupHelp string = `I'll manage students homeworks

<b>User Commands</b>
/joinstudent - Join school as student
/joinmodule - Join school as listener
/help - Help message

<b>Superuser Commands</b>
/startschool - Start school
/stopschool - Finish school
/report - School progress
/fullreport - School progress with homework list
`

	msgPrivateHelp string = `I'll manage students homeworks

<b>User Commands</b>
/start - Add user to database
/myreport - Your progress
/homeworks - Homeworks list
/help - Help message

<b>Superuser Commands</b>
/schools - Manage schools
/report - School progress
/fullreport - School progress with homework list

/users - Manage users
/setsuperuser - Set Superuser
/unsetsuperuser - Unset Superuser
`
)

func (handler *Handler) handleHelp(c telebot.Context) error {
	message := ""

	if c.Message().Private() {
		message = msgPrivateHelp
	} else {
		message = msgGroupHelp
	}

	handler.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Info("get account from database by telegram_id")
	account, err := handler.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		if err == store.ErrRecordNotFound {
			handler.logger.Info(err)
		} else {
			handler.logger.Error(err)
		}
	} else {
		handler.logger.WithFields(account.LogrusFields()).Info("account found")

		if account.Superuser {
			if c.Message().Private() {
				message += fmt.Sprintf(msgBotInfo, handler.version, handler.buildDate, runtime.Version())
			}
		}
	}

	return handler.editOrReply(c, message, nil)
}
