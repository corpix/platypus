package transmitters

import (
	"strings"

	"github.com/corpix/logger"
	"github.com/fatih/structs"

	"github.com/cryptounicorns/market-fetcher-http/errors"
	"github.com/cryptounicorns/market-fetcher-http/transmitters/transmitter"
	"github.com/cryptounicorns/market-fetcher-http/transmitters/transmitter/broadcast"
	"github.com/cryptounicorns/market-fetcher-http/transmitters/writers"
)

const (
	BroadcastTransmitterType = "broadcast"
)

func New(c Config, ws writers.Writers, w writers.Writer, e transmitter.ErrorHandler, l logger.Logger) (Transmitter, error) {
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
