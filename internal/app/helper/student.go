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

	if student.Active {
		blockCallback := &model.Callback{
			Created:     time.Now(),
			Type:        "student",
			TypeID:      callback.TypeID,
			Command:     "block",
			ListCommand: callback.ListCommand,
		}
		if err := hlpr.prepareCallback(blockCallback); err != nil {
			return "", nil, err
		}
		buttons = append(buttons, replyMarkup.Data("Block student", blockCallback.GetStringID()))
	} else {
		unblockCallback := &model.Callback{
			Created:     time.Now(),
			Type:        "student",
			TypeID:      callback.TypeID,
			Command:     "unblock",
			ListCommand: callback.ListCommand,
		}
		if err := hlpr.prepareCallback(unblockCallback); err != nil {
			return "", nil, err
		}
		buttons = append(buttons, replyMarkup.Data("Unblock student", unblockCallback.GetStringID()))
	}

	if student.FullCourse {
		studentCallback := &model.Callback{
			Created:     time.Now(),
			Type:        "student",
			TypeID:      callback.TypeID,
			Command:     "set_listener",
			ListCommand: callback.ListCommand,
		}
		if err := hlpr.prepareCallback(studentCallback); err != nil {
			return "", nil, err
		}
		buttons = append(buttons, replyMarkup.Data("Change to listener", studentCallback.GetStringID()))
	} else {
		listenerCallback := &model.Callback{
			Created:     time.Now(),
			Type:        "student",
			TypeID:      callback.TypeID,
			Command:     "set_student",
			ListCommand: callback.ListCommand,
		}
		if err := hlpr.prepareCallback(listenerCallback); err != nil {
			return "", nil, err
		}
		buttons = append(buttons, replyMarkup.Data("Change to student", listenerCallback.GetStringID()))
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
		Created:     time.Now(),
		Type:        "student",
		TypeID:      student.ID,
		Command:     "students_list",
		ListCommand: callback.ListCommand,
	}
	if err := hlpr.prepareCallback(backToStudentsListCallback); err != nil {
		return "", nil, err
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to Students List", backToStudentsListCallback.GetStringID())))

	backRow, err := hlpr.backToSchoolRow(replyMarkup, callback.ListCommand, student.School.ID)
	if err != nil {
		return "", nil, err
	}
	rows = append(rows, backRow)
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

	return fmt.Sprintf("School: %v\n\nChoose a student from the list below:", student.School.Title), replyMarkup, nil
}

// BlockStudent ...
func (hlpr *Helper) BlockStudent(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.TypeID)
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
	backRow, err := hlpr.backToStudentRow(replyMarkup, callback.ListCommand, student.ID)
	if err != nil {
		return "", nil, err
	}
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("Success! Student <b>%v</b> blocked", student.Account.Username), replyMarkup, nil
}

// UnblockStudent ...
func (hlpr *Helper) UnblockStudent(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.TypeID)
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
	backRow, err := hlpr.backToStudentRow(replyMarkup, callback.ListCommand, student.ID)
	if err != nil {
		return "", nil, err
	}
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("Success! Student <b>%v</b> unblocked", student.Account.Username), replyMarkup, nil
}

// SetStudent ...
func (hlpr *Helper) SetStudent(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.TypeID)
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
	backRow, err := hlpr.backToStudentRow(replyMarkup, callback.ListCommand, student.ID)
	if err != nil {
		return "", nil, err
	}
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("Success! Student <b>%v</b> updated. New type is %v", student.Account.Username, student.GetType()), replyMarkup, nil
}

// SetListener ...
func (hlpr *Helper) SetListener(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get student from database by id")
	student, err := hlpr.store.Student().FindByID(callback.TypeID)
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
	backRow, err := hlpr.backToStudentRow(replyMarkup, callback.ListCommand, student.ID)
	if err != nil {
		return "", nil, err
	}
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("Success! Student <b>%v</b> updated. New type is %v", student.Account.Username, student.GetType()), replyMarkup, nil
}

func (hlpr *Helper) backToStudentRow(replyMarkup *telebot.ReplyMarkup, listCommand string, studentID int64) (telebot.Row, error) {
	backToStudentCallback := &model.Callback{
		Created:     time.Now(),
		Type:        "student",
		TypeID:      studentID,
		Command:     "get",
		ListCommand: listCommand,
	}
	if err := hlpr.prepareCallback(backToStudentCallback); err != nil {
		return nil, err
	}

	backToStudentsListCallback := &model.Callback{
		Created:     time.Now(),
		Type:        "student",
		TypeID:      studentID,
		Command:     "students_list",
		ListCommand: listCommand,
	}
	if err := hlpr.prepareCallback(backToStudentsListCallback); err != nil {
		return nil, err
	}

	if listCommand == "get" {
		return replyMarkup.Row(
			replyMarkup.Data("<< Back to Student", backToStudentCallback.GetStringID()),
			replyMarkup.Data("<< Back to Students List", backToStudentsListCallback.GetStringID()),
		), nil
	}

	return replyMarkup.Row(replyMarkup.Data("<< Back to Students List", backToStudentsListCallback.GetStringID())), nil
}
