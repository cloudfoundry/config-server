package main

import (
	"os"
	"io"
	"net/http"
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

func config(w http.ResponseWriter, r *http.Request) {
	handler.handle(w, r)
}

func main() {
	port := os.Args[1]
	handler = ConfigHandler{ make(map[string]string) }

	http.HandleFunc("/v1/config", config)
	http.ListenAndServe(":" + port, nil)
}
