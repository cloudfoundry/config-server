package server

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dgrijalva/jwt-go"
)

type jwtTokenValidator struct {
	verificationKey *rsa.PublicKey
}

func NewJwtTokenValidator(jwtVerificationKeyPath string) (jwtTokenValidator, error) {
	bytes, err := ioutil.ReadFile(jwtVerificationKeyPath)
	if err != nil {
		return jwtTokenValidator{}, errors.WrapError(err, "Failed to read JWT Verification key")
	}

	if verificationKey, err := jwt.ParseRSAPublicKeyFromPEM(bytes); err != nil {
		return jwtTokenValidator{}, errors.WrapError(err, "Failed to parse RSA public key from PEM")
	} else {
		return jwtTokenValidator{verificationKey: verificationKey}, nil
	}
}

func (j jwtTokenValidator) Validate(token string) error {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return j.verificationKey, nil
	})

	if err != nil {
		return errors.WrapError(err, "Failed to validate token")
	} else {
		return nil
	}
}
