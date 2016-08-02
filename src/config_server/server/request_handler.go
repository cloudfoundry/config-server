package server

import (
	"config_server/store"
	"encoding/json"
	"net/http"
	"strings"
)

type requestHandler struct {
    store store.Store
}

func NewRequestHandler(store store.Store) http.Handler {
    return requestHandler { store }
}

func (handler requestHandler) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
    if handler.store == nil {
        http.Error(resWriter, "DB Store is nil", http.StatusInternalServerError)
    } else {
        handler.handleRequest(resWriter, req)
    }
}

func (handler requestHandler) handleRequest(resWriter http.ResponseWriter, req *http.Request) {

	paths := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(paths) != 3 {
        http.Error(resWriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	key := paths[len(paths)-1]

	switch req.Method {
	case "GET":
		handler.handleGet(key, resWriter)
	case "PUT":
        handler.handlePut(key, req, resWriter)
	default:
        http.Error(resWriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (handler requestHandler) handleGet(key string, resWriter http.ResponseWriter) {

	value, err := handler.store.Get(key)
	if err != nil {
        http.Error(resWriter, err.Error(), http.StatusNotFound)
		return
	}

	if value == "" {
        http.Error(resWriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
        respond(resWriter, value, http.StatusOK)
    }
}

func (handler requestHandler) handlePut(key string, req *http.Request, resWriter http.ResponseWriter) {

	type RequestBody struct {
		Value interface{}
	}
	var requestBody RequestBody

	if req.Body == nil {
        http.Error(resWriter, "Value cannot be empty", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
        http.Error(resWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	storeValue, err := store.StoreValue{Path: key, Value: requestBody.Value}.Json()

	if err != nil {
        http.Error(resWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	err = handler.store.Put(key, storeValue)
	if err != nil {
        http.Error(resWriter, err.Error(), http.StatusInternalServerError)
		return
	}

    resWriter.WriteHeader(http.StatusOK)
}

func respond(res http.ResponseWriter, message string, status int) {
	res.WriteHeader(status)

	_, err := res.Write([]byte(message))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}
