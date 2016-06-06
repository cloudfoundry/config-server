package store_test

import (
	. "github.com/onsi/ginkgo"
	"config_server/config"
)

var _ = Describe("DatabaseStore", func() {

	Describe("Given a properly initialized DatabaseStore", func() {
		var store DatabaseStore

		BeforeEach(func() {
			store = NewDatabaseStore(config.DBConfig{})
		})

		Context("calling put with valid data", func() {
			It("should not return error", func() {
			})
		})

		Context("calling get with valid key", func() {
			It("should return associated value", func() {
			})
		})
	})
})
