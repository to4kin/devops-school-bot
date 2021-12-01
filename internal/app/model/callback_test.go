package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
)

func TestCallback_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		c       func() *model.Callback
		isValid bool
	}{
		{
			name: "valid",
			c: func() *model.Callback {
				return model.TestAccountCallback(t)
			},
			isValid: true,
		},
		{
			name: "empty_type",
			c: func() *model.Callback {
				c := model.TestAccountCallback(t)
				c.Type = ""
				return c
			},
			isValid: false,
		},
		{
			name: "zero_type_id",
			c: func() *model.Callback {
				c := model.TestAccountCallback(t)
				c.TypeID = 0
				return c
			},
			isValid: false,
		},
		{
			name: "empty_command",
			c: func() *model.Callback {
				c := model.TestAccountCallback(t)
				c.Command = ""
				return c
			},
			isValid: false,
		},
		{
			name: "empty_list_command",
			c: func() *model.Callback {
				c := model.TestAccountCallback(t)
				c.ListCommand = ""
				return c
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.c().Validate())
			} else {
				assert.Error(t, tc.c().Validate())
			}
		})
	}
}
