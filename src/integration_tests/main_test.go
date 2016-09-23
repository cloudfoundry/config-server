package acceptance_tests

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "integration_tests/support"
)

var _ = Describe("Supported HTTP Methods", func() {

	BeforeEach(func() {
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

	Describe("Put", func() {
		It("errors when key has invalid characters", func() {
			resp, err := SendPutRequest("sm!urf/garg$amel/cat", "value")

			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(400))

			body, _ := ioutil.ReadAll(resp.Body)
			Expect(string(body)).To(ContainSubstring("Key must consist of alphanumeric, underscores, dashes, and forward slashes"))
		})

		Context("when key does NOT exist in server", func() {
			It("responds with status 204 when value is successfully stored", func() {
				resp, err := SendPutRequest("smurf", "blue")

				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(204))
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
})
