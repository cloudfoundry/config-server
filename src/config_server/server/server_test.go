package server_test

import (
	. "config_server/server"

	"config_server/server"
	"config_server/store"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	Describe("Given a server with no datastore", func() {
		var configServer ConfigServer
		BeforeEach(func() {
			configServer = server.NewServer(nil)
		})
		Context("starting the server", func() {
			It("should return an error", func() {
				err := configServer.Start(9000)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("DataStore can not be nil"))
			})
		})
	})

	Describe("Given a server with datastore", func() {
		var configServer ConfigServer
		BeforeEach(func() {
			configServer = server.NewServer(store.NewMemoryStore())
		})

		Context("starting the server", func() {
			It("should not return an error", func() {
				//				err := configServer.Start(9000)
				//				Expect(err).To(BeNil())
			})
		})
	})
})
