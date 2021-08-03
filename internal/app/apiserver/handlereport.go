package apiserver

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleReport(c telebot.Context) error {
	if c.Message().Private() {
		return nil
	}

	srv.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Debug("get account from database by telegram_id")
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		srv.logger.Error(err)

		if err == store.ErrRecordNotFound {
			return c.Reply(msgUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return nil
	}
	srv.logger.WithFields(account.LogrusFields()).Debug("account found")

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
		"account_id": account.ID,
		"school_id":  school.ID,
	}).Debug("get student from database by account_id and school_id")
	student, err := srv.store.Student().FindByAccountIDSchoolID(account.ID, school.ID)
	if err != nil {
		srv.logger.Error(err)

		if err == store.ErrRecordNotFound {
			return c.Reply(msgUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return nil
	}
	srv.logger.WithFields(student.LogrusFields()).Debug("student found")

	srv.logger.WithFields(logrus.Fields{
		"student_id": student.ID,
	}).Debug("get student homeworks from database by student_id")
	studentHomeworks, err := srv.store.Homework().FindByStudentID(student.ID)
	if err != nil {
		srv.logger.Error(err)

		if err == store.ErrRecordNotFound {
			c.Reply(msgHomeworkNotProvided, &telebot.SendOptions{ParseMode: "HTML"})
		}
		return nil
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(studentHomeworks),
	}).Debug("homeworks found")

	srv.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get all lessons from database by school_id")
	allLessons, err := srv.store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		srv.logger.Error(err)
		return nil
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(allLessons),
	}).Debug("lessons found")

	reportMessage := fmt.Sprintf(msgHomeworkReport, account.Username, school.Title)
	for _, lesson := range allLessons {
		counted := false
		for _, homework := range studentHomeworks {
			if homework.Lesson.ID == lesson.ID {
				counted = true
				if homework.Verify {
					reportMessage += fmt.Sprintf("%v - %v\n", iconHomeworkVerified, lesson.Title)
				} else {
					reportMessage += fmt.Sprintf("%v - %v\n", iconHomeworkNotVerified, lesson.Title)
				}
			}
		}

		if !counted {
			reportMessage += fmt.Sprintf("%v - %v\n", iconHomeworkNotProvided, lesson.Title)
		}
	}

	srv.logger.Debug("report sent")
	return c.Reply(reportMessage, &telebot.SendOptions{ParseMode: "HTML"})
}
