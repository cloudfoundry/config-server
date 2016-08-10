package types_test

import (
	. "config_server/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValueGeneratorFactoryConcrete", func() {
	var valueGeneratorFactory ValueGeneratorFactory

	Context("GetGenerator", func() {
		BeforeEach(func () {
			valueGeneratorFactory = NewValueGeneratorConcrete()
		})

		It("throws an error for unsupported value types", func() {
			_, err := valueGeneratorFactory.GetGenerator("bad_type")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Unsupported value type: bad_type"))
		})

        It("supports the password type", func() {
            generator, err := valueGeneratorFactory.GetGenerator("password")
            Expect(err).ToNot(HaveOccurred())
            Expect(generator).ToNot(BeNil())
        })
	})
})
