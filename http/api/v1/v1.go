package v1

import (
	"github.com/corpix/logger"
	"github.com/gorilla/mux"

	"github.com/cryptounicorns/market-fetcher-http/errors"
	"github.com/cryptounicorns/market-fetcher-http/feeds"
	"github.com/cryptounicorns/market-fetcher-http/http/api/config"
)

type V1 struct {
	Router  *mux.Router
	Tickers *Tickers
}

func (a *V1) Close() error {
	var (
		err error
	)

	err = a.Tickers.Close()
	if err != nil {
		return err
	}

	return nil
}

func New(c config.Config, r *mux.Router, f *feeds.Feeds, l logger.Logger) (*V1, error) {
	if r == nil {
		return nil, errors.NewErrNilArgument(r)
	}
	if f == nil {
		return nil, errors.NewErrNilArgument(f)
	}
	if l == nil {
		return nil, errors.NewErrNilArgument(l)
	}

	var (
		v1      = r.PathPrefix("/v1").Subrouter()
		tickers *Tickers
		err     error
	)

	tickers, err = NewTickers(
		c,
		v1,
		f.Tickers,
		l,
	)
	if err != nil {
		return nil, err
	}

	return &V1{
		Router:  v1,
		Tickers: tickers,
	}, nil
}
