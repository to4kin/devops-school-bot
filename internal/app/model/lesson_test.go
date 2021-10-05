package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
)

func TestLesson_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		l       func() *model.Lesson
		isValid bool
	}{
		{
			name: "valid",
			l: func() *model.Lesson {
				return model.TestLessonOne(t)
			},
			isValid: true,
		},
		{
			name: "empty_title",
			l: func() *model.Lesson {
				l := model.TestLessonOne(t)
				l.Title = ""
				return l
			},
			isValid: false,
		},
		{
			name: "empty_module",
			l: func() *model.Lesson {
				l := model.TestLessonOne(t)
				l.Module = &model.Module{}
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
