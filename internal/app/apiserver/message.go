package apiserver

var (
	msgVersion   string = "dev"
	msgBuildDate string = ""
	msgBotInfo   string = "<b>Bot information:</b>\nVersion: %v\nBuild date: %v\nBuilt with: %v"

	msgInternalError string = "Internal Server Error!\n\nSomething has gone wrong on the server"

	msgHelpCommand string = `I'll manage students homeworks

<b>Abbreviation</b>
[GC] - Commands available in <b>Group chats</b>
[PC] - Commands available in <b>Private chat</b> with bot
[All] - Commands available in <b>All chats</b>
[Admin] - Commands available only for bot <b>Administrators</b>

<b>Commands</b>
/start - <b>[GC][Admin]</b>: Start school. <b>[PC]</b>: Add user to database
/finish - <b>[GC][Admin]</b>: Finish school
/report - <b>[GC][Admin]</b>: School progress
/fullreport - <b>[GC][Admin]</b>: School progress with homework list
/join - <b>[GC]</b>: Join school as student
/myreport - <b>[GC]</b>: Your progress
/homeworks - <b>[GC]</b>: Homework list
/schools - <b>[PC][Admin]</b>: Manage schools
/users - <b>[PC][Admin]</b>: Manage users
/help - <b>[All]</b>: Help message
`

	msgWelcomeToSchool string = "<b>Welcome to %v!</b>\n\nI'll manage all your progress and provide the report if needed.\n" +
		sysHomeworkAdd + "\n\n" + sysHomeworkGuide

	msgUserCreated                 string = "Hello, <b>%v!</b>\nAccount created successfully!\n\n" + msgUserInfo
	msgUserExist                   string = "Hello, <b>%v!</b>\nAccount already exist!\n\n" + msgUserInfo
	msgUserInfo                    string = "Account info:\n\nFirst name: %v\nLast name: %v\nUsername: @%v\nSuperuser: %v"
	msgUserInsufficientPermissions string = `you have insufficient permissions, please contact teachers or mentors`
	msgUserNotJoined               string = `please join school first`
	msgUserAlreadyJoined           string = "you have already joined school <b>%v</b>\n\n" + sysHomeworkAdd

	msgSchoolNotFound        string = `school not started, please contact teachers or mentors`
	msgSchoolStarted         string = `school <b>%v</b> started`
	msgSchoolAlreadyStarted  string = `school <b>%v</b> already started`
	msgSchoolAlreadyFinished string = `school <b>%v</b> already finished`
	msgSchoolFinished        string = `school <b>%v</b> finished`
	msgSchoolInfo            string = "<b>%v</b>\n\nCreated: %v\nStudents: %v\nHomeworks: %v\nStatus: %v\n\nAccepted homeworks:\n%v"

	//msgStudentIsBlocked string = "Your student account was blocked!\n\nPlease contact mentors or teachers"
	msgStudentInfo string = "School: %v\n\nStudent info:\n\nFirst name: %v\nLast name: %v\nStatus: %v\n\n" + sysHomeworkGuide + "\n\nHomeworks:\n%v"

	msgHomeworkNotProvided string = "you haven't submitted your homework yet\n\n" + sysHomeworkAdd
	msgHomeworkReport      string = "Hello, @%v!\n\n" + sysHomeworkGuide + "\n\nYour progress in <b>%v</b>:\n"
	msgHomeworkInfo        string = "School: %v\n\nHomework info:\n\nTitle: %v"
	msgHomeworkList        string = "<b>Homework list</b>\n\n"

	msgReport string = "Academic performance\n\n<b><u>Name - Accepted/Not Provided - Type</u></b>\n"

	iconGreenCircle string = `🟢`
	//iconYellowCircle string = `🟡`
	iconRedCircle string = `🔴`

	sysHomeworkAdd   string = `To add homework, use the lesson hashtag along with the #homework, for example: <code>#homework #cicd</code>`
	sysHomeworkGuide string = `Homework guide:
` + iconGreenCircle + ` - homework is <b>accepted</b>
` + iconRedCircle + ` - homework is <b>NOT provided</b>`
)
