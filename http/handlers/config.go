package handlers

import (
	"github.com/cryptounicorns/platypus/http/handlers/handler/latest"
	"github.com/cryptounicorns/platypus/http/handlers/handler/stream"
	"github.com/cryptounicorns/platypus/http/handlers/handler/streams"
)

type Configs = []Config

type Config struct {
	Path   string
	Method string
	Type   string

	Stream  stream.Config
	Streams streams.Config
	Latest  latest.Config
}
