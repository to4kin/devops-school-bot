package handler

import (
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (handler *Handler) handleHomework(c telebot.Context) error {
	hlpr := helper.NewHelper(handler.store, handler.logger)

	if c.Message().Private() {
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

	handler.logger.WithFields(logrus.Fields{
		"chat_id": c.Message().Chat.ID,
	}).Info("get school by chat_id")
	school, err := handler.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		handler.logger.Error(err)

		if err == store.ErrRecordNotFound {
			return c.EditOrReply(helper.ErrSchoolNotStarted, &telebot.SendOptions{ParseMode: "HTML"})
		}

		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	handler.logger.WithFields(school.LogrusFields()).Info("school found")

	reportMessage, err := hlpr.GetLessonsReport(school)
	if err != nil && err != store.ErrRecordNotFound {
		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	if err == store.ErrRecordNotFound {
		reportMessage = helper.ErrReportNotFound
	}

	handler.logger.Info("report sent")
	return c.EditOrReply(reportMessage, &telebot.SendOptions{ParseMode: "HTML"})
}
