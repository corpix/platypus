package logger

import (
	"strings"

	"github.com/corpix/logger"
	logrusLogger "github.com/corpix/logger/logrus"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Level     string
	Formatter string
}

type Logger logger.Logger

func New(c Config) (Logger, error) {
	var (
		l   logrus.Level
		f   logrus.Formatter
		err error
	)
	l, err = logrus.ParseLevel(c.Level)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(c.Formatter) {
	case "text":
		f = &logrus.TextFormatter{}
	case "json":
		f = &logrus.JSONFormatter{}
	case "":
		f = &logrus.TextFormatter{}
	default:
		return nil, NewErrUnknownFormatter(c.Formatter)
	}

	log := logrus.New()
	log.Level = l
	log.Formatter = f

	return logrusLogger.New(log), nil
}
