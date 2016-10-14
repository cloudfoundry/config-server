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

			It("generates a unique id for new record", func() {
				store.Put("key1", "value1")
				value1, _ := store.Get("key1")
				store.Put("key2", "value2")
				value2, _ := store.Get("key2")

				Expect(value1.Id).To(Equal("0"))
				Expect(value2.Id).To(Equal("1"))
			})
		})

		Context("Get", func() {
			It("should return associated value", func() {
				store.Put("some_key", "some_value")
				returnedValue, err := store.Get("some_key")
				Expect(err).To(BeNil())
				Expect(returnedValue).To(Equal(Configuration{
					Id:    "0",
					Key:   "some_key",
					Value: "some_value",
				}))
			})
		})

		Context("Delete", func() {
			Context("Key exists", func() {
				BeforeEach(func() {
					store.Put("some_key", "some_value")
					value, err := store.Get("some_key")
					Expect(err).To(BeNil())
					Expect(value).To(Equal(Configuration{
						Id:    "0",
						Key:   "some_key",
						Value: "some_value",
					}))
				})

				It("removes value", func() {
					store.Delete("some_key")

					value, err := store.Get("some_key")
					Expect(err).To(BeNil())
					Expect(value).To(Equal(Configuration{}))
				})

				It("returns true", func() {
					deleted, err := store.Delete("some_key")
					Expect(err).To(BeNil())
					Expect(deleted).To(BeTrue())
				})
			})

			Context("Key does not exist", func() {
				It("returns false", func() {
					deleted, err := store.Delete("fake_key")
					Expect(deleted).To(BeFalse())
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
