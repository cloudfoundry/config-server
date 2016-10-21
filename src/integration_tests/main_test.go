package acceptance_tests

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	. "integration_tests/support"
	"net/http"
)

var _ = Describe("Supported HTTP Methods", func() {

	BeforeEach(func() {
		SetupDB()
		StartServer()
	})

	AfterEach(func() {
		StopServer()
	})

	Describe("GET", func() {

		Describe("Lookup by key", func() {
			It("errors when key has invalid characters", func() {
				resp, err := SendGetRequestByKey("sm!urf/garg$amel/cat")

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(400))

				body, _ := ioutil.ReadAll(resp.Body)
				Expect(string(body)).To(ContainSubstring("Key must consist of alphanumeric, underscores, dashes, and forward slashes"))
			})

			Context("when key does not exist in server", func() {
				It("responds with status 404", func() {
					resp, err := SendGetRequestByKey("smurf")

					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(Equal(404))
				})
			})

			Context("when key exists in server", func() {
				It("responds with status 200", func() {
					SendPutRequest("smurf", "blue")

					resp, err := SendGetRequestByKey("smurf")

					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(Equal(200))
				})

				It("sends back value along with key as json", func() {
					SendPutRequest("smurf", "blue")

					resp, err := SendGetRequestByKey("smurf")

					Expect(err).To(BeNil())

					resultMap := UnmarshalJsonString(resp.Body)

					Expect(resultMap["path"]).To(Equal("smurf"))
					Expect(resultMap["value"]).To(Equal("blue"))
				})

				It("handles keys with forward slashes", func() {
					key := "smurf/gar_gamel/c-at"

					SendPutRequest(key, "vroom")

					resp, err := SendGetRequestByKey(key)

					Expect(err).To(BeNil())

					resultMap := UnmarshalJsonString(resp.Body)
					Expect(resultMap["path"]).To(Equal(key))
					Expect(resultMap["value"]).To(Equal("vroom"))
				})
			})
		})

		Describe("Lookup by ID", func() {
			Context("when id does not exist in server", func() {
				It("responds with status 404", func() {
					resp, err := SendGetRequestByID("123")

					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(Equal(404))
				})
			})

			Context("when id exists in server", func() {
				It("responds with status 200", func() {
					putResponse, _ := SendPutRequest("smurf", "blue")
					config := UnmarshalJsonString(putResponse.Body)
					id := config["id"].(string)

					resp, err := SendGetRequestByID(id)

					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(Equal(200))
				})

				It("sends back value along with key as json", func() {
					putResponse, _ := SendPutRequest("annie", "diane")
					config := UnmarshalJsonString(putResponse.Body)
					id := config["id"].(string)

					resp, err := SendGetRequestByID(id)

					Expect(err).To(BeNil())

					resultMap := UnmarshalJsonString(resp.Body)

					Expect(resultMap["path"]).To(Equal("annie"))
					Expect(resultMap["value"]).To(Equal("diane"))
				})
			})
		})
	})

	Describe("PUT", func() {
		It("fails if content-type in the header is not set to application/json", func() {
			requestBytes := bytes.NewReader([]byte(`{"value":"smurf"`))
			req, _ := http.NewRequest("PUT", SERVER_URL+"/v1/data/blah", requestBytes)
			req.Header.Add("Authorization", "bearer "+ValidToken())

			resp, err := HTTPSClient.Do(req)
			Expect(resp.StatusCode).To(Equal(415))
			Expect(err).To(BeNil())

			body, _ := ioutil.ReadAll(resp.Body)
			Expect(string(body)).To(ContainSubstring("Unsupported Media Type - Accepts application/json only"))
		})

		It("errors when key has invalid characters", func() {
			resp, err := SendPutRequest("sm!urf/garg$amel/cat", "value")

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(400))

			body, _ := ioutil.ReadAll(resp.Body)
			Expect(string(body)).To(ContainSubstring("Key must consist of alphanumeric, underscores, dashes, and forward slashes"))
		})

		Context("when key does NOT exist in server", func() {
			It("responds with value & id", func() {
				resp, err := SendPutRequest("cross", "fit")

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))

				resultMap := UnmarshalJsonString(resp.Body)
				Expect(resultMap["path"]).To(Equal("cross"))
				Expect(resultMap["value"]).To(Equal("fit"))
			})

			It("responds with status 200 when value is successfully stored", func() {
				resp, err := SendPutRequest("smurf", "blue")

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
			})

			Context("when key is empty string", func() {
				It("is stored and responds with value & id", func() {
					resp, err := SendPutRequest("crossfit", "")

					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(Equal(200))

					resultMap := UnmarshalJsonString(resp.Body)
					Expect(resultMap["path"]).To(Equal("crossfit"))
					Expect(resultMap["value"]).To(Equal(""))
				})
			})
			Context("when key is nil", func() {
				It("is stored and responds with value & id", func() {
					resp, err := SendPutRequest("crossfit", nil)

					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(Equal(200))

					resultMap := UnmarshalJsonString(resp.Body)
					Expect(resultMap["path"]).To(Equal("crossfit"))
					Expect(resultMap["value"]).To(BeNil())
				})
			})
		})

		Context("when key exists in server", func() {
			It("updates the value", func() {
				SendPutRequest("smurf", "blue")

				getResp, _ := SendGetRequestByKey("smurf")

				resultMap := UnmarshalJsonString(getResp.Body)
				Expect(resultMap["path"]).To(Equal("smurf"))
				Expect(resultMap["value"]).To(Equal("blue"))

				SendPutRequest("smurf", "red")
				getResp, _ = SendGetRequestByKey("smurf")

				resultMap = UnmarshalJsonString(getResp.Body)
				Expect(resultMap["path"]).To(Equal("smurf"))
				Expect(resultMap["value"]).To(Equal("red"))
			})
		})
	})

	Describe("POST", func() {
		It("fails if content-type in the header is not set to application/json", func() {
			requestBytes := bytes.NewReader([]byte(`{"type":"password","parameters":{}}`))
			req, _ := http.NewRequest("POST", SERVER_URL+"/v1/data/blah", requestBytes)
			req.Header.Add("Authorization", "bearer "+ValidToken())

			resp, err := HTTPSClient.Do(req)
			Expect(resp.StatusCode).To(Equal(415))
			Expect(err).To(BeNil())

			body, _ := ioutil.ReadAll(resp.Body)
			Expect(string(body)).To(ContainSubstring("Unsupported Media Type - Accepts application/json only"))
		})

		It("generates a new id and password for a new key", func() {
			resp, _ := SendPostRequest("password-key", "password")
			result := UnmarshalJsonString(resp.Body)

			Expect(result["id"]).ToNot(BeNil())
			Expect(result["path"]).To(Equal("password-key"))
			Expect(result["value"]).To(MatchRegexp("[a-z0-9]{20}"))
		})

		It("generates a new id and certificate for a new key", func() {
			resp, _ := SendPostRequest("certificate-key", "certificate")
			result := UnmarshalJsonString(resp.Body)

			Expect(result["id"]).ToNot(BeNil())
			Expect(result["path"]).To(Equal("certificate-key"))

			value := result["value"].(map[string]interface{})
			cert, _ := ParseCertString(value["certificate"].(string))

			Expect(cert.DNSNames).Should(ContainElement("cnj"))
			Expect(cert.DNSNames).Should(ContainElement("deadlift"))
			Expect(cert.Subject.CommonName).To(Equal("burpees"))

			Expect(cert.Issuer.Organization).To(ContainElement("Internet Widgits Pty Ltd"))
			Expect(cert.Issuer.Country).To(ContainElement("AU"))
			Expect(cert.Issuer.Province).To(ContainElement("Some-State"))
		})
	})
})
