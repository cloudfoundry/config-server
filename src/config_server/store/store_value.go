package store

import (
	"encoding/json"
)

type StoreValue struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func (response StoreValue) Json() (string, error) {
	bytes, err := json.Marshal(&response)
	return string(bytes), err
}
