package store

import (
	"fmt"
	"strings"
)

const _PLACEHOLDER = "?"

// Why: ? vs $ inconsistency as bind parameter for mysql vs postgres
// Source: https://github.com/golang/go/issues/3602
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
