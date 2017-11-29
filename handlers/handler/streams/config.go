package streams

import (
	"github.com/corpix/effects/writer"

	"github.com/cryptounicorns/platypus/handlers/handler/stream"
)

type Config struct {
	Format string
	Inputs []stream.Config
	Wrap   string
	Writer writer.ConcurrentMultiWriterConfig
}
