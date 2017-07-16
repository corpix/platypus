package feeds

import (
	"github.com/corpix/logger"
	"github.com/corpix/queues"
	"github.com/corpix/queues/handler"
)

const (
	NsqFeedType   = queues.NsqQueueType
	KafkaFeedType = queues.KafkaQueueType
)

type Config queues.Config
type Feed interface {
	Consume(handler.Handler) error
	Close() error
}

func NewFromConfig(l logger.Logger, c Config) (Feed, error) {
	return queues.NewFromConfig(
		l,
		queues.Config(c),
	)
}
