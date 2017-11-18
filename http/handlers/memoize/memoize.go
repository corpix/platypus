package memoize

import (
	"github.com/cryptounicorns/platypus/http/handlers/cache"
	"github.com/cryptounicorns/platypus/http/handlers/consumer"
)

type Memoize struct {
	Cache    *cache.Cache
	Consumer *consumer.Consumer
}

func (m Memoize) Close() error {
	var (
		err error
	)

	err = m.Consumer.Close()
	if err != nil {
		return err
	}

	err = m.Cache.Close()
	if err != nil {
		return err
	}

	return nil
}
