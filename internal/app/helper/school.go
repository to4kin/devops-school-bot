package helper

import (
	"fmt"
	"sort"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

var (
	schoolText          string = "<b>%v</b>\n\nCreated: %v\nStudents: %d\nListeners: %d\nHomeworks: %d\nStatus: %v\n\nAccepted homeworks:\n%v"
	schoolsListText     string = "Choose a school from the list below:"
	schoolsNotFoundText string = "There are no schools in the database. Please add first"
	startSchoolText     string = "Success! School <b>%v</b> started"
	stopSchoolText      string = "Success! School <b>%v</b> finished"

	backToSchoolText      string = "<< Back to School"
	backToSchoolsListText string = "<< Back to Schools List"
)

// GetSchoolsList ...
func GetSchoolsList(str store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	var err error
	var schools []*model.School
	switch callback.ListCommand {
	case "start":
		schools, err = str.School().FindByActive(false)
	case "stop":
		schools, err = str.School().FindByActive(true)
	default:
		schools, err = str.School().FindAll()
	}

	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		return schoolsNotFoundText, nil, nil
	}

	var interfaceSlice []model.Interface = make([]model.Interface, len(schools))
	for i, v := range schools {
		interfaceSlice[i] = v
	}
	rows := rowsWithButtons(interfaceSlice, callback)

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(rows...)
	return schoolsListText, replyMarkup, nil
}

// GetSchool ...
func GetSchool(str store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	school, err := str.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	students, err := str.Student().FindByFullCourseSchoolID(true, school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	listeners, err := str.Student().FindByFullCourseSchoolID(false, school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	homeworks, err := str.Homework().FindBySchoolID(school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

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
		homeworksListCallback := &model.Callback{
			ID:          homeworks[0].ID,
			Type:        "homework",
			Command:     "homeworks_list",
			ListCommand: callback.ListCommand,
		}
		buttons = append(buttons, replyMarkup.Data("Homeworks", homeworksListCallback.ToString()))

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
	rows = append(rows, replyMarkup.Row(replyMarkup.Data(backToSchoolsListText, backToSchoolsListCallback.ToString())))
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
		schoolText,
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
func StartSchool(store store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	school, err := store.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	school.Active = true

	if err := store.School().Update(school); err != nil {
		return "", nil, err
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToSchoolRow(replyMarkup, callback, school.ID))

	return fmt.Sprintf(startSchoolText, school.Title), replyMarkup, nil
}

// StopSchool ...
func StopSchool(store store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	school, err := store.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	school.Active = false

	if err := store.School().Update(school); err != nil {
		return "", nil, err
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(backToSchoolRow(replyMarkup, callback, school.ID))

	return fmt.Sprintf(stopSchoolText, school.Title), replyMarkup, nil
}

// ReportSchool ...
func ReportSchool(str store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	school, err := str.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	reportMessage, err := GetReport(str, school)
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
func FullReportSchool(str store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	school, err := str.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	reportMessage, err := GetFullReport(str, school)
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
func GetSchoolHomeworks(str store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	school, err := str.School().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	reportMessage, err := GetLessonsReport(str, school)
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
			replyMarkup.Data(backToSchoolText, backToSchoolCallback.ToString()),
			replyMarkup.Data(backToSchoolsListText, backToSchoolsListCallback.ToString()),
		)
	}

	return replyMarkup.Row(replyMarkup.Data(backToSchoolsListText, backToSchoolsListCallback.ToString()))
}
