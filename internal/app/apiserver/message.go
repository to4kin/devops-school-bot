package apiserver

var (
	msgHelpCommand  string = "I'll manage students homeworks\n\n<b>Commands</b>"
	msgUserGroupCmd string = `
/join - will add you to school
/report - will return your progress in school
/help - will return this help message`
	msgSuperuserGroupCmd string = `
<b>Superuser only</b>
/start - will create and start school
/finish - will finish current school`

	msgWelcomeToSchool string = "<b>Welcome to %v!</b>\n\nI'll manage all your progress and provide the report if needed.\n" +
		sysHomeworkAdd + "\n\n" + sysHomeworkGuide

	msgUserInsufficientPermissions string = `you have insufficient permissions, please contact teachers or mentors`
	msgUserNotJoined               string = `please join the school first`
	msgUserAlreadyJoined           string = "you have already joined school <b>%v</b>\n\n" + sysHomeworkAdd

	msgSchoolNotFound        string = `school not started, please contact teachers or mentors`
	msgSchoolStarted         string = `school <b>%v</b> started`
	msgSchoolExist           string = `school <b>%v</b> exist and started`
	msgSchoolAlreadyFinished string = `school <b>%v</b> already finished`
	msgSchoolFinished        string = `school <b>%v</b> finished`

	msgHomeworkNotProvided string = "you haven't submitted your homework yet\n\n" + sysHomeworkAdd
	msgHomeworkReport      string = "Hello, @%v!\n\n" + sysHomeworkGuide + "\n\nYour progress in <b>%v</b>:\n"

	iconHomeworkVerified    string = `🟢`
	iconHomeworkNotVerified string = `🟡`
	iconHomeworkNotProvided string = `🔴`

	sysHomeworkAdd   string = `To add homework, use the lesson hashtag along with the #homework, for example: <code>#homework #cicd</code>`
	sysHomeworkGuide string = `Homework guide:
` + iconHomeworkVerified + ` - homework is <b>verified</b>
` + iconHomeworkNotVerified + ` - homework is <b>NOT verified</b>
` + iconHomeworkNotProvided + ` - homework is <b>NOT provided</b>`
)
