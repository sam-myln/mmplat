package logging

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// Logger returns the standard logrus logger.
func Logger() (logger *logrus.Logger) {
	logger = &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    &logrus.TextFormatter{
			TimestampFormat: time.DateTime,
			DisableColors:   true,
			FullTimestamp:   true,

		},
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.DebugLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	return
}


