package helper

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// GetUserReport ...
func (hlpr *Helper) GetUserReport(account *model.Account, school *model.School) (string, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"account_id": account.ID,
		"school_id":  school.ID,
	}).Debug("get student from database by account_id and school_id")
	student, err := hlpr.store.Student().FindByAccountIDSchoolID(account.ID, school.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			return ErrUserNotJoined, nil
		}

		return "", err
	}
	hlpr.logger.WithFields(student.LogrusFields()).Debug("student found")

	if !student.Active {
		return "Your student account is blocked!\n\nPlease contact mentors or teachers", nil
	}

	hlpr.logger.WithFields(logrus.Fields{
		"student_id": student.ID,
	}).Debug("get student homeworks from database by student_id")
	studentHomeworks, err := hlpr.store.Homework().FindByStudentID(student.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			return fmt.Sprintf("you haven't submitted your homework yet\n\n%v", sysHomeworkAdd), nil
		}
		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(studentHomeworks),
	}).Debug("student homeworks found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get lessons from database by school_id")
	allLessons, err := hlpr.store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(allLessons),
	}).Debug("lessons found")

	reportMessage := fmt.Sprintf(
		"Hello, @%v!\n\n"+msgStudentInfo,
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
func (hlpr *Helper) GetReport(school *model.School) (string, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get lessons from database by school_id")
	lessons, err := hlpr.store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(lessons),
	}).Debug("lessons found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get students from database by school_id")
	students, err := hlpr.store.Student().FindBySchoolID(school.ID)
	if err != nil {
		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Debug("students found")

	reportMessage, err := hlpr.prepareReportMsg(students, lessons)
	if err != nil {
		return "", err
	}

	return reportMessage, nil
}

// GetFullReport ...
func (hlpr *Helper) GetFullReport(school *model.School) (string, error) {
	fullReport, err := hlpr.GetLessonsReport(school)
	if err != nil {
		return "", err
	}

	reportMessage, err := hlpr.GetReport(school)
	if err != nil {
		return "", err
	}

	fullReport += "\n" + reportMessage

	return fullReport, nil
}

// GetLessonsReport ...
func (hlpr *Helper) GetLessonsReport(school *model.School) (string, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Debug("get lessons from database by school_id")
	lessons, err := hlpr.store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(lessons),
	}).Debug("lessons found")

	reportMessage := "<b>Homework list</b>\n"
	reportMessage += fmt.Sprintf("\n<b>Module: %v\n</b>", lessons[0].Module.Title)
	reportMessage += fmt.Sprintf("%v\n", lessons[0].Title)

	for i := 1; i < len(lessons); i++ {
		if lessons[i].Module.ID != lessons[i-1].Module.ID {
			reportMessage += fmt.Sprintf("\n<b>Module: %v\n</b>", lessons[i].Module.Title)
		}
		reportMessage += fmt.Sprintf("%v\n", lessons[i].Title)
	}

	return reportMessage, nil
}

func (hlpr *Helper) prepareReportMsg(students []*model.Student, lessons []*model.Lesson) (string, error) {
	reportMessage := "Academic performance\n\n<b><u>Name - Accepted/Not Provided - Type</u></b>\n<pre>"
	for _, student := range students {
		acceptedHomework := 0
		notProvidedHomework := 0

		homeworks, err := hlpr.store.Homework().FindByStudentID(student.ID)
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
