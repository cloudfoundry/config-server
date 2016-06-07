package store

import (
	"config_server/config"
	"database/sql"
)

type MysqlStore struct {
	config config.DBConfig
}

func NewMysqlStore(config config.DBConfig) MysqlStore {
	return MysqlStore{config}
}

func (store MysqlStore) Put(key string, value string) error {

	db, err := DBConnection(store.config)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, _ := db.Query("SELECT value FROM config c WHERE c.key = ?", key)
	if rows.Next() == false {
		_, err = db.Exec("INSERT INTO config VALUES(?,?)", key, value)
	} else {
		_, err = db.Exec("UPDATE config SET config.value = ? WHERE config.key = ?", value, key)
	}

	return err
}

func (store MysqlStore) Get(key string) (string, error) {

	var value string

	db, err := DBConnection(store.config)
	if err != nil {
		return value, err
 	}
	defer db.Close()

	err = db.QueryRow("SELECT value FROM config c WHERE c.key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return value, nil
	}

	return value, err
}
