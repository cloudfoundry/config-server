package types

import (
	"config_server/config"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/cloudfoundry/bosh-utils/errors"
)

type CertificateGenerator struct {
	config config.ServerConfig
	loader CertsLoader
}

type CertResponse struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
	CA          string `json:"ca"`
}

type CertParams struct {
	CommonName      string
	AlternativeName []string
}

func NewCertificateGenerator(config config.ServerConfig, loader CertsLoader) CertificateGenerator {
	return CertificateGenerator{config: config, loader: loader}
}

func (cfg CertificateGenerator) Generate(parameters interface{}) (interface{}, error) {
	params := parameters.(map[string]interface{})
	commonName := params["common_name"].(string)
	alternativeNames := []string{}

	if params["alternative_names"] != nil {
		for _, altName := range params["alternative_names"].([]interface{}) {
			alternativeNames = append(alternativeNames, altName.(string))
		}
	}

	cParams := CertParams{CommonName: commonName, AlternativeName: alternativeNames}
	return cfg.generateCert(cParams)
}

func (cfg CertificateGenerator) generateCert(cParams CertParams) (CertResponse, error) {
	var certResponse CertResponse

	rootCA, rootCAKey, err := cfg.loader.LoadCerts(cfg.config.CACertificateFilePath, cfg.config.CAPrivateKeyFilePath)
	if err != nil {
		return certResponse, errors.WrapError(err, "Loading certificates")
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return certResponse, errors.WrapError(err, "Generating Serial Number")
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return certResponse, errors.WrapError(err, "Generating Key")
	}

	now := time.Now()
	notAfter := now.Add(365 * 24 * time.Hour)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:      []string{"USA"},
			Organization: []string{"Cloud Foundry"},
			CommonName:   cParams.CommonName,
		},
		NotBefore:             now,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA: false,
	}

	for _, altName := range cParams.AlternativeName {
		template.DNSNames = append(template.DNSNames, altName)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, rootCA, &privateKey.PublicKey, rootCAKey)
	if err != nil {
		return certResponse, errors.WrapError(err, "Generating Certificate")
	}

	encodedCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	encodedPrivatekey := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	encodedRootCACert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootCA.Raw})

	certResponse = CertResponse{
		Certificate: string(encodedCert),
		PrivateKey:  string(encodedPrivatekey),
		CA:          string(encodedRootCACert),
	}

	return certResponse, nil
}
