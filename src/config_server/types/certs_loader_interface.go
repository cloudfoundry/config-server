package types

import (
	"crypto/rsa"
	"crypto/x509"
)

type CertsLoader interface {
	LoadCerts(certFile, keyFile string) (*x509.Certificate, *rsa.PrivateKey, error)
}
