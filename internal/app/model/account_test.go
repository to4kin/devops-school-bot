package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
)

func TestAccount_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		a       func() *model.Account
		isValid bool
	}{
		{
			name: "valid",
			a: func() *model.Account {
				return model.TestAccount(t)
			},
			isValid: true,
		},
		{
			name: "zero_telegram_id",
			a: func() *model.Account {
				a := model.TestAccount(t)
				a.TelegramID = int64(0)
				return a
			},
			isValid: false,
		},
		{
			name: "empty_first_name",
			a: func() *model.Account {
				a := model.TestAccount(t)
				a.FirstName = ""
				return a
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.a().Validate())
			} else {
				assert.Error(t, tc.a().Validate())
			}
		})
	}
}
