package xpm

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("XPM", func() {
	//var err error
	It("should encode xpm files", func() {
		xpm := NewXPMKeygroup()
		Expect(xpm).ToNot(BeNil())
	})

	/* 	It("should detect endianness", func() {
	   		exs, err := exs.NewExsFromFile("testdata/MC-202 bass.exs")
	   		Expect(err).To(BeNil())
	   		Expect(exs.BigEndian).To(BeFalse())
	   	})

	   	It("should detect size expanded file", func() {
	   		exs, err := exs.NewExsFromFile("testdata/Big News (slow sweeps).exs")
	   		Expect(err).To(BeNil())
	   		Expect(exs.IsSizeExpanded).To(BeTrue())
	   	}) */
})
