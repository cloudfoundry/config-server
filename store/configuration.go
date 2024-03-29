package store

import (
	"encoding/json"
)

type Configuration struct {
	ID                string
	Name              string
	Value             string
	ParameterChecksum string
}

func (rv Configuration) StringifiedJSON() (string, error) {
	var val map[string]interface{}

	err := json.Unmarshal([]byte(rv.Value), &val) //nolint:ineffassign,staticcheck

	val["id"] = rv.ID
	val["name"] = rv.Name
	bytes, err := json.Marshal(&val)

	return string(bytes), err
}
