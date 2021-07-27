package apiserver

var (
	MsgHelpPrivate string = `I'll manage students homeworks`
	MsgHelpGroup   string = `I'll manage students homeworks

<b>Commands</b>
/join - will add you to the active school
/homeworks - will return your progress in active school
/help - will return this help message`

	MsgWelcomeToSchool string = `<b>Welcome to DevOps School!</b>

I'll manage all your progress and provide the report if needed.
To add homework, use the work hashtag along with the #homework, for example: <code>#homework #cicd</code>`

	MsgNoActiveSchool       string = `no active school found`
	MsgUserNotJoined        string = `please join the school first`
	MsgUserAlreadyJoined    string = `you have already joined the school`
	MsgHomeworkAlreadyAdded string = `homework has already been added`
)
