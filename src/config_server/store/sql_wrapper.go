package store

import "database/sql"

type sqlWrapper struct {
}

func NewSqlWrapper() sqlWrapper {
	return sqlWrapper{}
}

func (w sqlWrapper)Open(driverName, dataSourceName string) (IDb, error) {
	db, err := sql.Open(driverName, dataSourceName)
	return NewDbWrapper(db), err
}