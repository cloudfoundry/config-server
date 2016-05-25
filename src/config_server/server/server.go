package server

import (
	"io"
	"net/http"
	"strconv"
)

var handler ConfigHandler

type ConfigHandler struct {
	db map[string]string
}

func (c ConfigHandler) handle(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		key := req.FormValue("key")
		result := c.db[key]
		io.WriteString(res, result)
	case "PUT":
		key := req.FormValue("key")
		value := req.FormValue("value")
		c.db[key] = value
	default:
		res.WriteHeader(http.StatusNotFound)
	}
}

func StartServer(port int) {
	handler = ConfigHandler{make(map[string]string)}

	http.HandleFunc("/v1/config", handler.handle)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
