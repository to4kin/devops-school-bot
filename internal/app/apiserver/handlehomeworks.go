package apiserver

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

func (srv *server) homeworkRespond(c telebot.Context, callback *model.Callback) error {
	srv.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get homework from database by id")
	homework, err := srv.store.Homework().FindByID(callback.ID)
	if err != nil {
		srv.logger.Error(err)
		return srv.respondAlert(c, msgInternalError)
	}
	srv.logger.WithFields(homework.LogrusFields()).Debug("homework found")

	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}

	backCallback := &model.Callback{
		Type: "homework",
		ID:   homework.ID,
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to homework list", "homeworks_list", backCallback.ToString())))

	schoolCallback := &model.Callback{
		Type: "school",
		ID:   homework.Student.School.ID,
	}

	toSchool := replyMarkup.Data("<< Back to school", "get", schoolCallback.ToString())
	toSchoolList := replyMarkup.Data("<< Back to school list", "schools_list", schoolCallback.ToString())
	rows = append(rows, replyMarkup.Row(toSchool, toSchoolList))

	replyMarkup.Inline(rows...)

	return c.EditOrSend(
		fmt.Sprintf(
			msgHomeworkInfo,
			homework.Student.School.Title,
			homework.Lesson.Title,
		),
		&telebot.SendOptions{ParseMode: "HTML"},
		replyMarkup,
	)
}

func (srv *server) homeworksNaviButtons(c telebot.Context, callback *model.Callback) error {
	srv.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get homework from database by id")
	homework, err := srv.store.Homework().FindByID(callback.ID)
	if err != nil {
		srv.logger.Error(err)
		return srv.respondAlert(c, msgInternalError)
	}
	srv.logger.WithFields(homework.LogrusFields()).Debug("homework found")

	srv.logger.WithFields(logrus.Fields{
		"school_id": homework.Student.School.ID,
	}).Debug("get all homeworks by school_id")
	homeworks, err := srv.store.Homework().FindBySchoolID(homework.Student.School.ID)
	if err != nil {
		srv.logger.Error(err)
		return srv.respondAlert(c, msgInternalError)
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Debug("homeworks found")

	var interfaceSlice []model.Interface = make([]model.Interface, len(homeworks))
	for i, v := range homeworks {
		interfaceSlice[i] = v
	}
	rows := naviButtons(interfaceSlice, callback)

	schoolCallback := &model.Callback{
		Type: "school",
		ID:   homework.Student.School.ID,
	}

	replyMarkup := &telebot.ReplyMarkup{}

	toSchool := replyMarkup.Data("<< Back to school", "get", schoolCallback.ToString())
	toSchoolList := replyMarkup.Data("<< Back to school list", "schools_list", schoolCallback.ToString())

	rows = append(rows, replyMarkup.Row(toSchool, toSchoolList))

	replyMarkup.Inline(rows...)
	return c.EditOrSend(fmt.Sprintf("School: %v\n\nChoose a homework from the list below:", homework.Student.School.Title), replyMarkup)
}