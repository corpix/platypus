package memoryttl

import (
	"time"
)

type Config struct {
	TTL        time.Duration
	Resolution time.Duration
}
