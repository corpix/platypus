package consumer

import (
	"github.com/corpix/formats"
	"github.com/cryptounicorns/queues"
)

type Meta struct {
	queues.Queue
	Config Config
	Format formats.Format
	stream <-chan Result
}
