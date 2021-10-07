package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v3"
)

// HandleLambda ...
func (handler *Handler) HandleLambda(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var u telebot.Update
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "error" + err.Error(),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	handler.logger.WithFields(logrus.Fields{
		"update_id": u.ID,
	}).Info("new message received")
	handler.bot.ProcessUpdate(u)
	return events.APIGatewayProxyResponse{
		Body:       "ok",
		StatusCode: http.StatusOK,
	}, nil
}
