package apiserver

import (
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	homeworkHashtag string = "#homework"
)

func (srv *server) handleOnText(c telebot.Context) error {
	if c.Message().Private() {
		return nil
	}

	text := strings.ToLower(c.Message().Text)

	if strings.Contains(text, homeworkHashtag) {
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
				logrus.Error(err)
				return c.Reply(msgUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
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
				return c.Reply(msgUserNotJoined, &telebot.SendOptions{ParseMode: "HTML"})
			}

			logrus.Error(err)
			return nil
		}
		logrus.Debug(student.ToString())

		for _, entity := range c.Message().Entities {
			switch entity.Type {
			case "hashtag":
				hashtag := text[entity.Offset : entity.Offset+entity.Length]
				if hashtag == homeworkHashtag {
					logrus.Debug("homework hashtag skipped: ", homeworkHashtag)
					continue
				}

				logrus.Debug("get lesson from database by title: ", hashtag)
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
							MessageID: int64(c.Message().ID),
							Verify:    false,
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
