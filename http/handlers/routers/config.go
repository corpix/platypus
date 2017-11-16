package routers

import (
	"github.com/cryptounicorns/platypus/http/handlers/routers/router/broadcast"
)

type Config struct {
	Type      string
	Broadcast broadcast.Config
}
