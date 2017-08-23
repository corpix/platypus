package channel

import (
	"github.com/corpix/loggers"

	"github.com/corpix/queues/consumer"
	"github.com/corpix/queues/errors"
	"github.com/corpix/queues/message"
	"github.com/corpix/queues/producer"
)

type Channel struct {
	config  Config
	log     loggers.Logger
	channel chan message.Message
}

func (q *Channel) Producer() (producer.Producer, error) {
	return NewProducer(q.channel)
}

func (q *Channel) Consumer() (consumer.Consumer, error) {
	return NewConsumer(q.channel)
}

func (q *Channel) Close() error {
	close(q.channel)

	return nil
}

func NewFromConfig(c Config, l loggers.Logger) (*Channel, error) {
	if l == nil {
		return nil, errors.NewErrNilArgument(l)
	}

	return &Channel{
		config: c,
		log:    l,
		channel: make(
			chan message.Message,
			c.Capacity,
		),
	}, nil
}
