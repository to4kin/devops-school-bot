package helper

import (
	"fmt"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

// GetSchoolsList returns all schools from database
// The message is populated with buttons
func (hlpr *Helper) GetSchoolsList(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.Info("get all schools from database")
	schools, err := hlpr.store.School().FindAll()
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		return "There are no schools in the database. Please add first", nil, nil
	}

	hlpr.logger.WithFields(logrus.Fields{
		"count": len(schools),
	}).Info("schools found")

	var interfaceSlice []model.Interface = make([]model.Interface, len(schools))
	for i, v := range schools {
		interfaceSlice[i] = v
	}
	rows, err := hlpr.rowsWithButtons(interfaceSlice, callback)
	if err != nil {
		return "", nil, err
	}

	replyMarkup := &telebot.ReplyMarkup{}
	replyMarkup.Inline(rows...)
	return "Choose a school from the list below:", replyMarkup, nil
}

// GetSchool returns School card, where
// Buttons Start school/Stop school depend on Active
// Button Students depend on count of students (>0)
// Button Listeners depend on count of Listeners (>0)
// Buttons Report/Full Report and Homeworks depend on cound of Homeworks (>0)
func (hlpr *Helper) GetSchool(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get school from database by id")
	school, err := hlpr.store.School().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Info("school found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Info("get all students from database by school_id")
	students, err := hlpr.store.Student().FindByFullCourseSchoolID(true, school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Info("students found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Info("get all listeners from database by school_id")
	listeners, err := hlpr.store.Student().FindByFullCourseSchoolID(false, school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(listeners),
	}).Info("listeners found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Info("get all homeworks from database by school_id")
	homeworks, err := hlpr.store.Homework().FindBySchoolID(school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Info("homeworks found")

	var buttons []telebot.Btn
	replyMarkup := &telebot.ReplyMarkup{}
	stopSchoolCallback := &model.Callback{
		Created:     time.Now(),
		Type:        callback.Type,
		TypeID:      callback.TypeID,
		Command:     "update_status",
		ListCommand: callback.ListCommand,
	}
	if err := hlpr.prepareCallback(stopSchoolCallback); err != nil {
		return "", nil, err
	}

	buttons = append(buttons, replyMarkup.Data(fmt.Sprintf("Set to %v", (&model.School{Active: !school.Active}).GetStatusText()), stopSchoolCallback.GetStringID()))

	if len(students) > 0 {
		studentsListCallback := &model.Callback{
			Created:     time.Now(),
			Type:        "student",
			TypeID:      students[0].ID,
			Command:     "students_list",
			ListCommand: callback.ListCommand,
		}
		if err := hlpr.prepareCallback(studentsListCallback); err != nil {
			return "", nil, err
		}
		buttons = append(buttons, replyMarkup.Data("Students", studentsListCallback.GetStringID()))
	}
	if len(listeners) > 0 {
		listenersListCallback := &model.Callback{
			Created:     time.Now(),
			Type:        "student",
			TypeID:      listeners[0].ID,
			Command:     "students_list",
			ListCommand: callback.ListCommand,
		}
		if err := hlpr.prepareCallback(listenersListCallback); err != nil {
			return "", nil, err
		}
		buttons = append(buttons, replyMarkup.Data("Listeners", listenersListCallback.GetStringID()))
	}
	if len(homeworks) > 0 {
		reportCallback := &model.Callback{
			Created:     time.Now(),
			Type:        callback.Type,
			TypeID:      callback.TypeID,
			Command:     "report",
			ListCommand: callback.ListCommand,
		}
		if err := hlpr.prepareCallback(reportCallback); err != nil {
			return "", nil, err
		}
		buttons = append(buttons, replyMarkup.Data("Report", reportCallback.GetStringID()))

		fullReportCallback := &model.Callback{
			Created:     time.Now(),
			Type:        callback.Type,
			TypeID:      callback.TypeID,
			Command:     "full_report",
			ListCommand: callback.ListCommand,
		}
		if err := hlpr.prepareCallback(fullReportCallback); err != nil {
			return "", nil, err
		}
		buttons = append(buttons, replyMarkup.Data("Full Report", fullReportCallback.GetStringID()))

		homeworksCallback := &model.Callback{
			Created:     time.Now(),
			Type:        callback.Type,
			TypeID:      callback.TypeID,
			Command:     "homeworks",
			ListCommand: callback.ListCommand,
		}
		if err := hlpr.prepareCallback(homeworksCallback); err != nil {
			return "", nil, err
		}
		buttons = append(buttons, replyMarkup.Data("Homeworks", homeworksCallback.GetStringID()))
	}

	var rows []telebot.Row
	div, mod := len(buttons)/2, len(buttons)%2
	for i := 0; i < div; i++ {
		rows = append(rows, replyMarkup.Row(buttons[i*2], buttons[i*2+1]))
	}
	if mod != 0 {
		rows = append(rows, replyMarkup.Row(buttons[div*2]))
	}

	backToSchoolsListCallback := &model.Callback{
		Created:     time.Now(),
		Type:        callback.Type,
		TypeID:      callback.TypeID,
		Command:     "schools_list",
		ListCommand: callback.ListCommand,
	}
	if err := hlpr.prepareCallback(backToSchoolsListCallback); err != nil {
		return "", nil, err
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to Schools List", backToSchoolsListCallback.GetStringID())))
	replyMarkup.Inline(rows...)

	lessons := make(map[string]int)
	for _, homework := range homeworks {
		if homework.Active {
			lessons[homework.Lesson.Title]++
		}
	}

	// NOTE: Go runs from a random offset for map iteration. We need a workaround if we need a sorted output for map
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

// UpdateSchoolStatus updates School status in database
func (hlpr *Helper) UpdateSchoolStatus(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get school from database by id")
	school, err := hlpr.store.School().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Info("school found")

	school.Active = !school.Active

	hlpr.logger.Info("change school status")
	if err := hlpr.store.School().Update(school); err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Info("status changed")

	replyMarkup := &telebot.ReplyMarkup{}
	backRow, err := hlpr.backToSchoolRow(replyMarkup, callback.ListCommand, school.ID)
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("Success! School <b>%v</b> status changed to %v", school.Title, school.GetStatusText()), replyMarkup, nil
}

// ReportSchool returns report for School
func (hlpr *Helper) ReportSchool(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get school from database by id")
	school, err := hlpr.store.School().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Info("school found")

	reportMessage, err := hlpr.GetReport(school)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		reportMessage = ErrReportNotFound
	}

	replyMarkup := &telebot.ReplyMarkup{}
	backRow, err := hlpr.backToSchoolRow(replyMarkup, callback.ListCommand, school.ID)
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("School <b>%v</b>\n\n%v", school.Title, reportMessage), replyMarkup, nil
}

// FullReportSchool return Dull Report for School
func (hlpr *Helper) FullReportSchool(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get school from database by id")
	school, err := hlpr.store.School().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Info("school found")

	reportMessage, err := hlpr.GetFullReport(school)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		reportMessage = ErrReportNotFound
	}

	replyMarkup := &telebot.ReplyMarkup{}
	backRow, err := hlpr.backToSchoolRow(replyMarkup, callback.ListCommand, school.ID)
	replyMarkup.Inline(backRow)

	return fmt.Sprintf("School <b>%v</b>\n\n%v", school.Title, reportMessage), replyMarkup, nil
}

// GetSchoolHomeworks returns Homeworks for School
func (hlpr *Helper) GetSchoolHomeworks(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.TypeID,
	}).Info("get school from database by id")
	school, err := hlpr.store.School().FindByID(callback.TypeID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(school.LogrusFields()).Info("school found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Info("get all homeworks from database by school_id")
	homeworks, err := hlpr.store.Homework().FindBySchoolID(school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Info("homeworks found")

	reportMessage, err := hlpr.GetLessonsReport(school)
	if err != nil && err != store.ErrRecordNotFound {
		return "", nil, err
	}

	if err == store.ErrRecordNotFound {
		reportMessage = ErrReportNotFound
	}

	replyMarkup := &telebot.ReplyMarkup{}
	//var interfaceSlice []model.Interface = make([]model.Interface, len(homeworks))
	//for i, v := range homeworks {
	//	interfaceSlice[i] = v
	//}
	//interfaceSlice = removeDuplicate(interfaceSlice)

	//homeworkCallback := &model.Callback{
	//	Type:        "homework",
	//	Command:     "get",
	//	ListCommand: "get",
	//}

	var rows []telebot.Row
	//rows := rowsWithButtons(interfaceSlice, homeworkCallback)
	backRow, err := hlpr.backToSchoolRow(replyMarkup, callback.ListCommand, school.ID)
	rows = append(rows, backRow)
	replyMarkup.Inline(rows...)

	return fmt.Sprintf("School <b>%v</b>\n\n%v", school.Title, reportMessage), replyMarkup, nil
}

func (hlpr *Helper) backToSchoolRow(replyMarkup *telebot.ReplyMarkup, listCommand string, schoolID int64) (telebot.Row, error) {
	backToSchoolCallback := &model.Callback{
		Created:     time.Now(),
		Type:        "school",
		TypeID:      schoolID,
		Command:     "get",
		ListCommand: listCommand,
	}

	if err := hlpr.prepareCallback(backToSchoolCallback); err != nil {
		return nil, err
	}

	backToSchoolsListCallback := &model.Callback{
		Created:     time.Now(),
		Type:        "school",
		TypeID:      schoolID,
		Command:     "schools_list",
		ListCommand: listCommand,
	}

	if err := hlpr.prepareCallback(backToSchoolsListCallback); err != nil {
		return nil, err
	}

	if listCommand == "get" {
		return replyMarkup.Row(
			replyMarkup.Data("<< Back to School", backToSchoolCallback.GetStringID()),
			replyMarkup.Data("<< Back to Schools List", backToSchoolsListCallback.GetStringID()),
		), nil
	}

	return replyMarkup.Row(replyMarkup.Data("<< Back to Schools List", backToSchoolsListCallback.GetStringID())), nil
}
