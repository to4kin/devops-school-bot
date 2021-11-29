package handler

import (
	"fmt"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

func (handler *Handler) handleHomework(c telebot.Context) error {
	if !c.Message().Private() {
		return c.EditOrReply(fmt.Sprintf(helper.ErrWrongChatType, "PRIVATE"), &telebot.SendOptions{ParseMode: "HTML"})
	}

	hlpr := helper.NewHelper(handler.store, handler.logger)
	callback := &model.Callback{
		ID:          0,
		Type:        "school",
		Command:     "homeworks",
		ListCommand: "homeworks",
	}

	replyMessage, replyMarkup, err := hlpr.GetSchoolsList(callback)
	if err != nil {
		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	return c.EditOrReply(replyMessage, replyMarkup)
}
