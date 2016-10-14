package server

import (
	"config_server/store"
	"config_server/types"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cloudfoundry/bosh-utils/errors"
	"regexp"
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
	key, err := handler.extractKeyFromURLPath(req.URL.Path)

	if err != nil {
		http.Error(resWriter, err.Error(), http.StatusBadRequest)
		return
	}

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

	emptyValue := store.Configuration{}

	if value == emptyValue {
		http.Error(resWriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		result, _ := value.StringifiedJSON()
		handler.respond(resWriter, result, http.StatusOK)
	}
}

func (handler requestHandler) handlePut(key string, req *http.Request, resWriter http.ResponseWriter) {
	value, err := handler.readPutRequest(req)

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
	generationType, parameters, err := handler.readPostRequest(req)

	if err != nil {
		http.Error(resWriter, err.Error(), http.StatusBadRequest)
		return
	}

	emptyValue := store.Configuration{}

	value, err := handler.store.Get(key)
	if value != emptyValue {
		result, _ := value.StringifiedJSON()
		handler.respond(resWriter, result, http.StatusOK)

	} else {
		generator, err := handler.valueGeneratorFactory.GetGenerator(generationType)
		if err != nil {
			http.Error(resWriter, "Unsupport type {put type here}", http.StatusBadRequest)
			return
		}

		generatedValue, err := generator.Generate(parameters)
		if err != nil {
			http.Error(resWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		err = handler.saveToStore(key, generatedValue)
		if err != nil {
			http.Error(resWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		configuration, err := handler.store.Get(key)
		if err != nil {
			http.Error(resWriter, err.Error(), http.StatusInternalServerError)
		}

		result, _ := configuration.StringifiedJSON()
		handler.respond(resWriter, result, http.StatusCreated)
	}
}

func (handler requestHandler) handleDelete(key string, req *http.Request, resWriter http.ResponseWriter) {
	deleted, err := handler.store.Delete(key)

	if err == nil {
		if deleted {
			handler.respond(resWriter, "", http.StatusNoContent)
		} else {
			handler.respond(resWriter, "", http.StatusNotFound)
		}
	} else {
		http.Error(resWriter, err.Error(), http.StatusInternalServerError)
	}
}

func (handler requestHandler) saveToStore(key string, value interface{}) error {
	configValue :=  make(map[string]interface{})
	configValue["value"] = value

	bytes, err := json.Marshal(&configValue)

	if err != nil {
		return err
	}

	err = handler.store.Put(key, string(bytes))
	if err != nil {
		return err
	}

	return nil
}

func (handler requestHandler) respond(res http.ResponseWriter, message string, status int) {
	res.WriteHeader(status)

	_, err := res.Write([]byte(message))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func (handler requestHandler) readPutRequest(req *http.Request) (interface{}, error) {
	jsonMap, err := handler.readJSONBody(req)

	if err != nil {
		return nil, err
	}

	value, keyExist := jsonMap["value"]
	if !keyExist {
		return nil, errors.Error("JSON request body shoud contain the key 'value'")
	}

	return value, nil
}

func (handler requestHandler) readPostRequest(req *http.Request) (string, interface{}, error) {
	jsonMap, err := handler.readJSONBody(req)

	if err != nil {
		return "", nil, err
	}

	generationType, keyExist := jsonMap["type"]
	if !keyExist {
		return "", nil, errors.Error("JSON request body shoud contain the key 'type'")
	}

	return generationType.(string), jsonMap["parameters"], nil
}

func (handler requestHandler) readJSONBody(req *http.Request) (map[string]interface{}, error) {
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

func (handler requestHandler) extractKeyFromURLPath(path string) (string, error) {
	paths := strings.Split(strings.Trim(path, "/"), "/")

	if len(paths) < 3 {
		return "", errors.Error("Request URL invalid, seems to be missing key")
	}

	tokens := paths[2:]

	var validKeyToken = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

	for _, token := range tokens {
		if !validKeyToken.MatchString(token) {
			return "", errors.Error("Key must consist of alphanumeric, underscores, dashes, and forward slashes")
		}
	}

	return strings.Join(tokens, "/"), nil
}
