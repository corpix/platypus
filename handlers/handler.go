package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/fatih/structs"

	"github.com/cryptounicorns/platypus/handlers/handler/latest"
	"github.com/cryptounicorns/platypus/handlers/handler/latests"
	"github.com/cryptounicorns/platypus/handlers/handler/stream"
	"github.com/cryptounicorns/platypus/handlers/handler/streams"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	// Run could be called multiple times
	// So if we have error, we could log it
	// and run handler again, ServeHTTP should
	// serve data as normal.
	// (in worst case it should return 500)
	Run(context.Context) error
	Close() error
}

func New(c Config, l loggers.Logger) (Handler, error) {
	var (
		t   = strings.ToLower(c.Type)
		log = prefixwrapper.New(
			fmt.Sprintf(
				"Handler(%s): ",
				t,
			),
			l,
		)
		handler Handler
		err     error
	)

	for _, v := range structs.New(c).Fields() {
		if strings.ToLower(v.Name()) != t {
			continue
		}

		switch t {
		case latest.Name:
			handler, err = latest.New(
				v.Value().(latest.Config),
				log,
			)
			if err != nil {
				return nil, err
			}
		case latests.Name:
			handler, err = latests.New(
				v.Value().(latests.Config),
				log,
			)
			if err != nil {
				return nil, err
			}
		case stream.Name:
			handler, err = stream.New(
				v.Value().(stream.Config),
				log,
			)
			if err != nil {
				return nil, err
			}
		case streams.Name:
			handler, err = streams.New(
				v.Value().(streams.Config),
				log,
			)
			if err != nil {
				return nil, err
			}
		default:
			continue
		}

		return handler, nil
	}

	return nil, NewErrUnknownHandlerType(c.Type)
}
