package server

import "encoding/json"

type ConfigResponse struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

func (response ConfigResponse) Json() (string, error) {
	bytes, err := json.Marshal(&response)
	return string(bytes), err
}
