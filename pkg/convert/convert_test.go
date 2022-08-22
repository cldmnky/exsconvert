package convert

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Convert", func() {
	Context("Helpers", func() {
		//var err error
		It("should normalize values", func() {
			v := normalizeValue(float64(25), float64(0), float64(50))
			Expect(v).To(Equal("0.500000"))
		})

		It("should convert gain to db", func() {
			v := convertGain(float64(-96))
			Expect(v).To(Equal("0.353000"))
			v = convertGain(float64(12))
			Expect(v).To(Equal("1.000000"))
		})
	})
})
