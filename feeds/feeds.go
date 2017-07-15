package feeds

import (
	"github.com/corpix/logger"
	"github.com/corpix/queues"
)

const (
	NsqFeedType   = queues.NsqQueueType
	KafkaFeedType = queues.KafkaQueueType
)

type Config queues.Config
type Feed queues.Queue

func NewFromConfig(l logger.Logger, c Config) (Feed, error) {
	return queues.NewFromConfig(
		l,
		queues.Config(c),
	)
}
