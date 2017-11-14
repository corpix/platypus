package http

import (
	"net/http"

	"github.com/corpix/loggers"
	"github.com/gorilla/mux"

	"github.com/cryptounicorns/platypus/http/endpoints"
)

type Server struct {
	Config
	endpoints endpoints.Endpoints
	router    *mux.Router
	log       loggers.Logger
}

func (s *Server) Serve() error {
	s.log.Printf(
		"Starting server on '%s'...",
		s.Config.Addr,
	)

	return http.ListenAndServe(
		s.Config.Addr,
		s.router,
	)
}

func (s *Server) Close() error {
	return s.endpoints.Close()
}

func New(c Config, l loggers.Logger) (*Server, error) {
	var (
		r   = mux.NewRouter()
		es  endpoints.Endpoints
		err error
	)

	es, err = endpoints.New(
		c.Endpoints,
		r,
		l,
	)
	if err != nil {
		return nil, err
	}

	return &Server{
		Config:    c,
		endpoints: es,
		router:    r,
		log:       l,
	}, nil
}
