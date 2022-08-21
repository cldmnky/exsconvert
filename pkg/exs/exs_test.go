package exs_test

import (
	"github.com/cldmnky/exsconvert/pkg/exs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exs", func() {
	//var err error
	It("should load exs files", func() {
		//_, err = exs.NewExsFromFile("testdata/MC-202 bass.exs")
		//Expect(err).To(BeNil())
		//_, err = exs.NewExsFromFile("testdata/Big News (slow sweeps).exs")
		//Expect(err).To(BeNil())
		//_, err = exs.NewExsFromFile("testdata/K3 Big.exs")
		//Expect(err).To(BeNil())
		//_, err = exs.NewExsFromFile("testdata/filter-DFAM-WFM-LP.exs")
		//Expect(err).To(BeNil())
		//_, err = exs.NewExsFromFile("testdata/MC-202 bass.exs")
		//Expect(err).To(BeNil())
		//_, err = exs.NewExsFromFile("testdata/80sThreats.exs")
		//Expect(err).To(BeNil())
		exs, err := exs.NewFromFile("testdata/Hi Hat 909 Clean.exs")
		Expect(err).To(BeNil())
		Expect(exs.BigEndian).To(BeFalse())
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
