package memoryttl

import (
	"github.com/cryptounicorns/platypus/time"
)

type Config struct {
	TTL        time.Duration
	Resolution time.Duration
}
