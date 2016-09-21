package acceptance_tests

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GET Requests", func() {
	var client *http.Client

	BeforeEach(func() {
		client = HTTPSClient()
	})

	It("should get status 404 when key on the server does not exist", func() {
		req, err := http.NewRequest("GET", "https://localhost:9000/v1/data/key", nil)
		req.Header.Add("Authorization", "bearer " + ValidToken())
		resp, err := client.Do(req)

		Expect(err).To(BeNil())
		Expect(resp.StatusCode).To(Equal(404))
	})
})
