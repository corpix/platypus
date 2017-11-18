package cache

import (
	"github.com/cryptounicorns/platypus/stores"
)

type Config struct {
	Key   string
	Store stores.Config
}
