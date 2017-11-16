package http

import (
	"github.com/cryptounicorns/platypus/http/handlers"
)

type Config struct {
	Addr     string
	Handlers handlers.Configs
}
