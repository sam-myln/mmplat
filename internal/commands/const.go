package commands

const (
	fmtCmdShort = "mmplat %s"
	fmtCmdLong  = `mmplat %s
Multimedia platform for file sharing and content streaming
`
	fmtCmdEx    = `mmplat --bind|--addr localhost:8080 --folders dir... [-r|--recursive] [--two-factor false] 
[--file-formats formats...] [--log info|debug|error] [--login-data files...]`

	cmdFlagNameFolders = "folders"
	cmdFlagNameFormats = "file-formats"
	cmdFlagNameLogLevel = "log"
	cmdFlagNameRecursive = "recursive"
	cmdFlagNameRecursiveShort = "r"
	cmdFlagNameAddrAddress = "addr"
	cmdFlagNameBindAddress = "bind"
	cmdFlagNameAllowGuest = "allow-guest"
	cmdFlagNameLoginData = "login-data"

	cmdFlagNamePredefinedVideo = "mp4,mkv,avi,ts"
	cmdFlagNamePredefinedAudio = "mp3,flc"
	cmdFlagNameLogLevelDefault = "error"
	cmdFlagNameBindAddressDefault = "localhost:8080"

	version = "1.0"
)
