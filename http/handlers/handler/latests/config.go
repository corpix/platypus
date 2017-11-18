package latests

import (
	"github.com/cryptounicorns/platypus/http/handlers/memoize"
)

type Config struct {
	Format  string
	Memoize []memoize.Config
	Wrap    string
}
