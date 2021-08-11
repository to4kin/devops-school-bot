package apiserver

var (
	msgVersion   string = "dev"
	msgBuildDate string = ""
	msgBotInfo   string = "<b>Bot information:</b>\nVersion: %v\nBuild date: %v\nBuilt with: %v"

	msgHelpCommand    string = "I'll manage students homeworks\n\n<b>Commands</b>"
	msgUserPrivateCmd string = `
/start - will add you to the database for future use`
	msgSuperuserPrivateCmd string = `
<b>Superuser only</b>
/schools - will provide the interface to manage schools`
	msgUserGroupCmd string = `
/join - will add you to school
/myreport - will return your progress in school
/help - will return this help message`
	msgSuperuserGroupCmd string = `
<b>Superuser only</b>
/start - will create and start school
/finish - will finish current school`

	msgWelcomeToSchool string = "<b>Welcome to %v!</b>\n\nI'll manage all your progress and provide the report if needed.\n" +
		sysHomeworkAdd + "\n\n" + sysHomeworkGuide

	msgUserCreated                 string = "Hello, <b>%v!</b>\nAccount created successfully:\n\n" + msgUserInfo
	msgUserExist                   string = "Hello, <b>%v!</b>\nAccount already exist:\n\n" + msgUserInfo
	msgUserInfo                    string = "FirstName: %v\nLastName: %v\nUsername: %v\nSuperuser: %v"
	msgUserInsufficientPermissions string = `you have insufficient permissions, please contact teachers or mentors`
	msgUserNotJoined               string = `please join the school first`
	msgUserAlreadyJoined           string = "you have already joined school <b>%v</b>\n\n" + sysHomeworkAdd

	msgSchoolNotFound        string = `school not started, please contact teachers or mentors`
	msgSchoolStarted         string = `school <b>%v</b> started`
	msgSchoolAlreadyStarted  string = `school <b>%v</b> already started`
	msgSchoolAlreadyFinished string = `school <b>%v</b> already finished`
	msgSchoolFinished        string = `school <b>%v</b> finished`
	msgSchoolInfo            string = "<b>%v</b>\n\nCreated: %v\nStudents: %v\nHomeworks: %v\nStatus: %v"

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
