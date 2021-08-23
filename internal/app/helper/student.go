package helper

import (
	"fmt"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	studentText        string = "School: %v\n\nStudent info:\n\nFirst name: %v\nLast name: %v\nStatus: %v\n\n" + SysHomeworkGuide + "\n\nHomeworks:\n%v"
	studentsListText   string = "School: %v\n\nChoose a student from the list below:"
	blockStudentText   string = "Success! Student <b>%v</b> blocked"
	unblockStudentText string = "Success! Student <b>%v</b> unblocked"

	backToStudentText      string = "<< Back to Student"
	backToStudentsListText string = "<< Back to Students List"
)

// GetStudent ...
func GetStudent(store store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	student, err := store.Student().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}

	if student.Active {
		blockCallback := &model.Callback{
			ID:          callback.ID,
			Type:        "student",
			Command:     "block",
			ListCommand: callback.ListCommand,
		}
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Block student", blockCallback.ToString())))
	} else {
		unblockCallback := &model.Callback{
			ID:          callback.ID,
			Type:        "student",
			Command:     "unblock",
			ListCommand: callback.ListCommand,
		}
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Unblock student", unblockCallback.ToString())))
	}

	backToStudentsListCallback := &model.Callback{
		ID:          student.ID,
		Type:        "student",
		Command:     "students_list",
		ListCommand: callback.ListCommand,
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data(backToStudentsListText, backToStudentsListCallback.ToString())))
	rows = append(rows, backToSchoolRow(replyMarkup, student.School.ID))
	replyMarkup.Inline(rows...)

	homeworks, _ := store.Homework().FindByStudentID(callback.ID)
	lessons, _ := store.Lesson().FindBySchoolID(student.School.ID)

	text := ""
	for _, lesson := range lessons {
		counted := false
		for _, homework := range homeworks {
			if homework.Lesson.ID == lesson.ID {
				counted = true
				text += fmt.Sprintf("%v - %v\n", iconGreenCircle, lesson.Title)
			}
		}

		if !counted {
			text += fmt.Sprintf("%v - %v\n", iconRedCircle, lesson.Title)
		}
	}

	return fmt.Sprintf(
			studentText,
			student.School.Title,
			student.Account.FirstName,
			student.Account.LastName,
			student.GetStatusText(),
			text,
		),
		replyMarkup,
		nil
}

// GetStudentsList ...
func GetStudentsList(store store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	student, err := store.Student().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	students, err := store.Student().FindBySchoolID(student.School.ID)
	if err != nil {
		return "", nil, err
	}

	replyMarkup := &telebot.ReplyMarkup{}
	var interfaceSlice []model.Interface = make([]model.Interface, len(students))
	for i, v := range students {
		interfaceSlice[i] = v
	}

	rows := rowsWithButtons(interfaceSlice, callback)
	rows = append(rows, backToSchoolRow(replyMarkup, student.School.ID))
	replyMarkup.Inline(rows...)

	return fmt.Sprintf(studentsListText, student.School.Title), replyMarkup, nil
}

// BlockStudent ...
func BlockStudent(store store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	student, err := store.Student().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	student.Active = false

	if err := store.Student().Update(student); err != nil {
		return "", nil, err
	}

	replyMarkup := &telebot.ReplyMarkup{}
	fillStudentReplyMarkup(replyMarkup, callback)

	return fmt.Sprintf(blockStudentText, student.Account.Username), replyMarkup, nil
}

// UnblockStudent ...
func UnblockStudent(store store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	student, err := store.Student().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	student.Active = true

	if err := store.Student().Update(student); err != nil {
		return "", nil, err
	}

	replyMarkup := &telebot.ReplyMarkup{}
	fillStudentReplyMarkup(replyMarkup, callback)

	return fmt.Sprintf(unblockStudentText, student.Account.Username), replyMarkup, nil
}

func fillStudentReplyMarkup(replyMarkup *telebot.ReplyMarkup, callback *model.Callback) {
	if callback.ListCommand == "get" {
		backToStudentCallback := *callback
		backToStudentCallback.Command = "get"

		backToStudentsListCallback := *callback
		backToStudentsListCallback.Command = "students_list"

		replyMarkup.Inline(replyMarkup.Row(
			replyMarkup.Data(backToStudentText, backToStudentCallback.ToString()),
			replyMarkup.Data(backToStudentsListText, backToStudentsListCallback.ToString()),
		))
	}
}
