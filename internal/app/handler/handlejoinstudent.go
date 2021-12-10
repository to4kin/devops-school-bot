package handler

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (handler *Handler) handleJoinStudent(c telebot.Context) error {
	if c.Message().Private() {
		return handler.editOrReply(c, fmt.Sprintf(helper.ErrWrongChatType, "SCHOOL"), nil)
	}

	handler.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Info("get account from database by telegram_id")
	account, err := handler.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		if err == store.ErrRecordNotFound {
			handler.logger.Info("account not found, will create a new one")
			account = &model.Account{
				Created:    time.Now(),
				TelegramID: int64(c.Sender().ID),
				FirstName:  c.Sender().FirstName,
				LastName:   c.Sender().LastName,
				Username:   c.Sender().Username,
				Superuser:  false,
			}

			if err := handler.store.Account().Create(account); err != nil {
				handler.logger.Error(err)
				return handler.editOrReply(c, helper.ErrInternal, nil)
			}

			handler.logger.WithFields(account.LogrusFields()).Info("account created")
		} else {
			handler.logger.Error(err)
			return handler.editOrReply(c, helper.ErrInternal, nil)
		}
	} else {
		handler.logger.WithFields(account.LogrusFields()).Info("account found")
	}

	handler.logger.WithFields(logrus.Fields{
		"chat_id": c.Message().Chat.ID,
	}).Info("get school by chat_id")
	school, err := handler.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		handler.logger.Error(err)
		if err == store.ErrRecordNotFound {
			return handler.editOrReply(c, helper.ErrSchoolNotStarted, nil)
		}

		return handler.editOrReply(c, helper.ErrInternal, nil)
	}
	handler.logger.WithFields(school.LogrusFields()).Info("school found")

	if !school.Active {
		handler.logger.WithFields(school.LogrusFields()).Info("school finished")
		return handler.editOrReply(c, helper.ErrSchoolNotStarted, nil)
	}

	handler.logger.WithFields(logrus.Fields{
		"account_id": account.ID,
		"school_id":  school.ID,
	}).Info("get student from database by account_id and school_id")
	student, err := handler.store.Student().FindByAccountIDSchoolID(account.ID, school.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			handler.logger.Info("student not found, will create a new one")
			student := &model.Student{
				Created:    time.Now(),
				Account:    account,
				School:     school,
				Active:     true,
				FullCourse: true,
			}

			if err := handler.store.Student().Create(student); err != nil {
				handler.logger.Error(err)
				return handler.editOrReply(c, helper.ErrInternal, nil)
			}

			handler.logger.WithFields(student.LogrusFields()).Info("student created")
			reportMessage := fmt.Sprintf(helper.MsgWelcomeToSchool, school.Title, student.GetType())
			if student.FullCourse {
				reportMessage += "\n\n" + helper.SysStudentGuide
			} else {
				reportMessage += "\n\n" + helper.SysListenerGuide
			}
			return handler.editOrReply(c, reportMessage, nil)
		}

		handler.logger.Error(err)
		return handler.editOrReply(c, helper.ErrInternal, nil)
	}
	handler.logger.WithFields(student.LogrusFields()).Info("student exist")
	return handler.editOrReply(c, fmt.Sprintf(helper.MsgUserAlreadyJoined, school.Title, student.GetType()), nil)
}
