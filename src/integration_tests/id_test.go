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
			response, _ := SendPutRequest("Dale", "Wick")
			resultMap := UnmarshalJsonString(response.Body)

			Expect(resultMap["id"]).ToNot(BeNil())
			Expect(len((resultMap["id"].(string))) > 0).To(BeTrue())
			Expect(resultMap["name"]).To(Equal("Dale"))
			Expect(resultMap["value"]).To(Equal("Wick"))
		})

		It("generates new id for different names", func() {
			response1, _ := SendPutRequest("Dale", "Wick")
			resultMap1 := UnmarshalJsonString(response1.Body)
			Expect(resultMap1["id"]).ToNot(BeNil())
			Expect(len((resultMap1["id"].(string))) > 0).To(BeTrue())

			response2, _ := SendPutRequest("Alan", "Donovan")
			resultMap2 := UnmarshalJsonString(response2.Body)
			Expect(resultMap2["id"]).ToNot(BeNil())
			Expect(len((resultMap2["id"].(string))) > 0).To(BeTrue())

			Expect(resultMap1["id"]).ToNot(Equal(resultMap2["id"]))
		})

		It("generates new id for existing name", func() {
			response1, _ := SendPutRequest("Dale", "Wick")
			resultMap1 := UnmarshalJsonString(response1.Body)

			response2, _ := SendPutRequest("Dale", "Wick")
			resultMap2 := UnmarshalJsonString(response2.Body)

			Expect(resultMap1["id"]).ToNot(Equal(resultMap2["id"]))
			Expect(resultMap1["name"]).To(Equal(resultMap2["name"]))
			Expect(resultMap1["value"]).To(Equal(resultMap2["value"]))
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
