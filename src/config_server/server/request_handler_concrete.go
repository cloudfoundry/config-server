package server

import (
	"config_server/store"
	"encoding/json"
	"net/http"
	"strings"
)

type requestHandlerImpl struct {
	store store.Store
}

func NewConcreteRequestHandler(datastore store.Store) RequestHandler {
	return requestHandlerImpl{datastore}
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

	respondSmurf(res, http.StatusOK, value)
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

	storeValue, err := store.StoreValue{Path: key, Value: requestBody.Value}.Json()

	if err != nil {
		respondSmurf(res, http.StatusInternalServerError, err.Error())
		return
	}

	err = handler.store.Put(key, storeValue)
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
