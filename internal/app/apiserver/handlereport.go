package apiserver

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleReport(c telebot.Context) error {
	if c.Message().Private() {
		return nil
	}

	logrus.Debug("get active school")
	school, err := srv.store.School().FindActive()
	if err != nil {
		logrus.Error(err)

		if err == store.ErrRecordNotFound {
			return c.Reply(MsgNoActiveSchool, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return nil
	}
	logrus.Debug(school.ToString())

	logrus.Debug("get account from database by telegram_id: ", c.Sender().ID)
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		logrus.Error(err)

		if err == store.ErrRecordNotFound {
			return c.Reply(MsgUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return nil
	}
	logrus.Debug(account.ToString())

	logrus.Debug("get student from database by account_id: ", account.ID, " and school_id: ", school.ID)
	student, err := srv.store.Student().FindByAccountIDSchoolID(account.ID, school.ID)
	if err != nil {
		logrus.Error(err)

		if err == store.ErrRecordNotFound {
			return c.Reply(MsgUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return nil
	}
	logrus.Debug(student.ToString())

	logrus.Debug("get student homeworks from database by student_id: ", student.ID)
	studentHomeworks, err := srv.store.Homework().FindByStudentID(student.ID)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	logrus.Debug("get all homeworks from database by school_id: ", school.ID)
	allLessons, err := srv.store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	reportMessage := fmt.Sprintf("Hello, @%v!\n\nYour progress in <b>DevOps School %v</b>:\n", account.Username, school.Title)
	for _, lesson := range allLessons {
		if err != nil {
			logrus.Error(err)
			return nil
		}

		counted := false
		for _, homework := range studentHomeworks {
			if homework.Lesson.ID == lesson.ID {
				counted = true
				if homework.Verify {
					reportMessage += fmt.Sprintf("✔️ - %v\n", lesson.Title)
				} else {
					reportMessage += fmt.Sprintf("⭕ - %v\n", lesson.Title)
				}
			}
		}

		if !counted {
			reportMessage += fmt.Sprintf("❌ - %v\n", lesson.Title)
		}

		logrus.Debug(lesson.ToString())
	}

	return c.Reply(reportMessage, &telebot.SendOptions{ParseMode: "HTML"})
}
