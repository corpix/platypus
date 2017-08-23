package consumer

import (
	"github.com/corpix/formats"
	"github.com/corpix/loggers"
	"github.com/corpix/queues"
	"github.com/corpix/queues/consumer"

	"github.com/cryptounicorns/market-fetcher-http/errors"
)

type Consumer struct {
	consumer consumer.Consumer
	*consumer.UnmarshalConsumer
}

func (c *Consumer) Close() error {
	var (
		err error
	)

	err = c.consumer.Close()
	if err != nil {
		return err
	}

	return c.UnmarshalConsumer.Close()
}

func New(q queues.Queue, t interface{}, f formats.Format, l loggers.Logger) (*Consumer, error) {
	if q == nil {
		return nil, errors.NewErrNilArgument(q)
	}
	if t == nil {
		return nil, errors.NewErrNilArgument(t)
	}
	if f == nil {
		return nil, errors.NewErrNilArgument(f)
	}
	if l == nil {
		return nil, errors.NewErrNilArgument(l)
	}

	var (
		c   consumer.Consumer
		uc  *consumer.UnmarshalConsumer
		err error
	)

	c, err = q.Consumer()
	if err != nil {
		return nil, err
	}

	uc, err = consumer.NewUnmarshalConsumer(
		t,
		c,
		f,
		l,
	)
	if err != nil {
		c.Close()
		return nil, err
	}

	return &Consumer{
		consumer:          c,
		UnmarshalConsumer: uc,
	}, nil
}
