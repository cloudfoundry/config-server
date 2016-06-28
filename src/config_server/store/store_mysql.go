package store

import (
	"database/sql"
)

type mysqlStore struct {
	dbProvider DbProvider
}

func NewMysqlStore(dbProvider DbProvider) Store {
	return mysqlStore{dbProvider}
}

func (ms mysqlStore) Put(key string, value string) error {

	db, err := ms.dbProvider.Db()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO config VALUES(?,?)", key, value)
	if err != nil {
		_, err = db.Exec("UPDATE config SET config.config_value = ? WHERE config.config_key = ?", value, key)
	}

	return err
}

func (ms mysqlStore) Get(key string) (string, error) {

	var value string

	db, err := ms.dbProvider.Db()
	if err != nil {
		return value, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT config_value FROM config c WHERE c.config_key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return value, nil
	}

	return value, err
}
