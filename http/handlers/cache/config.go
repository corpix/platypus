package cache

import (
	"github.com/cryptounicorns/platypus/http/handlers/stores"
)

type Config struct {
	Key   string
	Store stores.Config
}
