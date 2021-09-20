package apiserver

import (
	"regexp"
	"strings"
	"time"

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

	text := ""
	var entities []telebot.MessageEntity

	if c.Message().Text != "" {
		text = strings.ToLower(c.Message().Text)
		entities = c.Message().Entities
	} else {
		text = strings.ToLower(c.Message().Caption)
		entities = c.Message().CaptionEntities
	}

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

		if !school.Active {
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

		if !student.Active {
			srv.logger.WithFields(student.LogrusFields()).Debug("student is not active")
			return nil
		}

		for _, entity := range entities {
			switch entity.Type {
			case "hashtag":
				hashtag := text[entity.Offset : entity.Offset+entity.Length]
				if hashtag == homeworkHashtag {
					srv.logger.WithFields(logrus.Fields{
						"hashtag": homeworkHashtag,
					}).Debug("homework hashtag skipped")
					continue
				}

				reg, err := regexp.Compile("[^a-zA-Z0-9]+")
				if err != nil {
					srv.logger.Error(err)
					return nil
				}
				moduleTitle := reg.ReplaceAllString(hashtag, "")

				srv.logger.WithFields(logrus.Fields{
					"title": moduleTitle,
				}).Debug("get module from database by title")
				module, err := srv.store.Module().FindByTitle(moduleTitle)
				if err != nil {
					if err == store.ErrRecordNotFound {
						srv.logger.Debug("module not found, will create a new one")
						module = &model.Module{
							Title: moduleTitle,
						}

						if err := srv.store.Module().Create(module); err != nil {
							srv.logger.Error(err)
							continue
						}

						srv.logger.WithFields(module.LogrusFields()).Debug("module created")
					} else {
						srv.logger.Error(err)
						continue
					}
				} else {
					srv.logger.WithFields(module.LogrusFields()).Debug("module found")
				}

				srv.logger.WithFields(logrus.Fields{
					"title": hashtag,
				}).Debug("get lesson from database by title")
				lesson, err := srv.store.Lesson().FindByTitle(hashtag)
				if err != nil {
					if err == store.ErrRecordNotFound {
						srv.logger.Debug("lesson not found, will create a new one")
						lesson = &model.Lesson{
							Title:  hashtag,
							Module: module,
						}

						if err := srv.store.Lesson().Create(lesson); err != nil {
							srv.logger.Error(err)
							continue
						}

						srv.logger.WithFields(lesson.LogrusFields()).Debug("lesson created")
					} else {
						srv.logger.Error(err)
						continue
					}
				} else {
					srv.logger.WithFields(lesson.LogrusFields()).Debug("lesson found")
				}

				srv.logger.WithFields(logrus.Fields{
					"student_id": student.ID,
					"lesson_id":  lesson.ID,
				}).Debug("get homework from database by student_id and lesson_id")
				homework, err := srv.store.Homework().FindByStudentIDLessonID(student.ID, lesson.ID)
				if err != nil {
					if err == store.ErrRecordNotFound {
						srv.logger.Debug("homework not found, will create a new one")
						homework = &model.Homework{
							Created:   time.Now(),
							Student:   student,
							Lesson:    lesson,
							MessageID: int64(c.Message().ID),
							Verify:    false,
						}

						if err := srv.store.Homework().Create(homework); err != nil {
							srv.logger.Error(err)
							continue
						}

						srv.logger.WithFields(homework.LogrusFields()).Debug("homework created")
						continue
					} else {
						srv.logger.Error(err)
						continue
					}
				}
				srv.logger.WithFields(homework.LogrusFields()).Debug("homework found")
			}
		}
	}

	return nil
}
