package apiserver

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) handleSchools(c telebot.Context) error {
	if !c.Message().Private() {
		return nil
	}

	srv.logger.WithFields(logrus.Fields{
		"telegram_id": c.Sender().ID,
	}).Debug("get account from database by telegram_id")
	account, err := srv.store.Account().FindByTelegramID(int64(c.Sender().ID))
	if err != nil {
		srv.logger.Error(err)
		return nil
	}
	srv.logger.WithFields(account.LogrusFields()).Debug("account found")

	if !account.Superuser {
		srv.logger.WithFields(account.LogrusFields()).Debug("account has insufficient permissions")
		return c.EditOrSend(msgUserInsufficientPermissions, &telebot.SendOptions{ParseMode: "HTML"})
	}

	return srv.schoolsNaviButtons(c, 0)
}

func (srv *server) schoolRespond(c telebot.Context, school *model.School) error {
	srv.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get students by school_id")
	students, err := srv.store.Student().FindBySchoolID(school.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			srv.logger.Debug(err)
		} else {
			srv.logger.Error(err)
		}
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Debug("students found")

	srv.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get homeworks by school_id")
	homeworks, err := srv.store.Homework().FindBySchoolID(school.ID)
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

	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}
	status := ""
	if school.Finished {
		status = "Finished"
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Re-Activate school", school.Title, "school", "re_activate")))
	} else {
		status = "Active"
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Finish school", school.Title, "school", "finish")))
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to school list", school.Title, "school", "back_to_list")))
	replyMarkup.Inline(rows...)

	return c.EditOrSend(
		fmt.Sprintf(
			msgSchoolInfo,
			school.Title,
			fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d",
				school.Created.Day(), school.Created.Month(), school.Created.Year(),
				school.Created.Hour(), school.Created.Minute(), school.Created.Second()),
			len(students),
			len(homeworks),
			status,
		),
		&telebot.SendOptions{ParseMode: "HTML"},
		replyMarkup,
	)
}

func (srv *server) schoolsNaviButtons(c telebot.Context, page int) error {
	srv.logger.Debug("get all schools")
	schools, err := srv.store.School().FindAll()
	if err != nil {
		srv.logger.Error(err)
		return nil
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(schools),
	}).Debug("schools found")

	var buttons []telebot.Btn
	replyMarkup := &telebot.ReplyMarkup{}
	for _, school := range schools {
		buttons = append(buttons, replyMarkup.Data(school.Title, school.Title, "school", "get"))
	}

	var rows []telebot.Row
	div, mod := len(schools)/2, len(schools)%2
	btnNext := replyMarkup.Data("Next page >>", strconv.Itoa(page+1), "school", "next")
	btnPrevious := replyMarkup.Data("<< Previous page", strconv.Itoa(page-1), "school", "previous")

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

	replyMarkup.Inline(rows...)
	return c.EditOrSend("Choose a school from the list below:", replyMarkup)
}
