package helper

import (
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	maxRows = 3
)

// Helper ...
type Helper struct {
	logger *logrus.Logger
	store  store.Store
}

// NewHelper ...
func NewHelper(store store.Store, logger *logrus.Logger) *Helper {
	hlpr := &Helper{
		logger: logger,
		store:  store,
	}

	return hlpr
}

func rowsWithButtons(values []model.Interface, callback *model.Callback) []telebot.Row {
	page := 0
	for i, value := range values {
		if callback.ID == value.GetID() {
			page = i / (maxRows * 2)
			break
		}
	}

	var buttons []telebot.Btn
	replyMarkup := &telebot.ReplyMarkup{}
	for _, value := range values {
		valueCallback := &model.Callback{
			ID:          value.GetID(),
			Type:        callback.Type,
			Command:     callback.ListCommand,
			ListCommand: callback.ListCommand,
		}
		buttons = append(buttons, replyMarkup.Data(value.GetButtonTitle(), valueCallback.ToString()))
	}

	var rows []telebot.Row
	div, mod := len(values)/2, len(values)%2

	nextCallback := &model.Callback{
		Type:        callback.Type,
		Command:     "next",
		ListCommand: callback.ListCommand,
	}

	previousCallback := &model.Callback{
		Type:        callback.Type,
		Command:     "previous",
		ListCommand: callback.ListCommand,
	}

	if div > maxRows*(page+1) {
		for i := maxRows * page; i < maxRows*(page+1); i++ {
			rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
		}

		nextCallback.ID = values[maxRows*2*(page+1)].GetID()
		btnNext := replyMarkup.Data(">>", nextCallback.ToString())

		if page > 0 {
			previousCallback.ID = values[maxRows*2*(page-1)].GetID()
			btnPrevious := replyMarkup.Data("<<", previousCallback.ToString())

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
			previousCallback.ID = values[maxRows*2*(page-1)].GetID()
			btnPrevious := replyMarkup.Data("<<", previousCallback.ToString())

			rows = append(rows, replyMarkup.Row(btnPrevious))
		}
	}

	return rows
}

func removeDuplicate(slice []model.Interface) []model.Interface {
	allKeys := make(map[string]bool)
	list := []model.Interface{}
	for _, item := range slice {
		if _, value := allKeys[item.GetButtonTitle()]; !value {
			allKeys[item.GetButtonTitle()] = true
			list = append(list, item)
		}
	}
	return list
}
