package helper

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

// GetStudent ...
func (hlpr *Helper) GetStudent(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student found")

	var buttons []telebot.Btn
	replyMarkup := &telebot.ReplyMarkup{}

	if student.Active {
		blockCallback := &model.Callback{
			ID:          callback.ID,
			Type:        "student",
			Command:     "block",
			ListCommand: callback.ListCommand,
		}
		buttons = append(buttons, replyMarkup.Data("Block student", blockCallback.ToString()))
	} else {
		unblockCallback := &model.Callback{
			ID:          callback.ID,
			Type:        "student",
			Command:     "unblock",
			ListCommand: callback.ListCommand,
		}
		buttons = append(buttons, replyMarkup.Data("Unblock student", unblockCallback.ToString()))
	}

	if student.FullCourse {
		studentCallback := &model.Callback{
			ID:          callback.ID,
			Type:        "student",
			Command:     "set_listener",
			ListCommand: callback.ListCommand,
		}
		buttons = append(buttons, replyMarkup.Data("Change to listener", studentCallback.ToString()))
	} else {
		listenerCallback := &model.Callback{
			ID:          callback.ID,
			Type:        "student",
			Command:     "set_student",
			ListCommand: callback.ListCommand,
		}
		buttons = append(buttons, replyMarkup.Data("Change to student", listenerCallback.ToString()))
	}

	var rows []telebot.Row
	div, mod := len(buttons)/2, len(buttons)%2
	for i := 0; i < div; i++ {
		rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
	}
	if mod != 0 {
		rows = append(rows, replyMarkup.Row(buttons[div*2]))
	}

	backToStudentsListCallback := &model.Callback{
		ID:          student.ID,
		Type:        "student",
		Command:     "students_list",
		ListCommand: callback.ListCommand,
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to Students List", backToStudentsListCallback.ToString())))
	rows = append(rows, backToSchoolRow(replyMarkup, callback, student.School.ID))
	replyMarkup.Inline(rows...)

	reportMessage, err := hlpr.GetUserReport(student.Account, student.School)
	if err != nil {
		return "", nil, err
	}

	return reportMessage, replyMarkup, nil
}

// GetStudentsList ...
func (hlpr *Helper) GetStudentsList(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": student.School.ID,
	}).Info("get all students from database by school_id")
	students, err := hlpr.store.Student().FindByFullCourseSchoolID(student.FullCourse, student.School.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Info("students found")

	replyMarkup := &telebot.ReplyMarkup{}
	var interfaceSlice []model.Interface = make([]model.Interface, len(students))
	for i, v := range students {
		interfaceSlice[i] = v
	}

	rows := rowsWithButtons(interfaceSlice, callback)
	rows = append(rows, backToSchoolRow(replyMarkup, callback, student.School.ID))
	replyMarkup.Inline(rows...)

	return fmt.Sprintf("School: %v\n\nChoose a student from the list below:", student.School.Title), replyMarkup, nil
}

// BlockStudent ...
func (hlpr *Helper) BlockStudent(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student found")

	student.Active = false

	hlpr.logger.Info("block student")
	if err := hlpr.store.Student().Update(student); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student blocked")

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToStudentRow(replyMarkup, callback, student.ID))

	return fmt.Sprintf("Success! Student <b>%v</b> blocked", student.Account.Username), replyMarkup, nil
}

// UnblockStudent ...
func (hlpr *Helper) UnblockStudent(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student found")

	student.Active = true

	hlpr.logger.Info("unblock student")
	if err := hlpr.store.Student().Update(student); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student unblocked")

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToStudentRow(replyMarkup, callback, student.ID))

	return fmt.Sprintf("Success! Student <b>%v</b> unblocked", student.Account.Username), replyMarkup, nil
}

// SetStudent ...
func (hlpr *Helper) SetStudent(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student found")

	student.FullCourse = true

	hlpr.logger.Info("move student to full course")
	if err := hlpr.store.Student().Update(student); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student moved")

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToStudentRow(replyMarkup, callback, student.ID))

	return fmt.Sprintf("Success! Student <b>%v</b> updated. New type is %v", student.Account.Username, student.GetType()), replyMarkup, nil
}

// SetListener ...
func (hlpr *Helper) SetListener(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student found")

	student.FullCourse = false

	hlpr.logger.Info("move student to listeners")
	if err := hlpr.store.Student().Update(student); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student moved")

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToStudentRow(replyMarkup, callback, student.ID))

	return fmt.Sprintf("Success! Student <b>%v</b> updated. New type is %v", student.Account.Username, student.GetType()), replyMarkup, nil
}

func backToStudentRow(replyMarkup *telebot.ReplyMarkup, callback *model.Callback, studentID int64) telebot.Row {
	backToStudentCallback := &model.Callback{
		ID:          studentID,
		Type:        "student",
		Command:     "get",
		ListCommand: callback.ListCommand,
	}

	backToStudentsListCallback := &model.Callback{
		ID:          studentID,
		Type:        "student",
		Command:     "students_list",
		ListCommand: callback.ListCommand,
	}

	if callback.ListCommand == "get" {
		return replyMarkup.Row(
			replyMarkup.Data("<< Back to Student", backToStudentCallback.ToString()),
			replyMarkup.Data("<< Back to Students List", backToStudentsListCallback.ToString()),
		)
	}

	return replyMarkup.Row(replyMarkup.Data("<< Back to Students List", backToStudentsListCallback.ToString()))
}
