package types

import (
    "crypto/x509"
    "crypto/rsa"
)

type CertsLoader interface {
    LoadCerts(certFile, keyFile string) (*x509.Certificate, *rsa.PrivateKey, error)
}
