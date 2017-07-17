package datasources

import (
	"github.com/corpix/formats"
	"github.com/corpix/queues/message"
	"github.com/corpix/trade/market"

	"github.com/cryptounicorns/market-fetcher-http/feeds"
	"github.com/cryptounicorns/market-fetcher-http/logger"
	"github.com/cryptounicorns/market-fetcher-http/warehouse"
)

type Ticker struct {
	Feed   chan *market.Ticker
	Store  *warehouse.Warehouse
	format formats.Format
	log    logger.Logger
}

func (t *Ticker) Consume(m message.Message) {
	var (
		ticker = &market.Ticker{}
		err    error
	)

	err = t.format.Unmarshal(m, ticker)
	if err != nil {
		t.log.Error(err)
		return
	}

	t.Store.Set(
		[]string{
			ticker.Market,
			ticker.CurrencyPair.String(),
		},
		ticker,
	)
	select {
	case t.Feed <- ticker:
	default:
	}

}

func (t *Ticker) Close() error {
	close(t.Feed)
	return nil
}

func NewTicker(feed feeds.Feed, format formats.Format, log logger.Logger) (*Ticker, error) {
	ticker := &Ticker{
		Feed:   make(chan *market.Ticker),
		Store:  warehouse.New(),
		format: format,
		log:    log,
	}

	err := feed.Consume(ticker.Consume)
	if err != nil {
		return nil, err
	}

	return ticker, nil
}
