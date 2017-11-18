package broadcast

import (
	"github.com/corpix/pool"

	"github.com/cryptounicorns/platypus/time"
)

type Config struct {
	WriteTimeout time.Duration
	Pool         pool.Config
}
