package stream

import (
	"github.com/cryptounicorns/platypus/http/handlers/consumer"
	"github.com/cryptounicorns/platypus/http/handlers/routers"
	//"github.com/cryptounicorns/queues"
)

type Config struct {
	Format string

	Consumer consumer.Config
	Router   routers.Config
}
