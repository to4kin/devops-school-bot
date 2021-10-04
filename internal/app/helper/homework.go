package helper

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gopkg.in/tucnak/telebot.v3"
)

// GetHomework ...
func (hlpr *Helper) GetHomework(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get homework from database by id")
	homework, err := hlpr.store.Homework().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(homework.LogrusFields()).Debug("homework found")

	var rows []telebot.Row
	replyMarkup := &telebot.ReplyMarkup{}

	homeworksListCallback := &model.Callback{
		ID:          homework.ID,
		Type:        "homework",
		Command:     "homeworks_list",
		ListCommand: callback.ListCommand,
	}
	rows = append(rows, replyMarkup.Row(replyMarkup.Data("<< Back to Homeworks List", homeworksListCallback.ToString())))
	rows = append(rows, backToSchoolRow(replyMarkup, callback, homework.Student.School.ID))
	replyMarkup.Inline(rows...)

	return fmt.Sprintf(
			"School: <b>%v</b>\n\nHomework info:\n\nTitle: %v\nModule: %v",
			homework.Student.School.Title,
			homework.Lesson.Title,
			homework.Lesson.Module.Title,
		),
		replyMarkup,
		nil
}

// GetHomeworksList ...
func (hlpr *Helper) GetHomeworksList(callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"id": callback.ID,
	}).Debug("get homework from database by id")
	homework, err := hlpr.store.Homework().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(homework.LogrusFields()).Debug("homework found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": homework.Student.School.ID,
	}).Debug("get homeworks from database by school_id")
	homeworks, err := hlpr.store.Homework().FindBySchoolID(homework.Student.School.ID)
	if err != nil {
		return "", nil, err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(homeworks),
	}).Debug("homeworks found")

	replyMarkup := &telebot.ReplyMarkup{}
	var interfaceSlice []model.Interface = make([]model.Interface, len(homeworks))
	for i, v := range homeworks {
		interfaceSlice[i] = v
	}

	interfaceSlice = removeDuplicate(interfaceSlice)

	rows := rowsWithButtons(interfaceSlice, callback)
	rows = append(rows, backToSchoolRow(replyMarkup, callback, homework.Student.School.ID))
	replyMarkup.Inline(rows...)

	return fmt.Sprintf(
			"School: <b>%v</b>\n\nChoose a homework from the list below:",
			homework.Student.School.Title,
		),
		replyMarkup,
		nil
}

func removeDuplicate(slice []model.Interface) []model.Interface {
	allKeys := make(map[string]bool)
	list := []model.Interface{}
	for _, item := range slice {
		if _, value := allKeys[item.GetButtonTitle()]; !value {
			allKeys[item.GetButtonTitle()] = true
			list = append(list, item)
		}
	}
	return list
}
