package http

import (
	"github.com/gobwas/ws/wsutil"

	"github.com/cryptounicorns/platypus/transmitters/writers"
)

func NewTransmitterWriter() writers.Writer {
	return wsutil.WriteServerBinary
}
