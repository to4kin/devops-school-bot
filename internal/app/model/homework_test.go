package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/model"
)

func TestHomework_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		h       func() *model.Homework
		isValid bool
	}{
		{
			name: "valid",
			h: func() *model.Homework {
				return model.TestHomework(t)
			},
			isValid: true,
		},
		{
			name: "empty_student",
			h: func() *model.Homework {
				h := model.TestHomework(t)
				h.Student = &model.Student{}
				return h
			},
			isValid: false,
		},
		{
			name: "empty_lesson",
			h: func() *model.Homework {
				h := model.TestHomework(t)
				h.Lesson = &model.Lesson{}
				return h
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.h().Validate())
			} else {
				assert.Error(t, tc.h().Validate())
			}
		})
	}
}
