package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
)

func TestStudent_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		s       func() *model.Student
		isValid bool
	}{
		{
			name: "valid",
			s: func() *model.Student {
				return model.TestStudent(t)
			},
			isValid: true,
		},
		{
			name: "zero_telegram_id",
			s: func() *model.Student {
				s := model.TestStudent(t)
				s.TelegramID = int64(0)
				return s
			},
			isValid: false,
		},
		{
			name: "empty_first_name",
			s: func() *model.Student {
				s := model.TestStudent(t)
				s.FirstName = ""
				return s
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.s().Validate())
			} else {
				assert.Error(t, tc.s().Validate())
			}
		})
	}
}
