package broadcast

import (
	"time"

	"github.com/corpix/pool"
)

type Config struct {
	WriteTimeout time.Duration
	Pool         pool.Config
}
