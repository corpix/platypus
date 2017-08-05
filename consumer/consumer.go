package consumer

import (
	"github.com/corpix/formats"
	"github.com/corpix/logger"
	"github.com/corpix/queues"
	"github.com/corpix/queues/consumer"

	"github.com/cryptounicorns/market-fetcher-http/errors"
)

func New(q queues.Queue, t interface{}, f formats.Format, l logger.Logger) (*consumer.Consumer, error) {
	var (
		c   *consumer.Consumer
		err error
	)

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

	c, err = consumer.New(t, f, l)
	if err != nil {
		return nil, err
	}

	err = q.Consume(c.Handler)
	if err != nil {
		return nil, err
	}

	return c, nil
}
