package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"os"
	"time"
)

type Logger interface {
	Infof(format string, args ...interface{})
}

type FasthttpAdapter struct {
	logger fasthttp.Logger
}

func (f *FasthttpAdapter) Infof(format string, args ...interface{}) {
	f.logger.Printf(format, args...)
}

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

