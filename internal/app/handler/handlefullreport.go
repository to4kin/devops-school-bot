package handler

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (handler *Handler) handleFullReport(c telebot.Context) error {
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

	hlpr := helper.NewHelper(handler.store, handler.logger)

	if c.Message().Private() {
		callback := &model.Callback{
			ID:          0,
			Type:        "school",
			Command:     "full_report",
			ListCommand: "full_report",
		}

		replyMessage, replyMarkup, err := hlpr.GetSchoolsList(callback)
		if err != nil {
			handler.logger.Error(err)
			return handler.editOrReply(c, helper.ErrInternal, nil)
		}

		return handler.editOrReply(c, replyMessage, replyMarkup)
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

	reportMessage, err := hlpr.GetFullReport(school)
	if err != nil && err != store.ErrRecordNotFound {
		handler.logger.Error(err)
		return handler.editOrReply(c, helper.ErrInternal, nil)
	}

	if err == store.ErrRecordNotFound {
		reportMessage = helper.ErrReportNotFound
	}

	return handler.editOrReply(c, fmt.Sprintf("School <b>%v</b>\n\n%v", school.Title, reportMessage), nil)
}
