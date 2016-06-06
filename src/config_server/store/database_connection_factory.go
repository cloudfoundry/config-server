package store

import (
	"config_server/config"
	"database/sql"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"errors"
)

func DBConnection(config config.DBConfig) (*sql.DB, error) {
	var db *sql.DB

	connectionString, err := ConnectionString(config)
	if err != nil {
		return db, err
	}

	db, err = sql.Open(config.Adapter, connectionString)
	db.SetMaxOpenConns(config.ConnectionOptions.MaxOpenConnections)
	db.SetMaxIdleConns(config.ConnectionOptions.MaxIdleConnections)

	return db, err
}

func ConnectionString(config config.DBConfig) (string, error) {

	var connectionString string
	var err error

	switch config.Adapter {
	case "postgres":
		connectionString = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			config.User, config.Password, config.Name)
	case "mysql":
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			config.User, config.Password, config.Host, config.Port, config.Name)
	default:
		err = errors.New(fmt.Sprintf("Unsupported adapter: %s", config.Adapter))
	}

	return connectionString, err
}

