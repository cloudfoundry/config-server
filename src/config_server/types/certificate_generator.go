package types

import (
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "log"
    "time"
    "math/big"
    "config_server/config"
)

type certificateGenerator struct {
    config config.ServerConfig
    loader CertsLoader
}

type CertParams struct {
    CommonName string
    AlternativeName []string
}

func NewCertificateGenerator(config config.ServerConfig, loader CertsLoader) certificateGenerator {
    return certificateGenerator{config: config, loader: loader}
}

func (cfg certificateGenerator) Generate(parameters interface{}) (interface{}, error) {
    params := parameters.(map[string]interface{})
    commonName := params["common_name"].(string)
    alternativeNames := []string{}

    if params["AlternativeNames"] != nil {
        alternativeNames = params["alternative_names"].([]string)
    }

    cParams := CertParams{CommonName: commonName, AlternativeName: alternativeNames}
    return cfg.GenerateCert(cParams)
}

func (cfg certificateGenerator) GenerateCert(cParams CertParams) (string, error) {
    parent, key, err := cfg.loader.LoadCerts(cfg.config.CertificateFilePath, cfg.config.PrivateKeyFilePath)

    if err != nil {
        return "", err
    }

    notBefore := time.Now()
    notAfter := notBefore.Add(365*24*time.Hour)

    serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
    serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
    if err != nil {
        log.Printf("failed to generate serial number:\n %s", err)
        return "", err
    }

    template := x509.Certificate {
        SerialNumber: serialNumber,
        Subject: pkix.Name {
            Organization: []string{"Internet Widgits Pty Ltd"},
        },
        NotBefore: notBefore,
        NotAfter:  notAfter,

        KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        BasicConstraintsValid: true,
    }
    template.DNSNames = append(template.DNSNames, cParams.CommonName)

    for _, altName := range cParams.AlternativeName {
        template.DNSNames = append(template.DNSNames, altName)
    }

    derBytes, err := x509.CreateCertificate(rand.Reader, &template, parent, &key.PublicKey, key)
    if err != nil {
        log.Printf("Failed to create certificate:\n%s", err)
        return "", err
    }

    cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
    return string(cert), nil
}
