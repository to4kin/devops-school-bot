package helper

import (
	"fmt"

	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

var (
	report string = "Academic performance\n\n<b><u>Name - Accepted/Not Provided - Type</u></b>\n"

	homeworksListReport string = "<b>Homework list</b>\n\n"
	homeworkNotProvided string = "you haven't submitted your homework yet\n\n" + sysHomeworkAdd
	homeworkReport      string = "Hello, @%v!\n\n" + msgStudentInfo

	studentIsBlocked string = "Your student account is blocked!\n\nPlease contact mentors or teachers"
)

// GetUserReport ...
func GetUserReport(str store.Store, account *model.Account, school *model.School) (string, error) {
	student, err := str.Student().FindByAccountIDSchoolID(account.ID, school.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			return ErrUserNotJoined, nil
		}

		return "", err
	}

	if !student.Active {
		return studentIsBlocked, nil
	}

	studentHomeworks, err := str.Homework().FindByStudentID(student.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			return homeworkNotProvided, nil
		}
		return "", err
	}

	allLessons, err := str.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		return "", err
	}

	reportMessage := fmt.Sprintf(
		homeworkReport,
		account.Username,
		student.Account.FirstName,
		student.Account.LastName,
		student.GetType(),
		student.GetStatusText(),
	)

	if student.FullCourse {
		reportMessage += "\n\n" + SysStudentGuide
		reportMessage += fmt.Sprintf("\n\nYour progress in <b>%v</b>:\n", school.Title)

		for _, lesson := range allLessons {
			counted := false
			for _, homework := range studentHomeworks {
				if homework.Lesson.ID == lesson.ID {
					counted = true
					reportMessage += fmt.Sprintf("%v - %v\n", iconGreenCircle, lesson.Title)
				}
			}

			if !counted {
				reportMessage += fmt.Sprintf("%v - %v\n", iconRedCircle, lesson.Title)
			}
		}
	} else {
		reportMessage += "\n\n" + SysListenerGuide
		reportMessage += fmt.Sprintf("\n\nYour progress in <b>%v</b>:\n", school.Title)
		studentModules := []*model.Module{}

		for _, homework := range studentHomeworks {
			studentModules = appendModule(studentModules, homework.Lesson.Module)
		}

		for _, module := range studentModules {
			for _, lesson := range allLessons {
				if module.ID == lesson.Module.ID {
					counted := false
					for _, homework := range studentHomeworks {
						if homework.Lesson.ID == lesson.ID {
							counted = true
							reportMessage += fmt.Sprintf("%v - %v\n", iconGreenCircle, lesson.Title)
						}
					}

					if !counted {
						reportMessage += fmt.Sprintf("%v - %v\n", iconRedCircle, lesson.Title)
					}
				}
			}
		}
	}

	return reportMessage, nil
}

// GetReport ...
func GetReport(store store.Store, school *model.School) (string, error) {
	lessons, err := store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		return "", err
	}

	students, err := store.Student().FindBySchoolID(school.ID)
	if err != nil {
		return "", err
	}

	reportMessage, err := prepareReportMsg(store, students, lessons)
	if err != nil {
		return "", err
	}

	return reportMessage, nil
}

// GetFullReport ...
func GetFullReport(store store.Store, school *model.School) (string, error) {
	fullReport, err := GetLessonsReport(store, school)
	if err != nil {
		return "", err
	}

	reportMessage, err := GetReport(store, school)
	if err != nil {
		return "", err
	}

	fullReport += "\n" + reportMessage

	return fullReport, nil
}

// GetLessonsReport ...
func GetLessonsReport(store store.Store, school *model.School) (string, error) {
	lessons, err := store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		return "", err
	}

	reportMessage := homeworksListReport
	for i, lesson := range lessons {
		reportMessage += fmt.Sprintf("%d - %v\n", i+1, lesson.Title)
	}

	return reportMessage, nil
}

func prepareReportMsg(str store.Store, students []*model.Student, lessons []*model.Lesson) (string, error) {
	reportMessage := report + "<pre>"
	for _, student := range students {
		acceptedHomework := 0
		notProvidedHomework := 0

		homeworks, err := str.Homework().FindByStudentID(student.ID)
		if err != nil && err != store.ErrRecordNotFound {
			return "", err
		}

		if student.FullCourse {
			for _, lesson := range lessons {
				counted := false
				for _, homework := range homeworks {
					if homework.Lesson.ID == lesson.ID {
						counted = true
						acceptedHomework++
					}
				}

				if !counted {
					notProvidedHomework++
				}
			}
		} else {
			studentModules := []*model.Module{}

			for _, homework := range homeworks {
				studentModules = appendModule(studentModules, homework.Lesson.Module)
			}

			for _, module := range studentModules {
				for _, lesson := range lessons {
					if module.ID == lesson.Module.ID {
						counted := false
						for _, homework := range homeworks {
							if homework.Lesson.ID == lesson.ID {
								counted = true
								acceptedHomework++
							}
						}

						if !counted {
							notProvidedHomework++
						}
					}
				}
			}
		}

		reportMessage += fmt.Sprintf("%v %v - %d/%d - %v\n",
			student.Account.FirstName, student.Account.LastName, acceptedHomework, notProvidedHomework, student.GetType())
	}
	reportMessage += "</pre>"

	return reportMessage, nil
}

func appendModule(slice []*model.Module, homework *model.Module) []*model.Module {
	for _, elem := range slice {
		if elem.ID == homework.ID {
			return slice
		}
	}
	return append(slice, homework)
}
