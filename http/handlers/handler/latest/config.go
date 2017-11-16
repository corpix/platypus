package latest

import (
	"github.com/cryptounicorns/platypus/http/handlers/cache"
	"github.com/cryptounicorns/platypus/http/handlers/consumer"
)

type Config struct {
	Format string

	Consumer consumer.Config
	Cache    cache.Config
}
