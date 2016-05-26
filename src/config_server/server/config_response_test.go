package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"config_server/server"
)

var _ = Describe("Server/ConfigResponse", func() {
	Describe("Given a config response", func() {
		Context("with fields populated", func() {
			It("should generate a json string", func() {
				configResponse := server.ConfigResponse{Path:"key", Value:"result"}
				Expect(configResponse.Json()).To(Equal("{\"path\":\"key\",\"value\":\"result\"}"))
			})
		})
	})
})
