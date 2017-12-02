package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"

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
	)

	switch t {
	case latest.Name:
		return latest.New(
			c.Latest,
			log,
		)
	case latests.Name:
		return latests.New(
			c.Latests,
			log,
		)
	case stream.Name:
		return stream.New(
			c.Stream,
			log,
		)
	case streams.Name:
		return streams.New(
			c.Streams,
			log,
		)
	default:
		return nil, NewErrUnknownHandlerType(c.Type)
	}
}
