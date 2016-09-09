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

func (store MemoryStore) Delete(key string) (bool, error) {
	deleted := false
	value, _ := store.Get(key)

    // map contains key, delete
    if len(value) > 0 {
		delete(store.db, key)
		deleted = true
	}

	return deleted, nil
}
