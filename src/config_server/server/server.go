package server

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

var handler ConfigHandler

type ConfigHandler struct {
	db map[string]string
}

func StartServer(port int) {
	handler = ConfigHandler{make(map[string]string)}

	http.HandleFunc("/v1/config/", handler.handle)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (c ConfigHandler) handle(res http.ResponseWriter, req *http.Request) {
	paths := strings.Split(req.URL.Path, "/")

	if len(paths) != 4 { // We only accept /<version>/config/<key>
		res.WriteHeader(http.StatusNotFound)
	}

	key := paths[len(paths)-1]

	switch req.Method {
	case "GET":
		value := c.db[key]
		if value == "" {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		response, err := ConfigResponse{Path: key, Value: value}.Json()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			io.WriteString(res, err.Error())
		} else {
			io.WriteString(res, response)
		}

	case "PUT":
		value := req.FormValue("value")
		c.db[key] = value

	default:
		res.WriteHeader(http.StatusNotFound)
	}
}
