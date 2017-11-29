package stores

import (
	"fmt"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/fatih/structs"

	"github.com/corpix/stores/store/memory"
	"github.com/corpix/stores/store/memoryttl"
)

func New(c Config, l loggers.Logger) (Store, error) {
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
		case memory.Name:
			return memory.New(
				v.Value().(memory.Config),
				log,
			)
		case memoryttl.Name:
			return memoryttl.New(
				v.Value().(memoryttl.Config),
				log,
			)
		}
	}

	return nil, NewErrUnknownStoreType(c.Type)
}
