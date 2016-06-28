package server

import (
	"net/http"
	"strconv"
)

type configServerImpl struct {
	requestHandler RequestHandler
}

func NewServer(requestHandler RequestHandler) ConfigServer {
	return configServerImpl{
		requestHandler: requestHandler,
	}
}

func (server configServerImpl) Start(port int, certificateFilePath string, privateKeyFilePath string) error {
	http.HandleFunc("/v1/config/", server.requestHandler.HandleRequest)
	return http.ListenAndServeTLS(":"+strconv.Itoa(port), certificateFilePath, privateKeyFilePath, nil)
}
