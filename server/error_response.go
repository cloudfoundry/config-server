package server

import "encoding/json"

type errorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(err error) errorResponse {
	return errorResponse{Error: err.Error()}
}

func (e errorResponse) GenerateErrorMsg() string {
	response, err := json.Marshal(e)
	if err != nil {
		return `{"error": "Unknown Error"}`
	}

	return string(response)
}
