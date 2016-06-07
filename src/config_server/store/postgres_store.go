package store

import (
	"config_server/config"
	"database/sql"
)

type PostgresStore struct {
	config config.DBConfig
}

func NewPostgresStore(config config.DBConfig) PostgresStore {
	return PostgresStore{config}
}

func (store PostgresStore) Put(key string, value string) error {

	db, err := DBConnection(store.config)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, _ := db.Query("SELECT value FROM config WHERE key = $1", key)
	if rows.Next() == false {
		_, err = db.Exec("INSERT INTO config VALUES($1, $2)", key, value)
	} else {
		_, err = db.Exec("UPDATE config SET value=$1 WHERE key=$2", value, key)
	}

	return err
}

func (store PostgresStore) Get(key string) (string, error) {

	var value string

	db, err := DBConnection(store.config)
	if err != nil {
		return value, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT value FROM config WHERE key = $1", key).Scan(&value)
	if err == sql.ErrNoRows {
		return value, nil
	}

	return value, err
}
