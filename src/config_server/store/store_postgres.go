package store

import (
	"database/sql"
)

type postgresStore struct {
	dbProvider DbProvider
}

func NewPostgresStore(dbProvider DbProvider) Store {
	return postgresStore{dbProvider}
}

func (ps postgresStore) Put(key string, value string) error {

	db, err := ps.dbProvider.Db()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO config VALUES($1, $2)", key, value)
	if err != nil {
		_, err = db.Exec("UPDATE config SET config_value=$1 WHERE config_key=$2", value, key)
	}

	return err
}

func (ps postgresStore) Get(key string) (string, error) {

	var value string

	db, err := ps.dbProvider.Db()
	if err != nil {
		return value, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT config_value FROM config WHERE config_key = $1", key).Scan(&value)
	if err == sql.ErrNoRows {
		return value, nil
	}

	return value, err
}

func (ps postgresStore) Delete(key string) (bool, error) {

	db, err := ps.dbProvider.Db()
	if err != nil {
		return false, err
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM config WHERE config_key = $1", key)
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
