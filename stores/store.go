package stores

type Store interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Iter(func(key string, value interface{})) error
}
