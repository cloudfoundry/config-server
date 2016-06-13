package store

import (
	. "config_server/config"
	"strings"
)

func CreateStore(config ServerConfig) Store {

	var store Store

	if strings.EqualFold(config.Store, "database") {
		dbConfig := config.Database

		if (strings.EqualFold(dbConfig.Adapter, "postgres")) {
			dbProvider := NewConcreteDbProvider(NewSqlWrapper(), dbConfig)
			store = NewPostgresStore(dbProvider)

		} else if (strings.EqualFold(dbConfig.Adapter, "mysql")) {
			dbProvider := NewConcreteDbProvider(NewSqlWrapper(), dbConfig)
			store = NewMysqlStore(dbProvider)

		} else {
			panic("Unsupported adapter")
		}

	} else {
		store = NewMemoryStore()
	}

	return store
}
