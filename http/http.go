package http

import (
	"net/http"

	"github.com/corpix/loggers"
	"github.com/gorilla/mux"

	"github.com/cryptounicorns/platypus/http/handlers"
)

type Server struct {
	Config
	handlers handlers.Handlers
	router   *mux.Router
	log      loggers.Logger
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
	return s.handlers.Close()
}

func New(c Config, l loggers.Logger) (*Server, error) {
	var (
		r   = mux.NewRouter()
		hs  = make(handlers.Handlers, len(c.Handlers))
		h   handlers.Handler
		err error
	)

	for k, v := range c.Handlers {
		h, err = handlers.New(
			v,
			r,
			l,
		)
		if err != nil {
			return nil, err
		}

		hs[k] = h
	}

	return &Server{
		Config:   c,
		handlers: hs,
		router:   r,
		log:      l,
	}, nil
}
