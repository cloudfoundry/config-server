package server

import (
	"net/http"
	"strings"
	"encoding/json"
	"config_server/store"
)

type requestHandlerImpl struct {
	store          store.Store
	encoderDecoder store.EncoderDecoder
}

func NewConcreteRequestHandler(datastore store.Store) RequestHandler {
	return requestHandlerImpl{datastore, store.NewEncoderDecoder()}
}

func (handler requestHandlerImpl) HandleRequest(res http.ResponseWriter, req *http.Request) {
	if handler.store == nil {
		respondSmurf(res, http.StatusInternalServerError, "DB Store is nil")
		return
	}

	paths := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(paths) != 3 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	key := paths[len(paths)-1]

	switch req.Method {
	case "GET":
		handler.handleGet(key, res)
	case "PUT":
		handler.handlePut(key, req, res)
	default:
		res.WriteHeader(http.StatusNotFound)
	}
}

func (handler requestHandlerImpl) handleGet(key string, res http.ResponseWriter) {

	value, err := handler.store.Get(key)
	if err != nil {
		respondSmurf(res, http.StatusInternalServerError, err.Error())
		return
	}

	if value == "" {
		respondSmurf(res, http.StatusNotFound, "")
		return
	}

	decodedValue, encDecErr := handler.encoderDecoder.Decode(value)

	if encDecErr != nil {
		respondSmurf(res, http.StatusInternalServerError, encDecErr.Error())
		return
	}

	response, err := ConfigResponse{Path: key, Value: decodedValue}.Json()
	if err != nil {
		respondSmurf(res, http.StatusInternalServerError, err.Error())
		return
	}

	respondSmurf(res, http.StatusOK, response)
}

func (handler requestHandlerImpl) handlePut(key string, req *http.Request, res http.ResponseWriter) {

	type RequestBody struct {
		Value interface{}
	}
	var requestBody RequestBody

	if req.Body == nil {
		respondSmurf(res, http.StatusBadRequest, "Value cannot be empty")
		return
	}

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		respondSmurf(res, http.StatusInternalServerError, err.Error())
		return
	}

	encodedValue, encDecErr := handler.encoderDecoder.Encode(requestBody.Value)

	if encDecErr != nil {
		respondSmurf(res, http.StatusInternalServerError, encDecErr.Error())
		return
	}

	err = handler.store.Put(key, encodedValue)
	if err != nil {
		respondSmurf(res, http.StatusInternalServerError, err.Error())
		return
	}

	res.WriteHeader(http.StatusOK)
}

func respondSmurf(res http.ResponseWriter, status int, message string) {
	res.WriteHeader(status)
	_, err := res.Write([]byte(message))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}
