package apiserver

import (
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	maxRows = 3
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
		return c.Respond(&telebot.CallbackResponse{
			Text:      msgInternalError,
			ShowAlert: true,
		})
	}
	srv.logger.WithFields(account.LogrusFields()).Debug("account found")

	if !account.Superuser {
		srv.logger.WithFields(account.LogrusFields()).Debug("account has insufficient permissions")
		return c.EditOrSend(msgUserInsufficientPermissions, &telebot.SendOptions{ParseMode: "HTML"})
	}

	callbackData := strings.Split(c.Callback().Data[1:], "|")
	callbackUnique := callbackData[0]

	callback := &model.Callback{}
	callback.Unmarshal([]byte(callbackData[1]))

	srv.logger.WithFields(logrus.Fields{
		"callback_type":   callback.Type,
		"callback_id":     callback.ID,
		"callback_unique": callbackUnique,
	}).Debug("parse callback data")

	switch callback.Type {
	case "school":
		switch callbackUnique {
		case "schools_list", "next", "previous":
			return srv.schoolsNaviButtons(c, callback)

		case "re_activate":
			srv.logger.WithFields(logrus.Fields{
				"id": callback.ID,
			}).Debug("get school by id")
			school, err := srv.store.School().FindByID(callback.ID)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school found")
			school.Active = true
			if err := srv.store.School().Update(school); err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school re-activated")

			return srv.schoolRespond(c, callback)

		case "finish":
			srv.logger.WithFields(logrus.Fields{
				"id": callback.ID,
			}).Debug("get school by id")
			school, err := srv.store.School().FindByID(callback.ID)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school found")
			school.Active = false
			if err := srv.store.School().Update(school); err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school finished")

			return srv.schoolRespond(c, callback)

		case "get":
			return srv.schoolRespond(c, callback)
		}
	case "account":
		switch callbackUnique {
		case "accounts_list", "next", "previous":
			return srv.usersNaviButtons(c, callback)

		case "update":
			srv.logger.WithFields(logrus.Fields{
				"id": callback.ID,
			}).Debug("get account from database by id")
			account, err := srv.store.Account().FindByID(callback.ID)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(account.LogrusFields()).Debug("account found")

			if account.FirstName == c.Sender().FirstName &&
				account.LastName == c.Sender().LastName &&
				account.Username == c.Sender().Username {
				return c.Respond(&telebot.CallbackResponse{
					Text:      "User account is up to date.\nNothing to update!",
					ShowAlert: true,
				})
			}

			if err := srv.store.Account().Update(account); err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(account.LogrusFields()).Debug("account updated")

			return srv.userRespond(c, callback)

		case "set_superuser":
			srv.logger.WithFields(logrus.Fields{
				"id": callback.ID,
			}).Debug("get account from database by id")
			account, err := srv.store.Account().FindByID(callback.ID)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(account.LogrusFields()).Debug("account found")
			account.Superuser = true

			if err := srv.store.Account().Update(account); err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(account.LogrusFields()).Debug("account updated")

			return srv.userRespond(c, callback)

		case "unset_superuser":
			srv.logger.WithFields(logrus.Fields{
				"id": callback.ID,
			}).Debug("get account from database by id")
			account, err := srv.store.Account().FindByID(callback.ID)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(account.LogrusFields()).Debug("account found")
			account.Superuser = false

			if err := srv.store.Account().Update(account); err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(account.LogrusFields()).Debug("account updated")

			return srv.userRespond(c, callback)

		case "get":
			return srv.userRespond(c, callback)

		}
	case "student":
		switch callbackUnique {
		case "students_list", "next", "previous":
			return srv.studentsNaviButtons(c, callback)

		case "get":
			return srv.studentRespond(c, callback)

		}
	}

	return c.Respond(&telebot.CallbackResponse{
		Text:      msgInternalError,
		ShowAlert: true,
	})
}
