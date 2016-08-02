package server

import (
    "net/http"
    "strings"
    "errors"
)

type authenticationHandler struct {
    tokenValidator TokenValidator
    nextHandler http.Handler
}

func NewAuthenticationHandler(tokenValidator TokenValidator, nextHandler http.Handler) http.Handler {
    return authenticationHandler { tokenValidator, nextHandler }
}

func (handler authenticationHandler) ServeHTTP(resWriter http.ResponseWriter, req *http.Request) {
    err := handler.authenticate(req)
    if err == nil {
        handler.nextHandler.ServeHTTP(resWriter, req)
    } else {
        http.Error(resWriter, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
    }
}

func (handler authenticationHandler) authenticate(req *http.Request) error {
    authHeader := req.Header.Get("Authorization")
    if len(authHeader) == 0 {
        return errors.New("Missing Token")
    }

    jwtToken, err := handler.checkTokenFormat(authHeader)
    if err == nil {
        err = handler.tokenValidator.Validate(jwtToken)
    }

    return err
}

func (handler authenticationHandler) checkTokenFormat(token string) (string, error) {
    tokenParts := strings.Split(token, " ")
    if len(tokenParts) != 2 {
        return "", errors.New("Invalid token format")
    }

    tokenType, userToken := tokenParts[0], tokenParts[1]
    if !strings.EqualFold(tokenType, "bearer") {
        return "", errors.New("Invalid token type: " + tokenType)
    }

    return userToken, nil
}
