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

	http.HandleFunc("/v1/config/", server.handleRequest)
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (server ConfigServer) handleRequest(res http.ResponseWriter, req *http.Request) {
	paths := strings.Split(req.URL.Path, "/")

	if len(paths) != 4 { // We only accept /<version>/config/<key>
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
		res.WriteHeader(http.StatusInternalServerError)
		io.WriteString(res, err.Error())
		return
	}

	if value == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := ConfigResponse{Path: key, Value: value}.Json()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		io.WriteString(res, err.Error())
		return
	}

	res.WriteHeader(http.StatusOK)
	io.WriteString(res, response)
}

func (server ConfigServer) handlePut(key string, value string, res http.ResponseWriter) {
	err := server.store.Put(key, value)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		io.WriteString(res, err.Error())
		return
	}
	res.WriteHeader(http.StatusOK)
}
