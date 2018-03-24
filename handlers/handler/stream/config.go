package stream

import (
	"github.com/corpix/effects/writer"
	"github.com/cryptounicorns/queues"
)

type Config struct {
	Consumer queues.GenericConfig
	Writer   writer.ConcurrentMultiWriterConfig
}
