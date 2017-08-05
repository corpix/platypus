package config

import (
	"github.com/cryptounicorns/market-fetcher-http/consumer"
	"github.com/cryptounicorns/market-fetcher-http/stores"
	"github.com/cryptounicorns/market-fetcher-http/transmitters"
)

type Config struct {
	Consumer    consumer.Config
	Store       stores.Config
	Transmitter transmitters.Config
}
