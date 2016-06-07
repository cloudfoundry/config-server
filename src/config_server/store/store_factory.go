package store

import (
	. "config_server/config"
	"strings"
)

func CreateStore(dbConfig DBConfig) Store {

	var dataStore Store

	if dbConfig == (DBConfig{}) {
		dataStore = NewMemoryStore()
	} else if (strings.EqualFold(dbConfig.Adapter, "postgres")) {
		dataStore = NewPostgresStore(dbConfig)
	} else if (strings.EqualFold(dbConfig.Adapter, "mysql")) {
		dataStore = NewMysqlStore(dbConfig)
	} else {
		panic("Unsupported adapter")
	}

	return dataStore
}
