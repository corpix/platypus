package memory

import (
	"github.com/corpix/loggers"

	cmap "github.com/orcaman/concurrent-map"
)

type Memory struct {
	storage cmap.ConcurrentMap
	Config  Config
}

func (s *Memory) Set(key string, value interface{}) error {
	s.storage.Set(key, value)

	return nil
}

func (s *Memory) Get(key string) (interface{}, error) {
	v, _ := s.storage.Get(key)

	return v, nil
}

func (s *Memory) Remove(key string) error {
	s.storage.Remove(key)

	return nil
}

func (s *Memory) Iter(fn func(key string, value interface{})) error {
	s.storage.IterCb(fn)

	return nil
}

func (s *Memory) Close() error {
	return nil
}

func New(c Config, l loggers.Logger) (*Memory, error) {
	return &Memory{
		storage: cmap.New(),
		Config:  c,
	}, nil
}
