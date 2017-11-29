package handlers

import (
	"github.com/cryptounicorns/platypus/handlers/handler/latest"
	"github.com/cryptounicorns/platypus/handlers/handler/latests"
	"github.com/cryptounicorns/platypus/handlers/handler/stream"
	"github.com/cryptounicorns/platypus/handlers/handler/streams"
)

type Configs = []Config

type Config struct {
	Path   string `validate:"required"`
	Method string `validate:"required"`
	Type   string `validate:"required"`

	Latest  latest.Config
	Latests latests.Config

	Stream  stream.Config
	Streams streams.Config
}
