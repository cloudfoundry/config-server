package store_test

import (
	. "config_server/store"

	. "github.com/onsi/ginkgo"
	//	. "github.com/onsi/gomega"
)

var _ = Describe("DatabaseStore", func() {

	Describe("Given a properly initialized DatabaseStore", func() {
		var store DatabaseStore

		BeforeEach(func() {
			store = NewDatabaseStore()
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
