package handler

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (handler *Handler) handleSchools(c telebot.Context) error {
	if !c.Message().Private() {
		return handler.editOrReply(c, fmt.Sprintf(helper.ErrWrongChatType, "PRIVATE"), nil)
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

	callback := &model.Callback{
		ID:          0,
		Type:        "school",
		Command:     "get",
		ListCommand: "get",
	}

	hlpr := helper.NewHelper(handler.store, handler.logger)
	replyMessage, replyMarkup, err := hlpr.GetSchoolsList(callback)
	if err != nil {
		handler.logger.Error(err)
		return handler.editOrReply(c, helper.ErrInternal, nil)
	}

	return handler.editOrReply(c, replyMessage, replyMarkup)
}
