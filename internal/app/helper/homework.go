package helper

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

// GetHomework ...
func (hlpr *Helper) GetHomework(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Info("get homework from database by id")
	homework, err := hlpr.store.Homework().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(homework.LogrusFields()).Info("homework found")

	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}

	if homework.Active {
		disableHomeworkCallback := *callback
		disableHomeworkCallback.Command = "disable_homework"
		rows = append(rows, replyMarkup.Row(replyMarkup.Data(fmt.Sprintf("Disable all %v", homework.Lesson.Title), disableHomeworkCallback.ToString())))
	} else {
		enableHomeworkCallback := *callback
		enableHomeworkCallback.Command = "enable_homework"
		rows = append(rows, replyMarkup.Row(replyMarkup.Data(fmt.Sprintf("Enable all %v", homework.Lesson.Title), enableHomeworkCallback.ToString())))
	}

	homeworksListCallback := &model.Callback{
		ID:          homework.ID,
		Type:        "homework",
		Command:     "homeworks_list",
		ListCommand: callback.ListCommand,
	}
	schoolsListCallback := &model.Callback{
		ID:          homework.Student.School.ID,
		Type:        "school",
		Command:     "homeworks",
		ListCommand: "homeworks",
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to Homeworks List", homeworksListCallback.ToString())))
	rows = append(rows, backToSchoolRow(replyMarkup, schoolsListCallback, homework.Student.School.ID))
	replyMarkup.Inline(rows...)

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": homework.Student.School.ID,
		"lesson_id": homework.Lesson.ID,
	}).Info("get all homeworks from database by school_id and lesson_id")
	homeworks, err := hlpr.store.Homework().FindBySchoolIDLessonID(homework.Student.School.ID, homework.Lesson.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Info("homeworks found")

	disabledCound := 0
	enabledCound := 0
	for _, hw := range homeworks {
		if hw.Active {
			enabledCound++
		} else {
			disabledCound++
		}
	}

	return fmt.Sprintf(
			"School: <b>%v</b>\n\nHomework info:\n\nTitle: %v\nModule: %v\nEnabled: %d\nDisabled: %d",
			homework.Student.School.Title,
			homework.Lesson.Title,
			homework.Lesson.Module.Title,
			enabledCound,
			disabledCound,
		),
		replyMarkup,
		nil
}

// GetHomeworksList ...
func (hlpr *Helper) GetHomeworksList(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Info("get homework from database by id")
	homework, err := hlpr.store.Homework().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(homework.LogrusFields()).Info("homework found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": homework.Student.School.ID,
	}).Info("get homeworks from database by school_id")
	homeworks, err := hlpr.store.Homework().FindBySchoolID(homework.Student.School.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Info("homeworks found")

	reportMessage, err := hlpr.GetLessonsReport(homework.Student.School)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		reportMessage = ErrReportNotFound
	}

	replyMarkup := &telebot.ReplyMarkup{}
	var interfaceSlice []model.Interface = make([]model.Interface, len(homeworks))
	for i, v := range homeworks {
		interfaceSlice[i] = v
	}
	interfaceSlice = removeDuplicate(interfaceSlice)

	schoolsListCallback := &model.Callback{
		ID:          homework.Student.School.ID,
		Type:        "school",
		Command:     "homeworks",
		ListCommand: "homeworks",
	}

	rows := rowsWithButtons(interfaceSlice, callback)
	//var rows []telebot.Row
	rows = append(rows, backToSchoolRow(replyMarkup, schoolsListCallback, homework.Student.School.ID))
	replyMarkup.Inline(rows...)

	return fmt.Sprintf("School <b>%v</b>\n\n%v", homework.Student.School.Title, reportMessage), replyMarkup, nil
}

// DisableHomework ...
func (hlpr *Helper) DisableHomework(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Info("get homework from database by id")
	homework, err := hlpr.store.Homework().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(homework.LogrusFields()).Info("homework found")

	hlpr.logger.WithFields(logrus.Fields{
		"lesson_id": homework.Lesson.ID,
	}).Info("get all homeworks from database by lesson_id")
	homeworks, err := hlpr.store.Homework().FindBySchoolIDLessonID(homework.Student.School.ID, homework.Lesson.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Info("homeworks found")

	for _, hw := range homeworks {
		hw.Active = false

		hlpr.logger.Info("disable homework")
		if err := hlpr.store.Homework().Update(hw); err != nil {
			return "", nil, err
		}
		hlpr.logger.WithFields(hw.LogrusFields()).Info("homework disabled")
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToHomeworkRow(replyMarkup, callback, homework.ID))

	return fmt.Sprintf("Success! Homeworks with lesson <b>%v</b> were <b>DISABLED</b>", homework.Lesson.Title), replyMarkup, nil
}

// EnableHomework ...
func (hlpr *Helper) EnableHomework(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Info("get homework from database by id")
	homework, err := hlpr.store.Homework().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(homework.LogrusFields()).Info("homework found")

	hlpr.logger.WithFields(logrus.Fields{
		"lesson_id": homework.Lesson.ID,
	}).Info("get all homeworks from database by lesson_id")
	homeworks, err := hlpr.store.Homework().FindBySchoolIDLessonID(homework.Student.School.ID, homework.Lesson.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Info("homeworks found")

	for _, hw := range homeworks {
		hw.Active = true

		hlpr.logger.Info("enable homework")
		if err := hlpr.store.Homework().Update(hw); err != nil {
			return "", nil, err
		}
		hlpr.logger.WithFields(hw.LogrusFields()).Info("homework enabled")
	}
	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToHomeworkRow(replyMarkup, callback, homework.ID))

	return fmt.Sprintf("Success! Homeworks with lesson <b>%v</b> were <b>ENABLED</b>", homework.Lesson.Title), replyMarkup, nil
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
