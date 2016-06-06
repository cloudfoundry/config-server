package store

import (
	"config_server/config"
	"database/sql"
)

type DatabaseStore struct {
	config config.DBConfig
}

func NewDatabaseStore(config config.DBConfig) DatabaseStore {
	return DatabaseStore{config}
}

func (store DatabaseStore) Put(key string, value string) error {

	db, err := DBConnection(store.config)
	if err != nil {
		return err
	}
	defer db.Close()

	query := SQLReplacer(store.config.Adapter, "INSERT INTO config VALUES(?,?)")
	_, err = db.Exec(query, key, value)

	return err
}

func (store DatabaseStore) Get(key string) (string, error) {

	var value string

	db, err := DBConnection(store.config)
	if err != nil {
		return value, err
	}
	defer db.Close()

	query := SQLReplacer(store.config.Adapter, "SELECT value FROM config c WHERE c.key = ?")
	err = db.QueryRow(query, key).Scan(&value)

	if err == sql.ErrNoRows {
		return value, nil
	}

	return value, err
}
