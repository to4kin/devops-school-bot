package apiserver

var (
	msgHelpPrivate string = `I'll manage students homeworks`
	msgHelpGroup   string = `I'll manage students homeworks

<b>Commands</b>
/start - will create a user
/join - will add you to the active school
/report - will return your progress in active school
/help - will return this help message
`

	msgWelcomeToSchool string = `<b>Welcome to DevOps School!</b>

I'll manage all your progress and provide the report if needed.
To add homework, use the work hashtag along with the #homework, for example: <code>#homework #cicd</code>
`

	msgNoActiveSchool    string = `no active school found`
	msgUserNotJoined     string = `please join the school first`
	msgUserAlreadyJoined string = `you have already joined the school`

	msgHomeworkVerified    string = `🟢`
	msgHomeworkNotVerified string = `🟡`
	msgHomeworkNotProvided string = `🔴`
	msgHomeworkReport      string = `Hello, @%v!

Guide:
` + msgHomeworkVerified + ` - homework is <b>verified</b>
` + msgHomeworkNotVerified + ` - homework is <b>NOT verified</b>
` + msgHomeworkNotProvided + ` - homework is <b>NOT provided</b>

Your progress in <b>DevOps School %v</b>:
`
)
