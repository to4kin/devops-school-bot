package apiserver

import (
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleJoin(c telebot.Context) error {
	if c.Message().Private() {
		return nil
	}

	logrus.Debug("get school by chat_id: ", c.Message().Chat.ID)
	school, err := srv.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			logrus.Error(err)
			return c.Reply(msgNoActiveSchool, &telebot.SendOptions{ParseMode: "HTML"})
		}

		logrus.Error(err)
		return nil
	}
	logrus.Debug(school.ToString())

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

	logrus.Debug("get student from database by account_id: ", account.ID, " and school_id: ", school.ID)
	student, err := srv.store.Student().FindByAccountIDSchoolID(account.ID, school.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			logrus.Debug("student not found, will create a new one")
			student := &model.Student{
				Account: account,
				School:  school,
				Active:  true,
			}

			if err := srv.store.Student().Create(student); err != nil {
				logrus.Error(err)
				return nil
			}

			logrus.Debug(student.ToString())
			return c.Reply(msgWelcomeToSchool, &telebot.SendOptions{ParseMode: "HTML"})
		}

		logrus.Error(err)
		return nil
	}
	logrus.Debug(student.ToString())

	return c.Reply(msgUserAlreadyJoined, &telebot.SendOptions{ParseMode: "HTML"})
}
