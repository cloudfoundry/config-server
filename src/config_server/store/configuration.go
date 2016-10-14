package store

import (
	"encoding/json"
)


type Configuration struct {
	Id string
	Key string
	Value string
}

func (rv Configuration) StringifiedJSON() (string, error) {
	var val map[string]interface{}

	err := json.Unmarshal([]byte(rv.Value), &val)

	val["id"] = rv.Id
	val["path"] = rv.Key
	bytes, err := json.Marshal(&val)

	return string(bytes), err
}