package stores

import (
	"fmt"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/fatih/structs"

	"github.com/cryptounicorns/platypus/stores/store"
	"github.com/cryptounicorns/platypus/stores/store/memory"
	"github.com/cryptounicorns/platypus/stores/store/memoryttl"
)

const (
	MemoryStoreType    = "memory"
	MemoryTTLStoreType = "memoryttl"
)

func New(c Config, l loggers.Logger) (store.Store, error) {
	var (
		t   = strings.ToLower(c.Type)
		log = prefixwrapper.New(
			fmt.Sprintf(
				"Store(%s): ",
				t,
			),
			l,
		)
	)

	for _, v := range structs.New(c).Fields() {
		if strings.ToLower(v.Name()) != t {
			continue
		}

		switch t {
		case MemoryStoreType:
			return memory.New(
				v.Value().(memory.Config),
				log,
			)
		case MemoryTTLStoreType:
			return memoryttl.New(
				v.Value().(memoryttl.Config),
				log,
			)
		}
	}

	return nil, NewErrUnknownStoreType(c.Type)
}
