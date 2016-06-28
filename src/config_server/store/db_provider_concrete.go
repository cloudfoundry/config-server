package store

import (
	"config_server/config"
	"errors"
	"fmt"
	"github.com/BurntSushi/migration"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type concreteDbProvider struct {
	config config.DBConfig
	sql    ISql
}

func NewConcreteDbProvider(sql ISql, config config.DBConfig) DbProvider {
	return concreteDbProvider{config, sql}
}

func (p concreteDbProvider) Db() (IDb, error) {

	var db IDb

	connectionString, err := p.connectionString(p.config)
	if err != nil {
		return db, err
	}

	db, err = p.sql.Open(p.config.Adapter, connectionString, Migrations)
	if err != nil {
		return db, err
	}

	db.SetMaxOpenConns(p.config.ConnectionOptions.MaxOpenConnections)
	db.SetMaxIdleConns(p.config.ConnectionOptions.MaxIdleConnections)

	return db, err
}

func (p concreteDbProvider) connectionString(config config.DBConfig) (string, error) {

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

var Migrations = []migration.Migrator{
	InitialSchema,
}

func InitialSchema(tx migration.LimitedTx) error {
	_, err := tx.Exec("CREATE TABLE config (config_key VARCHAR(255) NOT NULL UNIQUE, config_value TEXT NOT NULL)")
	return err
}
