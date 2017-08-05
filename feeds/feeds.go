package feeds

import (
	"github.com/corpix/logger"
	"github.com/corpix/queues"
)

type Feeds struct {
	Tickers queues.Queue
}

func (fs *Feeds) Close() error {
	var (
		err error
	)

	err = fs.Tickers.Close()
	if err != nil {
		return err
	}

	return nil
}

func New(c Config, l logger.Logger) (*Feeds, error) {
	var (
		tickers queues.Queue
		err     error
	)

	tickers, err = queues.NewFromConfig(
		c.Tickers,
		l,
	)
	if err != nil {
		return nil, err
	}

	return &Feeds{
		Tickers: tickers,
	}, nil
}
