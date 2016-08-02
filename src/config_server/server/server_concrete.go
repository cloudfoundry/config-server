package server

import (
	"net/http"
	"strconv"
    "config_server/config"
    "config_server/store"
)

type configServer struct {
    config config.ServerConfig
}

func NewConfigServer(config config.ServerConfig) ConfigServer {
	return configServer { config }
}

func (cs configServer) Start() error {

    cs.configureHandler()

	return http.ListenAndServeTLS(":" + strconv.Itoa(cs.config.Port),
        cs.config.CertificateFilePath,
        cs.config.PrivateKeyFilePath, nil)
}

func (cs configServer) configureHandler() {

    jwtTokenValidator, err := NewJwtTokenValidator(cs.config.JwtVerificationKeyPath)
    if err != nil {
        panic("Unable to start server\n" + err.Error())
    }

    requestHandler := NewRequestHandler(store.CreateStore(cs.config))
    authenticationHandler := NewAuthenticationHandler(jwtTokenValidator, requestHandler)

    http.Handle("/v1/data/", authenticationHandler)
}
