package acceptance_tests

import (
	"os"
	"path/filepath"
	"io/ioutil"
	"crypto/tls"
	"crypto/x509"
	"net/http"
)

func HTTPSClient() *http.Client {
	sslCertPath := pathFor("ssl.crt")
	sslKeyPath := pathFor("ssl.key")
	rootCAPath := pathFor("ssl_root_ca.crt")

	cert, err := tls.LoadX509KeyPair(sslCertPath, sslKeyPath)
	if err != nil {
		panic(err.Error())
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(rootCAPath)
	if err != nil {
		panic(err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	client := &http.Client{Transport: transport}

	return client
}

func ValidToken() string {
	tokenPath := pathFor("uaa.token")
	dat, err := ioutil.ReadFile(tokenPath)

	if err != nil {
		panic(err.Error())
	}

	return string(dat)
}

func pathFor(fileName string) string {
	var path, rootDir string

	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err.Error())
	}

	path = filepath.Join(rootDir, "src", "integration_tests", "assets", fileName)

	return path
}