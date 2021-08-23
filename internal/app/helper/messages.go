package helper

var (
	// SysHomeworkAdd ...
	SysHomeworkAdd string = `To add homework, use the lesson hashtag along with the #homework, for example: <code>#homework #cicd</code>`
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
