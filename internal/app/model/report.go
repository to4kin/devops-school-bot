package model

import "fmt"

// Report object represents a report
type Report struct {
	Student     *Student
	Accepted    []*Homework
	NotProvided []*Lesson
}

// GetCSVHeader returns string header for CSV file
func (r *Report) GetCSVHeader() string {
	return "full_name;school;type;module;lesson;provided\n"
}

// GetCSVLine converts object to CSV line
func (r *Report) GetCSVLine() string {
	line := ""

	for _, accepted := range r.Accepted {
		line += fmt.Sprintf("%s;%s;%s;%s;%s;%s",
			r.Student.Account.GetFullName(),
			r.Student.School.Title,
			r.Student.GetType(),
			accepted.Lesson.Module.Title,
			accepted.Lesson.Title,
			"true\n",
		)
	}

	for _, notProvided := range r.NotProvided {
		line += fmt.Sprintf("%s;%s;%s;%s;%s;%s",
			r.Student.Account.GetFullName(),
			r.Student.School.Title,
			r.Student.GetType(),
			notProvided.Module.Title,
			notProvided.Title,
			"false\n",
		)
	}

	return line
}
