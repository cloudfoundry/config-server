package store

import (
	"database/sql"
)

type rowsWrapper struct {
	rows *sql.Rows
}

func NewRowsWrapper(rows *sql.Rows) rowsWrapper {
	return rowsWrapper{rows}
}

func (w rowsWrapper) Next() bool {
	return w.rows.Next()
}
