package store

import (
	"database/sql"
	"strconv"
)

type postgresStore struct {
	dbProvider DbProvider
}

func NewPostgresStore(dbProvider DbProvider) Store {
	return postgresStore{dbProvider}
}

func (ps postgresStore) Put(name string, value string) error {

	db, err := ps.dbProvider.Db()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO configurations (name, value) VALUES($1, $2)", name, value)
	if err != nil {
		_, err = db.Exec("UPDATE configurations SET value=$1 WHERE name=$2", value, name)
	}

	return err
}

func (ps postgresStore) GetByName(name string) (Configuration, error) {
	result := Configuration{}

	db, err := ps.dbProvider.Db()
	if err != nil {
		return result, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, name, value FROM configurations WHERE name = $1 ORDER BY id DESC LIMIT 1", name).Scan(&result.Id, &result.Name, &result.Value)
	if err == sql.ErrNoRows {
		return result, nil
	}

	return result, err
}

func (ps postgresStore) GetById(id string) (Configuration, error) {
	result := Configuration{}

	_, err := strconv.Atoi(id)
	if err != nil {
		return result, nil
	}

	db, err := ps.dbProvider.Db()
	if err != nil {
		return result, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, name, value FROM configurations WHERE id = $1", id).Scan(&result.Id, &result.Name, &result.Value)
	if err == sql.ErrNoRows {
		return result, nil
	}

	return result, err
}

func (ps postgresStore) Delete(name string) (bool, error) {

	db, err := ps.dbProvider.Db()
	if err != nil {
		return false, err
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM configurations WHERE name = $1", name)
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
