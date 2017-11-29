package latests

import (
	"github.com/cryptounicorns/platypus/handlers/handler/latest"
)

type Config struct {
	Format string
	Inputs []latest.Config
	Wrap   string
}
