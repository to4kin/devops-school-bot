package handler

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gopkg.in/tucnak/telebot.v3"
)

func (handler *Handler) handleCallback(c telebot.Context) error {
	handler.logger.WithFields(logrus.Fields{
		"callback_data": c.Callback().Data[1:],
	}).Info("handle callback")

	callbackID, err := strconv.ParseInt(c.Callback().Data[1:], 10, 64)
	if err != nil {
		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrOldCallbackData, &telebot.SendOptions{ParseMode: "HTML"})
	}

	handler.logger.WithFields(logrus.Fields{
		"callback_id": callbackID,
	}).Info("callback parsed")

	// handler.logger.WithFields(logrus.Fields{
	// 	"telegram_id": c.Sender().ID,
	// }).Info("get account from database by telegram_id")
	// account, err := handler.store.Account().FindByTelegramID(int64(c.Sender().ID))
	// if err != nil {
	// 	handler.logger.Error(err)
	// 	return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	// }
	// handler.logger.WithFields(account.LogrusFields()).Info("account found")

	//if !account.Superuser {
	//	handler.logger.WithFields(account.LogrusFields()).Info("account has insufficient permissions")
	//	return c.EditOrReply(helper.ErrInsufficientPermissions, &telebot.SendOptions{ParseMode: "HTML"})
	//}

	callback, err := handler.store.Callback().FindByID(callbackID)
	if err != nil {
		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	handler.logger.WithFields(callback.LogrusFields()).Info("callback found")

	replyMessage := ""
	replyMarkup := &telebot.ReplyMarkup{}

	hlpr := helper.NewHelper(handler.store, handler.logger)

	switch callback.Type {
	case "school":
		switch callback.Command {
		case "schools_list", "next", "previous":
			replyMessage, replyMarkup, err = hlpr.GetSchoolsList(callback)

		case "get":
			replyMessage, replyMarkup, err = hlpr.GetSchool(callback)

		case "update_status":
			replyMessage, replyMarkup, err = hlpr.UpdateSchoolStatus(callback)

		case "report":
			replyMessage, replyMarkup, err = hlpr.ReportSchool(callback)

		case "full_report":
			replyMessage, replyMarkup, err = hlpr.FullReportSchool(callback)

		case "homeworks":
			replyMessage, replyMarkup, err = hlpr.GetSchoolHomeworks(callback)

		}
	case "account":
		switch callback.Command {
		case "accounts_list", "next", "previous":
			replyMessage, replyMarkup, err = hlpr.GetUsersList(callback)

		case "get":
			replyMessage, replyMarkup, err = hlpr.GetUser(callback, c.Sender())

		case "update":
			replyMessage, replyMarkup, err = hlpr.UpdateUser(callback, c.Sender())

		case "set_superuser":
			replyMessage, replyMarkup, err = hlpr.SetSuperuser(callback)

		case "unset_superuser":
			replyMessage, replyMarkup, err = hlpr.UnsetSuperuser(callback)

		}
	case "student":
		switch callback.Command {
		case "students_list", "next", "previous":
			replyMessage, replyMarkup, err = hlpr.GetStudentsList(callback)

		case "get":
			replyMessage, replyMarkup, err = hlpr.GetStudent(callback)

		case "update_status":
			replyMessage, replyMarkup, err = hlpr.UpdateStudentStatus(callback)

		case "update_type":
			replyMessage, replyMarkup, err = hlpr.UpdateStudentType(callback)

		}
	case "homework":
		switch callback.Command {
		case "homeworks_list", "next", "previous":
			replyMessage, replyMarkup, err = hlpr.GetHomeworksList(callback)

		case "get":
			replyMessage, replyMarkup, err = hlpr.GetHomework(callback)

		case "update_status":
			replyMessage, replyMarkup, err = hlpr.UpdateHomeworkStatus(callback)

		}
	}

	if err != nil {
		handler.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	if replyMessage == "" {
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	return c.EditOrReply(replyMessage, &telebot.SendOptions{ParseMode: "HTML"}, replyMarkup)
}
