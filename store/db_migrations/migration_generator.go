package db_migrations

import (
	"strings"

	"github.com/BurntSushi/migration"
)

func GetMigrations(adapter string) []migration.Migrator {
	var migrations []string

	if strings.EqualFold(adapter, "mysql") {
		migrations = MysqlMigrations()
	} else if strings.EqualFold(adapter, "postgres") {
		migrations = PostgresMigrations()
	}

	var result []migration.Migrator

	for _, mig := range migrations {
		query := func(tx migration.LimitedTx) error {
			_, err := tx.Exec(mig)
			return err
		}

		result = append(result, query)
	}

	return result
}
