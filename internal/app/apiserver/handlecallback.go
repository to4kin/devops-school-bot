package apiserver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	maxRows = 3
)

func (srv *server) handleCallback(c telebot.Context) error {
	srv.logger.WithFields(logrus.Fields{
		"callback_data": c.Callback().Data[1:],
	}).Debug("handle callback")

	callbackData := strings.Split(c.Callback().Data[1:], "|")

	if len(callbackData) < 3 {
		srv.logger.WithFields(logrus.Fields{
			"callback_data": callbackData,
		}).Error("callback is not supported")
		return nil
	}

	callbackValue := callbackData[0]
	callbackType := callbackData[1]
	callbackAction := callbackData[2]

	switch callbackType {
	case "school":
		if callbackAction == "back_to_list" {
			return srv.handleSchools(c)
		}

		if callbackAction == "previous" {
			page, err := strconv.Atoi(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return nil
			}

			if page == 0 {
				return srv.handleSchools(c)
			}

			return srv.schoolNaviButtons(c, page)
		}

		if callbackAction == "next" {
			page, err := strconv.Atoi(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return nil
			}

			return srv.schoolNaviButtons(c, page)
		}

		if callbackAction == "re_activate" {
			srv.logger.WithFields(logrus.Fields{
				"title": callbackValue,
			}).Debug("get school by title")
			school, err := srv.store.School().FindByTitle(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return nil
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school found")

			if err := srv.store.School().ReActivate(school); err != nil {
				srv.logger.Error(err)
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school re-activated")

			return srv.schoolRespond(c, school)
		}

		if callbackAction == "finish" {
			srv.logger.WithFields(logrus.Fields{
				"title": callbackValue,
			}).Debug("get school by title")
			school, err := srv.store.School().FindByTitle(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return nil
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school found")

			if err := srv.store.School().Finish(school); err != nil {
				srv.logger.Error(err)
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school finished")

			return srv.schoolRespond(c, school)
		}

		if callbackAction == "get" {
			srv.logger.WithFields(logrus.Fields{
				"title": callbackValue,
			}).Debug("get school by title")
			school, err := srv.store.School().FindByTitle(callbackValue)
			if err != nil {
				srv.logger.Error(err)
				return nil
			}
			srv.logger.WithFields(school.LogrusFields()).Debug("school found")

			return srv.schoolRespond(c, school)
		}
	}

	return c.Edit("handleCallback")
}

func (srv *server) schoolRespond(c telebot.Context, school *model.School) error {
	srv.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get students by school_id")
	students, err := srv.store.Student().FindBySchoolID(school.ID)
	if err != nil {
		srv.logger.Error(err)
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Debug("students found")

	srv.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get homeworks by school_id")
	homeworks, err := srv.store.Homework().FindBySchoolID(school.ID)
	if err != nil {
		srv.logger.Error(err)
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
			fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
				school.Created.Year(), school.Created.Month(), school.Created.Day(),
				school.Created.Hour(), school.Created.Minute(), school.Created.Second()),
			len(students),
			len(homeworks),
			status,
		),
		&telebot.SendOptions{ParseMode: "HTML"},
		replyMarkup,
	)
}

func (srv *server) schoolNaviButtons(c telebot.Context, page int) error {
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
	btnNext := replyMarkup.Data("Next list >>", strconv.Itoa(page+1), "school", "next")
	btnPrevious := replyMarkup.Data("<< Previous list", strconv.Itoa(page-1), "school", "previous")

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
