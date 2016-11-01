package acceptance_tests

import (
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

	Describe("PUT", func() {
		It("generates new id for a name when adding new value", func() {
			SendPutRequest("Dale", "Wick")
			resp, _ := SendGetRequestByName("Dale")

			resultMap := UnmarshalJsonString(resp.Body)
			Expect(resultMap["id"]).ToNot(BeNil())
			Expect(len((resultMap["id"].(string))) > 0).To(BeTrue())
			Expect(resultMap["value"]).To(Equal("Wick"))
		})

		It("generates new id for different names", func() {
			SendPutRequest("Dale", "Wick")
			resp1, _ := SendGetRequestByName("Dale")
			resultMap1 := UnmarshalJsonString(resp1.Body)
			Expect(resultMap1["id"]).ToNot(BeNil())
			Expect(len((resultMap1["id"].(string))) > 0).To(BeTrue())

			SendPutRequest("Alan", "Donovan")
			resp2, _ := SendGetRequestByName("Alan")
			resultMap2 := UnmarshalJsonString(resp2.Body)
			Expect(resultMap2["id"]).ToNot(BeNil())
			Expect(len((resultMap2["id"].(string))) > 0).To(BeTrue())

			Expect(resultMap1["id"]).ToNot(Equal(resultMap2["id"]))
		})
	})

	Describe("POST", func() {
		It("generates a new id and password for a new name", func() {
			resp, _ := SendPostRequest("pass", "password")
			result := UnmarshalJsonString(resp.Body)

			Expect(result["id"]).ToNot(BeNil())
			Expect(result["value"]).ToNot(BeNil())
		})

		It("generates a new id and certificate for a new name", func() {
			resp, _ := SendPostRequest("cert", "certificate")
			result := UnmarshalJsonString(resp.Body)

			Expect(result["id"]).ToNot(BeNil())
			Expect(result["value"]).ToNot(BeNil())
		})
	})

})
