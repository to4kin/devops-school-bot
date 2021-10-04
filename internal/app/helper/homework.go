package helper

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

// GetHomework ...
func (hlpr *Helper) GetHomework(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get homework from database by id")
	homework, err := hlpr.store.Homework().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(homework.LogrusFields()).Debug("homework found")

	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}

	if homework.Active {
		disableHomeworkCallback := *callback
		disableHomeworkCallback.Command = "disable_homework"
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Disable homework", disableHomeworkCallback.ToString())))
	} else {
		enableHomeworkCallback := *callback
		enableHomeworkCallback.Command = "enable_homework"
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Enable homework", enableHomeworkCallback.ToString())))
	}

	homeworksListCallback := &model.Callback{
		ID:          homework.Student.School.ID,
		Type:        "school",
		Command:     "homeworks",
		ListCommand: "homeworks",
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to Homeworks List", homeworksListCallback.ToString())))
	replyMarkup.Inline(rows...)

	return fmt.Sprintf(
			"School: <b>%v</b>\n\nHomework info:\n\nTitle: %v\nModule: %v\nStatus: %v",
			homework.Student.School.Title,
			homework.Lesson.Title,
			homework.Lesson.Module.Title,
			homework.GetStatusText(),
		),
		replyMarkup,
		nil
}

// DisableHomework ...
func (hlpr *Helper) DisableHomework(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get homework from database by id")
	homework, err := hlpr.store.Homework().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(homework.LogrusFields()).Debug("homework found")

	hlpr.logger.WithFields(logrus.Fields{
		"lesson_id": homework.Lesson.ID,
	}).Debug("get all homeworks from database by lesson_id")
	homeworks, err := hlpr.store.Homework().FindByLessonID(homework.Lesson.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Debug("homeworks found")

	for _, hw := range homeworks {
		hw.Active = false

		hlpr.logger.Debug("disable homework")
		if err := hlpr.store.Homework().Update(hw); err != nil {
			return "", nil, err
		}
		hlpr.logger.WithFields(hw.LogrusFields()).Debug("homework disabled")
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToHomeworkRow(replyMarkup, callback, homework.ID))

	return fmt.Sprintf("Success! Homework <b>%v</b> was <b>DISABLED</b>", homework.Lesson.Title), replyMarkup, nil
}

// EnableHomework ...
func (hlpr *Helper) EnableHomework(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get homework from database by id")
	homework, err := hlpr.store.Homework().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(homework.LogrusFields()).Debug("homework found")

	hlpr.logger.WithFields(logrus.Fields{
		"lesson_id": homework.Lesson.ID,
	}).Debug("get all homeworks from database by lesson_id")
	homeworks, err := hlpr.store.Homework().FindByLessonID(homework.Lesson.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Debug("homeworks found")

	for _, hw := range homeworks {
		hw.Active = true

		hlpr.logger.Debug("enable homework")
		if err := hlpr.store.Homework().Update(hw); err != nil {
			return "", nil, err
		}
		hlpr.logger.WithFields(hw.LogrusFields()).Debug("homework enabled")
	}
	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToHomeworkRow(replyMarkup, callback, homework.ID))

	return fmt.Sprintf("Success! Homework <b>%v</b> was <b>ENABLED</b>", homework.Lesson.Title), replyMarkup, nil
}

func backToHomeworkRow(replyMarkup *telebot.ReplyMarkup, callback *model.Callback, homeworkID int64) telebot.Row {
	backToHomeworkCallback := &model.Callback{
		ID:          homeworkID,
		Type:        "homework",
		Command:     "get",
		ListCommand: callback.ListCommand,
	}

	backToHomeworkssListCallback := &model.Callback{
		ID:          homeworkID,
		Type:        "homework",
		Command:     "homeworks_list",
		ListCommand: callback.ListCommand,
	}

	if callback.ListCommand == "get" {
		return replyMarkup.Row(
			replyMarkup.Data("<< Back to Homework", backToHomeworkCallback.ToString()),
			replyMarkup.Data("<< Back to Homeworks List", backToHomeworkssListCallback.ToString()),
		)
	}

	return replyMarkup.Row(replyMarkup.Data("<< Back to Homeworks List", backToHomeworkssListCallback.ToString()))
}
