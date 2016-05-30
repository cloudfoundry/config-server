package store_test

import (
	. "config_server/store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MemoryStore", func() {

	Describe("Given a properly initialized MemoryStore", func() {
		var store MemoryStore

		BeforeEach(func() {
			store = NewMemoryStore()
		})

		Context("calling put with valid data", func() {
			It("should not return error", func() {
				err := store.Put("key", "value")
				Expect(err).To(BeNil())
			})
		})

		Context("calling get with valid key", func() {
			It("should return asoociated value", func() {
				store.Put("key", "value")
				returnedValue, err := store.Get("key")
				Expect(err).To(BeNil())
				Expect(returnedValue).To(Equal("value"))
			})
		})
	})
})
