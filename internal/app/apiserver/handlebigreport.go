package apiserver

import (
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleBigReport(c telebot.Context) error {
	if c.Message().Private() {
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

	srv.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get all lessons from database by school_id")
	lessons, err := srv.store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		srv.logger.Error(err)
		return nil
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(lessons),
	}).Debug("lessons found")

	srv.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get students from database by school_id")
	students, err := srv.store.Student().FindBySchoolID(school.ID)
	if err != nil {
		srv.logger.Error(err)
		return nil
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Debug("students found")

	reportMessage := srv.prepareHomeworkReportMsg(lessons)

	report, err := srv.prepareReportMsg(students, lessons)
	if err != nil {
		srv.logger.Error(err)
		return nil
	}

	reportMessage += "\n" + report

	srv.logger.Debug("report sent")
	return c.Reply(reportMessage, &telebot.SendOptions{ParseMode: "HTML"})
}
