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

		Context("Put", func() {
			It("should not return error when adding a string type value", func() {
				err := store.Put("key", "value")
				Expect(err).To(BeNil())
			})
		})

		Context("Get", func() {
			It("should return associated value", func() {
				store.Put("key", "value")
				returnedValue, err := store.Get("key")
				Expect(err).To(BeNil())
				Expect(returnedValue).To(Equal("value"))
			})
		})

		Context("Delete", func() {
			Context("Key exists", func() {
				BeforeEach(func() {
					store.Put("key", "value")
					value, err := store.Get("key")
					Expect(err).To(BeNil())
					Expect(value).To(Equal("value"))
				})

				It("removes value", func() {
					store.Delete("key")

					value, err := store.Get("key")
					Expect(err).To(BeNil())
					Expect(value).To(Equal(""))
				})

				It("returns true", func() {
					deleted, err := store.Delete("key")
					Expect(err).To(BeNil())
					Expect(deleted).To(BeTrue())
				})
			})

			Context("Key does not exist", func() {
				It("returns false", func() {
					deleted, err := store.Delete("key")
					Expect(deleted).To(BeFalse())
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
