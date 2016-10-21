package store

import (
	"encoding/json"
)

type Configuration struct {
	Id    string
	Name  string
	Value string
}

func (rv Configuration) StringifiedJSON() (string, error) {
	var val map[string]interface{}

	err := json.Unmarshal([]byte(rv.Value), &val)

	val["id"] = rv.Id
	val["name"] = rv.Name
	bytes, err := json.Marshal(&val)

	return string(bytes), err
}
