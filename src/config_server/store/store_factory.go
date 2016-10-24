package store

import (
	"config_server/config"
	"strings"

	"github.com/cloudfoundry/bosh-utils/errors"
)

func CreateStore(config config.ServerConfig) (Store, error) {
	var store Store

	if strings.EqualFold(config.Store, "database") {
		dbConfig := config.Database

		if strings.EqualFold(dbConfig.Adapter, "postgres") {
			dbProvider := NewConcreteDbProvider(NewSQLWrapper(), dbConfig)
			store = NewPostgresStore(dbProvider)

		} else if strings.EqualFold(dbConfig.Adapter, "mysql") {
			dbProvider := NewConcreteDbProvider(NewSQLWrapper(), dbConfig)
			store = NewMysqlStore(dbProvider)

		} else {
			return nil, errors.Errorf("Unsupported adapter '%s'", dbConfig.Adapter)
		}
	} else {
		store = NewMemoryStore()
	}

	return store, nil
}
