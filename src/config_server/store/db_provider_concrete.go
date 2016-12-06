package store

import (
	"fmt"
	"github.com/cloudfoundry/bosh-utils/errors"

	// blank import to load database drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"config_server/config"
	"config_server/store/db_migrations"
)

type concreteDbProvider struct {
	config config.DBConfig
	sql    ISql
	db     IDb
}

func NewConcreteDbProvider(sql ISql, config config.DBConfig) (DbProvider, error) {
	connectionString, err := connectionString(config)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to generate DB connection string")
	}

	db, err := sql.Open(config.Adapter, connectionString, db_migrations.GetMigrations(config.Adapter))
	if err != nil {
		return nil, errors.WrapError(err, "Failed to open connection to DB")
	}

	db.SetMaxOpenConns(config.ConnectionOptions.MaxOpenConnections)
	db.SetMaxIdleConns(config.ConnectionOptions.MaxIdleConnections)

	provider := concreteDbProvider{db: db}
	return provider, err
}

func (p concreteDbProvider) Db() (IDb, error) {
	if p.db == nil {
		return nil, errors.Error("Database not initialized")
	}
	return p.db, nil
}

func connectionString(config config.DBConfig) (string, error) {

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
		err = errors.Errorf("Unsupported adapter: %s", config.Adapter)
	}

	return connectionString, err
}
