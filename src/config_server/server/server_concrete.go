package server

import (
	"net/http"
	"strconv"
    "config_server/config"
    "config_server/store"
    "config_server/types"
	"github.com/cloudfoundry/bosh-utils/errors"
)

type configServer struct {
    config config.ServerConfig
}

func NewConfigServer(config config.ServerConfig) ConfigServer {
	return configServer { config }
}

func (cs configServer) Start() error {
    if err := cs.configureHandler(); err != nil {
		return err
	}

	return http.ListenAndServeTLS(":" + strconv.Itoa(cs.config.Port),
        cs.config.CertificateFilePath,
        cs.config.PrivateKeyFilePath, nil)
}

func (cs configServer) configureHandler() error {

    jwtTokenValidator, err := NewJwtTokenValidator(cs.config.JwtVerificationKeyPath)
    if err != nil {
        return errors.WrapError(err, "Failed to create JWT token validator")
    }
	requestHandler,err := NewRequestHandler(store.CreateStore(cs.config), types.NewValueGeneratorConcrete(cs.config))
	if err != nil {
		return errors.WrapError(err, "Failed to create Request Handler")
	}
    authenticationHandler := NewAuthenticationHandler(jwtTokenValidator, requestHandler)

    http.Handle("/v1/data/", authenticationHandler)

	return nil
}
