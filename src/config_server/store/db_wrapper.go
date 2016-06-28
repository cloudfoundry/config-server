package store

import (
	"database/sql"
)

type dbWrapper struct {
	db *sql.DB
}

func NewDbWrapper(db *sql.DB) dbWrapper {
	return dbWrapper{db}
}

func (w dbWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	return w.db.Exec(query, args...)
}

func (w dbWrapper) Query(query string, args ...interface{}) (IRows, error) {
	rows, err := w.db.Query(query, args...)
	return NewRowsWrapper(rows), err
}

func (w dbWrapper) QueryRow(query string, args ...interface{}) IRow {
	return NewRowWrapper(w.db.QueryRow(query, args...))
}

func (w dbWrapper) Close() {
	w.db.Close()
}

func (w dbWrapper) SetMaxOpenConns(n int) {
	w.db.SetMaxOpenConns(n)
}

func (w dbWrapper) SetMaxIdleConns(n int) {
	w.db.SetMaxIdleConns(n)
}
