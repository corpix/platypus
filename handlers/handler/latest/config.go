package latest

import (
	"github.com/corpix/stores"
	"github.com/cryptounicorns/queues"
)

type Config struct {
	Format   string
	Key      string
	Store    stores.Config
	Consumer queues.GenericConfig
}
