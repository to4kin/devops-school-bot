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
	// SysHomeworkAdd ...
	SysHomeworkAdd string = "To add homework, use the lesson hashtag along with the #homework, for example: <code>#homework #cicd</code>"
	// SysHomeworkGuide ...
	SysHomeworkGuide string = `Homework guide:
` + iconGreenCircle + ` - homework is <b>accepted</b>
` + iconRedCircle + ` - homework is <b>NOT provided</b>`
)

var (
	iconGreenCircle string = `🟢`
	//iconYellowCircle string = `🟡`
	iconRedCircle string = `🔴`
)
