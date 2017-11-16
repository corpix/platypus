package handlers

import (
	"fmt"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/fatih/structs"
	"github.com/gorilla/mux"

	"github.com/cryptounicorns/platypus/http/handlers/handler/latest"
	"github.com/cryptounicorns/platypus/http/handlers/handler/stream"
)

const (
	LatestType = "latest"
	StreamType = "stream"
)

type Handlers []Handler

func (es Handlers) Close() error {
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

func New(c Config, r *mux.Router, l loggers.Logger) (Handler, error) {
	var (
		t   = strings.ToLower(c.Type)
		log = prefixwrapper.New(
			fmt.Sprintf(
				"Handler(%s): ",
				t,
			),
			l,
		)
		h   Handler
		err error
	)

	for _, v := range structs.New(c).Fields() {
		if strings.ToLower(v.Name()) != t {
			continue
		}

		switch t {
		case LatestType:
			h, err = latest.New(
				v.Value().(latest.Config),
				log,
			)
			if err != nil {
				return nil, err
			}
		case StreamType:
			h, err = stream.New(
				v.Value().(stream.Config),
				log,
			)
			if err != nil {
				return nil, err
			}
		default:
			continue
		}

		r.
			Path(c.Path).
			Methods(c.Method).
			Handler(h)

		return h, nil
	}

	return nil, NewErrUnknownHandlerType(c.Type)
}
