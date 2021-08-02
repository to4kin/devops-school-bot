package apiserver

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleReport(c telebot.Context) error {
	logger := logrus.WithFields(logrus.Fields{
		"handler": "report",
	})

	if c.Message().Private() {
		return nil
	}

	logger.Debug("get school by chat_id: ", c.Message().Chat.ID)
	school, err := srv.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		logger.Error(err)

		if err == store.ErrRecordNotFound {
			return c.Reply(msgSchoolNotFound, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return nil
	}
	logger.Debug(school.ToString())

	logger.Debug("get account from database by telegram_id: ", c.Sender().ID)
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		logger.Error(err)

		if err == store.ErrRecordNotFound {
			return c.Reply(msgUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return nil
	}
	logger.Debug(account.ToString())

	logger.Debug("get student from database by account_id: ", account.ID, " and school_id: ", school.ID)
	student, err := srv.store.Student().FindByAccountIDSchoolID(account.ID, school.ID)
	if err != nil {
		logger.Error(err)

		if err == store.ErrRecordNotFound {
			return c.Reply(msgUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return nil
	}
	logger.Debug(student.ToString())

	logger.Debug("get student homeworks from database by student_id: ", student.ID)
	studentHomeworks, err := srv.store.Homework().FindByStudentID(student.ID)
	if err != nil {
		logger.Error(err)
		return nil
	}
	logger.Debug("found homeworks: ", len(studentHomeworks))

	logger.Debug("get all lessons from database by school_id: ", school.ID)
	allLessons, err := srv.store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		logger.Error(err)
		return nil
	}
	logger.Debug("found lessons: ", len(allLessons))

	reportMessage := fmt.Sprintf(msgHomeworkReport, account.Username, school.Title)
	for _, lesson := range allLessons {
		if err != nil {
			logger.Error(err)
			return nil
		}

		counted := false
		for _, homework := range studentHomeworks {
			if homework.Lesson.ID == lesson.ID {
				counted = true
				if homework.Verify {
					reportMessage += fmt.Sprintf("%v - %v\n", msgHomeworkVerified, lesson.Title)
				} else {
					reportMessage += fmt.Sprintf("%v - %v\n", msgHomeworkNotVerified, lesson.Title)
				}
			}
		}

		if !counted {
			reportMessage += fmt.Sprintf("%v - %v\n", msgHomeworkNotProvided, lesson.Title)
		}
	}

	return c.Reply(reportMessage, &telebot.SendOptions{ParseMode: "HTML"})
}
