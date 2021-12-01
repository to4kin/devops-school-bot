package helper

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

// GetUsersList ...
func (hlpr *Helper) GetUsersList(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	var err error
	var accounts []*model.Account
	switch callback.ListCommand {
	case "set_superuser":
		hlpr.logger.Info("get all users from database")
		accounts, err = hlpr.store.Account().FindBySuperuser(false)
	case "unset_superuser":
		// TODO: remove Sender from accounts,because user can remove superuser permissions for himself
		hlpr.logger.Info("get all superusers from database")
		accounts, err = hlpr.store.Account().FindBySuperuser(true)
	default:
		hlpr.logger.Info("get all accounts from database")
		accounts, err = hlpr.store.Account().FindAll()
	}

	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		return "There are no users in the database. Please add first", nil, nil
	}

	hlpr.logger.WithFields(logrus.Fields{
		"count": len(accounts),
	}).Info("accounts found")

	var interfaceSlice []model.Interface = make([]model.Interface, len(accounts))
	for i, v := range accounts {
		interfaceSlice[i] = v
	}
	rows, err := hlpr.rowsWithButtons(interfaceSlice, callback)
	if err != nil {
		return "", nil, err
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(rows...)
	return "Choose a user from the list below:", replyMarkup, nil
}

// GetUser ...
func (hlpr *Helper) GetUser(callback *model.Callback, sender *telebot.User) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get account from database by id")
	account, err := hlpr.store.Account().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(account.LogrusFields()).Info("account found")

	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}
	if sender.ID != account.TelegramID {
		if account.Superuser {
			unsetSuperuserCallback := &model.Callback{
				Created:     time.Now(),
				Type:        callback.Type,
				TypeID:      callback.TypeID,
				Command:     "unset_superuser",
				ListCommand: callback.ListCommand,
			}
			if err := hlpr.prepareCallback(unsetSuperuserCallback); err != nil {
				return "", nil, err
			}
			rows = append(rows, replyMarkup.Row(replyMarkup.Data("Unset Superuser", unsetSuperuserCallback.GetStringID())))

		} else {
			setSuperuserCallback := &model.Callback{
				Created:     time.Now(),
				Type:        callback.Type,
				TypeID:      callback.TypeID,
				Command:     "set_superuser",
				ListCommand: callback.ListCommand,
			}
			if err := hlpr.prepareCallback(setSuperuserCallback); err != nil {
				return "", nil, err
			}
			rows = append(rows, replyMarkup.Row(replyMarkup.Data("Set Superuser", setSuperuserCallback.GetStringID())))
		}
	} else {
		if account.FirstName != sender.FirstName || account.LastName != sender.LastName || account.Username != sender.Username {
			updateCallback := &model.Callback{
				Created:     time.Now(),
				Type:        callback.Type,
				TypeID:      callback.TypeID,
				Command:     "update",
				ListCommand: callback.ListCommand,
			}
			if err := hlpr.prepareCallback(updateCallback); err != nil {
				return "", nil, err
			}
			rows = append(rows, replyMarkup.Row(replyMarkup.Data("Update account", updateCallback.GetStringID())))
		}
	}

	backToUsersListCallback := &model.Callback{
		Created:     time.Now(),
		Type:        callback.Type,
		TypeID:      callback.TypeID,
		Command:     "accounts_list",
		ListCommand: callback.ListCommand,
	}
	if err := hlpr.prepareCallback(backToUsersListCallback); err != nil {
		return "", nil, err
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to Users List", backToUsersListCallback.GetStringID())))
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
func (hlpr *Helper) UpdateUser(callback *model.Callback, sender *telebot.User) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get account from database by id")
	account, err := hlpr.store.Account().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(account.LogrusFields()).Info("account found")

	account.FirstName = sender.FirstName
	account.LastName = sender.LastName
	account.Username = sender.Username

	hlpr.logger.Info("update account")
	if err := hlpr.store.Account().Update(account); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(account.LogrusFields()).Info("account updated")

	replyMarkup := &telebot.ReplyMarkup{}
	backRow, err := hlpr.backToUserRow(replyMarkup, callback.ListCommand, account.ID)
	if err != nil {
		return "", nil, err
	}
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("Success! Account <b>%v</b> updated", account.Username), replyMarkup, nil
}

// SetSuperuser ...
func (hlpr *Helper) SetSuperuser(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get account from database by id")
	account, err := hlpr.store.Account().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(account.LogrusFields()).Info("account found")

	account.Superuser = true

	hlpr.logger.Info("set superuser")
	if err := hlpr.store.Account().Update(account); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(account.LogrusFields()).Info("account updated")

	replyMarkup := &telebot.ReplyMarkup{}
	backRow, err := hlpr.backToUserRow(replyMarkup, callback.ListCommand, account.ID)
	if err != nil {
		return "", nil, err
	}
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("Success! Superuser access <b>ENABLED</b> for user <b>%v</b>", account.Username), replyMarkup, nil
}

// UnsetSuperuser ...
func (hlpr *Helper) UnsetSuperuser(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get account from database by id")
	account, err := hlpr.store.Account().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(account.LogrusFields()).Info("account found")

	account.Superuser = false

	hlpr.logger.Info("unset superuser")
	if err := hlpr.store.Account().Update(account); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(account.LogrusFields()).Info("account updated")

	replyMarkup := &telebot.ReplyMarkup{}
	backRow, err := hlpr.backToUserRow(replyMarkup, callback.ListCommand, account.ID)
	if err != nil {
		return "", nil, err
	}
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("Success! Superuser access <b>DISABLED</b> for user <b>%v</b>", account.Username), replyMarkup, nil
}

func (hlpr *Helper) backToUserRow(replyMarkup *telebot.ReplyMarkup, listCommand string, accountID int64) (telebot.Row, error) {
	backToUser := &model.Callback{
		Created:     time.Now(),
		Type:        "account",
		TypeID:      accountID,
		Command:     "get",
		ListCommand: listCommand,
	}
	if err := hlpr.prepareCallback(backToUser); err != nil {
		return nil, err
	}

	backToUsersList := &model.Callback{
		Created:     time.Now(),
		Type:        "account",
		TypeID:      accountID,
		Command:     "accounts_list",
		ListCommand: listCommand,
	}
	if err := hlpr.prepareCallback(backToUsersList); err != nil {
		return nil, err
	}

	if listCommand == "get" {
		return replyMarkup.Row(
			replyMarkup.Data("<< Back to User", backToUser.GetStringID()),
			replyMarkup.Data("<< Back to Users List", backToUsersList.GetStringID()),
		), nil
	}

	return replyMarkup.Row(replyMarkup.Data("<< Back to Users List", backToUsersList.GetStringID())), nil
}
