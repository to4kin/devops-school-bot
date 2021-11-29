package helper

var (
	// ErrInternal returns string for internal error
	ErrInternal string = "Internal server error!\n\nsomething has gone wrong on the server"
	// ErrInsufficientPermissions returns string for permissions error
	ErrInsufficientPermissions = "You have insufficient permissions, please contact teachers or mentors"
	// ErrWrongChatType returns string for wrong chat typer error. Please specify needed chat
	// type via fmt.
	ErrWrongChatType = "Wrong chat type, please use %v chat for this command"
	// ErrSchoolNotStarted returns string for not started school eeror
	ErrSchoolNotStarted = "School not started, please contact teachers or mentors"
	// ErrReportNotFound returns string for report not found error
	ErrReportNotFound = "Insufficient data for the report"
	// ErrUserNotJoined returns string for not joined user error
	ErrUserNotJoined = "Please join school first"

	// MsgWelcomeToSchool returns welcome string message
	MsgWelcomeToSchool string = "<b>Welcome to %v as %v!</b>\n\nI'll manage all your progress and provide the report if needed.\n" +
		sysHomeworkAdd
	// MsgUserAlreadyJoined returns already joined string message
	MsgUserAlreadyJoined string = "you have already joined school <b>%v</b> as %v\n\n" + sysHomeworkAdd

	// SysListenerGuide returns guide for listeners
	SysListenerGuide string = "Homework guide:\nPlease provide homeworks only for your module!"
	// SysStudentGuide returns guide for students
	SysStudentGuide string = `Homework guide:
` + iconGreenCircle + ` - homework is <b>accepted</b>
` + iconRedCircle + ` - homework is <b>NOT provided</b>`
)

var (
	iconGreenCircle string = `🟢`
	//iconYellowCircle string = `🟡`
	iconRedCircle string = `🔴`

	sysHomeworkAdd string = "To add homework, use the lesson hashtag along with the #homework, for example: <code>#homework #cicd</code>"
)
