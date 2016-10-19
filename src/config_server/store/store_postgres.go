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

	_, err = db.Exec("INSERT INTO configurations (config_key, value) VALUES($1, $2)", key, value)
	if err != nil {
		_, err = db.Exec("UPDATE configurations SET value=$1 WHERE config_key=$2", value, key)
	}

	return err
}

func (ps postgresStore) GetByKey(key string) (Configuration, error) {
	result := Configuration{}

	db, err := ps.dbProvider.Db()
	if err != nil {
		return result, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, config_key, value FROM configurations WHERE config_key = $1 ORDER BY id DESC LIMIT 1", key).Scan(&result.Id, &result.Key, &result.Value)
	if err == sql.ErrNoRows {
		return result, nil
	}

	return result, err
}

func (ps postgresStore) GetById(id string) (Configuration, error) {
	result := Configuration{}

	db, err := ps.dbProvider.Db()
	if err != nil {
		return result, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, config_key, value FROM configurations WHERE id = $1", id).Scan(&result.Id, &result.Key, &result.Value)
	if err == sql.ErrNoRows {
		return result, nil
	}

	return result, err
}

func (ps postgresStore) Delete(key string) (bool, error) {

	db, err := ps.dbProvider.Db()
	if err != nil {
		return false, err
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM configurations WHERE config_key = $1", key)
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
