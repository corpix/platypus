package routers

import (
	"fmt"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/fatih/structs"

	"github.com/cryptounicorns/platypus/http/handlers/routers/errors"
	"github.com/cryptounicorns/platypus/http/handlers/routers/router/broadcast"
	"github.com/cryptounicorns/platypus/http/handlers/routers/writer"
)

const (
	BroadcastRouterType = "broadcast"
)

func New(c Config, w writer.Iterator, e errors.Handler, l loggers.Logger) (Router, error) {
	var (
		t   = strings.ToLower(c.Type)
		log = prefixwrapper.New(
			fmt.Sprintf(
				"Router(%s): ",
				t,
			),
			l,
		)
	)

	for _, v := range structs.New(c).Fields() {
		if strings.ToLower(v.Name()) != t {
			continue
		}

		switch t {
		case BroadcastRouterType:
			return broadcast.New(
				v.Value().(broadcast.Config),
				w,
				e,
				log,
			)
		}
	}

	return nil, errors.NewErrUnknownRouterType(c.Type)
}
