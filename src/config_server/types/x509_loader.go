package types

import (
    "crypto/x509"
    "crypto/rsa"
    "io/ioutil"
    "log"
    "encoding/pem"
)

type x509Loader struct {}

func NewX509Loader() CertsLoader {
    return x509Loader{}
}

func (l x509Loader) LoadCerts(certFilePath, keyFilePath string) (*x509.Certificate, *rsa.PrivateKey, error) {
    cf, e := ioutil.ReadFile(certFilePath)
    if e != nil {
        log.Printf("Failed to load parent certificate file\n %v", e.Error())
        return nil, nil, e
    }

    kf, e := ioutil.ReadFile(keyFilePath)
    if e != nil {
        log.Printf("Failed to load parent key file:\n%v", e.Error())
        return nil, nil, e
    }

    cpb, _ := pem.Decode(cf)
    kpb, _ := pem.Decode(kf)
    crt, e := x509.ParseCertificate(cpb.Bytes)

    if e != nil {
        log.Printf("Failed to parse parent certificate:\n%v", e.Error())
        return nil, nil, e
    }

    key, e := x509.ParsePKCS1PrivateKey(kpb.Bytes)
    if e != nil {
        log.Printf("Failed to parse parent key:\n%v", e.Error())
        return nil, nil, e
    }

    return crt, key, nil
}
