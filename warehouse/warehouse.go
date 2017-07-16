package warehouse

import (
	"strings"

	cmap "github.com/orcaman/concurrent-map"
)

const (
	keyDelimiter = "/"
)

//

type Warehouse struct {
	data cmap.ConcurrentMap
}

func (w Warehouse) Set(key []string, value interface{}) {
	w.data.Set(
		strings.Join(key, keyDelimiter),
		value,
	)
}

func (w Warehouse) Get(key []string) (interface{}, bool) {
	return w.data.Get(
		strings.Join(key, keyDelimiter),
	)
}

func (w Warehouse) FindPrefix(key []string) map[string]interface{} {
	var (
		res       = map[string]interface{}{}
		joinedKey = strings.Join(key, keyDelimiter)
	)

	w.data.IterCb(
		func(k string, v interface{}) {
			if strings.HasPrefix(k, joinedKey) {
				res[k] = v
			}
		},
	)

	return res
}

//

func New() *Warehouse {
	return &Warehouse{cmap.New()}
}
