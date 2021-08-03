package apiserver

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleJoin(c telebot.Context) error {
	if c.Message().Private() {
		return nil
	}

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
		} else {
			srv.logger.Error(err)
			return nil
		}
	} else {
		srv.logger.WithFields(account.LogrusFields()).Debug("account found")
	}

	srv.logger.WithFields(logrus.Fields{
		"chat_id": c.Message().Chat.ID,
	}).Debug("get school by chat_id")
	school, err := srv.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		srv.logger.Error(err)
		if err == store.ErrRecordNotFound {
			return c.Reply(msgSchoolNotFound, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return nil
	}
	srv.logger.WithFields(school.LogrusFields()).Debug("school found")

	if school.Finished {
		srv.logger.WithFields(school.LogrusFields()).Debug("school already finished")
		return c.Reply(fmt.Sprintf(msgSchoolAlreadyFinished, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
	}

	srv.logger.WithFields(logrus.Fields{
		"account_id": account.ID,
		"school_id":  school.ID,
	}).Debug("get student from database by account_id and school_id")
	student, err := srv.store.Student().FindByAccountIDSchoolID(account.ID, school.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			srv.logger.Debug("student not found, will create a new one")
			student := &model.Student{
				Created: time.Now(),
				Account: account,
				School:  school,
				Active:  true,
			}

			if err := srv.store.Student().Create(student); err != nil {
				srv.logger.Error(err)
				return nil
			}

			srv.logger.WithFields(student.LogrusFields()).Debug("student created")
			return c.Reply(fmt.Sprintf(msgWelcomeToSchool, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
		}

		srv.logger.Error(err)
		return nil
	}
	srv.logger.WithFields(student.LogrusFields()).Debug("student exist")
	return c.Reply(fmt.Sprintf(msgUserAlreadyJoined, school.Title), &telebot.SendOptions{ParseMode: "HTML"})
}
