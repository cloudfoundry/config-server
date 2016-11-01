package store

import (
	"database/sql"
)

type RowsWrapper struct {
	rows *sql.Rows
}

func NewRowsWrapper(rows *sql.Rows) RowsWrapper {
	return RowsWrapper{rows}
}

func (w RowsWrapper) Next() bool {
	return w.rows.Next()
}
