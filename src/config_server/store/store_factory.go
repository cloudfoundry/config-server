package store

import (
	. "config_server/config"
)

func CreateStore(dbConfig DBConfig) Store {

	var dataStore Store

	if dbConfig == (DBConfig{}) {
		dataStore = NewMemoryStore()
	} else {
		dataStore = NewDatabaseStore(dbConfig)
	}

	return dataStore
}
