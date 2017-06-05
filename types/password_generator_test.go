package types_test

import (
	. "github.com/cloudfoundry/config-server/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PasswordGenerator", func() {

	Describe("passwordGenerator", func() {
		var generator ValueGenerator

		BeforeEach(func() {
			generator = NewPasswordGenerator()
		})

		Context("Generate", func() {
			It("generates a 20 character password", func() {
				password, err := generator.Generate(nil)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(password.(string))).To(Equal(20))
			})

			It("generates unique passwords", func() {
				password1, err := generator.Generate(nil)
				Expect(err).ToNot(HaveOccurred())

				password2, err := generator.Generate(nil)
				Expect(err).ToNot(HaveOccurred())

				Expect(password1).ToNot(Equal(password2))
			})

			It("only uses allowed characters", func() {
				for i := 0; i < 20; i++ { // arbitrary number
					password, err := generator.Generate(nil)
					Expect(err).ToNot(HaveOccurred())
					Expect(password).To(MatchRegexp("^[a-z0-9]{20}$"))
				}
			})

			Context("when allowed_characters are specified", func() {
				It("generates a password with allowed_characters", func() {
					params := map[string]interface{}{"allowed_characters": "0123456789"}
					for i := 0; i < 20; i++ { // arbitrary number
						password, err := generator.Generate(params)
						Expect(err).ToNot(HaveOccurred())
						Expect(password).To(MatchRegexp("^[0-9]{20}$"))
					}
				})

				Context("when allowed_characters are not string", func() {
					It("returns an error", func() {
						params := map[string]interface{}{"allowed_characters": []int{1234567890}}
						_, err := generator.Generate(params)
						Expect(err).To(HaveOccurred())
					})
				})
			})

			Context("with specified length", func() {
				It("generates a password with specified length", func() {
					params := map[string]interface{}{"length": 50}
					password, err := generator.Generate(params)
					Expect(err).ToNot(HaveOccurred())
					Expect(len(password.(string))).To(Equal(50))
				})

				Context("when lenght is too small", func() {
					It("generates a 20 character password", func() {
						params := map[string]interface{}{"length": 10}
						password, err := generator.Generate(params)
						Expect(err).ToNot(HaveOccurred())
						Expect(len(password.(string))).To(Equal(20))
					})
				})

				Context("when length is not int", func() {
					It("returns an error", func() {
						params := map[string]interface{}{"length": "10"}
						_, err := generator.Generate(params)
						Expect(err).To(HaveOccurred())
					})
				})
			})

		})
	})
})
