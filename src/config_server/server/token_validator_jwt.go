package server

import (
    "github.com/dgrijalva/jwt-go"
    "crypto/rsa"
    "io/ioutil"
)

type jwtTokenValidator struct {
    verificationKey *rsa.PublicKey
}

func NewJwtTokenValidator(jwtVerificationKeyPath string) (jwtTokenValidator, error) {

    bytes, err := ioutil.ReadFile(jwtVerificationKeyPath)
    if err != nil {
        return jwtTokenValidator{}, err
    }

    verificationKey, err := jwt.ParseRSAPublicKeyFromPEM(bytes)
    return jwtTokenValidator {verificationKey}, err
}

func (j jwtTokenValidator) Validate(token string) error {
    _, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        return j.verificationKey, nil
    })
    return err
}