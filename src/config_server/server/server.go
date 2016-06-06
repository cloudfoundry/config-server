package server

import (
	"config_server/store"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ConfigServer struct {
	store store.Store
}

func NewServer(store store.Store) ConfigServer {
	return ConfigServer{
		store: store,
	}
}

func (server ConfigServer) Start(port int) error {
	if server.store == nil {
		return errors.New("DataStore can not be nil")
	}

	http.HandleFunc("/v1/config/", server.HandleRequest)
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
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
		server.handlePut(key, req.FormValue("value"), res)
	default:
		res.WriteHeader(http.StatusNotFound)
	}
}

func (server ConfigServer) handleGet(key string, res http.ResponseWriter) {

	value, err := server.store.Get(key)
	if err != nil {
		server.respond(res, http.StatusInternalServerError, err.Error())
		return
	}

	if value == "" {
		server.respond(res, http.StatusNotFound, "")
		return
	}

	response, err := ConfigResponse{Path: key, Value: value}.Json()
	if err != nil {
		server.respond(res, http.StatusInternalServerError, err.Error())
		return
	}

	server.respond(res, http.StatusOK, response)
}

func (server ConfigServer) handlePut(key string, value string, res http.ResponseWriter) {

	if value == "" {
		server.respond(res, http.StatusBadRequest, "Value cannot be empty")
		return
	}

	err := server.store.Put(key, value)
	if err != nil {
		server.respond(res, http.StatusInternalServerError, err.Error())
		return
	}

	res.WriteHeader(http.StatusOK)
}

func (server ConfigServer) respond(res http.ResponseWriter, status int, message string) {
	res.WriteHeader(status)
	io.WriteString(res, message)
}
