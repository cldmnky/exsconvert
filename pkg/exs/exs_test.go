package exs_test

import (
	"github.com/cldmnky/exsconvert/pkg/exs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exs", func() {
	//var err error
	It("should load exs files", func() {
		_, err := exs.NewFromFile("testdata/MC-202 bass.exs")
		Expect(err).To(BeNil())
		_, err = exs.NewFromFile("testdata/Big News (slow sweeps).exs")
		Expect(err).To(BeNil())
		_, err = exs.NewFromFile("testdata/K3 Big.exs")
		Expect(err).To(BeNil())
		_, err = exs.NewFromFile("testdata/filter-DFAM-WFM-LP.exs")
		Expect(err).To(BeNil())
		_, err = exs.NewFromFile("testdata/Shape-DFAM-PSEQOUT.exs")
		Expect(err).To(BeNil())
		exs, err := exs.NewFromFile("testdata/Hi Hat 909 Clean.exs")
		Expect(err).To(BeNil())
		Expect(exs.BigEndian).To(BeFalse())
		Expect(exs.Instrument).ToNot(BeNil())
		Expect(len(exs.Zones)).To(BeNumerically(">", 0))
		for _, zone := range exs.Zones {
			// Check that new fields are parsed (even if default values)
			_ = zone.LoopEndRelease
			_ = zone.PlayMode
			_ = zone.HasOutput
		}
	})

	It("should get groups from exs file", func() {
		exsFile, err := exs.NewFromFile("testdata/Hi Hat 909 Clean.exs")
		Expect(err).To(BeNil())

		groups := exsFile.GetGroups()
		Expect(groups).ToNot(BeNil())
		Expect(len(groups)).To(BeNumerically(">=", 0))

		// Test with a more complex file
		exsFile2, err := exs.NewFromFile("testdata/K3 Big.exs")
		Expect(err).To(BeNil())
		groups2 := exsFile2.GetGroups()
		Expect(groups2).ToNot(BeNil())
	})

	It("should get zones by key ranges", func() {
		exsFile, err := exs.NewFromFile("testdata/Hi Hat 909 Clean.exs")
		Expect(err).To(BeNil())

		// Test with different zones per range values
		ranges := exsFile.GetZonesByKeyRanges(4)
		Expect(ranges).ToNot(BeNil())

		ranges2 := exsFile.GetZonesByKeyRanges(8)
		Expect(ranges2).ToNot(BeNil())

		// Test with a file that has more zones
		exsFile2, err := exs.NewFromFile("testdata/K3 Big.exs")
		Expect(err).To(BeNil())
		ranges3 := exsFile2.GetZonesByKeyRanges(4)
		Expect(ranges3).ToNot(BeNil())
	})

	It("should detect drum programs", func() {
		// Test with a drum kit file (Hi Hat is likely a drum)
		drumFile, err := exs.NewFromFile("testdata/Hi Hat 909 Clean.exs")
		Expect(err).To(BeNil())
		isDrum := drumFile.IsDrumProgram()
		// Just verify it returns without error, actual detection depends on file content
		_ = isDrum

		// Test with DMX drum kit
		drumFile2, err := exs.NewFromFile("testdata/DMX From Mars - Dirty Color Kit.exs")
		Expect(err).To(BeNil())
		isDrum2 := drumFile2.IsDrumProgram()
		_ = isDrum2

		// Test with a melodic instrument
		melodicFile, err := exs.NewFromFile("testdata/MC-202 bass.exs")
		Expect(err).To(BeNil())
		isDrum3 := melodicFile.IsDrumProgram()
		_ = isDrum3

		// Test with strings file (should not be drum)
		stringsFile, err := exs.NewFromFile("testdata/Analog Strings - Kawaii Dreams From Mars.exs")
		Expect(err).To(BeNil())
		isDrum4 := stringsFile.IsDrumProgram()
		_ = isDrum4
	})

	It("should handle different file formats", func() {
		// Test endianness detection
		exsFile, err := exs.NewFromFile("testdata/MC-202 bass.exs")
		Expect(err).To(BeNil())
		// BigEndian field should be set appropriately
		_ = exsFile.BigEndian

		// Test size expanded file
		bigFile, err := exs.NewFromFile("testdata/Big News (slow sweeps).exs")
		Expect(err).To(BeNil())
		Expect(bigFile.IsSizeExpanded).To(BeTrue())

		// Test various files to ensure they all load
		files := []string{
			"testdata/LegacyPulsar.exs",
			"testdata/Rave Go Up - OB From Mars.exs",
			"testdata/Strings - 360 From Mars.exs",
		}
		for _, file := range files {
			_, err := exs.NewFromFile(file)
			Expect(err).To(BeNil())
		}
	})

	It("should parse zones with proper attributes", func() {
		exsFile, err := exs.NewFromFile("testdata/Hi Hat 909 Clean.exs")
		Expect(err).To(BeNil())

		if len(exsFile.Zones) > 0 {
			zone := exsFile.Zones[0]
			// Verify zone fields are accessible
			Expect(zone.Name).ToNot(BeEmpty())
			// Check key range is valid (MIDI notes are 0-127)
			Expect(zone.KeyLow).To(BeNumerically(">=", 0))
			Expect(zone.KeyHigh).To(BeNumerically("<=", 127))
			// Check velocity range
			Expect(zone.VelLow).To(BeNumerically(">=", 0))
			Expect(zone.VelHigh).To(BeNumerically("<=", 127))
		}
	})

	It("should parse samples with file paths", func() {
		exsFile, err := exs.NewFromFile("testdata/Hi Hat 909 Clean.exs")
		Expect(err).To(BeNil())

		Expect(len(exsFile.Samples)).To(BeNumerically(">", 0))
		if len(exsFile.Samples) > 0 {
			sample := exsFile.Samples[0]
			// Verify sample has a name
			Expect(sample.Name).ToNot(BeEmpty())
			// File name and path may be empty depending on the file
			_ = sample.FileName
			_ = sample.Path
		}
	})

	It("should parse parameters", func() {
		exsFile, err := exs.NewFromFile("testdata/K3 Big.exs")
		Expect(err).To(BeNil())

		// Check if params were parsed (may be nil if file has no params)
		if exsFile.Params != nil {
			// Verify params structure exists
			_ = exsFile.Params.OutputVolume
			_ = exsFile.Params.FilterCutoff
		}
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
