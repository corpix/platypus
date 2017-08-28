package endpoint

import (
	"github.com/corpix/queues"

	"github.com/cryptounicorns/market-fetcher-http/consumer"
	"github.com/cryptounicorns/market-fetcher-http/stores"
	"github.com/cryptounicorns/market-fetcher-http/transmitters"
)

type Config struct {
	Path   string
	Method string

	Queue       queues.Config
	Consumer    consumer.Config
	Store       stores.Config
	Transmitter transmitters.Config
}
