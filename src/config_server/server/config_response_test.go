package server_test

import (
	"config_server/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server/ConfigResponse", func() {

	Describe("Given a config response", func() {

		Context("with fields populated", func() {
			It("should generate a json string", func() {
				configResponse := server.ConfigResponse{Path: "key", Value: "result"}
				Expect(configResponse.Json()).To(Equal("{\"path\":\"key\",\"value\":\"result\"}"))
			})
		})

		Context("without fields populated", func() {
			It("should generate a json string with empty values", func() {
				configResponse := server.ConfigResponse{}
				Expect(configResponse.Json()).To(Equal("{\"path\":\"\",\"value\":\"\"}"))
			})
		})
	})
})
