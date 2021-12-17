package handler

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	homeworkHashtag string = "#homework"
)

func (handler *Handler) handleOnText(c telebot.Context) error {
	if c.Message().Private() {
		handler.logger.Info("message is private, exiting...")
		return nil
	}

	text := ""
	var entities []telebot.MessageEntity

	// NOTE: multimedia messages have CaptionEntities and Caption, not Entities and Text
	// Also, if the message was edited - Entities and Text are under c.Update().EditedMessage
	if c.Update().EditedMessage != nil {
		text = strings.ToLower(c.Update().EditedMessage.Text)
		entities = c.Update().EditedMessage.Entities
	} else if c.Message().Text != "" {
		text = strings.ToLower(c.Message().Text)
		entities = c.Message().Entities
	} else {
		text = strings.ToLower(c.Message().Caption)
		entities = c.Message().CaptionEntities
	}

	// NOTE: Telegram uses UTF16 encoding for calculating Length and Offset
	// so when just ASCII text is used there are no problems at all, since ASCII always uses 1 byte for each character.
	//
	utfEncodedString := utf16.Encode([]rune(text))

	if strings.Contains(text, homeworkHashtag) {
		handler.logger.WithFields(logrus.Fields{
			"chat_id": c.Message().Chat.ID,
		}).Info("get school by chat_id")
		school, err := handler.store.School().FindByChatID(c.Message().Chat.ID)
		if err != nil {
			handler.logger.Error(err)
			return nil
		}
		handler.logger.WithFields(school.LogrusFields()).Info("school found")

		if !school.Active {
			handler.logger.WithFields(school.LogrusFields()).Info("school already finished")
			return nil
		}

		handler.logger.WithFields(logrus.Fields{
			"telegram_id": c.Sender().ID,
		}).Info("get account from database by telegram_id")
		account, err := handler.store.Account().FindByTelegramID(int64(c.Sender().ID))
		if err != nil {
			handler.logger.Error(err)
			return nil
		}
		handler.logger.WithFields(account.LogrusFields()).Info("account found")

		handler.logger.WithFields(logrus.Fields{
			"account_id": account.ID,
			"school_id":  school.ID,
		}).Info("get student from database by account_id and school_id")
		student, err := handler.store.Student().FindByAccountIDSchoolID(account.ID, school.ID)
		if err != nil {
			handler.logger.Error(err)
			return nil
		}
		handler.logger.WithFields(student.LogrusFields()).Info("student found")

		if !student.Active {
			handler.logger.WithFields(student.LogrusFields()).Info("student is not active")
			return nil
		}

		// NOTE: if c.Update().EditedMessage is not null - the message was edited
		if c.Update().EditedMessage != nil {
			handler.logger.WithFields(logrus.Fields{
				"message_id": c.Update().EditedMessage.ID,
			}).Info("message was edited, need to delete old homework first")
			if err := handler.store.Homework().DeleteByMessageIDStudentID(int64(c.Update().EditedMessage.ID), student.ID); err != nil {
				if err == store.ErrRecordNotFound {
					handler.logger.Info(err)
				} else {
					handler.logger.Error(err)
					return nil
				}
			} else {
				handler.logger.WithFields(logrus.Fields{
					"message_id": c.Update().EditedMessage.ID,
				}).Info("old message was deleted")
			}
		}

		for _, entity := range entities {
			switch entity.Type {
			case "hashtag":
				// NOTE: Telegram uses UTF16 encoding for calculating Length and Offset
				// so when just ASCII text is used there are no problems at all, since ASCII always uses 1 byte for each character.
				//
				//hashtag := text[entity.Offset : entity.Offset+entity.Length]
				runeString := utf16.Decode(utfEncodedString[entity.Offset : entity.Offset+entity.Length])
				hashtag := string(runeString)
				if hashtag == homeworkHashtag {
					handler.logger.WithFields(logrus.Fields{
						"hashtag": homeworkHashtag,
					}).Info("homework hashtag skipped")
					continue
				}

				handler.logger.WithFields(logrus.Fields{
					"hashtag": hashtag,
				}).Info("hashtag found")

				reg, err := regexp.Compile("[^a-zA-Z]+")
				if err != nil {
					handler.logger.Error(err)
					return nil
				}
				moduleTitle := reg.ReplaceAllString(hashtag, "")

				handler.logger.WithFields(logrus.Fields{
					"title": moduleTitle,
				}).Info("get module from database by title")
				module, err := handler.store.Module().FindByTitle(moduleTitle)
				if err != nil {
					if err == store.ErrRecordNotFound {
						handler.logger.Info("module not found, will create a new one")
						module = &model.Module{
							Title: moduleTitle,
						}

						if err := handler.store.Module().Create(module); err != nil {
							handler.logger.Error(err)
							continue
						}

						handler.logger.WithFields(module.LogrusFields()).Info("module created")
					} else {
						handler.logger.Error(err)
						continue
					}
				} else {
					handler.logger.WithFields(module.LogrusFields()).Info("module found")
				}

				handler.logger.WithFields(logrus.Fields{
					"title": hashtag,
				}).Info("get lesson from database by title")
				lesson, err := handler.store.Lesson().FindByTitle(hashtag)
				if err != nil {
					if err == store.ErrRecordNotFound {
						handler.logger.Info("lesson not found, will create a new one")
						lesson = &model.Lesson{
							Title:  hashtag,
							Module: module,
						}

						if err := handler.store.Lesson().Create(lesson); err != nil {
							handler.logger.Error(err)
							continue
						}

						handler.logger.WithFields(lesson.LogrusFields()).Info("lesson created")
					} else {
						handler.logger.Error(err)
						continue
					}
				} else {
					handler.logger.WithFields(lesson.LogrusFields()).Info("lesson found")
				}

				handler.logger.WithFields(logrus.Fields{
					"student_id": student.ID,
					"lesson_id":  lesson.ID,
				}).Info("get homework from database by student_id and lesson_id")
				homework, err := handler.store.Homework().FindByStudentIDLessonID(student.ID, lesson.ID)
				if err != nil {
					if err == store.ErrRecordNotFound {
						handler.logger.Info("homework not found, will create a new one")
						homework = &model.Homework{
							Created:   time.Now(),
							Student:   student,
							Lesson:    lesson,
							MessageID: int64(c.Message().ID),
							Verify:    false,
							Active:    true,
						}

						if err := handler.store.Homework().Create(homework); err != nil {
							handler.logger.Error(err)
							continue
						}

						handler.logger.WithFields(homework.LogrusFields()).Info("homework created")
						continue
					} else {
						handler.logger.Error(err)
						continue
					}
				}
				handler.logger.WithFields(homework.LogrusFields()).Info("homework found")
			}
		}
	}

	return nil
}
