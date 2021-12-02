package helper

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

// GetStudent ...
func (hlpr *Helper) GetStudent(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student found")

	var buttons []telebot.Btn
	replyMarkup := &telebot.ReplyMarkup{}

	statusCallback := &model.Callback{
		Created:     time.Now(),
		Type:        "student",
		TypeID:      callback.TypeID,
		Command:     "update_status",
		ListCommand: callback.ListCommand,
	}
	if err := hlpr.prepareCallback(statusCallback); err != nil {
		return "", nil, err
	}
	buttons = append(buttons, replyMarkup.Data(fmt.Sprintf("Set to %v", (&model.Student{Active: !student.Active}).GetStatusText()), statusCallback.GetStringID()))

	typeCallback := &model.Callback{
		Created:     time.Now(),
		Type:        "student",
		TypeID:      callback.TypeID,
		Command:     "update_type",
		ListCommand: callback.ListCommand,
	}
	if err := hlpr.prepareCallback(typeCallback); err != nil {
		return "", nil, err
	}
	buttons = append(buttons, replyMarkup.Data(fmt.Sprintf("Change to %v", (&model.Student{FullCourse: !student.FullCourse}).GetType()), typeCallback.GetStringID()))

	var rows []telebot.Row
	div, mod := len(buttons)/2, len(buttons)%2
	for i := 0; i < div; i++ {
		rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
	}
	if mod != 0 {
		rows = append(rows, replyMarkup.Row(buttons[div*2]))
	}

	backToStudentsListCallback := &model.Callback{
		Created:     time.Now(),
		Type:        "student",
		TypeID:      student.ID,
		Command:     "students_list",
		ListCommand: callback.ListCommand,
	}
	if err := hlpr.prepareCallback(backToStudentsListCallback); err != nil {
		return "", nil, err
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data(fmt.Sprintf("<< Back to %vs List", student.GetType()), backToStudentsListCallback.GetStringID())))

	backRow, err := hlpr.backToSchoolRow(replyMarkup, callback.ListCommand, student.School.ID)
	if err != nil {
		return "", nil, err
	}
	rows = append(rows, backRow)
	replyMarkup.Inline(rows...)

	reportMessage := fmt.Sprintf(
		"<b>Account info:</b>\nFirst name: %v\nLast name: %v\n\n",
		student.Account.FirstName,
		student.Account.LastName,
	)

	message, err := hlpr.GetStudentReport(student)
	if err != nil {
		return "", nil, err
	}

	reportMessage += message

	return reportMessage, replyMarkup, nil
}

// GetStudentsList ...
func (hlpr *Helper) GetStudentsList(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.TypeID)
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

	rows, err := hlpr.rowsWithButtons(interfaceSlice, callback)
	if err != nil {
		return "", nil, err
	}

	backRow, err := hlpr.backToSchoolRow(replyMarkup, callback.ListCommand, student.School.ID)
	if err != nil {
		return "", nil, err
	}
	rows = append(rows, backRow)
	replyMarkup.Inline(rows...)

	return fmt.Sprintf("School: %v\n\nChoose a %v from the list below:", student.School.Title, student.GetType()), replyMarkup, nil
}

// UpdateStudentStatus ...
func (hlpr *Helper) UpdateStudentStatus(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student found")

	student.Active = !student.Active

	hlpr.logger.Info("change student status")
	if err := hlpr.store.Student().Update(student); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student status changed")

	replyMarkup := &telebot.ReplyMarkup{}
	backRow, err := hlpr.backToStudentRow(replyMarkup, callback.ListCommand, student)
	if err != nil {
		return "", nil, err
	}
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("Success! %v <b>%v</b> status changed to %v", student.GetType(), student.Account.GetFullName(), student.GetStatusText()), replyMarkup, nil
}

// UpdateStudentType ...
func (hlpr *Helper) UpdateStudentType(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student found")

	student.FullCourse = !student.FullCourse

	hlpr.logger.Info("move student to " + student.GetType())
	if err := hlpr.store.Student().Update(student); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Info("student moved")

	replyMarkup := &telebot.ReplyMarkup{}
	backRow, err := hlpr.backToStudentRow(replyMarkup, callback.ListCommand, student)
	if err != nil {
		return "", nil, err
	}
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("Success! Student <b>%v</b> updated. New type is %v", student.Account.GetFullName(), student.GetType()), replyMarkup, nil
}

func (hlpr *Helper) backToStudentRow(replyMarkup *telebot.ReplyMarkup, listCommand string, student *model.Student) (telebot.Row, error) {
	backToStudentCallback := &model.Callback{
		Created:     time.Now(),
		Type:        "student",
		TypeID:      student.ID,
		Command:     "get",
		ListCommand: listCommand,
	}
	if err := hlpr.prepareCallback(backToStudentCallback); err != nil {
		return nil, err
	}

	backToStudentsListCallback := &model.Callback{
		Created:     time.Now(),
		Type:        "student",
		TypeID:      student.ID,
		Command:     "students_list",
		ListCommand: listCommand,
	}
	if err := hlpr.prepareCallback(backToStudentsListCallback); err != nil {
		return nil, err
	}

	if listCommand == "get" {
		return replyMarkup.Row(
			replyMarkup.Data(fmt.Sprintf("<< Back to %v", student.GetType()), backToStudentCallback.GetStringID()),
			replyMarkup.Data(fmt.Sprintf("<< Back to %vs List", student.GetType()), backToStudentsListCallback.GetStringID()),
		), nil
	}

	return replyMarkup.Row(replyMarkup.Data(fmt.Sprintf("<< Back to %vs List", student.GetType()), backToStudentsListCallback.GetStringID())), nil
}
