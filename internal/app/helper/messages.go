package helper

var (
	// ErrInternal ...
	ErrInternal string = "Internal server error!\n\nsomething has gone wrong on the server"
	// ErrInsufficientPermissions ...
	ErrInsufficientPermissions = "You have insufficient permissions, please contact teachers or mentors"
	// ErrWrongChatType ...
	ErrWrongChatType = "Wrong chat type, please use another chat for this command"
	// ErrSchoolNotStarted ...
	ErrSchoolNotStarted = "School not started, please contact teachers or mentors"
	// ErrReportNotFound ...
	ErrReportNotFound = "Insufficient data for the report"
	// ErrUserNotJoined ...
	ErrUserNotJoined = "Please join school first"

	// MsgWelcomeToSchool ...
	MsgWelcomeToSchool string = "<b>Welcome to %v as %v!</b>\n\nI'll manage all your progress and provide the report if needed.\n" +
		sysHomeworkAdd
	// MsgUserAlreadyJoined ...
	MsgUserAlreadyJoined string = "you have already joined school <b>%v</b> as %v\n\n" + sysHomeworkAdd

	// SysListenerGuide ...
	SysListenerGuide string = "Homework guide:\nPlease provide homeworks only for your module!"
	// SysStudentGuide ...
	SysStudentGuide string = `Homework guide:
` + iconGreenCircle + ` - homework is <b>accepted</b>
` + iconRedCircle + ` - homework is <b>NOT provided</b>`
)

var (
	msgStudentInfo string = "Student info:\n\nFirst name: %v\nLast name: %v\nType: %v\nStatus: %v"

	iconGreenCircle string = `🟢`
	//iconYellowCircle string = `🟡`
	iconRedCircle string = `🔴`

	sysHomeworkAdd string = "To add homework, use the lesson hashtag along with the #homework, for example: <code>#homework #cicd</code>"
)
