package transmitters

import (
	"fmt"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/fatih/structs"

	"github.com/cryptounicorns/platypus/transmitters/transmitter"
	"github.com/cryptounicorns/platypus/transmitters/transmitter/broadcast"
	"github.com/cryptounicorns/platypus/transmitters/writers"
)

const (
	BroadcastTransmitterType = "broadcast"
)

func New(c Config, ws writers.Writers, w writers.Writer, e transmitter.ErrorHandler, l loggers.Logger) (Transmitter, error) {
	var (
		t   = strings.ToLower(c.Type)
		log = prefixwrapper.New(
			fmt.Sprintf(
				"Transmitter(%s): ",
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
		case BroadcastTransmitterType:
			return broadcast.New(
				v.Value().(broadcast.Config),
				ws,
				w,
				e,
				log,
			)
		}
	}

	return nil, NewErrUnknownTransmitterType(c.Type)
}
