package support

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const SERVER_URL string = "https://localhost:9000"

type requestBody struct {
	Value interface{} `json:"value"`
	Name string `json:"name"`
}

func SendGetRequestByName(name string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, SERVER_URL+"/v1/data?name="+name, nil)
	req.Header.Add("Authorization", "bearer "+ValidToken())

	return HTTPSClient.Do(req)
}

func SendGetRequestByID(id string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, SERVER_URL+"/v1/data/"+id, nil)
	req.Header.Add("Authorization", "bearer "+ValidToken())

	return HTTPSClient.Do(req)
}

func SendPutRequest(name string, value interface{}) (*http.Response, error) {
	data := requestBody{
		Value: value,
		Name: name,
	}

	requestBytes, _ := json.Marshal(&data)

	req, _ := http.NewRequest(http.MethodPut, SERVER_URL+"/v1/data", bytes.NewReader(requestBytes))
	req.Header.Add("Authorization", "bearer "+ValidToken())
	req.Header.Add("Content-Type", "application/json")

	return HTTPSClient.Do(req)
}

func SendPostRequest(name string, valueType string) (*http.Response, error) {
	var requestBytes *bytes.Reader

	switch valueType {
	case "password":
		requestBytes = bytes.NewReader([]byte(`{"name":"` + name + `","type":"password","parameters":{}}`))
	case "certificate":
		requestBytes = bytes.NewReader([]byte(`{"name":"` + name + `","type":"certificate","parameters":{"common_name": "burpees", "alternative_names":["cnj", "deadlift"]}}`))
	}

	req, _ := http.NewRequest(http.MethodPost, SERVER_URL+"/v1/data", requestBytes)
	req.Header.Add("Authorization", "bearer "+ValidToken())
	req.Header.Add("Content-Type", "application/json")

	return HTTPSClient.Do(req)
}

func SendDeleteRequest(name string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodDelete, SERVER_URL+"/v1/data?name="+name, nil)
	req.Header.Add("Authorization", "bearer "+ValidToken())

	return HTTPSClient.Do(req)
}
