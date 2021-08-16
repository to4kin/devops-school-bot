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
			name: "valid_json",
			a: func() []byte {
				return []byte(`
				{
					"type": "account",
					"id": "99999"
				}`)
			},
			isValid: true,
		},
		{
			name: "invalid_json",
			a: func() []byte {
				return []byte(`
				{
					"type": "account",
					"id": "99999",
				}`)
			},
			isValid: false,
		},
		{
			name: "invalid_id",
			a: func() []byte {
				return []byte(`
				{
					"type": "account",
					"id": "invalid",
				}`)
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, callback.Unmarshal(tc.a()))
			} else {
				assert.Error(t, callback.Unmarshal(tc.a()))
			}
		})
	}
}
