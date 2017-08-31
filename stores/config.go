package stores

import (
	"github.com/cryptounicorns/market-fetcher-http/stores/store/memory"
	"github.com/cryptounicorns/market-fetcher-http/stores/store/memoryttl"
)

type Config struct {
	Type      string
	Memory    memory.Config
	MemoryTTL memoryttl.Config
}
