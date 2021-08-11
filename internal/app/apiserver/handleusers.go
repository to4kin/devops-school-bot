package apiserver

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleUsers(c telebot.Context) error {
	if !c.Message().Private() {
		return nil
	}

	return srv.userList(c, 0)
}

func (srv *server) userList(c telebot.Context, page int) error {
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

	return srv.usersNaviButtons(c, page)
}

func (srv *server) userRespond(c telebot.Context, account *model.Account, page int) error {
	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}
	if c.Sender().ID != account.TelegramID {
		if account.Superuser {
			rows = append(rows, replyMarkup.Row(replyMarkup.Data("Unset Superuser", strconv.FormatInt(account.TelegramID, 10), "account", "unset_superuser", strconv.Itoa(page))))
		} else {
			rows = append(rows, replyMarkup.Row(replyMarkup.Data("Set Superuser", strconv.FormatInt(account.TelegramID, 10), "account", "set_superuser", strconv.Itoa(page))))
		}
	} else {
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Update account", strconv.FormatInt(account.TelegramID, 10), "account", "update", strconv.Itoa(page))))
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to user list", strconv.FormatInt(account.TelegramID, 10), "account", "back_to_list", strconv.Itoa(page))))
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

func (srv *server) usersNaviButtons(c telebot.Context, page int) error {
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

	var buttons []telebot.Btn
	replyMarkup := &telebot.ReplyMarkup{}
	for _, account := range accounts {
		buttons = append(buttons, replyMarkup.Data(account.Username, strconv.FormatInt(account.TelegramID, 10), "account", "get", strconv.Itoa(page)))
	}

	var rows []telebot.Row
	div, mod := len(accounts)/2, len(accounts)%2
	btnNext := replyMarkup.Data("Next page >>", strconv.Itoa(page+1), "account", "next", strconv.Itoa(page))
	btnPrevious := replyMarkup.Data("<< Previous page", strconv.Itoa(page-1), "account", "previous", strconv.Itoa(page))

	if div >= maxRows*(page+1) {
		for i := maxRows * page; i < maxRows*(page+1); i++ {
			rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
		}
		if page > 0 {
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
			rows = append(rows, replyMarkup.Row(btnPrevious))
		}
	}

	replyMarkup.Inline(rows...)
	return c.EditOrSend("Choose a user from the list below:", replyMarkup)
}
