package consumer

import (
	"github.com/corpix/loggers"
	"github.com/cryptounicorns/queues"
	"github.com/cryptounicorns/queues/consumer"
)

type Stream struct {
	Consumer
	QueueConsumer consumer.Consumer
	Meta          *Meta
}

func (c *Stream) Close() error {
	var (
		err error
	)

	err = c.Consumer.Close()
	if err != nil {
		return err
	}

	err = c.QueueConsumer.Close()
	if err != nil {
		return err
	}

	return c.Meta.Close()
}

func NewStream(c Config, l loggers.Logger) (*Stream, error) {
	var (
		cr = &Stream{
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

	cr.Consumer, err = NewFormat(cr.QueueConsumer, c)
	if err != nil {
		return nil, err
	}

	cr.Meta.Stream, err = cr.Consumer.Consume()
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func NewStreams(c []Config, l loggers.Logger) ([]*Stream, error) {
	var (
		consumer  *Stream
		consumers = make(
			[]*Stream,
			len(c),
		)
		err error
	)

	for k, v := range c {
		consumer, err = NewStream(v, l)
		if err != nil {
			return nil, err
		}

		consumers[k] = consumer
	}

	return consumers, nil
}
