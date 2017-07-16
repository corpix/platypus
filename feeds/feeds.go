package feeds

import (
	"github.com/corpix/logger"
	"github.com/corpix/queues"
	"github.com/corpix/queues/handler"
)

type Config struct {
	Format string
	Ticker queues.Config
}

type Feed interface {
	Consume(handler.Handler) error
	Close() error
}

type Feeds struct {
	Ticker Feed
}

func NewFromConfig(c Config, l logger.Logger) (*Feeds, error) {
	var (
		ticker Feed
		err    error
	)

	ticker, err = queues.NewFromConfig(
		c.Ticker,
		l,
	)
	if err != nil {
		return nil, err
	}

	return &Feeds{
		Ticker: ticker,
	}, nil
}
