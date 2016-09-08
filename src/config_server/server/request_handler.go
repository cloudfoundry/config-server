package server

import (
	"config_server/store"
	"config_server/types"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry/bosh-utils/errors"
)

type requestHandler struct {
	store                 store.Store
	valueGeneratorFactory types.ValueGeneratorFactory
}

func NewRequestHandler(store store.Store, valueGeneratorFactory types.ValueGeneratorFactory) (http.Handler, error) {
	if store == nil {
		return nil, errors.Error("Data store must be set")
	}
	return requestHandler{
		store: store,
		valueGeneratorFactory: valueGeneratorFactory,
	}, nil
}

func (handler requestHandler) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
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
	case "POST":
		handler.handlePost(key, req, resWriter)
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

	err := handler.readRequestBody(req, resWriter, &requestBody)
	if err != nil {
		return
	}

	err = handler.saveToStore(key, requestBody.Value)
	if err != nil {
		http.Error(resWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	resWriter.WriteHeader(http.StatusNoContent)
}

func (handler requestHandler) handlePost(key string, req *http.Request, resWriter http.ResponseWriter) {

	type RequestBody struct {
		Type       string
		Parameters interface{}
	}
	var requestBody RequestBody

	err := handler.readRequestBody(req, resWriter, &requestBody)
	if err != nil {
		return
	}

	value, err := handler.store.Get(key)
	if value != "" {
		respond(resWriter, value, http.StatusOK)

	} else {
		generator, err := handler.valueGeneratorFactory.GetGenerator(requestBody.Type)
		if err != nil {
			http.Error(resWriter, "Unable to create generator", http.StatusInternalServerError)
			return
		}

		value, err := generator.Generate(requestBody.Parameters)
		if err != nil {
			http.Error(resWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		err = handler.saveToStore(key, value)
		if err != nil {
			http.Error(resWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		value, err = handler.store.Get(key)
		if err != nil {
			http.Error(resWriter, err.Error(), http.StatusInternalServerError)
		}

		respond(resWriter, value.(string), http.StatusCreated)
	}
}

func (handler requestHandler) readRequestBody(req *http.Request, resWriter http.ResponseWriter, value interface{}) error {
	var err error

	if req.Body == nil {
		err = errors.Error("Value cannot be empty")
		http.Error(resWriter, err.Error(), http.StatusBadRequest)
	} else {
		err = json.NewDecoder(req.Body).Decode(value)
		if err != nil {
			http.Error(resWriter, err.Error(), http.StatusInternalServerError)
		}
	}

	return err
}

func (handler requestHandler) saveToStore(key string, value interface{}) error {

	storeValue, err := store.StoreValue{Path: key, Value: value}.Json()
	if err != nil {
		return err
	}

	err = handler.store.Put(key, storeValue)
	if err != nil {
		return err
	}

	return nil
}

func respond(res http.ResponseWriter, message string, status int) {
	res.WriteHeader(status)

	_, err := res.Write([]byte(message))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}
