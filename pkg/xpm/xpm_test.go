package xpm

import (
	"encoding/xml"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("XPM", func() {
	Context("when converting an EXS file", func() {
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

	Context("when parsing XPM testdata files", func() {
		var testdataFiles []string

		BeforeEach(func() {
			files, err := filepath.Glob("testdata/*.xpm")
			Expect(err).To(BeNil())
			testdataFiles = files
		})

		It("should find testdata files", func() {
			Expect(testdataFiles).ToNot(BeEmpty())
			Expect(len(testdataFiles)).To(BeNumerically(">", 0))
		})

		It("should parse all testdata XPM files", func() {
			for _, file := range testdataFiles {
				By("parsing " + filepath.Base(file))
				data, err := os.ReadFile(file)
				Expect(err).To(BeNil(), "should read file: "+file)
				Expect(data).ToNot(BeEmpty(), "file should not be empty: "+file)

				var xpm MPCVObject
				err = xml.Unmarshal(data, &xpm)
				Expect(err).To(BeNil(), "should parse XML: "+file)
			}
		})

		It("should validate XPM structure in all testdata files", func() {
			for _, file := range testdataFiles {
				By("validating " + filepath.Base(file))
				data, err := os.ReadFile(file)
				Expect(err).To(BeNil())

				var xpm MPCVObject
				err = xml.Unmarshal(data, &xpm)
				Expect(err).To(BeNil())

				// Validate basic structure - type can be "Drum" or "Keygroup"
				Expect(xpm.Program.Type).To(Or(Equal("Drum"), Equal("Keygroup")), "type should be Drum or Keygroup: "+file)
				Expect(xpm.Program.ProgramName).ToNot(BeEmpty(), "program name should not be empty: "+file)
			}
		})

		It("should validate Instruments in testdata files", func() {
			for _, file := range testdataFiles {
				By("checking instruments in " + filepath.Base(file))
				data, err := os.ReadFile(file)
				Expect(err).To(BeNil())

				var xpm MPCVObject
				err = xml.Unmarshal(data, &xpm)
				Expect(err).To(BeNil())

				// Skip files that might not have instruments (like Empty.xpm)
				if filepath.Base(file) == "Empty.xpm" {
					continue
				}

				// Check for Instruments
				Expect(xpm.Program.Instruments.Instrument).ToNot(BeEmpty(), "should have at least one instrument: "+file)

				// Validate each instrument has basic properties
				for i, inst := range xpm.Program.Instruments.Instrument {
					// Instruments have Number attribute but not Name field
					Expect(inst.Number).ToNot(BeEmpty(), "instrument %d number should not be empty in %s", i, file)
				}
			}
		})

		It("should validate Layers in testdata files", func() {
			for _, file := range testdataFiles {
				By("checking layers in " + filepath.Base(file))
				data, err := os.ReadFile(file)
				Expect(err).To(BeNil())

				var xpm MPCVObject
				err = xml.Unmarshal(data, &xpm)
				Expect(err).To(BeNil())

				// Skip files that might not have layers
				if filepath.Base(file) == "Empty.xpm" {
					continue
				}

				// Count total layers across all instruments
				totalLayers := 0
				layersWithSamples := 0

				// Check for Layers in instruments
				for _, inst := range xpm.Program.Instruments.Instrument {
					totalLayers += len(inst.Layers.Layer)
					for _, layer := range inst.Layers.Layer {
						// Count layers with actual sample names
						if layer.SampleName != "" {
							layersWithSamples++
						}
					}
				}

				// At least some layers should have sample names (unless it's an empty/template file)
				if totalLayers > 0 {
					Expect(layersWithSamples).To(BeNumerically(">", 0), "at least one layer should have a sample name in %s", file)
				}
			}
		})
	})
})
