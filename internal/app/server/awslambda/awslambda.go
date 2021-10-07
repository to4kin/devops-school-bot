package awslambda

import (
	"github.com/aws/aws-lambda-go/lambda"
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

	lambda.Start(handler.HandleLambda)
	return nil
}
