package endpoints

import (
	"github.com/corpix/loggers"
	"github.com/gorilla/mux"

	"github.com/cryptounicorns/platypus/errors"
	"github.com/cryptounicorns/platypus/http/endpoints/endpoint"
)

type Endpoints []*endpoint.Endpoint

func (es Endpoints) Close() error {
	var (
		err error
	)

	for _, e := range es {
		err = e.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func New(c Config, r *mux.Router, l loggers.Logger) (Endpoints, error) {
	if c == nil {
		return nil, errors.NewErrNilArgument(c)
	}
	if r == nil {
		return nil, errors.NewErrNilArgument(r)
	}
	if l == nil {
		return nil, errors.NewErrNilArgument(l)
	}

	var (
		es  = make(Endpoints, len(c))
		err error
	)

	for k, ec := range c {
		es[k], err = endpoint.New(
			ec,
			r,
			l,
		)
		if err != nil {
			return nil, err
		}
	}

	return es, nil
}
