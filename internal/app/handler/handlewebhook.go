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

		handler.logger.WithFields(logrus.Fields{
			"update_id": u.ID,
		}).Info("new message received")
		handler.bot.ProcessUpdate(u)
		handler.respond(rw, r, http.StatusOK, nil)
	}
}

func (handler *Handler) error(rw http.ResponseWriter, r *http.Request, code int, err error) {
	handler.respond(rw, r, code, map[string]string{"error": err.Error()})
}

func (handler *Handler) respond(rw http.ResponseWriter, r *http.Request, code int, data interface{}) {
	rw.WriteHeader(code)
	if data != nil {
		json.NewEncoder(rw).Encode(data)
	}
}
