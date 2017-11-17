package cache

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/corpix/loggers"

	"github.com/cryptounicorns/platypus/http/handlers/stores"
)

type Cache struct {
	config Config
	stores.Store
	keyTemplate *template.Template
}

func (c *Cache) Set(value interface{}) (string, error) {
	var (
		buf = bytes.NewBuffer([]byte{})
		key string
		err error
	)

	err = c.keyTemplate.Execute(buf, value)
	if err != nil {
		return key, err
	}
	key = string(
		buf.Bytes(),
	)

	return key, c.Store.Set(key, value)
}

func New(c Config, l loggers.Logger) (*Cache, error) {
	var (
		s   stores.Store
		t   *template.Template
		err error
	)

	s, err = stores.New(c.Store, l)
	if err != nil {
		return nil, err
	}

	t, err = template.New("key").Parse(
		strings.TrimSpace(c.Key),
	)
	if err != nil {
		return nil, err
	}

	return &Cache{
		config:      c,
		Store:       s,
		keyTemplate: t,
	}, nil
}