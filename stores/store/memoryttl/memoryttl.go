package memoryttl

import (
	"time"

	"github.com/corpix/loggers"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/cryptounicorns/platypus/errors"
)

type MemoryTTL struct {
	storage   cmap.ConcurrentMap
	timeouted cmap.ConcurrentMap
	log       loggers.Logger
	done      chan bool
	Config    Config
}

func (s *MemoryTTL) Set(key string, value interface{}) error {
	s.log.Debug("Set ", key, value)

	s.storage.Set(
		key,
		value,
	)
	s.timeouted.Set(
		key,
		time.Now().Add(s.Config.TTL),
	)

	return nil
}

func (s *MemoryTTL) Get(key string) (interface{}, error) {
	v, _ := s.storage.Get(key)

	return v, nil
}

func (s *MemoryTTL) Remove(key string) error {
	s.log.Debug("Remove ", key)
	s.storage.Remove(key)
	s.timeouted.Remove(key)

	return nil
}

func (s *MemoryTTL) Iter(fn func(key string, value interface{})) error {
	s.storage.IterCb(fn)

	return nil
}

func (s *MemoryTTL) Close() error {
	close(s.done)

	return nil
}

func (s *MemoryTTL) cancellationLoop() {
	for {
		select {
		case <-s.done:
			return
		case <-time.After(s.Config.Resolution):
			for k, v := range s.timeouted.Items() {
				if time.Now().After(v.(time.Time)) {
					s.Remove(k)
				}
			}
		}
	}
}

func New(c Config, l loggers.Logger) (*MemoryTTL, error) {
	if l == nil {
		return nil, errors.NewErrNilArgument(l)
	}

	var (
		s = &MemoryTTL{
			storage:   cmap.New(),
			timeouted: cmap.New(),
			log:       l,
			done:      make(chan bool),
			Config:    c,
		}
	)

	go s.cancellationLoop()

	return s, nil
}