package helper

import (
	"fmt"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
	"gopkg.in/tucnak/telebot.v3"
)

// GetHomework ...
func GetHomework(store store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	homework, err := store.Homework().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

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
func GetHomeworksList(store store.Store, callback *model.Callback) (string, *telebot.ReplyMarkup, error) {
	homework, err := store.Homework().FindByID(callback.ID)
	if err != nil {
		return "", nil, err
	}

	homeworks, err := store.Homework().FindBySchoolID(homework.Student.School.ID)
	if err != nil {
		return "", nil, err
	}

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
