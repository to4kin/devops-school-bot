package handler

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gopkg.in/tucnak/telebot.v3"
)

// HandleCron ...
func (handler *Handler) HandleCron(fullreport bool) {
	handler.logger.WithFields(logrus.Fields{
		"fullreport": fullreport,
	}).Debug("cron job started")

	schools, err := handler.store.School().FindByActive(true)
	if err != nil {
		handler.logger.Error(err)
	}

	hlpr := helper.NewHelper(handler.store, handler.logger)
	for _, school := range schools {
		schoolChat, err := handler.bot.ChatByID(school.ChatID)
		if err != nil {
			handler.logger.Error(err)
			continue
		}

		var reportMessage string
		if fullreport {
			reportMessage, err = hlpr.GetFullReport(school)
		} else {
			reportMessage, err = hlpr.GetReport(school)
		}

		if err != nil {
			handler.logger.Error(err)
			continue
		}

		handler.logger.WithFields(logrus.Fields{
			"school_title": school.Title,
			"fullreport":   fullreport,
		}).Debug("sent report by cron")
		handler.bot.Send(schoolChat, fmt.Sprintf("School <b>%v</b>\n\n%v", school.Title, reportMessage), &telebot.SendOptions{ParseMode: "HTML"})
	}
}
