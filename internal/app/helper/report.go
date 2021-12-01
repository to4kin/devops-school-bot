package helper

import (
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// GetUserReport returns report for User
// Account should be a vald Account struct
// School can be nil, in this case report will be generated for all schools where User is a student
func (hlpr *Helper) GetUserReport(account *model.Account, school *model.School) (string, error) {

	students := []*model.Student{}
	if school != nil {
		hlpr.logger.WithFields(logrus.Fields{
			"account_id": account.ID,
			"school_id":  school.ID,
		}).Info("get student from database by account_id and school_id")
		st, err := hlpr.store.Student().FindByAccountIDSchoolID(account.ID, school.ID)
		if err != nil {
			if err == store.ErrRecordNotFound {
				return ErrUserNotJoined, nil
			}

			return "", err
		}
		hlpr.logger.WithFields(st.LogrusFields()).Info("student found")
		students = append(students, st)
	} else {
		hlpr.logger.WithFields(logrus.Fields{
			"account_id": account.ID,
		}).Info("get students from database by account_id")
		sts, err := hlpr.store.Student().FindByAccountID(account.ID)
		if err != nil {
			if err == store.ErrRecordNotFound {
				return ErrUserNotJoined, nil
			}

			return "", err
		}
		hlpr.logger.WithFields(logrus.Fields{
			"count": len(sts),
		}).Info("student found")

		students = sts
	}

	reportMessage := fmt.Sprintf(
		"Hello, @%v!\n\n<b>Student info:</b>\nFirst name: %v\nLast name: %v\n",
		account.Username,
		account.FirstName,
		account.LastName,
	)

	for _, student := range students {
		reportMessage += fmt.Sprintf("\nSchool: <b>%v</b>:\nType: %v\nStatus: %v\n\n", student.School.Title, student.GetType(), student.GetStatusText())

		if !student.Active {
			reportMessage += "Your student account is blocked!\nPlease contact mentors or teachers"
			continue
		}

		hlpr.logger.WithFields(logrus.Fields{
			"student_id": student.ID,
		}).Info("get student homeworks from database by student_id")
		studentHomeworks, err := hlpr.store.Homework().FindByStudentID(student.ID)
		if err != nil {
			if err == store.ErrRecordNotFound {
				reportMessage += fmt.Sprintf("you haven't submitted your homework yet\n\n%v", sysHomeworkAdd)
				continue
			}
			return "", err
		}
		hlpr.logger.WithFields(logrus.Fields{
			"count": len(studentHomeworks),
		}).Info("student homeworks found")

		hlpr.logger.WithFields(logrus.Fields{
			"school_id": student.School.ID,
		}).Info("get lessons from database by school_id")
		allLessons, err := hlpr.store.Lesson().FindBySchoolID(student.School.ID)
		if err != nil {
			return "", err
		}
		hlpr.logger.WithFields(logrus.Fields{
			"count": len(allLessons),
		}).Info("lessons found")

		reportMessage += "<b>Progress:</b>\n"
		if student.FullCourse {
			for _, lesson := range allLessons {
				counted := false
				for _, homework := range studentHomeworks {
					if homework.Lesson.ID == lesson.ID {
						counted = true
						reportMessage += fmt.Sprintf(
							"%v - %v [%v]\n",
							iconGreenCircle,
							lesson.Title,
							fmt.Sprintf(
								"<a href='%v'>Go To Message</a>",
								homework.GetURL(),
							),
						)
					}
				}

				if !counted {
					reportMessage += fmt.Sprintf("%v - %v\n", iconRedCircle, lesson.Title)
				}
			}
		} else {
			reportMessage += "You have joined the following modules:\n"
			studentModules := []*model.Module{}

			for _, homework := range studentHomeworks {
				studentModules = appendModule(studentModules, homework.Lesson.Module)
			}

			for _, module := range studentModules {
				reportMessage += fmt.Sprintf("- <b>%v</b>\n", module.Title)
			}

			reportMessage += "\n"

			for _, module := range studentModules {
				for _, lesson := range allLessons {
					if module.ID == lesson.Module.ID {
						counted := false
						for _, homework := range studentHomeworks {
							if homework.Lesson.ID == lesson.ID {
								counted = true
								reportMessage += fmt.Sprintf(
									"%v - %v [%v]\n",
									iconGreenCircle,
									lesson.Title,
									fmt.Sprintf(
										"<a href='%v'>Go To Message</a>",
										homework.GetURL(),
									),
								)
							}
						}

						if !counted {
							reportMessage += fmt.Sprintf("%v - %v\n", iconRedCircle, lesson.Title)
						}
					}
				}
			}
		}
	}

	return reportMessage, nil
}

// GetReport returns academic perfomance for all active stundents in school
func (hlpr *Helper) GetReport(school *model.School) (string, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Info("get lessons from database by school_id")
	lessons, err := hlpr.store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(lessons),
	}).Info("lessons found")

	hlpr.logger.WithFields(logrus.Fields{
		"school_id":   school.ID,
		"full_course": "true",
	}).Info("get students from database by school_id")
	students, err := hlpr.store.Student().FindByFullCourseSchoolID(true, school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Info("students found")

	reportMessage := "Academic perfomance:\n\nStudents Report:\n"

	if len(students) > 0 {
		message, err := hlpr.prepareReportMsg(students, lessons)
		if err != nil {
			return "", err
		}

		reportMessage += message
	} else {
		reportMessage += ErrReportNotFound + "\n"
	}

	hlpr.logger.WithFields(logrus.Fields{
		"school_id":   school.ID,
		"full_course": "false",
	}).Info("get listeners from database by school_id")
	listeners, err := hlpr.store.Student().FindByFullCourseSchoolID(false, school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(listeners),
	}).Info("listeners found")

	if len(listeners) > 0 {
		reportMessage += "\nListeners Report:\n"

		message, err := hlpr.prepareReportMsg(listeners, lessons)
		if err != nil {
			return "", err
		}

		reportMessage += message
	}

	return reportMessage, nil
}

// GetFullReport returns academic performance for all active students in school.
// Additionally adds the list of homeworks at the beginning
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

// GetLessonsReport returns the list of homeworks for school
// The list is populated only with active homeworks provided by students
func (hlpr *Helper) GetLessonsReport(school *model.School) (string, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Info("get lessons from database by school_id")
	lessons, err := hlpr.store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(lessons),
	}).Info("lessons found")

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

	type AccountReport struct {
		Account     *model.Account
		Accepted    int
		NotProvided int
	}

	reports := []*AccountReport{}

	reportMessage := "<b><u>Accepted/Not Provided - Name</u></b>\n"
	for _, student := range students {
		if !student.Active {
			continue
		}

		accountReport := &AccountReport{
			Account:     student.Account,
			Accepted:    0,
			NotProvided: 0,
		}

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
						accountReport.Accepted++
					}
				}

				if !counted {
					accountReport.NotProvided++
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
								accountReport.Accepted++
							}
						}

						if !counted {
							accountReport.NotProvided++
						}
					}
				}
			}
		}

		reports = append(reports, accountReport)
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Accepted > reports[j].Accepted
	})

	for _, report := range reports {
		reportMessage += fmt.Sprintf("%d/%d - %v\n",
			report.Accepted, report.NotProvided, report.Account.GetMention())

	}

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
