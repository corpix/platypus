package http

import (
	"net/http"

	"github.com/corpix/loggers"
	"github.com/gorilla/mux"

	"github.com/cryptounicorns/market-fetcher-http/errors"
	"github.com/cryptounicorns/market-fetcher-http/feeds"
	"github.com/cryptounicorns/market-fetcher-http/http/api"
)

type Server struct {
	Config
	Feeds feeds.Feeds
	log   loggers.Logger
}

func (s *Server) Serve() error {
	s.log.Printf(
		"Starting server on '%s'...",
		s.Config.Addr,
	)
	r := mux.NewRouter()

	_, err := api.New(
		s.Config.Api,
		r,
		s.Feeds,
		s.log,
	)
	if err != nil {
		return err
	}

	return http.ListenAndServe(
		s.Config.Addr,
		r,
	)
}

func New(c Config, f feeds.Feeds, l loggers.Logger) (*Server, error) {
	if f == nil {
		return nil, errors.NewErrNilArgument(f)
	}
	if l == nil {
		return nil, errors.NewErrNilArgument(l)
	}

	return &Server{
		Config: c,
		Feeds:  f,
		log:    l,
	}, nil
}
