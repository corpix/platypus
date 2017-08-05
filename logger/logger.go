package logger

import (
	"github.com/corpix/logger"
	"github.com/corpix/logger/target/logrus"
)

type Config logrus.Config

func New(c Config) (logger.Logger, error) {
	return logrus.NewFromConfig(logrus.Config(c))
}
