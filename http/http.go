package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cryptounicorns/market-fetcher-http/datasources"
	"github.com/cryptounicorns/market-fetcher-http/logger"
)

type Config struct {
	Addr string
}

type Server struct {
	config      Config
	datasources datasources.Datasources
	logger      logger.Logger
}

func (s *Server) Serve() error {
	r := mux.NewRouter()
	// FIXME: Handlers
	s.logger.Printf("Listening on '%s'", s.config.Addr)
	return http.ListenAndServe(s.config.Addr, r)
}

func New(c Config, d datasources.Datasources, l logger.Logger) *Server {
	return &Server{
		config:      c,
		datasources: d,
		logger:      l,
	}
}
