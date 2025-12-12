package command

const (
	fmtCmdShort = "mmplat %s"
	fmtCmdUse = "mmplat [flags] [PATH...]"
	fmtCmdLong  = `mmplat %s
Multimedia platform for file sharing and content streaming
`


	cmdFlagNameLogLevel = "log"
	cmdFlagNameFormat   = "format"
	cmdFlagNamePort     = "port"
	cmdFlagNameRecurse  = "recursive"
	cmdFlagNameTest     = "test"

	cmdFlagValuePredefinedVideo = "mp4,mkv,avi,ts"
	cmdFlagValuePredefinedAudio = "mp3,flc,flac,aac"
	cmdFlagValueBindPortDefault = ":8080"

	version = "1.1"
)
