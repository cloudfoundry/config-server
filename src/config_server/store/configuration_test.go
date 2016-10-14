package store_test

import (
	"config_server/store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration", func() {

	Describe("StringifiedJSON", func() {
		Context("When value is a string", func(){
			It("returns json string from the given db result", func(){
				configuration := store.Configuration{
					Id: "123",
					Key: "smurf",
					Value: `{"value": "blue"}`,
				}

				jsonString, _ := configuration.StringifiedJSON()

				Expect(jsonString).To(Equal(`{"id":"123","path":"smurf","value":"blue"}`))
			})
		})

		Context("When value is a number", func(){
			It("returns json string from the given db result", func(){
				configuration := store.Configuration{
					Id: "123",
					Key: "smurf",
					Value: `{"value": 123}`,
				}

				jsonString, _ := configuration.StringifiedJSON()

				Expect(jsonString).To(Equal(`{"id":"123","path":"smurf","value":123}`))
			})
		})

		Context("When value is complex", func(){
			It("returns json string from the given db result", func(){
				configuration := store.Configuration{
					Id: "123",
					Key: "smurf",
					Value: `{"value": {"smurf":"gargamel"}}`,
				}

				jsonString, _ := configuration.StringifiedJSON()

				Expect(jsonString).To(Equal(`{"id":"123","path":"smurf","value":{"smurf":"gargamel"}}`))
			})
		})

	})
})
