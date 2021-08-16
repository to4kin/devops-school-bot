package apiserver

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) studentRespond(c telebot.Context, callback *model.Callback) error {
	srv.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get student from database by id")
	student, err := srv.store.Student().FindByID(callback.ID)
	if err != nil {
		srv.logger.Error(err)
		return srv.respondAlert(c, msgInternalError)
	}
	srv.logger.WithFields(student.LogrusFields()).Debug("student found")

	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}

	status := ""
	if student.Active {
		status = iconGreenCircle
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Block student", "block", callback.ToString())))
	} else {
		status = iconRedCircle
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Unblock student", "unblock", callback.ToString())))
	}

	backCallback := &model.Callback{
		Type: callback.Type,
		ID:   student.School.ID,
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to student list", "students_list", backCallback.ToString())))
	replyMarkup.Inline(rows...)

	srv.logger.WithFields(logrus.Fields{
		"student_id": callback.ID,
	}).Debug("get homeworks by student_id")
	homeworks, err := srv.store.Homework().FindByStudentID(callback.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			srv.logger.Debug(err)
		} else {
			srv.logger.Error(err)
		}
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Debug("homeworks found")

	srv.logger.WithFields(logrus.Fields{
		"school_id": student.School.ID,
	}).Debug("get all lessons from database by school_id")
	lessons, err := srv.store.Lesson().FindBySchoolID(student.School.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			srv.logger.Debug(err)
		} else {
			srv.logger.Error(err)
		}
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(lessons),
	}).Debug("lessons found")

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

	return c.EditOrSend(
		fmt.Sprintf(
			msgStudentInfo,
			student.School.Title,
			student.Account.FirstName,
			student.Account.LastName,
			status,
			text,
		),
		&telebot.SendOptions{ParseMode: "HTML"},
		replyMarkup,
	)
}

func (srv *server) studentsNaviButtons(c telebot.Context, callback *model.Callback) error {
	srv.logger.WithFields(logrus.Fields{
		"school_id": callback.ID,
	}).Debug("get all students by school_id")
	students, err := srv.store.Student().FindBySchoolID(callback.ID)
	if err != nil {
		srv.logger.Error(err)
		return srv.respondAlert(c, msgInternalError)
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Debug("students found")

	page := 0
	for i, student := range students {
		if callback.ID == student.ID {
			page = i / (maxRows * 2)
			break
		}
	}

	var buttons []telebot.Btn
	replyMarkup := &telebot.ReplyMarkup{}
	for _, student := range students {
		studentCallback := &model.Callback{
			Type: callback.Type,
			ID:   student.ID,
		}

		text := fmt.Sprintf("%v %v", iconRedCircle, student.Account.Username)
		if student.Active {
			text = fmt.Sprintf("%v %v", iconGreenCircle, student.Account.Username)
		}

		buttons = append(buttons, replyMarkup.Data(text, "get", studentCallback.ToString()))
	}

	var rows []telebot.Row
	div, mod := len(students)/2, len(students)%2

	nextCallback := &model.Callback{
		Type: callback.Type,
	}

	previousCallback := &model.Callback{
		Type: callback.Type,
	}

	if div >= maxRows*(page+1) {
		for i := maxRows * page; i < maxRows*(page+1); i++ {
			rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
		}

		nextCallback.ID = students[maxRows*2*(page+1)].ID
		btnNext := replyMarkup.Data("Next page >>", "next", nextCallback.ToString())

		if page > 0 {
			previousCallback.ID = students[maxRows*2*(page-1)].ID
			btnPrevious := replyMarkup.Data("<< Previous page", "previous", previousCallback.ToString())

			rows = append(rows, replyMarkup.Row(btnPrevious, btnNext))
		} else {
			rows = append(rows, replyMarkup.Row(btnNext))
		}
	} else {
		for i := maxRows * page; i < div; i++ {
			rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
		}
		if mod != 0 {
			rows = append(rows, replyMarkup.Row(buttons[div*2]))
		}
		if page > 0 {
			previousCallback.ID = students[maxRows*2*(page-1)].ID
			btnPrevious := replyMarkup.Data("<< Previous page", "previous", previousCallback.ToString())

			rows = append(rows, replyMarkup.Row(btnPrevious))
		}
	}

	schoolCallback := &model.Callback{
		Type: "school",
		ID:   callback.ID,
	}

	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to school", "get", schoolCallback.ToString())))

	replyMarkup.Inline(rows...)
	return c.EditOrSend("Choose a student from the list below:", replyMarkup)
}
