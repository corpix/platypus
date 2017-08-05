package http

import (
	"github.com/gobwas/ws/wsutil"

	"github.com/cryptounicorns/market-fetcher-http/transmitters/writers"
)

func NewTransmitterWriter() writers.Writer {
	return wsutil.WriteServerBinary
}
