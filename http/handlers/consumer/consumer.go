package consumer

import (
	"github.com/corpix/loggers"
	"github.com/cryptounicorns/queues"
	"github.com/cryptounicorns/queues/consumer"
)

type Consumer struct {
	*Format
	QueueConsumer consumer.Consumer
	Meta          *Meta
}

func (c *Consumer) Stream() <-chan Result {
	return c.Meta.stream
}

func (c *Consumer) Close() error {
	var (
		err error
	)

	err = c.Format.Close()
	if err != nil {
		return err
	}

	err = c.QueueConsumer.Close()
	if err != nil {
		return err
	}

	return c.Meta.Close()
}

func New(c Config, l loggers.Logger) (*Consumer, error) {
	var (
		cr = &Consumer{
			Meta: &Meta{
				Config: c,
			},
		}
		err error
	)

	cr.Meta.Queue, err = queues.New(c.Queue, l)
	if err != nil {
		return nil, err
	}

	cr.QueueConsumer, err = cr.Meta.Queue.Consumer()
	if err != nil {
		return nil, err
	}

	cr.Format, err = NewFormat(cr.QueueConsumer, c)
	if err != nil {
		return nil, err
	}

	cr.Meta.stream, err = cr.Format.Consume()
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func NewConsumers(c []Config, l loggers.Logger) ([]*Consumer, error) {
	var (
		consumer  *Consumer
		consumers = make(
			[]*Consumer,
			len(c),
		)
		err error
	)

	for k, v := range c {
		consumer, err = New(v, l)
		if err != nil {
			return nil, err
		}

		consumers[k] = consumer
	}

	return consumers, nil
}
