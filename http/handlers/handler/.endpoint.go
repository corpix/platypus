package handler

import (
	"github.com/corpix/loggers"
	"github.com/cryptounicorns/platypus/http/handlers/stores"
	"github.com/cryptounicorns/queues"
	"github.com/gorilla/mux"
)

type Handler struct {
	Router  *mux.Router
	Handler *Handler
	Queue   queues.Queue
}

func (e *Handler) Close() error {
	var (
		err error
	)

	err = e.Queue.Close()
	if err != nil {
		return err
	}

	err = e.Handler.Close()
	if err != nil {
		return err
	}

	return nil
}

func New(c Config, r *mux.Router, l loggers.Logger) (*Handler, error) {
	var (
		log = prefixwrapper.New(
			fmt.Sprintf(
				"Handler(%s, %s): ",
				c.Method,
				c.Path,
			),
			l,
		)

		s   stores.Store
		q   queues.Queue
		h   *Handler
		err error
	)

	q, err = queues.New(c.Queue, l)
	if err != nil {
		q.Close()
		return nil, err
	}

	s, err = stores.New(
		c.Store,
		log,
	)
	if err != nil {
		return nil, err
	}

	h, err = NewHandler(c, q, l)
	if err != nil {
		q.Close()
		h.Close()
		return nil, err
	}

	r.
		Path(c.Path).
		Methods(c.Method).
		HandlerFunc(h.Handle)

	return &Handler{
		Router:  r,
		Handler: h,
		Queue:   q,
	}, nil
}
