package store

import (
	"config_server/config"
)

func CreateStore(config config.DBConfig) Store {

	var dataStore Store

	if &config == nil {
		dataStore = NewMemoryStore()
	} else {
		dataStore = NewDatabaseStore(config)
	}
	
	return dataStore
}
