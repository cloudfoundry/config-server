package server

import (
	"config_server/store"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"encoding/json"
)

type ConfigServer struct {
	store store.Store
}

func NewServer(store store.Store) ConfigServer {
	return ConfigServer{
		store: store,
	}
}

func (server ConfigServer) Start(port int, certificateFilePath string, privateKeyFilePath string) error {
	if server.store == nil {
		return errors.New("DataStore can not be nil")
	}

	http.HandleFunc("/v1/config/", server.HandleRequest)
	return http.ListenAndServeTLS(":"+strconv.Itoa(port), certificateFilePath, privateKeyFilePath, nil)
}

func (server ConfigServer) HandleRequest(res http.ResponseWriter, req *http.Request) {

	paths := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	if len(paths) != 3 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	key := paths[len(paths)-1]

	switch req.Method {
	case "GET":
		server.handleGet(key, res)
	case "PUT":
		server.handlePut(key, req, res)
	default:
		res.WriteHeader(http.StatusNotFound)
	}
}

func (server ConfigServer) handleGet(key string, res http.ResponseWriter) {

	value, err := server.store.Get(key)
	if err != nil {
		respond(res, http.StatusInternalServerError, err.Error())
		return
	}

	if value == "" {
		respond(res, http.StatusNotFound, "")
		return
	}

	response, err := ConfigResponse{Path: key, Value: value}.Json()
	if err != nil {
		respond(res, http.StatusInternalServerError, err.Error())
		return
	}

	respond(res, http.StatusOK, response)
}

func (server ConfigServer) handlePut(key string, req *http.Request, res http.ResponseWriter) {

	type RequestBody struct {
		Value string
	}
	var requestBody RequestBody

	if req.Body == nil {
		respond(res, http.StatusBadRequest, "Value cannot be empty")
		return
	}

	err := json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		respond(res, http.StatusInternalServerError, err.Error())
		return
	}

	if requestBody.Value == "" {
		respond(res, http.StatusBadRequest, "Value cannot be empty")
		return
	}

	err = server.store.Put(key, requestBody.Value)
	if err != nil {
		respond(res, http.StatusInternalServerError, err.Error())
		return
	}

	res.WriteHeader(http.StatusOK)
}

func respond(res http.ResponseWriter, status int, message string) {
	res.WriteHeader(status)
	_, err := res.Write([]byte(message))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}
