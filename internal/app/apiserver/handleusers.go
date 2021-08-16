package apiserver

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleUsers(c telebot.Context) error {
	if !c.Message().Private() {
		return nil
	}

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

	calbback := &model.Callback{
		Type: "account",
		ID:   0,
	}

	return srv.usersNaviButtons(c, calbback)
}

func (srv *server) userRespond(c telebot.Context, callback *model.Callback) error {
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

	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}
	if c.Sender().ID != account.TelegramID {
		if account.Superuser {
			rows = append(rows, replyMarkup.Row(replyMarkup.Data("Unset Superuser", "unset_superuser", callback.ToString())))
		} else {
			rows = append(rows, replyMarkup.Row(replyMarkup.Data("Set Superuser", "set_superuser", callback.ToString())))
		}
	} else {
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Update account", "update", callback.ToString())))
	}

	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to user list", "accounts_list", callback.ToString())))
	replyMarkup.Inline(rows...)

	return c.EditOrSend(
		fmt.Sprintf(
			msgUserInfo,
			account.FirstName,
			account.LastName,
			account.Username,
			account.Superuser,
		),
		&telebot.SendOptions{ParseMode: "HTML"},
		replyMarkup,
	)
}

func (srv *server) usersNaviButtons(c telebot.Context, callback *model.Callback) error {
	srv.logger.Debug("get all accounts")
	accounts, err := srv.store.Account().FindAll()
	if err != nil {
		srv.logger.Error(err)
		return c.Respond(&telebot.CallbackResponse{
			Text:      msgInternalError,
			ShowAlert: true,
		})
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(accounts),
	}).Debug("accounts found")

	page := 0
	for i, account := range accounts {
		if callback.ID == account.ID {
			page = i / (maxRows * 2)
			break
		}
	}

	var buttons []telebot.Btn
	replyMarkup := &telebot.ReplyMarkup{}
	for _, account := range accounts {
		accountCallback := &model.Callback{
			Type: callback.Type,
			ID:   account.ID,
		}
		buttons = append(buttons, replyMarkup.Data(account.Username, "get", accountCallback.ToString()))
	}

	var rows []telebot.Row
	div, mod := len(accounts)/2, len(accounts)%2

	nextCallback := &model.Callback{
		Type: callback.Type,
	}

	previousCallback := &model.Callback{
		Type: callback.Type,
	}

	if div >= maxRows*(page+1) {
		for i := maxRows * page; i < maxRows*(page+1); i++ {
			rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
		}

		nextCallback.ID = accounts[maxRows*2*(page+1)].ID
		btnNext := replyMarkup.Data("Next page >>", "next", nextCallback.ToString())

		if page > 0 {
			previousCallback.ID = accounts[maxRows*2*(page-1)].ID
			btnPrevious := replyMarkup.Data("<< Previous page", "previous", previousCallback.ToString())

			rows = append(rows, replyMarkup.Row(btnPrevious, btnNext))
		} else {
			rows = append(rows, replyMarkup.Row(btnNext))
		}
	} else {
		for i := maxRows * page; i < div; i++ {
			rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
		}
		if mod != 0 {
			rows = append(rows, replyMarkup.Row(buttons[div*2]))
		}
		if page > 0 {
			previousCallback.ID = accounts[maxRows*2*(page-1)].ID
			btnPrevious := replyMarkup.Data("<< Previous page", "previous", previousCallback.ToString())

			rows = append(rows, replyMarkup.Row(btnPrevious))
		}
	}

	replyMarkup.Inline(rows...)
	return c.EditOrSend("Choose a user from the list below:", replyMarkup)
}
