package apiserver

import (
	"database/sql"
	"net/http"

	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store/sqlstore"
	"gopkg.in/tucnak/telebot.v3"
)

func Start(config *Config) error {
	db, err := newDb(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	srv := newServer(store)

	srv.bot, err = telebot.NewBot(telebot.Settings{
		Token:   config.TelegramBot.Token,
		Verbose: config.TelegramBot.Verbose,
	})

	if err != nil {
		return err
	}

	srv.configureBotHandler()

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDb(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
