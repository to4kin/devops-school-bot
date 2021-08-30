package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
)

func TestModule_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		l       func() *model.Module
		isValid bool
	}{
		{
			name: "valid",
			l: func() *model.Module {
				return model.TestModule(t)
			},
			isValid: true,
		},
		{
			name: "empty_title",
			l: func() *model.Module {
				l := model.TestModule(t)
				l.Title = ""
				return l
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.l().Validate())
			} else {
				assert.Error(t, tc.l().Validate())
			}
		})
	}
}
