package helper

import (
	"time"

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

func (hlpr *Helper) rowsWithButtons(values []model.Interface, callback *model.Callback) ([]telebot.Row, error) {
	page := 0

	hlpr.logger.Debug("get current page in list")
	for i, value := range values {
		if callback.TypeID == value.GetID() {
			page = i / (maxRows * 2)
			break
		}
	}
	hlpr.logger.WithFields(logrus.Fields{
		"page": page,
	}).Debug("page found")

	hlpr.logger.Debug("prepare callbacks for elements")
	var buttons []telebot.Btn
	replyMarkup := &telebot.ReplyMarkup{}
	for _, value := range values {
		valueCallback := &model.Callback{
			Created:     time.Now(),
			Type:        callback.Type,
			TypeID:      value.GetID(),
			Command:     callback.ListCommand,
			ListCommand: callback.ListCommand,
		}

		hlpr.prepareCallback(valueCallback)
		buttons = append(buttons, replyMarkup.Data(value.GetButtonTitle(), valueCallback.GetStringID()))
	}

	var rows []telebot.Row
	div, mod := len(values)/2, len(values)%2

	nextCallback := &model.Callback{
		Created:     time.Now(),
		Type:        callback.Type,
		Command:     "next",
		ListCommand: callback.ListCommand,
	}

	previousCallback := &model.Callback{
		Created:     time.Now(),
		Type:        callback.Type,
		Command:     "previous",
		ListCommand: callback.ListCommand,
	}

	if div > maxRows*(page+1) {
		for i := maxRows * page; i < maxRows*(page+1); i++ {
			rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
		}

		hlpr.logger.Debug("prepare callback for 'next button'")
		nextCallback.TypeID = values[maxRows*2*(page+1)].GetID()

		hlpr.prepareCallback(nextCallback)
		btnNext := replyMarkup.Data(">>", nextCallback.GetStringID())

		if page > 0 {
			hlpr.logger.Debug("prepare callback for 'previous button'")
			previousCallback.TypeID = values[maxRows*2*(page-1)].GetID()

			hlpr.prepareCallback(previousCallback)
			btnPrevious := replyMarkup.Data("<<", previousCallback.GetStringID())

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
			hlpr.logger.Debug("prepare callback for 'previous button'")
			previousCallback.TypeID = values[maxRows*2*(page-1)].GetID()

			hlpr.prepareCallback(previousCallback)
			btnPrevious := replyMarkup.Data("<<", previousCallback.GetStringID())

			rows = append(rows, replyMarkup.Row(btnPrevious))
		}
	}

	return rows, nil
}

func (hlpr *Helper) prepareCallback(callback *model.Callback) error {
	hlpr.logger.WithFields(callback.LogrusFields()).Debug("check if callback already exist")
	c, err := hlpr.store.Callback().FindByCallback(callback)
	if err != nil {
		if err == store.ErrRecordNotFound {
			hlpr.logger.WithFields(callback.LogrusFields()).Debug("insert callback into database")
			err := hlpr.store.Callback().Create(callback)
			if err != nil {
				hlpr.logger.Error(err)
				return err
			}
			hlpr.logger.WithFields(logrus.Fields{
				"new_callback_id": callback.ID,
			}).Debug("callback successfully updated")
			return nil
		}

		hlpr.logger.Error(err)
		return err
	}
	hlpr.logger.WithFields(c.LogrusFields()).Debug("callback found")

	callback.ID = c.ID
	callback.Created = c.Created
	return nil
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
