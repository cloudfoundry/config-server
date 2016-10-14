package store

type Store interface {
	Put(key string, value string) error
	Get(key string) (Configuration, error)
	Delete(key string) (bool, error)
}