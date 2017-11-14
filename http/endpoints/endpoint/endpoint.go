package endpoint

import (
	"github.com/corpix/loggers"
	"github.com/cryptounicorns/queues"
	"github.com/gorilla/mux"
)

type Endpoint struct {
	Router  *mux.Router
	Handler *Handler
	Queue   queues.Queue
}

func (e *Endpoint) Close() error {
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

func New(c Config, r *mux.Router, l loggers.Logger) (*Endpoint, error) {
	var (
		q   queues.Queue
		h   *Handler
		err error
	)

	q, err = queues.New(
		c.Queue,
		l,
	)
	if err != nil {
		return nil, err
	}

	h, err = NewHandler(
		c,
		q,
		l,
	)
	if err != nil {
		return nil, err
	}

	r.
		Path(c.Path).
		Methods(c.Method).
		HandlerFunc(h.Handle)

	return &Endpoint{
		Router:  r,
		Handler: h,
		Queue:   q,
	}, nil
}
