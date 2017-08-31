package stores

import (
	"strings"

	"github.com/corpix/loggers"
	"github.com/fatih/structs"

	"github.com/cryptounicorns/market-fetcher-http/errors"
	"github.com/cryptounicorns/market-fetcher-http/stores/store"
	"github.com/cryptounicorns/market-fetcher-http/stores/store/memory"
	"github.com/cryptounicorns/market-fetcher-http/stores/store/memoryttl"
)

const (
	MemoryStoreType    = "memory"
	MemoryTTLStoreType = "memoryttl"
)

func New(c Config, l loggers.Logger) (store.Store, error) {
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
		case MemoryTTLStoreType:
			return memoryttl.New(
				v.Value().(memoryttl.Config),
				l,
			)
		}
	}

	return nil, NewErrUnknownStoreType(c.Type)
}
