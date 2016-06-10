package store

import (
	. "config_server/config"
	"strings"
)

func CreateStore(dbConfig DBConfig) Store {

	var store Store

	if dbConfig == (DBConfig{}) {
		store = NewMemoryStore()

	} else if (strings.EqualFold(dbConfig.Adapter, "postgres")) {
		dbProvider := NewConcreteDbProvider(NewSqlWrapper(), dbConfig)
		store = NewPostgresStore(dbProvider)

	} else if (strings.EqualFold(dbConfig.Adapter, "mysql")) {
		dbProvider := NewConcreteDbProvider(NewSqlWrapper(), dbConfig)
		store = NewMysqlStore(dbProvider)

	} else {
		panic("Unsupported adapter")
	}

	return store
}
