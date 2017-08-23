package memory

import (
	"github.com/corpix/loggers"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/cryptounicorns/market-fetcher-http/errors"
)

type Memory struct {
	storage cmap.ConcurrentMap
	log     loggers.Logger

	Config
}

func (s Memory) Set(key string, value interface{}) error {
	s.log.Debug("Set", key, value)
	s.storage.Set(key, value)

	return nil
}

func (s Memory) Get(key string) (interface{}, error) {
	v, _ := s.storage.Get(key)

	return v, nil
}

func (s Memory) Iter(fn func(key string, value interface{})) error {
	s.storage.IterCb(
		func(key string, value interface{}) {
			fn(
				key,
				value,
			)
		},
	)

	return nil
}

//

func New(c Config, l loggers.Logger) (*Memory, error) {
	if l == nil {
		return nil, errors.NewErrNilArgument(l)
	}

	return &Memory{
		storage: cmap.New(),
		log:     l,
		Config:  c,
	}, nil
}
