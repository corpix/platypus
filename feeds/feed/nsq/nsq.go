package nsq

import (
	"github.com/corpix/queues/queue/nsq"
)

type Config nsq.Config

func WrapConfig(c Config) nsq.Config { return nsq.Config(c) }
