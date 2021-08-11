package apiserver

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
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

	if len(callbackData) < 4 {
		srv.logger.WithFields(logrus.Fields{
			"callback_data": callbackData,
		}).Error("callback is not supported")
		return c.Respond(&telebot.CallbackResponse{
			Text:      msgInternalError,
			ShowAlert: true,
		})
	}

	callbackValue := callbackData[0]
	callbackType := callbackData[1]
	callbackAction := callbackData[2]

	callbackFromPage, err := strconv.Atoi(callbackData[3])
	if err != nil {
		srv.logger.Error(err)
		return c.Respond(&telebot.CallbackResponse{
			Text:      msgInternalError,
			ShowAlert: true,
		})
	}

	srv.logger.WithFields(logrus.Fields{
		"callback_value":     callbackValue,
		"callback_type":      callbackType,
		"callback_action":    callbackAction,
		"callback_from_page": callbackFromPage,
	}).Debug("parse callback data")

	switch callbackType {
	case "school":
		if callbackAction == "back_to_list" {
			return srv.schoolList(c, callbackFromPage)
		}

		if callbackAction == "previous" {
			page, err := strconv.Atoi(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}

			if page == 0 {
				return srv.schoolList(c, 0)
			}

			return srv.schoolsNaviButtons(c, page)
		}

		if callbackAction == "next" {
			page, err := strconv.Atoi(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}

			return srv.schoolsNaviButtons(c, page)
		}

		if callbackAction == "re_activate" {
			srv.logger.WithFields(logrus.Fields{
				"title": callbackValue,
			}).Debug("get school by title")
			school, err := srv.store.School().FindByTitle(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school found")

			if err := srv.store.School().ReActivate(school); err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school re-activated")

			return srv.schoolRespond(c, school, callbackFromPage)
		}

		if callbackAction == "finish" {
			srv.logger.WithFields(logrus.Fields{
				"title": callbackValue,
			}).Debug("get school by title")
			school, err := srv.store.School().FindByTitle(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school found")

			if err := srv.store.School().Finish(school); err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school finished")

			return srv.schoolRespond(c, school, callbackFromPage)
		}

		if callbackAction == "get" {
			srv.logger.WithFields(logrus.Fields{
				"title": callbackValue,
			}).Debug("get school by title")
			school, err := srv.store.School().FindByTitle(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school found")

			return srv.schoolRespond(c, school, callbackFromPage)
		}
	case "account":
		if callbackAction == "back_to_list" {
			return srv.userList(c, callbackFromPage)
		}

		if callbackAction == "previous" {
			page, err := strconv.Atoi(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}

			if page == 0 {
				return srv.userList(c, 0)
			}

			return srv.usersNaviButtons(c, page)
		}

		if callbackAction == "next" {
			page, err := strconv.Atoi(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}

			return srv.usersNaviButtons(c, page)
		}

		if callbackAction == "update" {
			srv.logger.WithFields(logrus.Fields{
				"telegram_id": callbackValue,
			}).Debug("get account by telegram_id")
			telegramID, err := strconv.ParseInt(callbackValue, 10, 64)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			account, err := srv.store.Account().FindByTelegramID(telegramID)
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

			return srv.userRespond(c, account, callbackFromPage)
		}

		if callbackAction == "set_superuser" {
			srv.logger.WithFields(logrus.Fields{
				"telegram_id": callbackValue,
			}).Debug("get account by telegram_id")
			telegramID, err := strconv.ParseInt(callbackValue, 10, 64)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			account, err := srv.store.Account().FindByTelegramID(telegramID)
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

			return srv.userRespond(c, account, callbackFromPage)
		}

		if callbackAction == "unset_superuser" {
			srv.logger.WithFields(logrus.Fields{
				"telegram_id": callbackValue,
			}).Debug("get account by telegram_id")
			telegramID, err := strconv.ParseInt(callbackValue, 10, 64)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			account, err := srv.store.Account().FindByTelegramID(telegramID)
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

			return srv.userRespond(c, account, callbackFromPage)
		}

		if callbackAction == "get" {
			srv.logger.WithFields(logrus.Fields{
				"telegram_id": callbackValue,
			}).Debug("get account by telegram_id")
			telegramID, err := strconv.ParseInt(callbackValue, 10, 64)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			account, err := srv.store.Account().FindByTelegramID(telegramID)
			if err != nil {
				srv.logger.Error(err)
				return c.Respond(&telebot.CallbackResponse{
					Text:      msgInternalError,
					ShowAlert: true,
				})
			}
			srv.logger.WithFields(account.LogrusFields()).Debug("account found")

			return srv.userRespond(c, account, callbackFromPage)
		}
	}

	return c.Edit("handleCallback")
}
