package streams

import (
	"github.com/cryptounicorns/platypus/http/handlers/consumer"
	"github.com/cryptounicorns/platypus/http/handlers/routers"
)

type Config struct {
	Format string

	Consumers []consumer.Config
	Wrap      string
	Router    routers.Config
}
