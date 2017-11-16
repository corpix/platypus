package handlers

import (
	"github.com/cryptounicorns/platypus/http/handlers/handler/latest"
	"github.com/cryptounicorns/platypus/http/handlers/handler/stream"
)

type Configs = []Config

type Config struct {
	Path   string
	Method string
	Type   string
	Format string

	Stream stream.Config
	Latest latest.Config
}
