package stores

import (
	"github.com/cryptounicorns/market-fetcher-http/stores/store/memory"
)

type Config struct {
	Type   string
	Memory memory.Config
}
