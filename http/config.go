package http

import (
	"github.com/cryptounicorns/platypus/http/endpoints"
)

type Config struct {
	Addr      string
	Endpoints endpoints.Config
}
