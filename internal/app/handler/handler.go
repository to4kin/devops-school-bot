package handler

import (
	_ "github.com/golang-migrate/migrate/v4/source/file" // for migrations by file
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/configuration"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

// Handler ...
type Handler struct {
	logger *logrus.Logger
	store  store.Store
	bot    *telebot.Bot
}

// NewHandler ...
func NewHandler(config *configuration.Config, store store.Store) (*Handler, error) {
	// configure logger
	logger := configuration.ConfigureLogger()
	if level, err := logrus.ParseLevel(config.LogLevel); err != nil {
		logger.Error(err)
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
		logger.WithFields(logrus.Fields{
			"log_level": level,
		}).Info("log level was updated")
	}

	// configure telebot
	bot, err := telebot.NewBot(telebot.Settings{
		Token:       config.TelegramBot.Token,
		Verbose:     config.TelegramBot.Verbose,
		Synchronous: true,
	})
	if err != nil {
		return nil, err
	}

	handler := &Handler{
		logger: logger,
		store:  store,
		bot:    bot,
	}

	handler.configureBotHandlers()

	return handler, nil
}
