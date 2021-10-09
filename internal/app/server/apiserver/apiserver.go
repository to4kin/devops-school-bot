package apiserver

import (
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/configuration"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/handler"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/server"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/sqlstore"
)

// Server ...
type Server struct {
	config *configuration.Config
}

// New ...
func New(config *configuration.Config) *Server {
	return &Server{
		config: config,
	}
}

// Start ...
func (srv *Server) Start() error {
	db, err := server.NewPostgres(srv.config.Database.URL, srv.config.Database.Migrations)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)

	handler, err := handler.NewHandler(srv.config, store)
	if err != nil {
		return err
	}

	if srv.config.Apiserver.Cron.Enabled {
		cron := gocron.NewScheduler(time.UTC)
		if _, err := cron.Cron(srv.config.Apiserver.Cron.Schedule).Do(handler.HandleCron, srv.config.Apiserver.Cron.Fullreport); err != nil {
			return err
		}

		cron.StartAsync()
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handler.HandleWebHook()).Methods("POST")

	return http.ListenAndServe(srv.config.Apiserver.BindAddr, router)
}
