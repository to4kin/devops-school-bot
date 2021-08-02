package apiserver

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleFinish(c telebot.Context) error {
	logger := logrus.WithFields(logrus.Fields{
		"handler": "finish",
	})

	if c.Message().Private() {
		return nil
	}

	logger.Debug("get school by chat_id: ", c.Message().Chat.ID)
	school, err := srv.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		logger.Error(err)
		return nil
	}
	logger.Debug(school.ToString())

	if school.Finished {
		return nil
	}

	if err := srv.store.School().Finish(school); err != nil {
		logger.Error(err)
	}
	logger.Debug(school.ToString())

	return nil
}
