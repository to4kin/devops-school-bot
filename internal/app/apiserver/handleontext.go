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
		srv.logger.WithFields(logrus.Fields{
			"chat_id": c.Message().Chat.ID,
		}).Debug("get school by chat_id")
		school, err := srv.store.School().FindByChatID(c.Message().Chat.ID)
		if err != nil {
			srv.logger.Error(err)
			return nil
		}
		srv.logger.WithFields(school.LogrusFields()).Debug("school found")

		if school.Finished {
			srv.logger.WithFields(school.LogrusFields()).Debug("school already finished")
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

		srv.logger.WithFields(logrus.Fields{
			"account_id": account.ID,
			"school_id":  school.ID,
		}).Debug("get student from database by account_id and school_id")
		student, err := srv.store.Student().FindByAccountIDSchoolID(account.ID, school.ID)
		if err != nil {
			srv.logger.Error(err)
			return nil
		}
		srv.logger.WithFields(student.LogrusFields()).Debug("student found")

		for _, entity := range c.Message().Entities {
			switch entity.Type {
			case "hashtag":
				hashtag := text[entity.Offset : entity.Offset+entity.Length]
				if hashtag == homeworkHashtag {
					srv.logger.WithFields(logrus.Fields{
						"hashtag": homeworkHashtag,
					}).Debug("homework hashtag skipped")
					continue
				}

				srv.logger.WithFields(logrus.Fields{
					"title": hashtag,
				}).Debug("get lesson from database by title")
				lesson, err := srv.store.Lesson().FindByTitle(hashtag)
				if err != nil {
					if err == store.ErrRecordNotFound {
						srv.logger.Debug("lesson not found, will create a new one")
						lesson = &model.Lesson{
							Title: hashtag,
						}

						if err := srv.store.Lesson().Create(lesson); err != nil {
							srv.logger.Error(err)
							return nil
						}
					} else {
						srv.logger.Error(err)
						return nil
					}
				}
				srv.logger.WithFields(lesson.LogrusFields()).Debug("lesson found")

				srv.logger.WithFields(logrus.Fields{
					"student_id": student.ID,
					"lesson_id":  lesson.ID,
				}).Debug("get homework from database by student_id and lesson_id")
				homework, err := srv.store.Homework().FindByStudentIDLessonID(student.ID, lesson.ID)
				if err != nil {
					if err == store.ErrRecordNotFound {
						srv.logger.Debug("homework not found, will create a new one")
						homework = &model.Homework{
							Student:   student,
							Lesson:    lesson,
							MessageID: int64(c.Message().ID),
							Verify:    false,
						}

						if err := srv.store.Homework().Create(homework); err != nil {
							srv.logger.Error(err)
							return nil
						}
					} else {
						srv.logger.Error(err)
						return nil
					}
				}
				srv.logger.WithFields(homework.LogrusFields()).Debug("homework found")
			}
		}
	}

	return nil
}
