package store

import (
	"io"
)

type Store interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Remove(key string) error
	Iter(func(key string, value interface{})) error
	io.Closer
}
