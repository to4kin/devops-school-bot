package apiserver

import (
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/helper"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleCallback(c telebot.Context) error {
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

	switch callback.Type {
	case "school":
		switch callback.Command {
		case "schools_list", "next", "previous":
			replyMessage, replyMarkup, err = helper.GetSchoolsList(srv.store, callback)

		case "get":
			replyMessage, replyMarkup, err = helper.GetSchool(srv.store, callback)

		case "start":
			replyMessage, replyMarkup, err = helper.StartSchool(srv.store, callback)

		case "stop":
			replyMessage, replyMarkup, err = helper.StopSchool(srv.store, callback)

		case "report":
			replyMessage, replyMarkup, err = helper.ReportSchool(srv.store, callback)

		case "full_report":
			replyMessage, replyMarkup, err = helper.FullReportSchool(srv.store, callback)

		case "homeworks":
			replyMessage, replyMarkup, err = helper.GetSchoolHomeworks(srv.store, callback)

		}
	case "account":
		switch callback.Command {
		case "accounts_list", "next", "previous":
			replyMessage, replyMarkup, err = helper.GetUsersList(srv.store, callback)

		case "get":
			replyMessage, replyMarkup, err = helper.GetUser(srv.store, callback, c.Sender())

		case "update":
			replyMessage, replyMarkup, err = helper.UpdateUser(srv.store, callback, c.Sender())

		case "set_superuser":
			replyMessage, replyMarkup, err = helper.SetSuperuser(srv.store, callback)

		case "unset_superuser":
			replyMessage, replyMarkup, err = helper.UnsetSuperuser(srv.store, callback)

		}
	case "student":
		switch callback.Command {
		case "students_list", "next", "previous":
			replyMessage, replyMarkup, err = helper.GetStudentsList(srv.store, callback)

		case "get":
			replyMessage, replyMarkup, err = helper.GetStudent(srv.store, callback)

		case "block":
			replyMessage, replyMarkup, err = helper.BlockStudent(srv.store, callback)

		case "unblock":
			replyMessage, replyMarkup, err = helper.UnblockStudent(srv.store, callback)

		}
	case "homework":
		switch callback.Command {
		case "homeworks_list", "next", "previous":
			replyMessage, replyMarkup, err = helper.GetHomeworksList(srv.store, callback)

		case "get":
			replyMessage, replyMarkup, err = helper.GetHomework(srv.store, callback)

		}
	}

	if err != nil {
		return c.EditOrReply(helper.ErrInternal, &telebot.SendOptions{ParseMode: "HTML"})
	}

	return c.EditOrReply(replyMessage, &telebot.SendOptions{ParseMode: "HTML"}, replyMarkup)
}
