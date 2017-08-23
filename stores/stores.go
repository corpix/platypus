package stores

import (
	"strings"

	"github.com/corpix/loggers"
	"github.com/fatih/structs"

	"github.com/cryptounicorns/market-fetcher-http/errors"
	"github.com/cryptounicorns/market-fetcher-http/stores/store/memory"
)

const (
	MemoryStoreType = "memory"
)

func New(c Config, l loggers.Logger) (Store, error) {
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
