package stores

import (
	"github.com/cryptounicorns/platypus/http/handlers/stores/store/memory"
	"github.com/cryptounicorns/platypus/http/handlers/stores/store/memoryttl"
)

type Config struct {
	Type      string
	Memory    memory.Config
	MemoryTTL memoryttl.Config
}
