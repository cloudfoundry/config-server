package store

import (
	"github.com/BurntSushi/migration"
)

type sqlWrapper struct {
}

func NewSqlWrapper() sqlWrapper {
	return sqlWrapper{}
}

func (w sqlWrapper) Open(driverName, dataSourceName string, migrations []migration.Migrator) (IDb, error) {
	db, err := migration.Open(driverName, dataSourceName, migrations)
	return NewDbWrapper(db), err
}
