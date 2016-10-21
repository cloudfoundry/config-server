package store

import "strconv"

type MemoryStore struct {
	db map[string]Configuration
}

var dbCounter int

func NewMemoryStore() MemoryStore {
	dbCounter = 0
	return MemoryStore{db: make(map[string]Configuration)}
}

func (store MemoryStore) Put(name string, value string) error {
	val, ok := store.db[name]

	if ok == false {
		store.db[name] = Configuration{
			Name:  name,
			Value: value,
			Id:    strconv.Itoa(dbCounter),
		}
		dbCounter++
	} else {
		val.Value = value
		store.db[name] = val
	}

	return nil
}

func (store MemoryStore) GetByName(name string) (Configuration, error) {
	return store.db[name], nil
}

func (store MemoryStore) GetById(id string) (Configuration, error) {
	result := Configuration{}

	for _, config := range store.db {
		if config.Id == id {
			result = config
			break
		}
	}

	return result, nil
}

func (store MemoryStore) Delete(name string) (bool, error) {
	deleted := false
	result, _ := store.GetByName(name)

	// map contains name, delete
	if len(result.Value) > 0 {
		delete(store.db, name)
		deleted = true
	}

	return deleted, nil
}
