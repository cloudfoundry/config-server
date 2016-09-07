package types

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"

	"github.com/cloudfoundry/bosh-utils/errors"
)

type x509Loader struct{}

func NewX509Loader() CertsLoader {
	return x509Loader{}
}

func (l x509Loader) LoadCerts(certFilePath, keyFilePath string) (*x509.Certificate, *rsa.PrivateKey, error) {
	crt, err := l.parseCertificate(certFilePath)
	if err != nil {
		return nil, nil, err
	}

	key, err := l.parsePrivateKey(keyFilePath)
	if err != nil {
		return nil, nil, err
	}

	return crt, key, nil
}

func (x509Loader) parseCertificate(certFilePath string) (*x509.Certificate, error) {
	cf, e := ioutil.ReadFile(certFilePath)
	if e != nil {
		err := errors.Error("Failed to load certificate file")
		log.Printf(err.Error())
		return nil, err
	}

	cpb, _ := pem.Decode(cf)
	crt, e := x509.ParseCertificate(cpb.Bytes)

	if e != nil {
		err := errors.WrapError(e, "Failed to parse certificate")
		log.Printf(err.Error())
		return nil, err
	}

	return crt, nil
}

func (x509Loader) parsePrivateKey(keyFilePath string) (*rsa.PrivateKey, error) {
	kf, e := ioutil.ReadFile(keyFilePath)
	if e != nil {
		err := errors.Error("Failed to load private key file")
		log.Printf(err.Error())
		return nil, err
	}

	kpb, _ := pem.Decode(kf)

	key, e := x509.ParsePKCS1PrivateKey(kpb.Bytes)
	if e != nil {
		err := errors.WrapError(e, "Failed to parse private key")
		log.Printf(err.Error())
		return nil, err
	}
	return key, nil
}
