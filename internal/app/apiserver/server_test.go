package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store/teststore"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	token string
)

func TestMain(m *testing.M) {
	if token == "" {
		token = "1949550059:AAHTvp0Zm5ABVDKL8LVHAYkS-PEEGGZnEJE"
	}

	os.Exit(m.Run())
}

func TestServer_BotWebHookHandler(t *testing.T) {
	srv := newServer(teststore.New())

	bot, err := telebot.NewBot(telebot.Settings{Token: token})
	assert.NoError(t, err)
	assert.NotNil(t, bot)

	srv.bot = bot

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"update_id": 10000,
				"message": map[string]interface{}{
					"date": 1441645532,
					"chat": map[string]interface{}{
						"last_name":  "Test Lastname",
						"type":       "private",
						"id":         1111111,
						"first_name": "Test Firstname",
						"username":   "Testusername",
					},
					"message_id": 1365,
					"from": map[string]interface{}{
						"last_name":  "Test Lastname",
						"id":         1111111,
						"first_name": "Test Firstname",
						"username":   "Testusername",
					},
					"text": "/start #hashtag",
					"entities": []map[string]interface{}{
						{
							"type":   "hashtag",
							"offset": 7,
							"length": 8,
						},
						{
							"type":   "bot_command",
							"offset": 0,
							"length": 6,
						},
					},
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid_payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/", b)
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
