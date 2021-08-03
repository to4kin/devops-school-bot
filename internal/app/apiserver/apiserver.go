package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // for migrations by file
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/sqlstore"
	"gopkg.in/tucnak/telebot.v3"
)

// Start ...
func Start(config *Config) error {
	db, err := newDb(config.Database.URL, config.Database.Migrations)
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

	srv.configureLogger(config.LogLevel)
	srv.configureBotHandler()

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDb(databaseURL string, migrations string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrations,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return db, nil
}
