package datasources

import (
	"github.com/corpix/formats"
	"github.com/corpix/logger"

	"github.com/cryptounicorns/market-fetcher-http/feeds"
)

type Datasources struct {
	*Ticker
}

func (d *Datasources) Close() {
	d.Ticker.Close()
}

func New(f *feeds.Feeds, fmts formats.Format, log logger.Logger) (*Datasources, error) {
	var (
		t   *Ticker
		err error
	)

	t, err = NewTicker(
		f.Ticker,
		fmts,
		log,
	)
	if err != nil {
		return nil, err
	}

	return &Datasources{
		Ticker: t,
	}, nil
}
