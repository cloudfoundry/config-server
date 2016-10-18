package acceptance_tests

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "integration_tests/support"
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
		It("errors when key has invalid characters", func() {
			resp, err := SendGetRequest("sm!urf/garg$amel/cat")

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(400))

			body, _ := ioutil.ReadAll(resp.Body)
			Expect(string(body)).To(ContainSubstring("Key must consist of alphanumeric, underscores, dashes, and forward slashes"))
		})

		Context("when key does not exist in server", func() {
			It("responds with status 404", func() {
				resp, err := SendGetRequest("smurf")

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(404))
			})
		})

		Context("when key exists in server", func() {
			It("responds with status 200", func() {
				SendPutRequest("smurf", "blue")

				resp, err := SendGetRequest("smurf")

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
			})

			It("sends back value along with key as json", func() {
				SendPutRequest("smurf", "blue")

				resp, err := SendGetRequest("smurf")

				Expect(err).To(BeNil())

				resultMap := UnmarshalJsonString(resp.Body)

				Expect(resultMap["path"]).To(Equal("smurf"))
				Expect(resultMap["value"]).To(Equal("blue"))
			})

			It("handles keys with forward slashes", func() {
				key := "smurf/gar_gamel/c-at"

				SendPutRequest(key, "vroom")

				resp, err := SendGetRequest(key)

				Expect(err).To(BeNil())

				resultMap := UnmarshalJsonString(resp.Body)
				Expect(resultMap["path"]).To(Equal(key))
				Expect(resultMap["value"]).To(Equal("vroom"))
			})
		})
	})

	Describe("PUT", func() {
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

			Context("when key is empty string", func(){
				It("is stored and responds with value & id", func() {
					resp, err := SendPutRequest("crossfit", "")

					Expect(err).To(BeNil())
					Expect(resp.StatusCode).To(Equal(200))

					resultMap := UnmarshalJsonString(resp.Body)
					Expect(resultMap["path"]).To(Equal("crossfit"))
					Expect(resultMap["value"]).To(Equal(""))
				})
			})
			Context("when key is nil", func(){
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

				getResp, _ := SendGetRequest("smurf")

				resultMap := UnmarshalJsonString(getResp.Body)
				Expect(resultMap["path"]).To(Equal("smurf"))
				Expect(resultMap["value"]).To(Equal("blue"))

				SendPutRequest("smurf", "red")
				getResp, _ = SendGetRequest("smurf")

				resultMap = UnmarshalJsonString(getResp.Body)
				Expect(resultMap["path"]).To(Equal("smurf"))
				Expect(resultMap["value"]).To(Equal("red"))
			})
		})
	})

	Describe("POST", func() {
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
