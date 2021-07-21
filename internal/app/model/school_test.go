package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
)

func TestSchool_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		s       func() *model.School
		isValid bool
	}{
		{
			name: "valid",
			s: func() *model.School {
				return model.TestSchool(t)
			},
			isValid: true,
		},
		{
			name: "empty_title",
			s: func() *model.School {
				s := model.TestSchool(t)
				s.Title = ""
				return s
			},
			isValid: false,
		},
		{
			name: "active_and_finished_at_the_same_time",
			s: func() *model.School {
				s := model.TestSchool(t)
				s.Finished = true
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
