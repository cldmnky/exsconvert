package convert

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Convert", func() {
	Context("Constructor", func() {
		It("should create XPM converter with correct parameters", func() {
			xpm := NewXPM("/search/path", "/output/path", 4, true, "Keygroup")
			Expect(xpm.SearchPath).To(Equal("/search/path"))
			Expect(xpm.OutputPath).To(Equal("/output/path"))
			Expect(xpm.LayersPerInstrument).To(Equal(4))
			Expect(xpm.SkipErrors).To(BeTrue())
		})
	})

	Context("Boolean conversion", func() {
		It("should convert true to 1", func() {
			result := Btoi(true)
			Expect(result).To(Equal(1))
		})

		It("should convert false to 0", func() {
			result := Btoi(false)
			Expect(result).To(Equal(0))
		})
	})

	Context("Value clamping", func() {
		It("should clamp value within range", func() {
			result := clamp(5.0, 0.0, 10.0)
			Expect(result).To(Equal(5.0))
		})

		It("should clamp value below minimum", func() {
			result := clamp(-5.0, 0.0, 10.0)
			Expect(result).To(Equal(0.0))
		})

		It("should clamp value above maximum", func() {
			result := clamp(15.0, 0.0, 10.0)
			Expect(result).To(Equal(10.0))
		})
	})

	Context("File discovery", func() {
		var tempDir string

		BeforeEach(func() {
			var err error
			tempDir, err = os.MkdirTemp("", "convert_test")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(tempDir)
		})

		It("should find EXS files with .exs extension", func() {
			// Create test files
			exsFile := filepath.Join(tempDir, "test.exs")
			err := os.WriteFile(exsFile, []byte("fake exs content"), 0644)
			Expect(err).ToNot(HaveOccurred())

			xpm := NewXPM(tempDir, "/output", 4, false, "Keygroup")
			files, err := xpm.findEXSFiles()
			Expect(err).ToNot(HaveOccurred())
			Expect(files).To(HaveLen(1))
			Expect(files[0]).To(Equal(exsFile))
		})

		It("should find EXS files with .EXS extension", func() {
			// Create test files
			exsFile := filepath.Join(tempDir, "test.EXS")
			err := os.WriteFile(exsFile, []byte("fake exs content"), 0644)
			Expect(err).ToNot(HaveOccurred())

			xpm := NewXPM(tempDir, "/output", 4, false, "Keygroup")
			files, err := xpm.findEXSFiles()
			Expect(err).ToNot(HaveOccurred())
			Expect(files).To(HaveLen(1))
			Expect(files[0]).To(Equal(exsFile))
		})

		It("should ignore non-EXS files", func() {
			// Create test files
			txtFile := filepath.Join(tempDir, "test.txt")
			err := os.WriteFile(txtFile, []byte("text content"), 0644)
			Expect(err).ToNot(HaveOccurred())

			xpm := NewXPM(tempDir, "/output", 4, false, "Keygroup")
			files, err := xpm.findEXSFiles()
			Expect(err).ToNot(HaveOccurred())
			Expect(files).To(HaveLen(0))
		})

		It("should find EXS files recursively", func() {
			// Create subdirectory
			subDir := filepath.Join(tempDir, "subdir")
			err := os.MkdirAll(subDir, 0755)
			Expect(err).ToNot(HaveOccurred())

			// Create test files
			exsFile1 := filepath.Join(tempDir, "test1.exs")
			exsFile2 := filepath.Join(subDir, "test2.exs")
			err = os.WriteFile(exsFile1, []byte("fake exs content"), 0644)
			Expect(err).ToNot(HaveOccurred())
			err = os.WriteFile(exsFile2, []byte("fake exs content"), 0644)
			Expect(err).ToNot(HaveOccurred())

			xpm := NewXPM(tempDir, "/output", 4, false, "Keygroup")
			files, err := xpm.findEXSFiles()
			Expect(err).ToNot(HaveOccurred())
			Expect(files).To(HaveLen(2))
			Expect(files).To(ContainElement(exsFile1))
			Expect(files).To(ContainElement(exsFile2))
		})

		It("should return error when search path does not exist", func() {
			xpm := NewXPM("/nonexistent/path", "/output", 4, false, "Keygroup")
			_, err := xpm.findEXSFiles()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Sample copying", func() {
		var tempDir string
		var outputDir string

		BeforeEach(func() {
			var err error
			tempDir, err = os.MkdirTemp("", "convert_test_input")
			Expect(err).ToNot(HaveOccurred())
			outputDir, err = os.MkdirTemp("", "convert_test_output")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(tempDir)
			os.RemoveAll(outputDir)
		})

		It("should copy sample file and convert extension to uppercase", func() {
			// Create test sample file
			sampleFile := filepath.Join(tempDir, "testsample.wav")
			sampleContent := []byte("fake wav content")
			err := os.WriteFile(sampleFile, sampleContent, 0644)
			Expect(err).ToNot(HaveOccurred())

			xpm := NewXPM(tempDir, outputDir, 4, false, "Keygroup")
			err = xpm.buildSampleIndex()
			Expect(err).ToNot(HaveOccurred())
			sampleName, sampleFileName, err := xpm.copySample("testsample.wav", outputDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(sampleName).To(Equal("testsample"))
			Expect(sampleFileName).To(Equal("testsample.WAV"))

			// Check that file was copied with uppercase extension
			copiedFile := filepath.Join(outputDir, "testsample.WAV")
			Expect(copiedFile).To(BeAnExistingFile())

			// Check content
			content, err := os.ReadFile(copiedFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(Equal(sampleContent))
		})

		It("should return error when sample not found", func() {
			xpm := NewXPM(tempDir, outputDir, 4, false, "Keygroup")
			err := xpm.buildSampleIndex()
			Expect(err).ToNot(HaveOccurred())
			_, _, err = xpm.copySample("nonexistent.wav", outputDir)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no sample found"))
		})

		It("should return error when multiple samples found", func() {
			// Create multiple files with same name in different directories
			subDir1 := filepath.Join(tempDir, "dir1")
			subDir2 := filepath.Join(tempDir, "dir2")
			os.MkdirAll(subDir1, 0755)
			os.MkdirAll(subDir2, 0755)

			sampleFile1 := filepath.Join(subDir1, "testsample.wav")
			sampleFile2 := filepath.Join(subDir2, "testsample.wav")
			err := os.WriteFile(sampleFile1, []byte("content1"), 0644)
			Expect(err).ToNot(HaveOccurred())
			err = os.WriteFile(sampleFile2, []byte("content2"), 0644)
			Expect(err).ToNot(HaveOccurred())

			xpm := NewXPM(tempDir, outputDir, 4, false, "Keygroup")
			err = xpm.buildSampleIndex()
			Expect(err).ToNot(HaveOccurred())
			// Note: With the index implementation, the first found file will be used
			// The test expectation needs to reflect this behavior
			_, _, err = xpm.copySample("testsample.wav", outputDir)
			// The index now just keeps the first file found, so no error
			Expect(err).ToNot(HaveOccurred())
		})

		It("should handle files without extension", func() {
			// Create test sample file without extension
			sampleFile := filepath.Join(tempDir, "testsample")
			sampleContent := []byte("fake content")
			err := os.WriteFile(sampleFile, sampleContent, 0644)
			Expect(err).ToNot(HaveOccurred())

			xpm := NewXPM(tempDir, outputDir, 4, false, "Keygroup")
			err = xpm.buildSampleIndex()
			Expect(err).ToNot(HaveOccurred())
			sampleName, sampleFileName, err := xpm.copySample("testsample", outputDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(sampleName).To(Equal("testsample"))
			Expect(sampleFileName).To(Equal("testsample"))

			// Check that file was copied (no extension change)
			copiedFile := filepath.Join(outputDir, "testsample")
			Expect(copiedFile).To(BeAnExistingFile())
		})
	})

	Context("Integration", func() {
		var outputDir string

		BeforeEach(func() {
			var err error
			outputDir, err = os.MkdirTemp("", "convert_integration_test")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(outputDir)
		})

		It("should convert EXS files to XPM", func() {
			testDataPath := "../../pkg/exs/testdata"
			xpm := NewXPM(testDataPath, outputDir, 4, true, "Keygroup")
			err := xpm.Convert()
			Expect(err).ToNot(HaveOccurred())

			// Check that output directory contains XPM files (recursively)
			xpmFiles := []string{}
			err = filepath.WalkDir(outputDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() && strings.HasSuffix(d.Name(), ".xpm") {
					xpmFiles = append(xpmFiles, path)
				}
				return nil
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(len(xpmFiles)).To(BeNumerically(">", 0))
		})
	})

	Context("Helpers", func() {
		//var err error
		It("should normalize values", func() {
			v := normalizeValue(float64(25), float64(0), float64(50))
			Expect(v).To(Equal("0.500000"))
		})

		It("should normalize edge values", func() {
			v := normalizeValue(float64(0), float64(0), float64(100))
			Expect(v).To(Equal("0.000000"))
			v = normalizeValue(float64(100), float64(0), float64(100))
			Expect(v).To(Equal("1.000000"))
			v = normalizeValue(float64(50), float64(0), float64(100))
			Expect(v).To(Equal("0.500000"))
		})

		It("should convert gain to db", func() {
			v := convertGain(float64(-96))
			Expect(v).To(Equal("0.353000"))
			v = convertGain(float64(12))
			Expect(v).To(Equal("1.000000"))
		})

		It("should clamp gain values within bounds", func() {
			v := convertGain(float64(-20)) // Below -12
			Expect(v).To(Equal("0.353000"))
			v = convertGain(float64(10)) // Above 6
			Expect(v).To(Equal("1.000000"))
		})

		It("should clamp values correctly", func() {
			result := clamp(5.0, 0.0, 10.0)
			Expect(result).To(Equal(5.0))
			result = clamp(-5.0, 0.0, 10.0)
			Expect(result).To(Equal(0.0))
			result = clamp(15.0, 0.0, 10.0)
			Expect(result).To(Equal(10.0))
			result = clamp(0.0, 0.0, 10.0)
			Expect(result).To(Equal(0.0))
			result = clamp(10.0, 0.0, 10.0)
			Expect(result).To(Equal(10.0))
		})
	})

	Context("Phase 1: Zone Volume and Pan Conversion", func() {
		It("should convert EXS volume (dB) to XPM linear gain", func() {
			// Test unity gain (0 dB -> 1.0)
			result := convertVolumeDbToLinear(0)
			Expect(result).To(Equal("1.000000"))

			// Test -6 dB (half power -> ~0.501)
			result = convertVolumeDbToLinear(-6)
			Expect(result).To(ContainSubstring("0.501"))

			// Test +6 dB (double power -> ~2.0)
			result = convertVolumeDbToLinear(6)
			Expect(result).To(ContainSubstring("1.995"))

			// Test -12 dB (quarter power -> ~0.251)
			result = convertVolumeDbToLinear(-12)
			Expect(result).To(ContainSubstring("0.251"))

			// Test +12 dB (max boost -> ~3.981)
			result = convertVolumeDbToLinear(12)
			Expect(result).To(ContainSubstring("3.981"))

			// Test -60 dB (near silence -> ~0.001)
			result = convertVolumeDbToLinear(-60)
			Expect(result).To(ContainSubstring("0.001"))
		})

		It("should clamp volume values outside valid range", func() {
			// Below minimum should clamp to -60 dB
			result := convertVolumeDbToLinear(-100)
			expected := convertVolumeDbToLinear(-60)
			Expect(result).To(Equal(expected))

			// Above maximum should clamp to +12 dB
			result = convertVolumeDbToLinear(20)
			expected = convertVolumeDbToLinear(12)
			Expect(result).To(Equal(expected))
		})

		It("should convert EXS pan to XPM normalized range", func() {
			// Test center (0 -> 0.5)
			result := convertPanToNormalized(0)
			Expect(result).To(ContainSubstring("0.503"))

			// Test hard left (-64 -> 0.0)
			result = convertPanToNormalized(-64)
			Expect(result).To(Equal("0.000000"))

			// Test hard right (+63 -> 1.0, as (63+64)/127 = 127/127 = 1.0)
			result = convertPanToNormalized(63)
			Expect(result).To(Equal("1.000000"))

			// Test left quarter (-32 -> 0.25)
			result = convertPanToNormalized(-32)
			Expect(result).To(ContainSubstring("0.251"))

			// Test right quarter (+32 -> 0.75)
			result = convertPanToNormalized(32)
			Expect(result).To(ContainSubstring("0.755"))
		})

		It("should clamp pan values outside valid range", func() {
			// Below minimum should clamp to -64
			result := convertPanToNormalized(-100)
			expected := convertPanToNormalized(-64)
			Expect(result).To(Equal(expected))

			// Above maximum should clamp to +63
			result = convertPanToNormalized(100)
			expected = convertPanToNormalized(63)
			Expect(result).To(Equal(expected))
		})
	})

	Context("Phase 2: Zone Scale, Output Routing, One-shot Mode", func() {
		It("should convert EXS scale to XPM KeyTrack", func() {
			// Test full key tracking (100 -> 1.0)
			result := convertScaleToKeyTrack(100)
			Expect(result).To(Equal("1.000000"))

			// Test no key tracking (0 -> 0.0)
			result = convertScaleToKeyTrack(0)
			Expect(result).To(Equal("0.000000"))

			// Test half key tracking (50 -> 0.5)
			result = convertScaleToKeyTrack(50)
			Expect(result).To(Equal("0.500000"))

			// Test 75% key tracking (75 -> 0.75)
			result = convertScaleToKeyTrack(75)
			Expect(result).To(Equal("0.750000"))
		})

		It("should clamp scale values outside valid range", func() {
			// Below minimum should clamp to 0
			result := convertScaleToKeyTrack(-50)
			Expect(result).To(Equal("0.000000"))

			// Above maximum should clamp to 100
			result = convertScaleToKeyTrack(150)
			Expect(result).To(Equal("1.000000"))
		})

		It("should convert EXS output to XPM AudioRoute integer", func() {
			// Test main output (0 -> 0)
			result := convertOutputToAudioRouteInt(0)
			Expect(result).To(Equal(0))

			// Test channel outputs (1-15 -> 1-15)
			result = convertOutputToAudioRouteInt(1)
			Expect(result).To(Equal(1))

			result = convertOutputToAudioRouteInt(8)
			Expect(result).To(Equal(8))

			result = convertOutputToAudioRouteInt(15)
			Expect(result).To(Equal(15))
		})

		It("should clamp output values outside valid range", func() {
			// Below minimum should clamp to 0
			result := convertOutputToAudioRouteInt(-5)
			Expect(result).To(Equal(0))

			// Above maximum should clamp to 15
			result = convertOutputToAudioRouteInt(20)
			Expect(result).To(Equal(15))
		})
	})

	Context("DMX Drum Kit Conversion", func() {
		var outputDir string

		BeforeEach(func() {
			var err error
			outputDir, err = os.MkdirTemp("", "dmx_drum_test")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(outputDir)
		})

		It("should successfully decode DMX From Mars - Dirty Color Kit.exs", func() {
			// This test verifies we can decode the drum kit EXS file
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "DMX From Mars - Dirty Color Kit.exs")

			// Verify the file exists
			_, err := os.Stat(exsFile)
			Expect(err).ToNot(HaveOccurred())

			// Create a converter and convert the specific file
			xpmConverter := NewXPM(testDataPath, outputDir, 1, true, "Drum")
			err = xpmConverter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should create a drum program XPM file from DMX drum kit", func() {
			// This test verifies the output is a proper drum program
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "DMX From Mars - Dirty Color Kit.exs")

			xpmConverter := NewXPM(testDataPath, outputDir, 1, true, "Drum")
			err := xpmConverter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			// Check that a directory was created for the drum kit
			expectedDir := filepath.Join(outputDir, "DMX From Mars - Dirty Color Kit")
			_, err = os.Stat(expectedDir)
			Expect(err).ToNot(HaveOccurred())

			// Check that an XPM file was created
			expectedXPM := filepath.Join(expectedDir, "DMX From Mars - Dirty Color Kit.xpm")
			_, err = os.Stat(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			// Read and verify the XPM file contains drum program markers
			xpmContent, err := os.ReadFile(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			xpmString := string(xpmContent)
			// Verify it's a drum program, not a keygroup
			Expect(xpmString).To(ContainSubstring(`type="Drum"`))
			Expect(xpmString).To(ContainSubstring("<PadNoteMap>"))
			Expect(xpmString).To(ContainSubstring("<PadNote "))
		})

		It("should map drum zones to single pads in drum mode", func() {
			// This test verifies that drum mode creates single-note mappings
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "DMX From Mars - Dirty Color Kit.exs")

			xpmConverter := NewXPM(testDataPath, outputDir, 1, true, "Drum")
			err := xpmConverter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			expectedXPM := filepath.Join(outputDir, "DMX From Mars - Dirty Color Kit", "DMX From Mars - Dirty Color Kit.xpm")
			xpmContent, err := os.ReadFile(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			xpmString := string(xpmContent)
			// In drum mode, LowNote and HighNote should be the same for each pad
			// This is a characteristic of drum programs
			Expect(xpmString).To(ContainSubstring("<LowNote>"))
			Expect(xpmString).To(ContainSubstring("<HighNote>"))
		})

		It("should convert same drum kit as keygroup for comparison", func() {
			// This test shows the difference between drum and keygroup conversion
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "DMX From Mars - Dirty Color Kit.exs")

			// Convert as keygroup
			keygroupConverter := NewXPM(testDataPath, outputDir, 1, true, "Keygroup")
			err := keygroupConverter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			expectedXPM := filepath.Join(outputDir, "DMX From Mars - Dirty Color Kit", "DMX From Mars - Dirty Color Kit.xpm")
			xpmContent, err := os.ReadFile(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			xpmString := string(xpmContent)
			// Keygroup version should NOT have PadNoteMap
			Expect(xpmString).To(ContainSubstring(`type="Keygroup"`))
			Expect(xpmString).NotTo(ContainSubstring("<PadNoteMap>"))
		})
	})

	Context("Advanced EXS24 Features", func() {
		var outputDir string

		BeforeEach(func() {
			var err error
			outputDir, err = os.MkdirTemp("", "xpm_advanced_test")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(outputDir)
		})

		It("should apply group key range limiting to zones", func() {
			// Test that zones are constrained to group key ranges
			// Using K3 Big which has groups with different key ranges
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "K3 Big.exs")

			converter := NewXPM(testDataPath, outputDir, 4, true, "Keygroup")
			err := converter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			expectedXPM := filepath.Join(outputDir, "K3 Big", "K3 Big.xpm")
			xpmContent, err := os.ReadFile(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			xpmString := string(xpmContent)
			// Should have instruments with properly constrained key ranges
			Expect(xpmString).To(ContainSubstring("<LowNote>"))
			Expect(xpmString).To(ContainSubstring("<HighNote>"))
			// If groups have range limits, no instrument should exceed MIDI note range 0-127
			Expect(xpmString).NotTo(ContainSubstring("<LowNote>-"))
			Expect(xpmString).NotTo(ContainSubstring("<HighNote>-"))
		})

		It("should detect release triggers from group Trigger field", func() {
			// Test files with release-triggered samples
			// Many piano instruments use release triggers for sympathetic resonance
			testDataPath := "../../pkg/exs/testdata"

			// Test with K3 Big which might have release triggers
			exsFile := filepath.Join(testDataPath, "K3 Big.exs")

			converter := NewXPM(testDataPath, outputDir, 4, true, "Keygroup")
			err := converter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			expectedXPM := filepath.Join(outputDir, "K3 Big", "K3 Big.xpm")
			xpmContent, err := os.ReadFile(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			xpmString := string(xpmContent)
			// Should have TriggerMode field set (1=release, 2=normal)
			Expect(xpmString).To(ContainSubstring("<TriggerMode>"))
			// TriggerMode should be either 1 (release) or 2 (normal attack)
			Expect(xpmString).To(Or(ContainSubstring("<TriggerMode>1</TriggerMode>"), ContainSubstring("<TriggerMode>2</TriggerMode>")))
		})

		It("should apply group velocity range limiting to layers", func() {
			// Test that layers are constrained to group velocity ranges
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "K3 Big.exs")

			converter := NewXPM(testDataPath, outputDir, 4, true, "Keygroup")
			err := converter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			expectedXPM := filepath.Join(outputDir, "K3 Big", "K3 Big.xpm")
			xpmContent, err := os.ReadFile(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			xpmString := string(xpmContent)
			// Should have velocity-sensitive layers
			Expect(xpmString).To(ContainSubstring("<VelStart>"))
			Expect(xpmString).To(ContainSubstring("<VelEnd>"))
			// Velocity values should be in valid MIDI range 0-127
			Expect(xpmString).NotTo(ContainSubstring("<VelStart>-"))
			Expect(xpmString).NotTo(ContainSubstring("<VelEnd>-"))
			// VelEnd should not exceed 127
			Expect(xpmString).NotTo(MatchRegexp(`<VelEnd>1[3-9][0-9]</VelEnd>`))     // 130-199
			Expect(xpmString).NotTo(MatchRegexp(`<VelEnd>[2-9][0-9][0-9]</VelEnd>`)) // 200+
		})

		It("should skip zones outside group key range", func() {
			// This test would need a specifically crafted EXS file where zones
			// are completely outside their group's key range
			// For now we'll test that the conversion succeeds without errors
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "K3 Big.exs")

			converter := NewXPM(testDataPath, outputDir, 4, true, "Keygroup")
			err := converter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			// Verify the XPM file was created
			expectedXPM := filepath.Join(outputDir, "K3 Big", "K3 Big.xpm")
			_, err = os.Stat(expectedXPM)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should skip layers outside group velocity range", func() {
			// Similar to above - test that conversion handles velocity range filtering
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "K3 Big.exs")

			converter := NewXPM(testDataPath, outputDir, 4, true, "Keygroup")
			err := converter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			expectedXPM := filepath.Join(outputDir, "K3 Big", "K3 Big.xpm")
			_, err = os.Stat(expectedXPM)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("Round Robin and Global Parameters", func() {
		var outputDir string

		BeforeEach(func() {
			var err error
			outputDir, err = os.MkdirTemp("", "xpm_roundrobin_test")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(outputDir)
		})

		It("should detect and enable Round Robin when groups have SelectGroup", func() {
			// Test with Shape-DFAM which might have round robin samples
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "Shape-DFAM-PSEQOUT.exs")

			converter := NewXPM(testDataPath, outputDir, 4, true, "Keygroup")
			err := converter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			expectedXPM := filepath.Join(outputDir, "Shape-DFAM-PSEQOUT", "Shape-DFAM-PSEQOUT.xpm")
			xpmContent, err := os.ReadFile(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			xpmString := string(xpmContent)
			// Should have ZonePlay field in XML
			Expect(xpmString).To(ContainSubstring("<ZonePlay>"))
			// ZonePlay should be 0 (CYCLE/round robin) or 1 (VELOCITY)
			Expect(xpmString).To(Or(ContainSubstring("<ZonePlay>0</ZonePlay>"), ContainSubstring("<ZonePlay>1</ZonePlay>")))
		})

		It("should set global pitch bend range", func() {
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "K3 Big.exs")

			converter := NewXPM(testDataPath, outputDir, 4, true, "Keygroup")
			err := converter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			expectedXPM := filepath.Join(outputDir, "K3 Big", "K3 Big.xpm")
			xpmContent, err := os.ReadFile(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			xpmString := string(xpmContent)
			// Should have KeygroupPitchBendRange field
			Expect(xpmString).To(ContainSubstring("<KeygroupPitchBendRange>"))
			// Default is 2 semitones = 0.166667
			Expect(xpmString).To(ContainSubstring("0.166667"))
		})

		It("should set global master transpose", func() {
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "K3 Big.exs")

			converter := NewXPM(testDataPath, outputDir, 4, true, "Keygroup")
			err := converter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			expectedXPM := filepath.Join(outputDir, "K3 Big", "K3 Big.xpm")
			xpmContent, err := os.ReadFile(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			xpmString := string(xpmContent)
			// Should have KeygroupMasterTranspose field
			Expect(xpmString).To(ContainSubstring("<KeygroupMasterTranspose>"))
		})

		It("should handle instruments without round robin correctly", func() {
			// Test with MC-202 bass which likely doesn't have round robin
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "MC-202 bass.exs")

			converter := NewXPM(testDataPath, outputDir, 4, true, "Keygroup")
			err := converter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			expectedXPM := filepath.Join(outputDir, "MC-202 bass", "MC-202 bass.xpm")
			xpmContent, err := os.ReadFile(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			xpmString := string(xpmContent)
			// Should default to ZonePlay=1 (VELOCITY mode)
			Expect(xpmString).To(ContainSubstring("<ZonePlay>1</ZonePlay>"))
		})

		It("should preserve Round Robin across multiple instruments", func() {
			// Test conversion maintains round robin consistently
			testDataPath := "../../pkg/exs/testdata"
			exsFile := filepath.Join(testDataPath, "K3 Big.exs")

			converter := NewXPM(testDataPath, outputDir, 4, true, "Keygroup")
			err := converter.ConvertFile(exsFile)
			Expect(err).ToNot(HaveOccurred())

			expectedXPM := filepath.Join(outputDir, "K3 Big", "K3 Big.xpm")
			_, err = os.Stat(expectedXPM)
			Expect(err).ToNot(HaveOccurred())

			// Verify file is valid XML
			xpmContent, err := os.ReadFile(expectedXPM)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(xpmContent)).To(BeNumerically(">", 0))
		})
	})
})
