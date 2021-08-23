package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
)

func TestCallback_Unmarshal(t *testing.T) {
	callback := &model.Callback{}

	testCases := []struct {
		name    string
		a       func() string
		isValid bool
	}{
		{
			name: "valid_string",
			a: func() string {
				return string("99999|account|accounts|get")
			},
			isValid: true,
		},
		{
			name: "invalid_small_string",
			a: func() string {
				return string("99999|account|accounts")
			},
			isValid: false,
		},
		{
			name: "invalid_big_string",
			a: func() string {
				return string("99999|account|accounts|get|somethingelse")
			},
			isValid: false,
		},
		{
			name: "invalid_id",
			a: func() string {
				return string("invalid|account|accounts|get")
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
