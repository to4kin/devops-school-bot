package helper

import (
	"fmt"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

// GetUsersList ...
func GetUsersList(str store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	var err error
	var accounts []*model.Account
	switch callback.ListCommand {
	case "set_superuser":
		accounts, err = str.Account().FindBySuperuser(false)
	case "unset_superuser":
		accounts, err = str.Account().FindBySuperuser(true)
	default:
		accounts, err = str.Account().FindAll()
	}

	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		return "There are no users in the database. Please add first", nil, nil
	}

	var interfaceSlice []model.Interface = make([]model.Interface, len(accounts))
	for i, v := range accounts {
		interfaceSlice[i] = v
	}
	rows := rowsWithButtons(interfaceSlice, callback)

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(rows...)
	return "Choose a user from the list below:", replyMarkup, nil
}

// GetUser ...
func GetUser(store store.Store, callback *model.Callback, sender *telebot.User) (string, *telebot.ReplyMarkup, error) {
	account, err := store.Account().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}
	if sender.ID != account.TelegramID {
		if account.Superuser {
			unsetSuperuserCallback := *callback
			unsetSuperuserCallback.Command = "unset_superuser"
			rows = append(rows, replyMarkup.Row(replyMarkup.Data("Unset Superuser", unsetSuperuserCallback.ToString())))
		} else {
			setSuperuserCallback := *callback
			setSuperuserCallback.Command = "set_superuser"
			rows = append(rows, replyMarkup.Row(replyMarkup.Data("Set Superuser", setSuperuserCallback.ToString())))
		}
	} else {
		if account.FirstName != sender.FirstName || account.LastName != sender.LastName || account.Username != sender.Username {
			updateCallback := *callback
			updateCallback.Command = "update"
			rows = append(rows, replyMarkup.Row(replyMarkup.Data("Update account", updateCallback.ToString())))
		}
	}

	backToUsersListCallback := *callback
	backToUsersListCallback.Command = "accounts_list"
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to Users List", backToUsersListCallback.ToString())))
	replyMarkup.Inline(rows...)

	return fmt.Sprintf(
		"Account info:\n\nFirst name: %v\nLast name: %v\nUsername: @%v\nSuperuser: %v",
		account.FirstName,
		account.LastName,
		account.Username,
		account.Superuser,
	), replyMarkup, nil
}

// UpdateUser ...
func UpdateUser(store store.Store, callback *model.Callback, sender *telebot.User) (string, *telebot.ReplyMarkup, error) {
	account, err := store.Account().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	account.FirstName = sender.FirstName
	account.LastName = sender.LastName
	account.Username = sender.Username

	if err := store.Account().Update(account); err != nil {
		return "", nil, err
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToUserRow(replyMarkup, callback, account.ID))

	return fmt.Sprintf("Success! Account <b>%v</b> updated", account.Username), replyMarkup, nil
}

// SetSuperuser ...
func SetSuperuser(store store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	account, err := store.Account().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	account.Superuser = true

	if err := store.Account().Update(account); err != nil {
		return "", nil, err
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToUserRow(replyMarkup, callback, account.ID))

	return fmt.Sprintf("Success! Superuser access <b>ENABLED</b> for user <b>%v</b>", account.Username), replyMarkup, nil
}

// UnsetSuperuser ...
func UnsetSuperuser(store store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	account, err := store.Account().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	account.Superuser = false

	if err := store.Account().Update(account); err != nil {
		return "", nil, err
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToUserRow(replyMarkup, callback, account.ID))

	return fmt.Sprintf("Success! Superuser access <b>DISABLED</b> for user <b>%v</b>", account.Username), replyMarkup, nil
}

func backToUserRow(replyMarkup *telebot.ReplyMarkup, callback *model.Callback, accountID int64) telebot.Row {
	backToUser := &model.Callback{
		ID:          accountID,
		Type:        "account",
		Command:     "get",
		ListCommand: callback.ListCommand,
	}

	backToUsersList := &model.Callback{
		ID:          accountID,
		Type:        "account",
		Command:     "accounts_list",
		ListCommand: callback.ListCommand,
	}

	if callback.ListCommand == "get" {
		return replyMarkup.Row(
			replyMarkup.Data("<< Back to User", backToUser.ToString()),
			replyMarkup.Data("<< Back to Users List", backToUsersList.ToString()),
		)
	}

	return replyMarkup.Row(replyMarkup.Data("<< Back to Users List", backToUsersList.ToString()))
}
