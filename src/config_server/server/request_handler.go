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
		http.Error(resWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
	case "DELETE":
		handler.handleDelete(key, req, resWriter)
	default:
		http.Error(resWriter, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (handler requestHandler) handleGet(key string, resWriter http.ResponseWriter) {

	value, err := handler.store.Get(key)
	if err != nil {
		http.Error(resWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if value == "" {
		http.Error(resWriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		respond(resWriter, value, http.StatusOK)
	}
}

func (handler requestHandler) handlePut(key string, req *http.Request, resWriter http.ResponseWriter) {
	value, err := readPutRequest(req)

	if err != nil {
		http.Error(resWriter, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.saveToStore(key, value)

	if err != nil {
		http.Error(resWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	resWriter.WriteHeader(http.StatusNoContent)
}

func (handler requestHandler) handlePost(key string, req *http.Request, resWriter http.ResponseWriter) {
	generationType, parameters, err := readPostRequest(req)

	if err != nil {
		http.Error(resWriter, err.Error(), http.StatusBadRequest)
		return
	}

	value, err := handler.store.Get(key)
	if value != "" {
		respond(resWriter, value, http.StatusOK)

	} else {
		generator, err := handler.valueGeneratorFactory.GetGenerator(generationType)
		if err != nil {
			http.Error(resWriter, "Unsupport type {put type here}", http.StatusBadRequest)
			return
		}

		value, err := generator.Generate(parameters)
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

func (handler requestHandler) handleDelete(key string, req *http.Request, resWriter http.ResponseWriter) {
	deleted, err := handler.store.Delete(key)

	if err == nil {
		if deleted {
			respond(resWriter, "", http.StatusNoContent)
		} else {
			respond(resWriter, "", http.StatusNotFound)
		}
	} else {
		http.Error(resWriter, err.Error(), http.StatusInternalServerError)
	}
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

func readPutRequest(req *http.Request) (interface{}, error) {
	jsonMap, err := readJSONBody(req)

	if err != nil {
		return nil, err
	}

	value, keyExist := jsonMap["value"]
	if !keyExist {
		return nil, errors.Error("JSON request body shoud contain the key 'value'")
	}

	return value, nil
}

func readPostRequest(req *http.Request) (string, interface{}, error) {
	jsonMap, err := readJSONBody(req)

	if err != nil {
		return "", nil, err
	}

	generationType, keyExist := jsonMap["type"]
	if !keyExist {
		return "", nil, errors.Error("JSON request body shoud contain the key 'type'")
	}

	return generationType.(string), jsonMap["parameters"], nil
}

func readJSONBody(req *http.Request) (map[string]interface{}, error) {
	if req == nil {
		return nil, errors.Error("Request can't be nil")
	}

	if req.Body == nil {
		return nil, errors.Error("Request can't be empty")
	}

	var f interface{}
	if err := json.NewDecoder(req.Body).Decode(&f); err != nil {
		return nil, errors.Error("Request Body should be JSON string")
	}

	return f.(map[string]interface{}), nil
}
