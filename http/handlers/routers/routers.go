package routers

import (
	"fmt"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/fatih/structs"

	"github.com/cryptounicorns/platypus/http/handlers/routers/router"
	"github.com/cryptounicorns/platypus/http/handlers/routers/router/broadcast"
)

const (
	BroadcastRouterType = "broadcast"
)

func New(c Config, ws router.Writers, w router.Writer, e router.ErrorHandler, l loggers.Logger) (Router, error) {
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
				ws,
				w,
				e,
				log,
			)
		}
	}

	return nil, NewErrUnknownRouterType(c.Type)
}
