package logger

import (
	"github.com/corpix/logger"
	"github.com/corpix/logger/target/logrus"
)

type Config logrus.Config
type Logger logger.Logger

func New(c Config) (Logger, error) {
	return logrus.NewFromConfig(logrus.Config(c))
}
