package store

import (
	_ "github.com/lib/pq"
	"config_server/config"
	"errors"
)

type DatabaseStore struct {
	config config.DBConfig
}

func NewDatabaseStore(config config.DBConfig) DatabaseStore {
	return DatabaseStore{config}
}

func (store DatabaseStore) Put(key string, value string) error {

	if &key == nil || &value == nil {
		return errors.New("Key/Value cannot be nil")
	}

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

	rows, err := db.Query(SQLReplacer(store.config.Adapter, "SELECT value FROM config c  WHERE c.key = ?"), key)
	if err != nil {
		return value, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return value, err
		}
	}

	return value, nil
}
