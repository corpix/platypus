package feeds

import (
	"github.com/corpix/loggers"
	"github.com/corpix/queues"
)

type Feeds map[string]queues.Queue

func (fs Feeds) Close() error {
	var (
		err error
	)

	for _, v := range fs {
		err = v.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func New(c Config, l loggers.Logger) (Feeds, error) {
	var (
		qs  = Feeds{}
		q   queues.Queue
		err error
	)

	for k, v := range c {
		q, err = queues.NewFromConfig(v, l)
		if err != nil {
			return nil, err
		}

		qs[k] = q
	}

	return qs, nil
}
