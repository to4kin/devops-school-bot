package apiserver

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
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
	callback.ID = student.School.ID
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to student list", "students_list", callback.ToString())))
	replyMarkup.Inline(rows...)

	return c.EditOrSend(
		fmt.Sprintf(
			msgUserInfo,
			student.Account.FirstName,
			student.Account.LastName,
			student.Account.Username,
			student.Account.Superuser,
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
		buttons = append(buttons, replyMarkup.Data(student.Account.Username, "get", studentCallback.ToString()))
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
