package queues

import (
	"strings"

	"github.com/corpix/logger"
	"github.com/fatih/structs"

	"github.com/corpix/queues/handler"
	"github.com/corpix/queues/message"
	"github.com/corpix/queues/queue/kafka"
	"github.com/corpix/queues/queue/nsq"
)

//

const (
	KafkaQueueType = "kafka"
	NsqQueueType   = "nsq"
)

//

type Config struct {
	Type  string
	Kafka kafka.Config
	Nsq   nsq.Config
}

//

type Queue interface {
	Produce(message.Message) error
	Consume(handler.Handler) error
	Close() error
}

//

func NewFromConfig(c Config, l logger.Logger) (Queue, error) {
	var (
		t = strings.ToLower(c.Type)
	)

	for _, v := range structs.New(c).Fields() {
		if strings.ToLower(v.Name()) != t {
			continue
		}

		switch t {
		case KafkaQueueType:
			return kafka.NewFromConfig(
				v.Value().(kafka.Config),
				l,
			)
		case NsqQueueType:
			return nsq.NewFromConfig(
				v.Value().(nsq.Config),
				l,
			)
		}
	}

	return nil, NewErrUnknownQueueType(c.Type)
}