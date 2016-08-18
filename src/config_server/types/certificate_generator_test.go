package types_test

import (
    . "config_server/types"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "config_server/config"
    "config_server/types/fakes"
    "crypto/x509"
    "encoding/pem"
    "time"
)

func getCert(generator ValueGenerator, certString string) (*x509.Certificate, error) {
    block, _ := pem.Decode([]byte(certString))
    crt, err := x509.ParseCertificate(block.Bytes)

    return crt, err
}

var _ = Describe("CertificateGenerator", func() {

    Describe("certificateGenerator", func() {
        var (
            fakeLoader *fakes.FakeCertsLoader
            generator ValueGenerator
        )

        mockCertValue := `-----BEGIN CERTIFICATE-----
MIIDtzCCAp+gAwIBAgIJAMiFskqEjVfoMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwIBcNMTYwODE4MTUzNzE5WhgPMjI5MDA2MDIxNTM3MTla
MEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJ
bnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDIpbxZT/eoHNnJIyitzzoGv9iKw9wEFrRTBZ4I/L9Tu7zO+2TogNnL
JaJ/dGAg3x4pIjMwCxIfEK30/R6RzkJXPaCNfI2m8Jdm4EVKVoXKmpxBibjETUpM
wBoU57rc3nfiPRZXEe4rNJlFOtpMEQ62nN13ixTrdC4aBBtQn7BmdeEsxUHjfHVg
jk16S7uAbQyfO0pA1PNSguxNe+X53jja2iuIjwtU5iV0KmQ+haJA2eHlKDxVQiwZ
zbxHrc8xtQ5rZZlv5RdF4sx3G3G7kqNt+WEUElCd2R5LfpWmBQs0rIYIJS19B+Bd
KQoRF5TPaQGw7tlcU04d9cokNDmL4UCXAgMBAAGjgacwgaQwHQYDVR0OBBYEFHqH
VbSJvj/RzwYUZPl0oWAh9WH5MHUGA1UdIwRuMGyAFHqHVbSJvj/RzwYUZPl0oWAh
9WH5oUmkRzBFMQswCQYDVQQGEwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8G
A1UEChMYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkggkAyIWySoSNV+gwDAYDVR0T
BAUwAwEB/zANBgkqhkiG9w0BAQUFAAOCAQEAXWgmomCv/R2CnRFayBXqSyok0cBt
j624xYjqS/ilQmJ7xw0FNOCST7bLSJpTZBrN5fb/5Gzb1YZIHw3AIl+epKBYdPvw
Z2IWam9o3AkdJb3lfsbfG40DCqclqMA8ufyJMU8MdEt4ekpswOdFTqQBjJSB8OgS
6AnxtQ/NwUqZTzmFgfIZqTTtAOBpuuwOWbsF8eOJUKW3cz4yQ3wzVI55wbceM1tI
tEZ5+1W5gwhBaWc9orhRJ+PAqaDsMdoJtS5Q6K7XhujinNbJsC9duFsVVfDeWBk/
s3OzXXaX6jGwhORFMqjszMKank/52HGOQKe/NbDR4YiUakj5bHtCkRuYEg==
-----END CERTIFICATE-----`

        mockKeyValue := `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAyKW8WU/3qBzZySMorc86Br/YisPcBBa0UwWeCPy/U7u8zvtk
6IDZyyWif3RgIN8eKSIzMAsSHxCt9P0ekc5CVz2gjXyNpvCXZuBFSlaFypqcQYm4
xE1KTMAaFOe63N534j0WVxHuKzSZRTraTBEOtpzdd4sU63QuGgQbUJ+wZnXhLMVB
43x1YI5Neku7gG0MnztKQNTzUoLsTXvl+d442toriI8LVOYldCpkPoWiQNnh5Sg8
VUIsGc28R63PMbUOa2WZb+UXReLMdxtxu5KjbflhFBJQndkeS36VpgULNKyGCCUt
fQfgXSkKEReUz2kBsO7ZXFNOHfXKJDQ5i+FAlwIDAQABAoIBADXnyPVZtZh1v9D4
4Cnp/ZXonr2RJx/2FZYkcTPK64AMdQpKoe6RB3I7dy+0zUqnhsvYQ1ZJ8yFjcsjZ
5AeaGVqk1OiOKORLyPE7mYICQsmOxiIZZQlgFnEOPzOUmjBSmHSworrbt1fmNrNR
v2omPwSymhFOzV0Aho44wjnj3Rl5YerEg6ai+bf6pkTUn27WknvpZKwNlWy0/MMk
VCz5oXoqrIDLVaB0tISu81XrG0wzlQru0/hcTJguBUfNdAFwoV6VZBrr/yNNLlwN
ylSSkNHI/zQIzGJ7BRO7mpNzQSkkpgJfKnitlOn6VC/rnQfi8XAX1JRJe43SAMVm
fqvwDFECgYEA8mNsxMDyvSd/Lg1q8jAb2V7K9qfErWPOwqRTrJ1n1M9483rpUp0U
9zjFhLFFvaGraBkqjuJpw/li0dCiJ4uBroMLKCRuG2liMBEImN4sq1vfnd4ITxag
3YRdNN3NBC7LSLBoslDDZbVZILqVJigL4umck6DdJND2+mB/foMyr/8CgYEA0+o+
H07naHkz546RbM4+xhSw0OiqIAI9NPTKilQEieEu9coKTQaWlgokL5vtn9v1iOzW
glI6m7xJ3EqaA/4C1FaU/5NRf+u8PjjLftC8W6mQJQxPaYp8O4DojXJ7GsbPPL8N
eCpcrddwymEcoEyqaI1n7T37uZ383oQxXRY472kCgYBrHakZojMjiGrZzTAv6zbD
bvK+4hE2lt7ugXRA2ibikvVelDi8O5LiVgJjy3uIfAOls0ltb31SD8mt80dVtn8O
wfaSPNcz4fR0TXLBK54N9EH0IAUH/nYErtQJ4uMRMCTB8MOz2aEgN2412n7DJ/to
wdhiHTrdZENhDngJTq19vQKBgE1Jm0qT3nYN7k4/gu5p8h8QIMSwdouiBWyw8sWs
itM2m4ftHgClrHogTL5IYkvxTwWXS0zQbKur6kw7sRxPofyLb2Ae/JRpB4ix8hXY
TzonB3lbmgLmDRUumKIt5JQJej+vdlwjJxYIAwzsiiI0k9n56LNv7cTux/DnyZEn
r7AxAoGBAOMTpMVT/zC6pskArQ9W9nHd3mgh5xxEQ7MQdIuZxZ7ILCUEW++J7eBP
46FRnqoruu+Ytfg//53hC+w3S7Y1SjGK0DEFkYRwYbmgS/KKSsco5tvRVTembSYL
JQnj8h8DPalW3Dn7oQXZhjCCeY7qK+z+KvgqDwTyv8HpP6Eetwhm
-----END RSA PRIVATE KEY-----`

        BeforeEach(func() {
            config := config.ServerConfig{CertificateFilePath:"blah",
                PrivateKeyFilePath: "blah", Debug: false}
            fakeLoader = new(fakes.FakeCertsLoader)
            generator = NewCertificateGenerator(config, fakeLoader)

            cpb, _ := pem.Decode([]byte(mockCertValue))
            kpb, _ := pem.Decode([]byte(mockKeyValue))
            crt, _ := x509.ParseCertificate(cpb.Bytes)
            key, _ := x509.ParsePKCS1PrivateKey(kpb.Bytes)

            fakeLoader.LoadCertsReturns(crt, key, nil)
        })

        Context("Generate", func() {
            It("generates a certificate", func() {
                certString, err := generator.Generate(map[string]interface{}{"common_name": "test"})
                Expect(err).To(BeNil())

                certificate, err := getCert(generator, certString.(string))

                Expect(err).To(BeNil())
                Expect(certificate).ToNot(BeNil())
            })

            It("sets common name and alternative name as passed in", func() {
                altNames := []string{"alt1", "alt2"}
                certString, _ := generator.Generate(map[string]interface{}{"common_name": "test", "alternative_names": altNames})
                certificate, _ := getCert(generator, certString.(string))

                Expect(certificate.DNSNames).Should(ContainElement("test"))
                Expect(certificate.DNSNames).Should(ContainElement("alt1"))
                Expect(certificate.DNSNames).Should(ContainElement("alt2"))
            })

            It("should set expiry for the cert in 1 year", func() {
                certString, _ := generator.Generate(map[string]interface{}{"common_name": "test"})
                certificate, _ := getCert(generator, certString.(string))

                oneYearFromToday := time.Now().UTC().Add(365*24*time.Hour)

                Expect(certificate.NotAfter).Should(BeTemporally("~", oneYearFromToday, 5*time.Second))
            })

            It("should be signed by the parent CA", func() {
                certString, err := generator.Generate(map[string]interface{}{"common_name": "test"})

                roots := x509.NewCertPool()
                success := roots.AppendCertsFromPEM([]byte(mockCertValue))
                Expect(success).To(BeTrue())

                block, _ := pem.Decode([]byte(certString.(string)))
                Expect(block).ToNot(BeNil())

                cert, err := x509.ParseCertificate(block.Bytes)
                Expect(err).To(BeNil())

                opts := x509.VerifyOptions{
                    Roots:   roots,
                }

                _, err = cert.Verify(opts)

                Expect(err).To(BeNil())
            })

            It("is not a CA", func() {
                certString, _ := generator.Generate(map[string]interface{}{"common_name": "test"})
                certificate, _ := getCert(generator, certString.(string))

                Expect(certificate.IsCA).To(BeFalse())
            })
        })
    })
})