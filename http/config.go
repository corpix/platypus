package http

import (
	"github.com/cryptounicorns/market-fetcher-http/http/endpoints"
)

type Config struct {
	Addr      string
	Endpoints endpoints.Config
}
