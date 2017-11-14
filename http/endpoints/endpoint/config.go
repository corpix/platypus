package endpoint

import (
	"github.com/cryptounicorns/queues"

	"github.com/cryptounicorns/platypus/consumer"
	"github.com/cryptounicorns/platypus/stores"
	"github.com/cryptounicorns/platypus/transmitters"
)

type Config struct {
	Path   string
	Method string

	Queue       queues.Config
	Consumer    consumer.Config
	Store       stores.Config
	Transmitter transmitters.Config
}
