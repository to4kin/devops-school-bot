package apiserver

import (
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleOnText(c telebot.Context) error {
	if c.Message().Private() {
		return nil
	}

	text := strings.ToLower(c.Message().Text)

	if strings.Contains(text, "#homework") {
		logrus.Debug("get active school")
		school, err := srv.store.School().FindActive()
		if err != nil {
			if err == store.ErrRecordNotFound {
				logrus.Error(err)
				return c.Reply(MsgNoActiveSchool, &telebot.SendOptions{ParseMode: "HTML"})
			}

			logrus.Error(err)
			return nil
		}
		logrus.Debug(school.ToString())

		logrus.Debug("get account from database by telegram_id: ", c.Sender().ID)
		account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
		if err != nil {
			if err == store.ErrRecordNotFound {
				logrus.Error(err)
				return c.Reply(MsgUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
			}

			logrus.Error(err)
			return nil
		}
		logrus.Debug(account.ToString())

		logrus.Debug("get student from database by account_id: ", account.ID, " and school_id: ", school.ID)
		student, err := srv.store.Student().FindByAccountIDSchoolID(account.ID, school.ID)
		if err != nil {
			if err == store.ErrRecordNotFound {
				logrus.Error(err)
				return c.Reply(MsgUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
			}

			logrus.Error(err)
			return nil
		}
		logrus.Debug(student.ToString())

		for _, entity := range c.Message().Entities {
			switch entity.Type {
			case "hashtag":
				hashtag := text[entity.Offset : entity.Offset+entity.Length]
				if hashtag == "#homework" {
					logrus.Debug("special hashtag skipped")
					break
				}

				logrus.Debug("get lesson from database by title: ")
				lesson, err := srv.store.Lesson().FindByTitle(hashtag)
				if err != nil {
					if err == store.ErrRecordNotFound {
						logrus.Debug("lesson not found, will create a new one")
						lesson = &model.Lesson{
							Title: hashtag,
						}

						if err := srv.store.Lesson().Create(lesson); err != nil {
							logrus.Error(err)
							return nil
						}
					} else {
						logrus.Error(err)
						return nil
					}
				}
				logrus.Debug(lesson.ToString())

				logrus.Debug("get homework from database by student_id: ", student.ID, " and lesson_id: ", lesson.ID)
				homework, err := srv.store.Homework().FindByStudentIDLessonID(student.ID, lesson.ID)
				if err != nil {
					if err == store.ErrRecordNotFound {
						logrus.Debug("homework not found, will create a new one")
						homework = &model.Homework{
							Student:   student,
							Lesson:    lesson,
							ChatID:    int64(c.Message().Chat.ID),
							MessageID: int64(c.Message().ID),
						}

						if err := srv.store.Homework().Create(homework); err != nil {
							logrus.Error(err)
							return nil
						}
					} else {
						logrus.Error(err)
						return nil
					}
				}
				logrus.Debug(homework.ToString())
			}
		}
	}

	return nil
}
