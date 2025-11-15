package commands

const (
	fmtCmdShort = "mmplat %s"
	fmtCmdUse = "mmplat [flags] [PATH...]"
	fmtCmdLong  = `mmplat %s
Multimedia platform for file sharing and content streaming
`
	fmtCmdEx = `
mmplat --test
mmplat --test -r dir1 dir2
mmplat --port 8080 -r --file-formats txt,json
mmplat --log-level debug media
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
