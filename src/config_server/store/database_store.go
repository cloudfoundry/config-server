package store

import "errors"

type DatabaseStore struct {
}

func NewDatabaseStore() DatabaseStore {
	return DatabaseStore{}
}

func (store DatabaseStore) Put(key string, value string) error {
	return errors.New("Not implemented")
}
func (store DatabaseStore) Get(key string) (string, error) {
	return "", errors.New("Not implemented")
}
