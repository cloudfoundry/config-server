package store

import (
	"fmt"
	"strings"
	"config_server/config"
)

const _PLACEHOLDER = "?"

// ? vs $ inconsistency as parameter problem for mysql vs postgres
// Source: https://github.com/golang/go/issues/3602
// SQLReplacer replaces the characters "$" with the placeholder parameter for
// the given SQL engine.
func SQLReplacer(adapter string, query string) string {
	switch adapter {
	case "mysql", "sqlite":
		if _PLACEHOLDER != "?" && strings.Contains(query, _PLACEHOLDER) {
			return strings.Replace(query, _PLACEHOLDER, "?", -1)
		}
	case "postgres":
		for nParam := 1; strings.Contains(query, _PLACEHOLDER); nParam++ {
			query = strings.Replace(query, _PLACEHOLDER, fmt.Sprintf("$%d", nParam), 1)
		}
	default:
		panic("adapter not supported: " + adapter)
	}

	return query
}

func SQLConnectionString(config config.DBConfig) string {

	var connectionString string

	adapter := strings.ToLower(config.Adapter)

	switch adapter {
	case "postgres":
		connectionString = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			config.User, config.Password, config.Name)
	case "mysql":
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			config.User, config.Password, config.Host, config.Port, config.Name)
	}

	return connectionString
}