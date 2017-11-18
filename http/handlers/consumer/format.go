package consumer

import (
	"github.com/corpix/formats"
	"github.com/cryptounicorns/queues/consumer"
	"github.com/cryptounicorns/queues/result"
)

type Format struct {
	consumer consumer.Consumer
	format   formats.Format
	done     chan struct{}
}

func (c *Format) Consume() (<-chan Result, error) {
	var (
		stream         = make(chan Result)
		consumerStream <-chan result.Result
		err            error
	)

	consumerStream, err = c.consumer.Consume()
	if err != nil {
		return nil, err
	}

	go func() {
		var (
			r Result
		)

		for cr := range consumerStream {
			select {
			case <-c.done:
				return
			default:
				if cr.Err == nil {
					r.Err = c.format.Unmarshal(
						cr.Value,
						&r.Value,
					)
				} else {
					r.Err = cr.Err
				}

				stream <- r
			}
		}
	}()

	return stream, nil
}

func (c *Format) Close() error {
	close(c.done)
	return nil
}

func NewFormat(cr consumer.Consumer, c Config) (*Format, error) {
	var (
		f   formats.Format
		err error
	)

	f, err = formats.New(c.Format)
	if err != nil {
		return nil, err
	}

	return &Format{
		consumer: cr,
		format:   f,
		done:     make(chan struct{}),
	}, nil
}
