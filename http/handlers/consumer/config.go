package consumer

import (
	"github.com/cryptounicorns/queues"
)

type Config struct {
	Format string

	queues.Config
}
