package logging

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// NewLogger returns the standard logrus logger.
func NewLogger() (logger *logrus.Logger) {
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

