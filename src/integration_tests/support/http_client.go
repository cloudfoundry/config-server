package support

import (
	"encoding/json"
	"net/http"
	"bytes"
)

const SERVER_URL string = "https://localhost:9000"

type requestBody struct {
	Value interface{} `json:"value"`
}

func SendGetRequest(key string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", SERVER_URL+"/v1/data/"+key, nil)
	req.Header.Add("Authorization", "bearer "+ValidToken())

	return HTTPSClient.Do(req)
}

func SendPutRequest(key string, value interface{}) (*http.Response, error) {
	data := requestBody{
		Value: value,
	}

	requestBytes, _ := json.Marshal(&data)

	req, _ := http.NewRequest("PUT", SERVER_URL+"/v1/data/"+key, bytes.NewReader(requestBytes))
	req.Header.Add("Authorization", "bearer "+ValidToken())

	return HTTPSClient.Do(req)
}
