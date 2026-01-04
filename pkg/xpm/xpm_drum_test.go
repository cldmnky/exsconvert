package xpm_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cldmnky/exsconvert/pkg/xpm"
)

var _ = Describe("Drum Program XPM Generation", func() {
	Context("When creating a new drum program", func() {
		It("should create a valid drum program structure", func() {
			drumXPM := xpm.NewXPMDrum()

			Expect(drumXPM).NotTo(BeNil())
			Expect(drumXPM.Program.Type).To(Equal("Drum"))
			Expect(drumXPM.Program.ProgramName).To(Equal(""))
		})

		It("should have correct version information", func() {
			drumXPM := xpm.NewXPMDrum()

			Expect(drumXPM.Version.FileVersion).To(Equal("2.1"))
			Expect(drumXPM.Version.Application).To(Equal("MPC-V"))
		})

		It("should initialize 128 instruments for drum pads", func() {
			drumXPM := xpm.NewXPMDrum()

			Expect(len(drumXPM.Program.Instruments.Instrument)).To(Equal(128))
		})

		It("should have PadNoteMap with 128 entries", func() {
			drumXPM := xpm.NewXPMDrum()

			Expect(len(drumXPM.Program.PadNoteMap.PadNote)).To(Equal(128))
		})

		It("should map pad numbers 1-128 to MIDI notes 0-127", func() {
			drumXPM := xpm.NewXPMDrum()

			for i := 0; i < 128; i++ {
				padNote := drumXPM.Program.PadNoteMap.PadNote[i]
				// Numbers are formatted as strings like "1", "2", ... "128"
				Expect(padNote.Number).To(Equal(fmt.Sprintf("%d", i+1)))
				Expect(padNote.Note).To(Equal(i))
			}
		})

		It("should set drum-specific defaults", func() {
			drumXPM := xpm.NewXPMDrum()

			// No pitch bend for drums
			Expect(drumXPM.Program.KeygroupPitchBendRange).To(Equal("0.000000"))

			// Each instrument should have polyphony of 1 (monophonic per pad)
			for _, instrument := range drumXPM.Program.Instruments.Instrument {
				Expect(instrument.Polyphony).To(Equal(1))
			}
		})

		It("should set TriggerMode to 0 (one-shot) for all instruments", func() {
			drumXPM := xpm.NewXPMDrum()

			for _, instrument := range drumXPM.Program.Instruments.Instrument {
				Expect(instrument.TriggerMode).To(Equal(0))
			}
		})

		It("should not set LowNote/HighNote for drum instruments", func() {
			drumXPM := xpm.NewXPMDrum()

			// In drum mode, instruments don't use LowNote/HighNote
			// The note mapping is defined in PadNoteMap instead
			for _, instrument := range drumXPM.Program.Instruments.Instrument {
				Expect(instrument.LowNote).To(Equal(0))
				Expect(instrument.HighNote).To(Equal(0))
			}
		})
	})

	Context("When comparing Drum vs Keygroup programs", func() {
		It("should have different Type attributes", func() {
			drumXPM := xpm.NewXPMDrum()
			keygroupXPM := xpm.NewXPMKeygroup()

			Expect(drumXPM.Program.Type).To(Equal("Drum"))
			Expect(keygroupXPM.Program.Type).To(Equal("Keygroup"))
		})

		It("drum should have PadNoteMap while keygroup may not need it", func() {
			drumXPM := xpm.NewXPMDrum()

			Expect(len(drumXPM.Program.PadNoteMap.PadNote)).To(Equal(128))
			// Keygroup doesn't necessarily need PadNoteMap
		})

		It("drum should have zero pitch bend range", func() {
			drumXPM := xpm.NewXPMDrum()
			keygroupXPM := xpm.NewXPMKeygroup()

			Expect(drumXPM.Program.KeygroupPitchBendRange).To(Equal("0.000000"))
			Expect(keygroupXPM.Program.KeygroupPitchBendRange).NotTo(Equal("0.000000"))
		})
	})
})
