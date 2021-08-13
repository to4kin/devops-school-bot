package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
)

func TestCallback_UnmarshalJSON(t *testing.T) {
	callback := &model.Callback{}

	testCases := []struct {
		name    string
		a       func() []byte
		isValid bool
	}{
		{
			name: "valid_account",
			a: func() []byte {
				return []byte(`
				{
					"action": "get",
					"type": "account",
					"content": {
						"id": "99999",
						"created": "2006-01-02T15:04:05Z",
						"telegram_id": "99999",
						"first_name": "FirstName",
						"last_name": "LastName",
						"username": "Username",
						"superuser": "false"
					}
				}`)
			},
			isValid: true,
		},
		{
			name: "valid_school",
			a: func() []byte {
				return []byte(`
				{
					"action": "get",
					"type": "school",
					"content": {
						"id": "99999",
						"created": "2006-01-02T15:04:05Z",
						"title": "Title",
						"chat_id": "99999",
						"active": "true"
					}
				}`)
			},
			isValid: true,
		},
		{
			name: "valid_student",
			a: func() []byte {
				return []byte(`
				{
					"action": "get",
					"type": "student",
					"content": {
						"id": "99999",
						"created": "2006-01-02T15:04:05Z",
						"account": {
							"id": "99999",
							"created": "2006-01-02T15:04:05Z",
							"telegram_id": "99999",
							"first_name": "FirstName",
							"last_name": "LastName",
							"username": "Username",
							"superuser": "false"
						},
						"school": {
							"id": "99999",
							"created": "2006-01-02T15:04:05Z",
							"title": "Title",
							"chat_id": "99999",
							"active": "true"
						},
						"active": "true"
					}
				}`)
			},
			isValid: true,
		},
		{
			name: "unsupported_type",
			a: func() []byte {
				return []byte(`
				{
					"action": "get",
					"type": "unsupported_type",
					"content": {
						"id": "99999",
						"created": "2006-01-02T15:04:05Z",
						"telegram_id": "99999",
						"first_name": "FirstName",
						"last_name": "LastName",
						"username": "Username",
						"superuser": "false"
					}
				}`)
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, callback.UnmarshalJSON(tc.a()))
			} else {
				assert.Error(t, callback.UnmarshalJSON(tc.a()))
			}
		})
	}
}
