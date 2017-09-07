package transmitters

import (
	"strings"

	"github.com/corpix/loggers"
	"github.com/fatih/structs"

	"github.com/cryptounicorns/platypus/errors"
	"github.com/cryptounicorns/platypus/transmitters/transmitter"
	"github.com/cryptounicorns/platypus/transmitters/transmitter/broadcast"
	"github.com/cryptounicorns/platypus/transmitters/writers"
)

const (
	BroadcastTransmitterType = "broadcast"
)

func New(c Config, ws writers.Writers, w writers.Writer, e transmitter.ErrorHandler, l loggers.Logger) (Transmitter, error) {
	var (
		t = strings.ToLower(c.Type)
	)

	if ws == nil {
		return nil, errors.NewErrNilArgument(ws)
	}
	if w == nil {
		return nil, errors.NewErrNilArgument(w)
	}
	if e == nil {
		return nil, errors.NewErrNilArgument(e)
	}
	if l == nil {
		return nil, errors.NewErrNilArgument(l)
	}

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
				l,
			)
		}
	}

	return nil, NewErrUnknownTransmitterType(c.Type)
}
