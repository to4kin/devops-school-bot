package handler

import (
	"encoding/json"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file" // for migrations by file
	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/configuration"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

// Handler ...
type Handler struct {
	logger    *logrus.Logger
	store     store.Store
	bot       *telebot.Bot
	buildDate string
	version   string
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
		OnError: func(err error, ctx telebot.Context) {
			logger.Error(err)
		},
	})
	if err != nil {
		return nil, err
	}

	handler := &Handler{
		logger:    logger,
		store:     store,
		bot:       bot,
		buildDate: config.BuildDate,
		version:   config.Version,
	}

	handler.configureBotHandlers()

	return handler, nil
}

// Respond error to http request
func (handler *Handler) error(rw http.ResponseWriter, r *http.Request, code int, err error) {
	handler.respond(rw, r, code, map[string]string{"error": err.Error()})
}

// Respond ok to http request
func (handler *Handler) respond(rw http.ResponseWriter, r *http.Request, code int, data interface{}) {
	rw.WriteHeader(code)
	if data != nil {
		json.NewEncoder(rw).Encode(data)
	}
}

// Respond to telegram chat
func (handler *Handler) editOrReply(c telebot.Context, replyMessage string, replyMarkup *telebot.ReplyMarkup) error {
	handler.logger.WithFields(logrus.Fields{
		"message": replyMessage,
	}).Info("send message to user")
	return c.EditOrReply(replyMessage, &telebot.SendOptions{ParseMode: "HTML"}, replyMarkup)
}
