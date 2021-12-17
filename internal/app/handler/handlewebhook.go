package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v3"
)

// HandleWebHook ...
func (handler *Handler) HandleWebHook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var u telebot.Update
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			handler.error(rw, r, http.StatusBadRequest, err)
			return
		}

		if u.Callback != nil {
			handler.logger.WithFields(logrus.Fields{
				"update_id":         u.ID,
				"private":           u.Callback.Message.Private(),
				"callback_raw_data": u.Callback.Data,
			}).Info("new callback received")
			handler.bot.ProcessUpdate(u)
		}

		if u.Message != nil {
			handler.logger.WithFields(logrus.Fields{
				"update_id": u.ID,
				"private":   u.Message.Private(),
				"message":   u.Message.Text,
			}).Info("new message received")
			handler.bot.ProcessUpdate(u)
		}

		if u.EditedMessage != nil {
			handler.logger.WithFields(logrus.Fields{
				"update_id": u.ID,
				"private":   u.EditedMessage.Private(),
				"message":   u.EditedMessage.Text,
			}).Info("edited message received")
			handler.bot.ProcessUpdate(u)
		}

		handler.respond(rw, r, http.StatusOK, nil)
	}
}
