package main

import (
	"io"
	"net/http"
)

func dummy(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w,"dummy")
}

func main() {
	http.HandleFunc("/v1/config/dummy", dummy)
	http.ListenAndServe(":8000", nil)
}

