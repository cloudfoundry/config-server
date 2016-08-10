package types_test

import (
	. "config_server/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SecretGenerator", func() {

    Describe("secretGenerator", func() {
        var generator ValueGenerator

        BeforeEach(func() {
            generator = NewSecretGenerator()
        })

        Context("Generate", func() {
            It("generates a 20 character secret", func() {
                secret, err := generator.Generate(nil)
                Expect(err).ToNot(HaveOccurred())
                Expect(len(secret.(string))).To(Equal(20))
            })

            It("generates unique secrets", func() {
                secret1, err := generator.Generate(nil)
                Expect(err).ToNot(HaveOccurred())

                secret2, err := generator.Generate(nil)
                Expect(err).ToNot(HaveOccurred())

                Expect(secret1).ToNot(Equal(secret2))
            })

            It("only uses allowed characters", func() {
                for i := 0; i < 20; i++ { // arbitrary number
                    secret, err := generator.Generate(nil)
                    Expect(err).ToNot(HaveOccurred())
                    Expect(secret).To(MatchRegexp("^[A-Za-z0-9!@#$%^&*()\\-_=\\+,\\.?/:;\\{\\}\\[\\]`~]{20}$"))
                }
            })
        })
    })
})
