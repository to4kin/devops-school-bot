package helper

import (
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

// GetSchoolsList ...
func (hlpr *Helper) GetSchoolsList(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	var err error
	var schools []*model.School
	switch callback.ListCommand {
	case "start":
		hlpr.logger.Debug("get all inactive schools from database")
		schools, err = hlpr.store.School().FindByActive(false)
	case "stop":
		hlpr.logger.Debug("get all active schools from database")
		schools, err = hlpr.store.School().FindByActive(true)
	default:
		hlpr.logger.Debug("get all schools from database")
		schools, err = hlpr.store.School().FindAll()
	}

	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		return "There are no schools in the database. Please add first", nil, nil
	}

	hlpr.logger.WithFields(logrus.Fields{
		"count": len(schools),
	}).Debug("schools found")

	var interfaceSlice []model.Interface = make([]model.Interface, len(schools))
	for i, v := range schools {
		interfaceSlice[i] = v
	}
	rows := rowsWithButtons(interfaceSlice, callback)

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(rows...)
	return "Choose a school from the list below:", replyMarkup, nil
}

// GetSchool ...
func (hlpr *Helper) GetSchool(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get school from database by id")
	school, err := hlpr.store.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Debug("school found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get all students from database by school_id")
	students, err := hlpr.store.Student().FindByFullCourseSchoolID(true, school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Debug("students found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get all listeners from database by school_id")
	listeners, err := hlpr.store.Student().FindByFullCourseSchoolID(false, school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(listeners),
	}).Debug("listeners found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get all homeworks from database by school_id")
	homeworks, err := hlpr.store.Homework().FindBySchoolID(school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Debug("homeworks found")

	var buttons []telebot.Btn
	replyMarkup := &telebot.ReplyMarkup{}
	if school.Active {
		stopSchoolCallback := *callback
		stopSchoolCallback.Command = "stop"
		buttons = append(buttons, replyMarkup.Data("Stop school", stopSchoolCallback.ToString()))
	} else {
		startSchoolCallback := *callback
		startSchoolCallback.Command = "start"
		buttons = append(buttons, replyMarkup.Data("Start school", startSchoolCallback.ToString()))
	}
	if len(students) > 0 {
		studentsListCallback := &model.Callback{
			ID:          students[0].ID,
			Type:        "student",
			Command:     "students_list",
			ListCommand: callback.ListCommand,
		}
		buttons = append(buttons, replyMarkup.Data("Students", studentsListCallback.ToString()))
	}
	if len(listeners) > 0 {
		listenersListCallback := &model.Callback{
			ID:          listeners[0].ID,
			Type:        "student",
			Command:     "students_list",
			ListCommand: callback.ListCommand,
		}
		buttons = append(buttons, replyMarkup.Data("Listeners", listenersListCallback.ToString()))
	}
	if len(homeworks) > 0 {
		reportCallback := *callback
		reportCallback.Command = "report"
		buttons = append(buttons, replyMarkup.Data("Report", reportCallback.ToString()))

		fullReportCallback := *callback
		fullReportCallback.Command = "full_report"
		buttons = append(buttons, replyMarkup.Data("Full Report", fullReportCallback.ToString()))
	}

	var rows []telebot.Row
	div, mod := len(buttons)/2, len(buttons)%2
	for i := 0; i < div; i++ {
		rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
	}
	if mod != 0 {
		rows = append(rows, replyMarkup.Row(buttons[div*2]))
	}

	backToSchoolsListCallback := *callback
	backToSchoolsListCallback.Command = "schools_list"
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to Schools List", backToSchoolsListCallback.ToString())))
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

	homeworksList := ""
	for _, k := range keys {
		homeworksList += fmt.Sprintf("%v - %d\n", k, lessons[k])
	}

	return fmt.Sprintf(
		"<b>%v</b>\n\nCreated: %v\nStudents: %d\nListeners: %d\nHomeworks: %d\nStatus: %v\n\nAccepted homeworks:\n%v",
		school.Title,
		fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d",
			school.Created.Day(), school.Created.Month(), school.Created.Year(),
			school.Created.Hour(), school.Created.Minute(), school.Created.Second()),
		len(students),
		len(listeners),
		len(homeworks),
		school.GetStatusText(),
		homeworksList,
	), replyMarkup, nil
}

// StartSchool ...
func (hlpr *Helper) StartSchool(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get school from database by id")
	school, err := hlpr.store.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Debug("school found")

	school.Active = true

	hlpr.logger.Debug("start school")
	if err := hlpr.store.School().Update(school); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Debug("school started")

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToSchoolRow(replyMarkup, callback, school.ID))

	return fmt.Sprintf("Success! School <b>%v</b> started", school.Title), replyMarkup, nil
}

// StopSchool ...
func (hlpr *Helper) StopSchool(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get school from database by id")
	school, err := hlpr.store.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Debug("school found")

	school.Active = false

	hlpr.logger.Debug("stop school")
	if err := hlpr.store.School().Update(school); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Debug("school stopped")

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToSchoolRow(replyMarkup, callback, school.ID))

	return fmt.Sprintf("Success! School <b>%v</b> finished", school.Title), replyMarkup, nil
}

// ReportSchool ...
func (hlpr *Helper) ReportSchool(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get school from database by id")
	school, err := hlpr.store.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Debug("school found")

	reportMessage, err := hlpr.GetReport(school)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		reportMessage = ErrReportNotFound
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToSchoolRow(replyMarkup, callback, school.ID))

	return fmt.Sprintf("School <b>%v</b>\n\n%v", school.Title, reportMessage), replyMarkup, nil
}

// FullReportSchool ...
func (hlpr *Helper) FullReportSchool(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get school from database by id")
	school, err := hlpr.store.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Debug("school found")

	reportMessage, err := hlpr.GetFullReport(school)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		reportMessage = ErrReportNotFound
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToSchoolRow(replyMarkup, callback, school.ID))

	return fmt.Sprintf("School <b>%v</b>\n\n%v", school.Title, reportMessage), replyMarkup, nil
}

// GetSchoolHomeworks ...
func (hlpr *Helper) GetSchoolHomeworks(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get school from database by id")
	school, err := hlpr.store.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Debug("school found")

	reportMessage, err := hlpr.GetLessonsReport(school)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		reportMessage = ErrReportNotFound
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToSchoolRow(replyMarkup, callback, school.ID))

	return fmt.Sprintf("School <b>%v</b>\n\n%v", school.Title, reportMessage), replyMarkup, nil
}

func backToSchoolRow(replyMarkup *telebot.ReplyMarkup, callback *model.Callback, schoolID int64) telebot.Row {
	backToSchoolCallback := &model.Callback{
		ID:          schoolID,
		Type:        "school",
		Command:     "get",
		ListCommand: callback.ListCommand,
	}

	backToSchoolsListCallback := &model.Callback{
		ID:          schoolID,
		Type:        "school",
		Command:     "schools_list",
		ListCommand: callback.ListCommand,
	}

	if callback.ListCommand == "get" {
		return replyMarkup.Row(
			replyMarkup.Data("<< Back to School", backToSchoolCallback.ToString()),
			replyMarkup.Data("<< Back to Schools List", backToSchoolsListCallback.ToString()),
		)
	}

	return replyMarkup.Row(replyMarkup.Data("<< Back to Schools List", backToSchoolsListCallback.ToString()))
}
