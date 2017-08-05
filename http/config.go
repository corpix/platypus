package http

import (
	api "github.com/cryptounicorns/market-fetcher-http/http/api/config"
)

type Config struct {
	Addr string
	Api  api.Config
}
