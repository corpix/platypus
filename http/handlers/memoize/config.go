package memoize

import (
	"github.com/cryptounicorns/platypus/http/handlers/cache"
	"github.com/cryptounicorns/platypus/http/handlers/consumer"
)

type Config struct {
	Consumer consumer.Config
	Cache    cache.Config
}
