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

func (ms mysqlStore) Put(name string, value string) error {

	db, err := ms.dbProvider.Db()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO configurations (name, value) VALUES(?,?)", name, value)

	if err != nil {
		_, err = db.Exec("UPDATE configurations SET value = ? WHERE name = ?", value, name)
	}

	return err
}

func (ms mysqlStore) GetByName(name string) (Configuration, error) {
	result := Configuration{}

	db, err := ms.dbProvider.Db()
	if err != nil {
		return result, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, name, value FROM configurations WHERE name = ? ORDER BY id DESC LIMIT 1", name).Scan(&result.ID, &result.Name, &result.Value)
	if err == sql.ErrNoRows {
		return result, nil
	}

	return result, err
}

func (ms mysqlStore) GetByID(id string) (Configuration, error) {
	result := Configuration{}

	db, err := ms.dbProvider.Db()
	if err != nil {
		return result, err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, name, value FROM configurations WHERE id = ?", id).Scan(&result.ID, &result.Name, &result.Value)
	if err == sql.ErrNoRows {
		return result, nil
	}

	return result, err
}

func (ms mysqlStore) Delete(name string) (bool, error) {

	db, err := ms.dbProvider.Db()
	if err != nil {
		return false, err
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM configurations WHERE name = ?", name)
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
