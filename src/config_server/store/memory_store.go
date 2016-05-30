package store

type MemoryStore struct {
	db map[string]string
}

func NewMemoryStore() MemoryStore {
	return MemoryStore{make(map[string]string)}
}

func (store MemoryStore) Put(key string, value string) error {
	store.db[key] = value
	return nil
}

func (store MemoryStore) Get(key string) (string, error) {
	return store.db[key], nil
}
