package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
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
			name: "empty_account",
			s: func() *model.Student {
				s := model.TestStudent(t)
				s.Account = &model.Account{}
				return s
			},
			isValid: false,
		},
		{
			name: "empty_school",
			s: func() *model.Student {
				s := model.TestStudent(t)
				s.School = &model.School{}
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
