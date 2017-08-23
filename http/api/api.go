package api

import (
	"github.com/corpix/loggers"
	"github.com/gorilla/mux"

	"github.com/cryptounicorns/market-fetcher-http/errors"
	"github.com/cryptounicorns/market-fetcher-http/feeds"
	"github.com/cryptounicorns/market-fetcher-http/http/api/config"
	"github.com/cryptounicorns/market-fetcher-http/http/api/v1"
)

type Api struct {
	Router *mux.Router
	V1     *v1.V1
}

func (a *Api) Close() error {
	return a.V1.Close()
}

func New(c config.Config, r *mux.Router, f *feeds.Feeds, l loggers.Logger) (*Api, error) {
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
		api      = r.PathPrefix("/api").Subrouter()
		version1 *v1.V1
		err      error
	)

	version1, err = v1.New(
		c,
		api,
		f,
		l,
	)
	if err != nil {
		return nil, err
	}

	return &Api{
		Router: api,
		V1:     version1,
	}, nil
}
