package awslambda

import (
	"github.com/aws/aws-lambda-go/lambda"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/configuration"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/handler"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/server"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store/sqlstore"
)

// Server ...
type AWSLambda struct {
	config *configuration.Config
}

// New ...
func New(config *configuration.Config) *AWSLambda {
	return &AWSLambda{
		config: config,
	}
}

// Start ...
func (awslambda *AWSLambda) Start() error {
	db, err := server.NewPostgres(awslambda.config.Database.URL, awslambda.config.Database.Migrations)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)

	handler, err := handler.NewHandler(awslambda.config, store)
	if err != nil {
		return err
	}

	lambda.Start(handler.HandleLambda)
	return nil
}
