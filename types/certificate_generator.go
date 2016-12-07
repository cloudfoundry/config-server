package types

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"

	"github.com/cloudfoundry/bosh-utils/errors"
	"gopkg.in/yaml.v2"
)

type CertificateGenerator struct {
	loader CertsLoader
}

type CertResponse struct {
	Certificate string `json:"certificate" yaml:"certificate"`
	PrivateKey  string `json:"private_key" yaml:"private_key"`
	CA          string `json:"ca"          yaml:"ca"`
}

type certParams struct {
	CommonName       string   `yaml:"common_name"`
	AlternativeNames []string `yaml:"alternative_names"`
	IsCA             bool     `yaml:"is_ca"`
	CAName           string   `yaml:"ca"`
}

func NewCertificateGenerator(loader CertsLoader) CertificateGenerator {
	return CertificateGenerator{loader: loader}
}

func (cfg CertificateGenerator) Generate(parameters interface{}) (interface{}, error) {
	var params certParams
	err := objToStruct(parameters, &params)
	if err != nil {
		return nil, errors.Error("Failed to generate certificate, parameters are invalid.")
	}

	return cfg.generateCertificate(params)
}

func (cfg CertificateGenerator) generateCertificate(cParams certParams) (CertResponse, error) {
	var certResponse CertResponse

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
		BasicConstraintsValid: true,
		IsCA: cParams.IsCA,
	}

	var certificateRaw []byte
	var rootCARaw []byte

	if cParams.IsCA {
		template.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageCRLSign

		certificateRaw, err = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			return certResponse, errors.WrapError(err, "Generating CA certificate")
		}

		rootCARaw = certificateRaw
	} else {
		template.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		template.ExtKeyUsage = append(template.ExtKeyUsage, x509.ExtKeyUsageServerAuth)

		if cfg.loader == nil {
			panic("Expected CertificateGenerator to have Loader set")
		}
		rootCA, rootPKey, err := cfg.loader.LoadCerts(cParams.CAName)
		if err != nil {
			return certResponse, errors.WrapError(err, "Loading certificates")
		}

		for _, altName := range cParams.AlternativeNames {
			possibleIP := net.ParseIP(altName)
			if possibleIP == nil {
				template.DNSNames = append(template.DNSNames, altName)
			} else {
				template.IPAddresses = append(template.IPAddresses, possibleIP)
			}
		}

		certificateRaw, err = x509.CreateCertificate(rand.Reader, &template, rootCA, &privateKey.PublicKey, rootPKey)
		if err != nil {
			return certResponse, errors.WrapError(err, "Generating certificate")
		}
		rootCARaw = rootCA.Raw
	}

	encodedCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificateRaw})
	encodedPrivatekey := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	encodedRootCACert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootCARaw})

	certResponse = CertResponse{
		Certificate: string(encodedCert),
		PrivateKey:  string(encodedPrivatekey),
		CA:          string(encodedRootCACert),
	}

	return certResponse, nil
}

func objToStruct(input interface{}, str interface{}) error {
	valBytes, err := yaml.Marshal(input)
	if err != nil {
		return errors.WrapErrorf(err, "Expected input to be serializable")
	}

	err = yaml.Unmarshal(valBytes, str)
	if err != nil {
		return errors.WrapErrorf(err, "Expected input to be deserializable")
	}

	return nil
}
