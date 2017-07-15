package kafka

import (
	"github.com/corpix/queues/queue/kafka"
)

type Config kafka.Config

func WrapConfig(c Config) kafka.Config { return kafka.Config(c) }
