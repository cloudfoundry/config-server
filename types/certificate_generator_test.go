package types_test

import (
	"crypto/x509"
	"encoding/pem"
	"strings"
	"time"

	. "github.com/cloudfoundry/config-server/types"

	"github.com/cloudfoundry/config-server/types/typesfakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func parseCertString(certString string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certString))
	crt, err := x509.ParseCertificate(block.Bytes)

	return crt, err
}

func getCertResp(generator ValueGenerator, certParams map[interface{}]interface{}) CertResponse {
	certResp, err := generator.Generate(certParams)
	Expect(err).To(BeNil())

	return certResp.(CertResponse)
}

var _ = Describe("CertificateGenerator", func() {

	var (
		fakeLoader *typesfakes.FakeCertsLoader
		generator  ValueGenerator
		fakeRootCA *x509.Certificate
	)

	mockCertValue := `-----BEGIN CERTIFICATE-----
MIIDpTCCAo2gAwIBAgIRAOQkH+coubqiMByb4L5w+JcwDQYJKoZIhvcNAQELBQAw
bDELMAkGA1UEBhMCTkExDzANBgNVBAgTBk5hcm5pYTEUMBIGA1UEBxMLU3ByaW5n
ZmllbGQxFjAUBgNVBAoTDUZ1dHVyYW1hIENvcnAxHjAcBgNVBAMTFUZ1dHVyYW1h
IENvcnAgUm9vdCBDQTAeFw0xNzAyMDIxNTM0NDdaFw0zNzAxMjgxNTM0NDdaMGwx
CzAJBgNVBAYTAk5BMQ8wDQYDVQQIEwZOYXJuaWExFDASBgNVBAcTC1NwcmluZ2Zp
ZWxkMRYwFAYDVQQKEw1GdXR1cmFtYSBDb3JwMR4wHAYDVQQDExVGdXR1cmFtYSBD
b3JwIFJvb3QgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDEGay0
G0tFjbAMtI3UA1yZdKP5QWpwYE5ABHwgBrTQyOCvh6IkRiOO10uF5MXzYmbxz+NJ
QfBqO9riRG5TsVOrVUnSjkxc2WOtqxZVHJLmXWT9880UaiBzPDaDOKpB+cYaI8a6
JyxTJJj74nH77MLGuaMd0xEQnqSFTCZ3o1iGMHUMg21XMTimS5zYO2HIJOA7duYC
q15JrGRzIYtG/AdjBGNBW3UnY1rH9CUWn79JgMHVhAAO9T1HJOCQMWhIS1TABxRB
CS6Bb/EN+CJcYLEkAIfppiHAPBO+zATKSRIDsU48cdDi23rjQJCnT3sfpyLtPjea
vtbf1beX+nkXALd1AgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8E
BTADAQH/MB0GA1UdDgQWBBQr/j++cO+ZuBJeac5qgOyVTXViGjANBgkqhkiG9w0B
AQsFAAOCAQEAeoS9zrdLmxbDr9bi6YoGPhNcBK+xjCWfvoDpG/+3TJ6MeuNFkfEN
mddi80Dk+DZUln15ld38XNONETvT/CIgnynQbQV2nMQVAdDQ/e/FnV4QS1T6AuZc
CJvQJeFZecT8f7Tu/a4afB/dBWPA8tS4O4UmYVq6BWMADaAI07lgL+rvZgODJKm2
uSkkQiRM6sYKwl3Hj0PfBhxsU/cNKKwfCAGpduZbhjjr4c0llLoMQL1Hg70U5Kh/
Zu6OJgbh0gkhgmQh3cDqpqGnuU0A2fP07U/jScHxi3uF5Gogf0JICdMdsneNmyEd
OfmZ+QiPmUcWxip9HEiV95FKdMMKrZeBzQ==
-----END CERTIFICATE-----`

	mockKeyValue := `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAxBmstBtLRY2wDLSN1ANcmXSj+UFqcGBOQAR8IAa00Mjgr4ei
JEYjjtdLheTF82Jm8c/jSUHwajva4kRuU7FTq1VJ0o5MXNljrasWVRyS5l1k/fPN
FGogczw2gziqQfnGGiPGuicsUySY++Jx++zCxrmjHdMREJ6khUwmd6NYhjB1DINt
VzE4pkuc2DthyCTgO3bmAqteSaxkcyGLRvwHYwRjQVt1J2Nax/QlFp+/SYDB1YQA
DvU9RyTgkDFoSEtUwAcUQQkugW/xDfgiXGCxJACH6aYhwDwTvswEykkSA7FOPHHQ
4tt640CQp097H6ci7T43mr7W39W3l/p5FwC3dQIDAQABAoIBAFbrrnJyuq2MvSEU
kt0L1GqsmdXs2foPpf+YPVCQyGrW70i+jO8ZK7+vzGj/24Ii9tBuPieFk8Py3fV3
MZPlQNnrxvoOWui47wbWk+KL9M0pTo+GFjvbQqDxapRFdKojxVxmyTpQDMSZrznO
yIlLC5a8SAH7QVAlkXUIKPDUtXNclRpXP3td4C8oheTBOy/n3GcNwoI3+m6RtzEh
lfbaC0EVftDt9acTeUz69l7CEdmQkL5wlbWzUsEGDu9Jec8FByjcN51dT0Q3W1cj
EWEQft9J1mXk0Q4MQp9OEpPBpItv/vR8kZmx/mlgsCA1wdZRb1Z+baMvqyLrDnhe
sZODTs0CgYEAxljIqjYCHbrHdDV/bZNSU3gvBtp0fTx9qpk65l1PMe3ZT4bvNgzQ
UEVzSSKpq5chL0s9vgVDUOTc8ROwLR6nIJGhUcrfYup1lUDi++m/bc705BaX8C+j
e4s3n0v7eqN4Nmb6jzjThBdzQSUYHktTFx1T+zENccYQBPtRDzCgWc8CgYEA/Rm5
fqiLOmzRlizPdzt3gNBE7F4VA9bVVQSFZ3C9ChhEXI+37wFBH6yC20khzxoamf7h
klPsjx24SfaXn5H1cAv/i7iGBfyivyatfZ8t4QLIESO1FIByzL8ho3l9SrImXxkN
+rF2utlPnFRvcjnZITOCYtJt3ELCJs24eZxyn3sCgYBpp5WMhaRu7bWdQ3oThmxO
JjD2t5thsr1GCMWFPKMY95fQcxItennkqHoWtS4oRyLYLH5BIFCRYLgIevJXtoJU
KP8DsMt5x2bHEH9YrVTZS2rLrPVWbinpf2kro6/bzgQVBpnlfOG+9Tbhtr64NGuY
XnkDz0dYGaci4DR3oPFppQKBgCJXJy/kNl/K+/TgR5Xp36D+oRtg+ID42SEb5+3N
Ahkoib31UgN/rBJcGbUfCxNKe2mBh6GO+X3BjouahaAE8cQuHJIcizKswHgnC0sM
GfF5qyOIML2DYBAfrrS7eyBzY3lrsNXe1jkr2v1fB2/8IVW742j9HCLSe/0/1FPL
wlGZAoGAY+huOY40xWTAS4XqZVF9a22zbeCeeJxXSH/CGUrw11+xg/2h8hKEJp/e
+72iuoCAThbJQDGagAm1/W3sulGFmfP30YJuTOkwfipIQZsq44b3hPTDbONVeEh2
Iv9rGAlhC9tfzJgH11fZqwYJHCR58VmCV4sLvUSisqzdd3j0RGE=
-----END RSA PRIVATE KEY-----`

	BeforeEach(func() {
		fakeLoader = new(typesfakes.FakeCertsLoader)
		generator = NewCertificateGenerator(fakeLoader)

		cpb, _ := pem.Decode([]byte(mockCertValue))
		kpb, _ := pem.Decode([]byte(mockKeyValue))
		fakeRootCA, _ = x509.ParseCertificate(cpb.Bytes) //nolint:errcheck
		key, _ := x509.ParsePKCS1PrivateKey(kpb.Bytes)   //nolint:errcheck

		fakeLoader.LoadCertsReturns(fakeRootCA, key, nil)
	})

	Describe("Generate", func() {
		var params map[interface{}]interface{}
		BeforeEach(func() {
			params = map[interface{}]interface{}{"common_name": "bosh.io"}
		})

		Context("when passed parameters types are NOT correct", func() {
			It("returns an error when CommonName is not of type string", func() {
				params["common_name"] = []int{1}
				_, err := generator.Generate(params)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Failed to generate certificate, parameters are invalid: Expected input to be deserializable: yaml: unmarshal errors:\n  line 2: cannot unmarshal !!seq into string"))
			})

			It("returns an error when AlternativeName is not of type []string", func() {
				params["alternative_names"] = "smurf"
				_, err := generator.Generate(params)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Failed to generate certificate, parameters are invalid: Expected input to be deserializable: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `smurf` into []string"))
			})

			It("returns an error when ca is not of type string", func() {
				params["ca"] = []int{1}
				_, err := generator.Generate(params)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Failed to generate certificate, parameters are invalid: Expected input to be deserializable: yaml: unmarshal errors:\n  line 2: cannot unmarshal !!seq into string"))
			})
		})

		Context("when passed parameters types are correct", func() {
			var params map[interface{}]interface{}
			BeforeEach(func() {
				params = map[interface{}]interface{}{}
			})

			Context("when 'is_ca' is TRUE", func() {
				BeforeEach(func() {
					params["is_ca"] = true
				})

				Context("when 'ca' is NOT set", func() {
					var certificate *x509.Certificate
					var certResp CertResponse

					BeforeEach(func() {
						certResp = getCertResp(generator, params)
						certificate, _ = parseCertString(certResp.Certificate) //nolint:errcheck
					})

					It("set the CA field to itself", func() {
						Expect(strings.Trim(certResp.CA, "\n")).To(Equal(strings.Trim(certResp.Certificate, "\n")))
					})

					It("generates a root CA", func() {
						Expect(certificate.IsCA).To(BeTrue())
					})

					It("sets KeyUsage and ExtKeyUsage", func() {
						Expect(certificate.KeyUsage).To(Equal(x509.KeyUsageCertSign | x509.KeyUsageCRLSign))
						Expect(certificate.ExtKeyUsage).To(BeEmpty())
					})

					It("sets Issuer, Country & default Org", func() {
						Expect(certificate.Issuer.Country).To(Equal([]string{"US"}))
						Expect(certificate.Issuer.Organization).To(Equal([]string{"Cloud Foundry"}))
						Expect(certificate.Issuer.CommonName).To(Equal(""))
					})

					It("sets the SKI and AKI", func() {
						Expect(certificate.SubjectKeyId).ToNot(BeNil())
						Expect(certificate.SubjectKeyId).To(Equal(certificate.AuthorityKeyId))
					})

					Context("when organization set", func() {
						BeforeEach(func() {
							params = map[interface{}]interface{}{
								"is_ca":        true,
								"organization": "Hi Five BOSH",
							}
							certResp = getCertResp(generator, params)
							certificate, _ = parseCertString(certResp.Certificate) //nolint:errcheck
						})

						It("sets Issuer, Country & default Org", func() {
							Expect(certificate.Issuer.Country).To(Equal([]string{"US"}))
							Expect(certificate.Issuer.Organization).To(Equal([]string{"Hi Five BOSH"}))
							Expect(certificate.Issuer.CommonName).To(Equal(""))
						})
					})
				})

				Context("when 'ca' is NOT empty", func() {
					var certificate *x509.Certificate
					var certResp CertResponse
					BeforeEach(func() {
						params["ca"] = "smurf-cert"

						certResp = getCertResp(generator, params)
						certificate, _ = parseCertString(certResp.Certificate) //nolint:errcheck
					})

					It("generates an intermediate CA cert", func() {
						Expect(certificate.IsCA).To(BeTrue())
					})

					It("set the CA field to the signing CA", func() {
						Expect(strings.Trim(certResp.CA, "\n")).To(Equal(mockCertValue))
					})

					It("sets KeyUsage and ExtKeyUsage", func() {
						Expect(certificate.KeyUsage).To(Equal(x509.KeyUsageCertSign | x509.KeyUsageCRLSign))
						Expect(certificate.ExtKeyUsage).To(BeEmpty())
					})

					It("sets Issuer Country & Org", func() {
						Expect(certificate.Issuer.Country).To(Equal([]string{"NA"}))
						Expect(certificate.Issuer.Organization).To(Equal([]string{"Futurama Corp"}))
						Expect(certificate.Issuer.CommonName).To(Equal("Futurama Corp Root CA"))
					})

					It("should be signed by the root CA", func() {
						certString := certResp.Certificate

						roots := x509.NewCertPool()
						success := roots.AppendCertsFromPEM([]byte(mockCertValue))
						Expect(success).To(BeTrue())

						block, _ := pem.Decode([]byte(certString))
						Expect(block).ToNot(BeNil())

						cert, err := x509.ParseCertificate(block.Bytes)
						Expect(err).To(BeNil())

						opts := x509.VerifyOptions{
							Roots: roots,
						}

						_, err = cert.Verify(opts)

						Expect(err).To(BeNil())
					})

					It("should set the AKI from the root CA", func() {
						Expect(certificate.AuthorityKeyId).ToNot(BeNil())
						Expect(certificate.SubjectKeyId).ToNot(BeNil())
						Expect(certificate.AuthorityKeyId).To(Equal(fakeRootCA.SubjectKeyId))
					})
				})
			})

			Context("when 'is_ca' is FALSE", func() {

				Context("when 'ca' is empty", func() {

					It("should throw an error", func() {
						_, err := generator.Generate(params)
						Expect(err).ToNot(BeNil())
						Expect(err.Error()).To(Equal("Missing required CA name"))
					})
				})

				Context("when 'ca' is NOT empty", func() {
					BeforeEach(func() {
						params["ca"] = "smurf-ca"
						params["common_name"] = "bosh.io"
					})
					It("generates a certificate", func() {
						certResp := getCertResp(generator, params)
						certificate, err := parseCertString(certResp.Certificate)

						Expect(err).To(BeNil())
						Expect(certificate).ToNot(BeNil())
					})

					It("sets KeyUsage", func() {
						altNames := []interface{}{"cloudfoundry.com", "example.com"}
						params["alternative_names"] = altNames
						certResp := getCertResp(generator, params)
						certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

						Expect(certificate.KeyUsage).To(Equal(x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature))
					})

					It("sets common name and alternative name as passed in", func() {
						altNames := []interface{}{"cloudfoundry.com", "example.com"}
						params["alternative_names"] = altNames
						certResp := getCertResp(generator, params)
						certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

						Expect(certificate.Subject.CommonName).Should(Equal("bosh.io"))

						Expect(certificate.DNSNames).ShouldNot(ContainElement("bosh.io"))
						Expect(certificate.DNSNames).Should(ContainElement("cloudfoundry.com"))
						Expect(certificate.DNSNames).Should(ContainElement("example.com"))
					})

					It("should work if CN was also included in SAN", func() {
						altNames := []interface{}{"bosh.io", "cloudfoundry.com", "example.com"}
						params["alternative_names"] = altNames
						certResp := getCertResp(generator, params)
						certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

						Expect(certificate.Subject.CommonName).Should(Equal("bosh.io"))

						Expect(certificate.DNSNames).Should(ContainElement("bosh.io"))
						Expect(certificate.DNSNames).Should(ContainElement("cloudfoundry.com"))
						Expect(certificate.DNSNames).Should(ContainElement("example.com"))
					})

					It("should set expiry for the cert in 1 year", func() {
						certResp := getCertResp(generator, params)
						certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

						oneYearFromToday := time.Now().UTC().Add(365 * 24 * time.Hour)

						Expect(certificate.NotAfter).Should(BeTemporally("~", oneYearFromToday, 5*time.Second))
					})

					It("should be signed by the parent CA", func() {
						certResp := getCertResp(generator, params)
						certString := certResp.Certificate

						roots := x509.NewCertPool()
						success := roots.AppendCertsFromPEM([]byte(mockCertValue))
						Expect(success).To(BeTrue())

						block, _ := pem.Decode([]byte(certString))
						Expect(block).ToNot(BeNil())

						cert, err := x509.ParseCertificate(block.Bytes)
						Expect(err).To(BeNil())

						opts := x509.VerifyOptions{
							Roots: roots,
						}

						_, err = cert.Verify(opts)

						Expect(err).To(BeNil())
					})

					It("is not a CA", func() {
						certResp := getCertResp(generator, params)
						certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

						Expect(certificate.IsCA).To(BeFalse())
					})

					It("generates a 3072-bit private key by default", func() {
						certResp := getCertResp(generator, params)

						Expect(certResp.PrivateKey).NotTo(BeEmpty())

						block, _ := pem.Decode([]byte(certResp.PrivateKey))
						key, _ := x509.ParsePKCS1PrivateKey(block.Bytes) //nolint:errcheck

						Expect(key.PublicKey.N.BitLen()).To(Equal(3072))
					})

					It("generates a private key with custom key_length when specified", func() {
						params["key_length"] = 2048
						certResp := getCertResp(generator, params)

						Expect(certResp.PrivateKey).NotTo(BeEmpty())

						block, _ := pem.Decode([]byte(certResp.PrivateKey))
						key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
						Expect(err).To(BeNil())
						Expect(key.PublicKey.N.BitLen()).To(Equal(2048))
					})

					It("returns an error when key_length is not a standard size", func() {
						params["key_length"] = 1234
						_, err := generator.Generate(params)

						Expect(err).ToNot(BeNil())
						Expect(err.Error()).To(ContainSubstring("Invalid key_length: 1234"))
						Expect(err.Error()).To(ContainSubstring("Must be one of: [2048 3072 4096]"))
					})

					It("should have the public keys of the private key and certificate match", func() {
						certResp := getCertResp(generator, params)
						certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

						block, _ := pem.Decode([]byte(certResp.PrivateKey))
						key, _ := x509.ParsePKCS1PrivateKey(block.Bytes) //nolint:errcheck
						Expect(certificate.PublicKey).To(Equal(&key.PublicKey))
					})

					It("set the CA field to the signing CA", func() {
						certResp := getCertResp(generator, params)
						Expect(strings.Trim(certResp.CA, "\n")).To(Equal(mockCertValue))
					})

					It("set the AKI as the CA's SKI", func() {
						certResp := getCertResp(generator, params)
						certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck
						Expect(certificate.SubjectKeyId).ToNot(BeNil())
						Expect(certificate.AuthorityKeyId).To(Equal(fakeRootCA.SubjectKeyId))
					})

					Context("when ExtKeyUsage is NOT empty", func() {
						Context("when it is client_auth", func() {
							It("should include the x509.ExtKeyUsageClientAuth flag in the key", func() {
								altNames := []interface{}{"cloudfoundry.com", "example.com"}
								params["alternative_names"] = altNames
								params["extended_key_usage"] = []string{"client_auth"}
								certResp := getCertResp(generator, params)
								certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

								Expect(certificate.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}))
							})
						})

						Context("when it is server_auth", func() {
							It("should include the x509.ExtKeyUsageServerAuth flag in the key", func() {
								altNames := []interface{}{"cloudfoundry.com", "example.com"}
								params["alternative_names"] = altNames
								params["extended_key_usage"] = []string{"server_auth"}
								certResp := getCertResp(generator, params)
								certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

								Expect(certificate.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}))
							})
						})

						Context("when multiple auth types are set", func() {
							It("should include the x509.ExtKeyUsageServerAuth flag in the key", func() {
								altNames := []interface{}{"cloudfoundry.com", "example.com"}
								params["alternative_names"] = altNames
								params["extended_key_usage"] = []string{"client_auth", "server_auth"}
								certResp := getCertResp(generator, params)
								certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

								Expect(certificate.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}))
							})
						})

						Context("when it is neither server or client auth", func() {
							It("returns an error", func() {
								altNames := []interface{}{"cloudfoundry.com", "example.com"}
								params["alternative_names"] = altNames
								params["extended_key_usage"] = []string{"something not supported"}
								_, err := generator.Generate(params)

								Expect(err).ToNot(BeNil())
								Expect(err.Error()).To(Equal("Unsupported extended key usage value: something not supported"))
							})
						})
					})

					Context("when ExtKeyUsage is empty", func() {
						It("should include the x509.ExtKeyUsageServerAuth flag in the key", func() {
							altNames := []interface{}{"cloudfoundry.com", "example.com"}
							params["alternative_names"] = altNames
							certResp := getCertResp(generator, params)
							certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

							Expect(certificate.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}))
						})
					})

					Context("when 'duration' is NOT empty", func() {
						It("should set expiry for the cert in more than year", func() {
							params["duration"] = 365
							params["common_name"] = "bosh.io"
							certResp := getCertResp(generator, params)
							certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

							tenYearsFromToday := time.Now().UTC().Add(365 * 24 * time.Hour)

							Expect(certificate.NotAfter).Should(BeTemporally("~", tenYearsFromToday, 5*time.Second))
						})

						It("should set expiry for the cert correctly when duration is very large", func() {
							params["duration"] = 7000
							params["common_name"] = "bosh.io"
							certResp := getCertResp(generator, params)
							certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

							longTimeFromToday := time.Now().UTC().Add(7000 * 24 * time.Hour)

							Expect(certificate.NotAfter).Should(BeTemporally("~", longTimeFromToday, 5*time.Second))
						})

						It("should set expiry for the cert correctly when duration is invalid", func() {
							params["duration"] = -100
							params["common_name"] = "bosh.io"
							certResp := getCertResp(generator, params)
							certificate, _ := parseCertString(certResp.Certificate) //nolint:errcheck

							oneYearFromToday := time.Now().UTC().Add(365 * 24 * time.Hour)

							Expect(certificate.NotAfter).Should(BeTemporally("~", oneYearFromToday, 5*time.Second))
						})
					})
				})
			})
		})

		Context("when passed parameters use unsupported keys", func() {
			var params map[interface{}]interface{}
			BeforeEach(func() {
				params = map[interface{}]interface{}{
					"is_ca":              true,
					"extended_key_usage": []string{"random", "values"},
					"ext_key_usage":      []string{"random", "values"},
				}
			})

			It("returns an error", func() {
				_, err := generator.Generate(params)

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Failed to generate certificate, parameters are invalid: Unsupported parameter 'ext_key_usage'"))
			})
		})
	})
})
