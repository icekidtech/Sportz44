package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// New returns a structured logger configured for the given environment level.
func New(level string) *logrus.Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	switch level {
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	case "error":
		l.SetLevel(logrus.ErrorLevel)
	default:
		l.SetLevel(logrus.InfoLevel)
	}
	l.SetFormatter(&logrus.JSONFormatter{})
	return l
}
