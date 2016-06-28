package store_test

import (
	. "config_server/store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EncoderDecoder", func() {

	Describe("given a properly initalized EncoderDecoder object", func() {
		var encoderDecoder EncoderDecoder

		BeforeEach(func() {
			encoderDecoder = NewEncoderDecoder()
		})

		Context("Encoding", func() {
			It("should be able to encode a string value", func() {
				encodedStr, err := encoderDecoder.Encode("test")
				Expect(err).To(BeNil())
				Expect(encodedStr).To(Equal("G/+BAwEBA3ZhbAH/ggABAQEFVmFsdWUBEAAAABP/ggEGc3RyaW5nDAYABHRlc3QA"))
			})

			It("should be able to encode an integer value", func() {
				encodedStr, err := encoderDecoder.Encode(1)
				Expect(err).To(BeNil())
				Expect(encodedStr).To(Equal("G/+BAwEBA3ZhbAH/ggABAQEFVmFsdWUBEAAAAAz/ggEDaW50BAIAAgA="))

			})
		})

		Context("Decoding", func() {
			It("should be able to decode a string value", func() {
				decodedInterface, err := encoderDecoder.Decode("G/+BAwEBA3ZhbAH/ggABAQEFVmFsdWUBEAAAABP/ggEGc3RyaW5nDAYABHRlc3QA")
				Expect(err).To(BeNil())
				Expect(decodedInterface).To(Equal("test"))
			})

			It("should be able to decode an integer value", func() {
				decodedInterface, err := encoderDecoder.Decode("G/+BAwEBA3ZhbAH/ggABAQEFVmFsdWUBEAAAAAz/ggEDaW50BAIAAgA=")
				Expect(err).To(BeNil())
				Expect(decodedInterface).To(Equal(1))
			})
		})
	})
})
