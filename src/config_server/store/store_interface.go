package store

type Store interface {
	Put(key string, value string) error
	GetByKey(key string) (Configuration, error)
	GetById(id string) (Configuration, error)
	Delete(key string) (bool, error)
}
