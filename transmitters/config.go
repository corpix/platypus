package transmitters

import (
	"github.com/cryptounicorns/market-fetcher-http/transmitters/transmitter/broadcast"
)

type Config struct {
	Type      string
	Broadcast broadcast.Config
}
