package apiserver

import (
	"fmt"
	"sort"

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
		return srv.respondAlert(c, msgInternalError)
	}
	srv.logger.WithFields(account.LogrusFields()).Debug("account found")

	if !account.Superuser {
		srv.logger.WithFields(account.LogrusFields()).Debug("account has insufficient permissions")
		return c.EditOrSend(msgUserInsufficientPermissions, &telebot.SendOptions{ParseMode: "HTML"})
	}

	calbback := &model.Callback{
		Type: "school",
		ID:   0,
	}

	return srv.schoolsNaviButtons(c, calbback)
}

func (srv *server) schoolRespond(c telebot.Context, callback *model.Callback) error {
	srv.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get school by id")
	school, err := srv.store.School().FindByID(callback.ID)
	if err != nil {
		srv.logger.Error(err)
		return srv.respondAlert(c, msgInternalError)
	}
	srv.logger.WithFields(school.LogrusFields()).Debug("school found")

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
	if school.Active {
		status = iconGreenCircle
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Finish school", "finish", callback.ToString())))
	} else {
		status = iconRedCircle
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Re-Activate school", "re_activate", callback.ToString())))
	}
	if len(students) > 0 {
		studentCallback := &model.Callback{
			Type: "student",
			ID:   students[0].ID,
		}
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Students", "students_list", studentCallback.ToString())))
	}
	if len(homeworks) > 0 {
		homeworkCallback := &model.Callback{
			Type: "homework",
			ID:   homeworks[0].ID,
		}
		rows = append(rows, replyMarkup.Row(replyMarkup.Data("Homeworks", "homeworks_list", homeworkCallback.ToString())))
	}

	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to school list", "schools_list", callback.ToString())))
	replyMarkup.Inline(rows...)

	lessons := make(map[string]int)
	for _, homework := range homeworks {
		lessons[homework.Lesson.Title]++
	}

	// TODO: Go runs from a random offset for map iteration. We need a workaround if we need a sorted output for map
	var keys []string
	for k := range lessons {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	text := ""
	for _, k := range keys {
		text += fmt.Sprintf("%v - %d\n", k, lessons[k])
	}

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
			text,
		),
		&telebot.SendOptions{ParseMode: "HTML"},
		replyMarkup,
	)
}

func (srv *server) schoolsNaviButtons(c telebot.Context, callback *model.Callback) error {
	srv.logger.Debug("get all schools")
	schools, err := srv.store.School().FindAll()
	if err != nil {
		srv.logger.Error(err)
		return srv.respondAlert(c, msgInternalError)
	}
	srv.logger.WithFields(logrus.Fields{
		"count": len(schools),
	}).Debug("schools found")

	var interfaceSlice []model.Interface = make([]model.Interface, len(schools))
	for i, v := range schools {
		interfaceSlice[i] = v
	}
	rows := naviButtons(interfaceSlice, callback)

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(rows...)
	return c.EditOrSend("Choose a school from the list below:", replyMarkup)
}
