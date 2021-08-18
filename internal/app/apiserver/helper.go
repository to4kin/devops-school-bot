package apiserver

import (
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

func naviButtons(values []model.Interface, callback *model.Callback) []telebot.Row {
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
			Type: callback.Type,
			ID:   value.GetID(),
		}
		buttons = append(buttons, replyMarkup.Data(value.GetButtonTitle(), "get", valueCallback.ToString()))
	}

	var rows []telebot.Row
	div, mod := len(values)/2, len(values)%2

	nextCallback := &model.Callback{
		Type: callback.Type,
	}

	previousCallback := &model.Callback{
		Type: callback.Type,
	}

	if div > maxRows*(page+1) {
		for i := maxRows * page; i < maxRows*(page+1); i++ {
			rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
		}

		nextCallback.ID = values[maxRows*2*(page+1)].GetID()
		btnNext := replyMarkup.Data("Next page >>", "next", nextCallback.ToString())

		if page > 0 {
			previousCallback.ID = values[maxRows*2*(page-1)].GetID()
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
			previousCallback.ID = values[maxRows*2*(page-1)].GetID()
			btnPrevious := replyMarkup.Data("<< Previous page", "previous", previousCallback.ToString())

			rows = append(rows, replyMarkup.Row(btnPrevious))
		}
	}

	return rows
}
