package apiserver

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) studentsList(c telebot.Context, schoolID int64, page int) error {
	srv.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Debug("get account from database by telegram_id")
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		srv.logger.Error(err)
		return c.Respond(&telebot.CallbackResponse{
			Text:      msgInternalError,
			ShowAlert: true,
		})
	}
	srv.logger.WithFields(account.LogrusFields()).Debug("account found")

	if !account.Superuser {
		srv.logger.WithFields(account.LogrusFields()).Debug("account has insufficient permissions")
		return c.EditOrSend(msgUserInsufficientPermissions, &telebot.SendOptions{ParseMode: "HTML"})
	}

	return srv.studentsNaviButtons(c, schoolID, page)
}

func (srv *server) studentRespond(c telebot.Context, student *model.Student, fromPage int) error {
	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}
	rows = append(
		rows,
		replyMarkup.Row(
			replyMarkup.Data(
				"<< Back to student list",
				strconv.FormatInt(student.ID, 10),
				"student",
				"back_to_list",
				strconv.Itoa(fromPage),
				strconv.FormatInt(student.School.ID, 10),
			),
		),
	)
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

func (srv *server) studentsNaviButtons(c telebot.Context, schoolID int64, page int) error {
	srv.logger.WithFields(logrus.Fields{
		"school_id": schoolID,
	}).Debug("get all students by school_id")
	students, err := srv.store.Student().FindBySchoolID(schoolID)
	if err != nil {
		srv.logger.Error(err)
		return c.Respond(&telebot.CallbackResponse{
			Text:      msgInternalError,
			ShowAlert: true,
		})
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Debug("students found")

	var buttons []telebot.Btn
	replyMarkup := &telebot.ReplyMarkup{}
	for _, student := range students {
		buttons = append(buttons, replyMarkup.Data(student.Account.Username, strconv.FormatInt(student.ID, 10), "student", "get", strconv.Itoa(page)))
	}

	var rows []telebot.Row
	div, mod := len(students)/2, len(students)%2
	btnNext := replyMarkup.Data("Next page >>", strconv.Itoa(page+1), "student", "next", strconv.Itoa(page), strconv.FormatInt(schoolID, 10))
	btnPrevious := replyMarkup.Data("<< Previous page", strconv.Itoa(page-1), "student", "previous", strconv.Itoa(page), strconv.FormatInt(schoolID, 10))

	if div >= maxRows*(page+1) {
		for i := maxRows * page; i < maxRows*(page+1); i++ {
			rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
		}
		if page > 0 {
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
			rows = append(rows, replyMarkup.Row(btnPrevious))
		}
	}

	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to school", strconv.FormatInt(schoolID, 10), "school", "get", strconv.Itoa(page))))

	replyMarkup.Inline(rows...)
	return c.EditOrSend("Choose a student from the list below:", replyMarkup)
}
