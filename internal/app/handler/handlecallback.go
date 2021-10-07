package handler

import (
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *Handler) handleCallback(c telebot.Context) error {
	srv.logger.WithFields(logrus.Fields{
		"callback_data": c.Callback().Data[1:],
	}).Debug("handle callback")

	srv.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Debug("get account from database by telegram_id")
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		srv.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}
	srv.logger.WithFields(account.LogrusFields()).Debug("account found")

	if !account.Superuser {
		srv.logger.WithFields(account.LogrusFields()).Debug("account has insufficient permissions")
		return c.EditOrReply(helper.ErrInsufficientPermissions, &telebot.SendOptions{ParseMode: "HTML"})
	}

	callback := &model.Callback{}
	callback.Unmarshal(c.Callback().Data[1:])

	srv.logger.WithFields(logrus.Fields{
		"callback_id":           callback.ID,
		"callback_type":         callback.Type,
		"callback_command":      callback.Command,
		"callback_list_command": callback.ListCommand,
	}).Debug("parse callback data")

	replyMessage := ""
	replyMarkup := &telebot.ReplyMarkup{}

	hlpr := helper.NewHelper(srv.store, srv.logger)

	switch callback.Type {
	case "school":
		switch callback.Command {
		case "schools_list", "next", "previous":
			replyMessage, replyMarkup, err = hlpr.GetSchoolsList(callback)

		case "get":
			replyMessage, replyMarkup, err = hlpr.GetSchool(callback)

		case "start":
			replyMessage, replyMarkup, err = hlpr.StartSchool(callback)

		case "stop":
			replyMessage, replyMarkup, err = hlpr.StopSchool(callback)

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

		case "block":
			replyMessage, replyMarkup, err = hlpr.BlockStudent(callback)

		case "unblock":
			replyMessage, replyMarkup, err = hlpr.UnblockStudent(callback)

		case "set_student":
			replyMessage, replyMarkup, err = hlpr.SetStudent(callback)

		case "set_listener":
			replyMessage, replyMarkup, err = hlpr.SetListener(callback)

		}
	case "homework":
		switch callback.Command {
		case "homeworks_list", "next", "previous":
			replyMessage, replyMarkup, err = hlpr.GetHomeworksList(callback)

		case "get":
			replyMessage, replyMarkup, err = hlpr.GetHomework(callback)

		case "disable_homework":
			replyMessage, replyMarkup, err = hlpr.DisableHomework(callback)

		case "enable_homework":
			replyMessage, replyMarkup, err = hlpr.EnableHomework(callback)

		}
	}

	if err != nil {
		srv.logger.Error(err)
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	if replyMessage == "" {
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	return c.EditOrReply(replyMessage, &telebot.SendOptions{ParseMode: "HTML"}, replyMarkup)
}
