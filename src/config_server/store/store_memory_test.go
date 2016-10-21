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
				value1, _ := store.GetByName("key1")
				store.Put("key2", "value2")
				value2, _ := store.GetByName("key2")

				Expect(value1.Id).To(Equal("0"))
				Expect(value2.Id).To(Equal("1"))
			})
		})

		Context("GetByName", func() {
			It("should return associated value", func() {
				store.Put("some_name", "some_value")
				returnedValue, err := store.GetByName("some_name")
				Expect(err).To(BeNil())
				Expect(returnedValue).To(Equal(Configuration{
					Id:    "0",
					Name:  "some_name",
					Value: "some_value",
				}))
			})
		})

		Context("GetById", func() {
			It("should return associated value", func() {
				store.Put("some_name", "some_value")

				configuration, err := store.GetById("0")
				Expect(err).To(BeNil())
				Expect(configuration).To(Equal(Configuration{
					Id:    "0",
					Name:  "some_name",
					Value: "some_value",
				}))
			})
		})

		Context("Delete", func() {
			Context("Name exists", func() {
				BeforeEach(func() {
					store.Put("some_name", "some_value")
					value, err := store.GetByName("some_name")
					Expect(err).To(BeNil())
					Expect(value).To(Equal(Configuration{
						Id:    "0",
						Name:  "some_name",
						Value: "some_value",
					}))
				})

				It("removes value", func() {
					store.Delete("some_name")

					value, err := store.GetByName("some_name")
					Expect(err).To(BeNil())
					Expect(value).To(Equal(Configuration{}))
				})

				It("returns true", func() {
					deleted, err := store.Delete("some_name")
					Expect(err).To(BeNil())
					Expect(deleted).To(BeTrue())
				})
			})

			Context("Name does not exist", func() {
				It("returns false", func() {
					deleted, err := store.Delete("fake_key")
					Expect(deleted).To(BeFalse())
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
