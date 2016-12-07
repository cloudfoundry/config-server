package server

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/cloudfoundry/bosh-utils/errors"
	"github.com/dgrijalva/jwt-go"
)

type JwtTokenValidator struct {
	verificationKey *rsa.PublicKey
}

func NewJwtTokenValidator(jwtVerificationKeyPath string) (JwtTokenValidator, error) {
	bytes, err := ioutil.ReadFile(jwtVerificationKeyPath)
	if err != nil {
		return JwtTokenValidator{}, errors.WrapError(err, "Failed to read JWT Verification key")
	}

	verificationKey, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
	if err != nil {
		return JwtTokenValidator{}, errors.WrapError(err, "Failed to parse RSA public key from PEM")
	}

	return JwtTokenValidator{verificationKey: verificationKey}, nil
}

func (j JwtTokenValidator) Validate(token string) error {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return j.verificationKey, nil
	})

	if err != nil {
		return errors.WrapError(err, "Failed to validate token")
	}
	return nil
}
