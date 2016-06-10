package store

import "database/sql"

type rowWrapper struct {
	row *sql.Row
}

func NewRowWrapper(row *sql.Row) rowWrapper {
	return rowWrapper {row}
}

func (w rowWrapper) Scan(dest ...interface{}) error {
	return w.row.Scan(dest...)
}