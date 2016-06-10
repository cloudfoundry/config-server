package store

type ISql interface {
	Open(driverName, dataSourceName string) (IDb, error)
}
