package stores

import (
	"github.com/cryptounicorns/platypus/stores/store/memory"
	"github.com/cryptounicorns/platypus/stores/store/memoryttl"
)

type Config struct {
	Type      string
	Memory    memory.Config
	MemoryTTL memoryttl.Config
}
