package apiserver

import (
	"net/http"
	"time"

	_ "net/http/pprof" // for profiling

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/configuration"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/handler"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/server"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/sqlstore"
)

// APIServer ...
type APIServer struct {
	config *configuration.Config
}

// New ...
func New(config *configuration.Config) *APIServer {
	return &APIServer{
		config: config,
	}
}

// Start ...
func (apiserver *APIServer) Start() error {
	db, err := server.NewPostgres(apiserver.config.Database.URL, apiserver.config.Database.Migrations)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)

	handler, err := handler.NewHandler(apiserver.config, store)
	if err != nil {
		return err
	}

	if apiserver.config.Apiserver.Cron.Enabled {
		cron := gocron.NewScheduler(time.UTC)
		if _, err := cron.Cron(apiserver.config.Apiserver.Cron.Schedule).Do(handler.HandleCron, apiserver.config.Apiserver.Cron.Fullreport); err != nil {
			return err
		}

		cron.StartAsync()
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handler.HandleWebHook()).Methods("POST")

	return http.ListenAndServe(apiserver.config.Apiserver.BindAddr, router)
}
