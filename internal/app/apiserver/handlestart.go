package apiserver

import (
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleStart(c telebot.Context) error {
	logger := logrus.WithFields(logrus.Fields{
		"handler": "start",
	})

	if c.Message().Private() {
		return nil
	}

	logger.Debug("get school by chat_id: ", c.Message().Chat.ID)
	school, err := srv.store.School().FindByChatID(c.Message().Chat.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			logger.Debug("school not found, will create a new one")
			school = &model.School{
				Title:    c.Message().Chat.Title,
				ChatID:   c.Message().Chat.ID,
				Active:   true,
				Finished: false,
			}

			if err := srv.store.School().Create(school); err != nil {
				logger.Error(err)
				return nil
			}
		} else {
			logger.Error(err)
			return nil
		}
	}
	logger.Debug(school.ToString())
	return nil
}
