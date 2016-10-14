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

	_, err = db.Exec("INSERT INTO configurations (config_key, value) VALUES(?,?)", key, value)

	if err != nil {
		_, err = db.Exec("UPDATE configurations SET value = ? WHERE config_key = ?", value, key)
	}

	return err
}

func (ms mysqlStore) Get(key string) (Configuration, error) {
	result := Configuration{}

	db, err := ms.dbProvider.Db()
	if err != nil {
		return result, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, config_key, value FROM configurations WHERE config_key = ? ORDER BY id DESC LIMIT 1", key).Scan(&result.Id, &result.Key, &result.Value)
	if err == sql.ErrNoRows {
		return result, nil
	}

	return result, err
}

func (ms mysqlStore) Delete(key string) (bool, error) {

	db, err := ms.dbProvider.Db()
	if err != nil {
		return false, err
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM configurations WHERE config_key = ?", key)
	if (err != nil) || (result == nil) {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	deleted := rows > 0
	return deleted, err
}
