package helper

import (
	"errors"
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/model"
	"gitlab.devops.telekom.de/tvpp/prototypes/devops-school-bot/internal/app/store"
)

// GetUserReport returns report for User
// Account should be a vald Account struct
// School can be nil, in this case report will be generated for all schools where User is a student
// NOTE: This report is the same as GetStudentReport but for all joined schools by this user
func (hlpr *Helper) GetUserReport(account *model.Account) (string, error) {

	hlpr.logger.WithFields(logrus.Fields{
		"account_id": account.ID,
	}).Info("get students from database by account_id")
	students, err := hlpr.store.Student().FindByAccountID(account.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			return ErrUserNotJoined, nil
		}

		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Info("student found")

	reportMessage := ""

	for _, student := range students {
		if !student.Active {
			reportMessage += "Your student account is blocked!\nPlease contact mentors or teachers"
			continue
		}

		message, err := hlpr.GetStudentReport(student)
		if err != nil {
			return "", err
		}

		reportMessage += message + "\n"
	}

	return reportMessage, nil
}

// GetStudentReport returns student progress in School:
// 	Account info:
// 	First name: Ivan
// 	Last name: Ivanov
//
// 	School: DevOps School Test:
// 	Type: Student
// 	Status: 🟢Active
//
// 	Progress:
// 	🟢 - #cicd1 [Go To Message (https://t.me/c/1534814897/279)]
// 	🟢 - #cicd2 [Go To Message (https://t.me/c/1534814897/293)]
func (hlpr *Helper) GetStudentReport(student *model.Student) (string, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"school_id": student.School.ID,
	}).Info("get lessons from database by school_id")
	lessons, err := hlpr.store.Lesson().FindBySchoolID(student.School.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(lessons),
	}).Info("lessons found")

	reportMessage := fmt.Sprintf("School: <b>%v</b>:\nType: %v\nStatus: %v\n\n", student.School.Title, student.GetType(), student.GetStatusText())

	if err == store.ErrRecordNotFound {
		reportMessage += ErrReportNotFound
	} else {
		message, err := hlpr.prepareDetailedReportMsg(student, lessons)
		if err != nil {
			return "", err
		}

		reportMessage += message
	}

	return reportMessage, nil
}

// GetReport returns academic perfomance for all active stundents in school:
// 	School DevOps School Test
//
// 	Academic perfomance:
//
// 	Students Report:
// 	Accepted/Not Provided - Name
// 	1/1 - Ivan Petrov
// 	1/1 - Sergey Ivanov
//
// 	Way to go, Ivan Petrov, you're crushing it!
//
// 	Listeners Report:
// 	Accepted/Not Provided - Name
// 	2/0 - Ivan Ivanov
//
// 	Way to go, Ivan Ivanov, you're crushing it!
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
	reportMessage += "<b><u>Accepted/Not Provided - Name</u></b>\n"

	if len(students) > 0 {
		message, err := hlpr.prepareGeneralReportMsg(students, lessons)
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
		reportMessage += "<b><u>Accepted/Not Provided - Name</u></b>\n"

		message, err := hlpr.prepareGeneralReportMsg(listeners, lessons)
		if err != nil {
			return "", err
		}

		reportMessage += message
	}

	return reportMessage, nil
}

// GetFullReport returns academic performance for all active students in school.
// Additionally adds the list of homeworks at the beginning
// 	School DevOps School Test
//
// 	Homework list
//
// 	Module: cicd
// 	#cicd1
// 	#cicd2
//
// 	Academic perfomance:
//
// 	Students Report:
// 	Accepted/Not Provided - Name
// 	1/1 - Ivan Petrov
// 	1/1 - Sergey Ivanov
//
// 	Way to go, Ivan Petrov, you're crushing it!
//
// 	Listeners Report:
// 	Accepted/Not Provided - Name
// 	2/0 - Ivan Ivanov
//
// 	Way to go, Ivan Ivanov, you're crushing it!
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
// The list is populated only with active homeworks provided by students:
// 	School DevOps School Test
//
// 	Homework list
//
// 	Module: cicd
// 	#cicd1
// 	#cicd2
func (hlpr *Helper) GetLessonsReport(school *model.School) (string, error) {
	hlpr.logger.WithFields(logrus.Fields{
		"school_id": school.ID,
	}).Info("get lessons from database by school_id")
	lessons, err := hlpr.store.Lesson().FindBySchoolID(school.ID)
	if err != nil {
		if err == store.ErrRecordNotFound {
			hlpr.logger.Info(err)
		}

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

// GetCSVReport returns csv string
func (hlpr *Helper) GetCSVReport(school *model.School) (string, error) {
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
	students, err := hlpr.store.Student().FindBySchoolID(school.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return "", err
	}
	hlpr.logger.WithFields(logrus.Fields{
		"count": len(students),
	}).Info("students found")

	reports := []*model.Report{}
	for _, student := range students {
		report, err := hlpr.prepareReport(student, lessons)
		if err != nil {
			return "", err
		}

		reports = append(reports, report)
	}

	if len(reports) == 0 {
		return "", errors.New("")
	}

	csvReport := reports[0].GetCSVHeader()
	for _, report := range reports {
		csvReport += report.GetCSVLine()
	}

	return csvReport, nil
}

func (hlpr *Helper) prepareGeneralReportMsg(students []*model.Student, lessons []*model.Lesson) (string, error) {
	reports := []*model.Report{}
	for _, student := range students {
		if !student.Active {
			continue
		}

		report, err := hlpr.prepareReport(student, lessons)
		if err != nil {
			return "", err
		}

		reports = append(reports, report)
	}

	sort.Slice(reports, func(i, j int) bool {
		return len(reports[i].Accepted) > len(reports[j].Accepted)
	})

	reportMessage := ""
	for _, report := range reports {
		reportMessage += fmt.Sprintf("%d/%d - %s\n",
			len(report.Accepted), len(report.NotProvided), report.Student.Account.GetMention())

	}

	reportMessage += fmt.Sprintf("\nWay to go, %s, you're crushing it!\n", reports[0].Student.Account.GetMention())

	return reportMessage, nil
}

func (hlpr *Helper) prepareDetailedReportMsg(student *model.Student, lessons []*model.Lesson) (string, error) {
	report, err := hlpr.prepareReport(student, lessons)
	if err != nil {
		return "", err
	}

	reportMessage := "<b>Progress:</b>\n"
	for _, accepted := range report.Accepted {
		reportMessage += fmt.Sprintf(
			"%v - %v [%v]\n",
			iconGreenCircle,
			accepted.Lesson.Title,
			fmt.Sprintf(
				"<a href='%v'>Go To Message</a>",
				accepted.GetURL(),
			))
	}

	for _, notProvided := range report.NotProvided {
		reportMessage += fmt.Sprintf("%v - %v\n", iconRedCircle, notProvided.Title)
	}

	return reportMessage, nil
}

func (hlpr *Helper) prepareReport(student *model.Student, lessons []*model.Lesson) (*model.Report, error) {
	studentReport := &model.Report{
		Student:     student,
		Accepted:    []*model.Homework{},
		NotProvided: []*model.Lesson{},
	}

	homeworks, err := hlpr.store.Homework().FindByStudentID(student.ID)
	if err != nil && err != store.ErrRecordNotFound {
		return nil, err
	}

	if student.FullCourse {
		for _, lesson := range lessons {
			counted := false
			for _, homework := range homeworks {
				if homework.Lesson.ID == lesson.ID {
					counted = true
					studentReport.Accepted = append(studentReport.Accepted, homework)
				}
			}

			if !counted {
				studentReport.NotProvided = append(studentReport.NotProvided, lesson)
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
							studentReport.Accepted = append(studentReport.Accepted, homework)
						}
					}

					if !counted {
						studentReport.NotProvided = append(studentReport.NotProvided, lesson)
					}
				}
			}
		}
	}

	return studentReport, nil
}

func appendModule(slice []*model.Module, homework *model.Module) []*model.Module {
	for _, elem := range slice {
		if elem.ID == homework.ID {
			return slice
		}
	}
	return append(slice, homework)
}
