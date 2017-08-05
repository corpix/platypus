package stores

import (
	"strings"

	"github.com/corpix/logger"
	"github.com/fatih/structs"

	"github.com/cryptounicorns/market-fetcher-http/errors"
	"github.com/cryptounicorns/market-fetcher-http/stores/store/memory"
)

const (
	MemoryStoreType = "memory"
)

type Store interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Iter(func(key string, value interface{})) error
}

func New(c Config, l logger.Logger) (Store, error) {
	var (
		t = strings.ToLower(c.Type)
	)

	if l == nil {
		return nil, errors.NewErrNilArgument(l)
	}

	for _, v := range structs.New(c).Fields() {
		if strings.ToLower(v.Name()) != t {
			continue
		}

		switch t {
		case MemoryStoreType:
			return memory.New(
				v.Value().(memory.Config),
				l,
			)
		}
	}

	return nil, NewErrUnknownStoreType(c.Type)
}
