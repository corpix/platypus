package transmitters

import (
	"github.com/cryptounicorns/platypus/transmitters/transmitter/broadcast"
)

type Config struct {
	Type      string
	Broadcast broadcast.Config
}
