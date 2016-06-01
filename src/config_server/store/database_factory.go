package store

import (
	"config_server/config"
	"database/sql"
	"errors"
	"strings"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

func DBConnection(config config.DBConfig) (*sql.DB, error) {

	if &config == nil {
		return nil, errors.New("ServerConfig.Database required")
	}


	var dataSourceName string
	var adapter = strings.ToLower(config.Adapter)

	dataSourceName = SQLConnectionString(config)

	//max_connections: 12
	//pool_timeout: 25

	return sql.Open(adapter, dataSourceName)
}

