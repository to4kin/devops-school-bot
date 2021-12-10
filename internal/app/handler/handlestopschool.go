package handler

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	msgSchoolFinished string = `school <b>%v</b> finished`
)

func (handler *Handler) handleStopSchool(c telebot.Context) error {
	if c.Message().Private() {
		return handler.editOrReply(c, fmt.Sprintf(helper.ErrWrongChatType, "SCHOOL"), nil)
	}

	handler.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Info("get account from database by telegram_id")
	account, err := handler.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		handler.logger.Error(err)

		if err == store.ErrRecordNotFound {
			return handler.editOrReply(c, helper.ErrInsufficientPermissions, nil)
		}

		return handler.editOrReply(c, helper.ErrInternal, nil)
	}
	handler.logger.WithFields(account.LogrusFields()).Info("account found")

	if !account.Superuser {
		handler.logger.WithFields(account.LogrusFields()).Info("account has insufficient permissions")
		return handler.editOrReply(c, helper.ErrInsufficientPermissions, nil)
	}

	handler.logger.WithFields(logrus.Fields{
		"chat_id": c.Message().Chat.ID,
	}).Info("get school by chat_id")
	school, err := handler.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		handler.logger.Error(err)
		return handler.editOrReply(c, helper.ErrSchoolNotStarted, nil)
	}
	handler.logger.WithFields(school.LogrusFields()).Info("school found")

	if !school.Active {
		handler.logger.WithFields(school.LogrusFields()).Info("school already finished")
		return handler.editOrReply(c, fmt.Sprintf(msgSchoolFinished, school.Title), nil)
	}

	school.Active = false
	if err := handler.store.School().Update(school); err != nil {
		handler.logger.Error(err)
		return handler.editOrReply(c, helper.ErrInternal, nil)
	}
	handler.logger.WithFields(school.LogrusFields()).Info("school finished")
	return handler.editOrReply(c, fmt.Sprintf(msgSchoolFinished, school.Title), nil)
}
