package xpm

import (
	"encoding/xml"
	"fmt"
	"os"
)

var programPads = `{
    &quot;ProgramPads&quot;: {
        &quot;Universal&quot;: {
            &quot;value0&quot;: "true"
        },
        &quot;Type&quot;: {
            &quot;value0&quot;: 1
        },
        &quot;universalPad&quot;: 32512,
        &quot;pads&quot;: {
            &quot;value0&quot;: 0,
            &quot;value1&quot;: 0,
            &quot;value2&quot;: 0,
            &quot;value3&quot;: 0,
            &quot;value4&quot;: 0,
            &quot;value5&quot;: 0,
            &quot;value6&quot;: 0,
            &quot;value7&quot;: 0,
            &quot;value8&quot;: 0,
            &quot;value9&quot;: 0,
            &quot;value10&quot;: 0,
            &quot;value11&quot;: 0,
            &quot;value12&quot;: 0,
            &quot;value13&quot;: 0,
            &quot;value14&quot;: 0,
            &quot;value15&quot;: 0,
            &quot;value16&quot;: 0,
            &quot;value17&quot;: 0,
            &quot;value18&quot;: 0,
            &quot;value19&quot;: 0,
            &quot;value20&quot;: 0,
            &quot;value21&quot;: 0,
            &quot;value22&quot;: 0,
            &quot;value23&quot;: 0,
            &quot;value24&quot;: 0,
            &quot;value25&quot;: 0,
            &quot;value26&quot;: 0,
            &quot;value27&quot;: 0,
            &quot;value28&quot;: 0,
            &quot;value29&quot;: 0,
            &quot;value30&quot;: 0,
            &quot;value31&quot;: 0,
            &quot;value32&quot;: 0,
            &quot;value33&quot;: 0,
            &quot;value34&quot;: 0,
            &quot;value35&quot;: 0,
            &quot;value36&quot;: 0,
            &quot;value37&quot;: 0,
            &quot;value38&quot;: 0,
            &quot;value39&quot;: 0,
            &quot;value40&quot;: 0,
            &quot;value41&quot;: 0,
            &quot;value42&quot;: 0,
            &quot;value43&quot;: 0,
            &quot;value44&quot;: 0,
            &quot;value45&quot;: 0,
            &quot;value46&quot;: 0,
            &quot;value47&quot;: 0,
            &quot;value48&quot;: 0,
            &quot;value49&quot;: 0,
            &quot;value50&quot;: 0,
            &quot;value51&quot;: 0,
            &quot;value52&quot;: 0,
            &quot;value53&quot;: 0,
            &quot;value54&quot;: 0,
            &quot;value55&quot;: 0,
            &quot;value56&quot;: 0,
            &quot;value57&quot;: 0,
            &quot;value58&quot;: 0,
            &quot;value59&quot;: 0,
            &quot;value60&quot;: 0,
            &quot;value61&quot;: 0,
            &quot;value62&quot;: 0,
            &quot;value63&quot;: 0,
            &quot;value64&quot;: 0,
            &quot;value65&quot;: 0,
            &quot;value66&quot;: 0,
            &quot;value67&quot;: 0,
            &quot;value68&quot;: 0,
            &quot;value69&quot;: 0,
            &quot;value70&quot;: 0,
            &quot;value71&quot;: 0,
            &quot;value72&quot;: 0,
            &quot;value73&quot;: 0,
            &quot;value74&quot;: 0,
            &quot;value75&quot;: 0,
            &quot;value76&quot;: 0,
            &quot;value77&quot;: 0,
            &quot;value78&quot;: 0,
            &quot;value79&quot;: 0,
            &quot;value80&quot;: 0,
            &quot;value81&quot;: 0,
            &quot;value82&quot;: 0,
            &quot;value83&quot;: 0,
            &quot;value84&quot;: 0,
            &quot;value85&quot;: 0,
            &quot;value86&quot;: 0,
            &quot;value87&quot;: 0,
            &quot;value88&quot;: 0,
            &quot;value89&quot;: 0,
            &quot;value90&quot;: 0,
            &quot;value91&quot;: 0,
            &quot;value92&quot;: 0,
            &quot;value93&quot;: 0,
            &quot;value94&quot;: 0,
            &quot;value95&quot;: 0,
            &quot;value96&quot;: 0,
            &quot;value97&quot;: 0,
            &quot;value98&quot;: 0,
            &quot;value99&quot;: 0,
            &quot;value100&quot;: 0,
            &quot;value101&quot;: 0,
            &quot;value102&quot;: 0,
            &quot;value103&quot;: 0,
            &quot;value104&quot;: 0,
            &quot;value105&quot;: 0,
            &quot;value106&quot;: 0,
            &quot;value107&quot;: 0,
            &quot;value108&quot;: 0,
            &quot;value109&quot;: 0,
            &quot;value110&quot;: 0,
            &quot;value111&quot;: 0,
            &quot;value112&quot;: 0,
            &quot;value113&quot;: 0,
            &quot;value114&quot;: 0,
            &quot;value115&quot;: 0,
            &quot;value116&quot;: 0,
            &quot;value117&quot;: 0,
            &quot;value118&quot;: 0,
            &quot;value119&quot;: 0,
            &quot;value120&quot;: 0,
            &quot;value121&quot;: 0,
            &quot;value122&quot;: 0,
            &quot;value123&quot;: 0,
            &quot;value124&quot;: 0,
            &quot;value125&quot;: 0,
            &quot;value126&quot;: 0,
            &quot;value127&quot;: 0
        },
        &quot;UnusedPads&quot;: {
            &quot;value0&quot;: 1
        }
    }
}`

func NewXPMKeygroup() *MPCVObject {
	xpm := &MPCVObject{}
	xpm.Version = Version{
		FileVersion:        "2.1",
		Application:        "MPC-V",
		ApplicationVersion: "2.11.3.5",
		Platform:           "Linux",
	}
	//buff := new(bytes.Buffer)
	//xml.EscapeText(buff, []byte(programPads))
	xpm.Program = Program{
		Type:        "Keygroup",
		ProgramName: "",
		ProgramPads: ProgramPadsContent{
			Content: programPads,
		},
		CueBusEnable: "False",
		AudioRoute: AudioRoute{
			AudioRoute:              2,
			AudioRouteSubIndex:      0,
			AudioRouteChannelBitmap: 3,
			InsertsEnabled:          "True",
		},
		Send1:                      "0.000000",
		Send2:                      "0.000000",
		Send3:                      "0.000000",
		Send4:                      "0.000000",
		Volume:                     "0.707946",
		Mute:                       "False",
		Solo:                       "False",
		Pan:                        "0.500000",
		AutomationFilter:           1,
		Pitch:                      "0.000000",
		TuneCoarse:                 0,
		TuneFine:                   0,
		Mono:                       "False",
		ProgramPolyphony:           0,
		PortamentoTime:             "0.000000",
		PortamentoLegato:           "False",
		PortamentoQuantized:        "False",
		ProgramXfaderRoute:         0,
		KeygroupMasterTranspose:    "0.500000",
		KeygroupNumKeygroups:       0,
		KeygroupPitchBendRange:     "0.340000",
		KeygroupWheelToLfo:         "0.940000",
		KeygroupAftertouchToFilter: "0.410000",
	}

	instruments := make([]Instrument, 128)
	for i := 0; i < 128; i++ {
		instruments[i] = Instrument{
			Number:       fmt.Sprintf("%d", i+1),
			CueBusEnable: "False",
			AudioRoute: AudioRoute{
				AudioRoute:              0,
				AudioRouteSubIndex:      0,
				AudioRouteChannelBitmap: 3,
				InsertsEnabled:          "True",
			},
			Send1:                    "0.000000",
			Send2:                    "0.000000",
			Send3:                    "0.000000",
			Send4:                    "0.000000",
			Volume:                   "0.707946",
			Mute:                     "False",
			Solo:                     "False",
			Pan:                      "0.500000",
			AutomationFilter:         1,
			TuneCoarse:               0,
			TuneFine:                 0,
			Mono:                     "False",
			Polyphony:                0,
			FilterKeytrack:           "0.000000",
			LowNote:                  0,
			HighNote:                 127,
			IgnoreBaseNote:           "False",
			ZonePlay:                 0,
			MuteGroup:                0,
			MuteTarget1:              0,
			MuteTarget2:              0,
			MuteTarget3:              0,
			MuteTarget4:              0,
			SimultTarget1:            0,
			SimultTarget2:            0,
			SimultTarget3:            0,
			SimultTarget4:            0,
			LfoPitch:                 "0.000000",
			LfoCutoff:                "0.000000",
			LfoVolume:                "0.000000",
			LfoPan:                   "0.000000",
			OneShot:                  "False",
			FilterType:               3,
			Cutoff:                   "0.240000",
			Resonance:                "0.030000",
			FilterEnvAmt:             "0.330000",
			AfterTouchToFilter:       "0.000000",
			VelocityToStart:          "0.000000",
			VelocityToFilterAttack:   "0.000000",
			VelocityToFilter:         "0.000000",
			VelocityToFilterEnvelope: "0.250000",
			FilterAttack:             "0.000000",
			FilterDecay:              "0.640000",
			FilterSustain:            "0.007800",
			FilterRelease:            "0.000000",
			FilterHold:               "0.000000",
			FilterDecayType:          "True",
			FilterADEnvelope:         "True",
			VolumeHold:               "0.000000",
			VolumeDecayType:          "True",
			VolumeADEnvelope:         "True",
			VolumeAttack:             "0.000000",
			VolumeDecay:              "0.040000",
			VolumeSustain:            "1.000000",
			VolumeRelease:            "0.000000",
			VelocityToPitch:          "0.000000",
			VelocityToVolumeAttack:   "0.000000",
			VelocitySensitivity:      "0.310000",
			VelocityToPan:            "0.000000",
			LFO: LFO{
				Type:  "sine",
				Rate:  "0.500000",
				Sync:  0,
				Reset: "False",
			},
			WarpTempo:         "97.272003",
			WarpEnable:        "False",
			BpmLock:           "True",
			StretchPercentage: 100,
		}

		// create Layers
		layers := make([]Layer, 4)
		for j := 0; j < 4; j++ {
			layers[j] = Layer{
				Number:                   fmt.Sprintf("%d", j+1),
				Active:                   "True",
				Volume:                   "1.000000",
				Pan:                      "0.500000",
				TuneCoarse:               0,
				TuneFine:                 0,
				VelStart:                 0,
				VelEnd:                   127,
				SampleStart:              0,
				SampleEnd:                0,
				Loop:                     "False",
				LoopStart:                0,
				LoopEnd:                  0,
				LoopCrossfadeLength:      0,
				LoopTune:                 0,
				Mute:                     "False",
				RootNote:                 0,
				KeyTrack:                 "False",
				SampleName:               "",
				SampleFile:               "",
				SliceIndex:               129,
				Direction:                0,
				Offset:                   0,
				SliceStart:               0,
				SliceEnd:                 0,
				SliceLoopStart:           0,
				SliceLoop:                0,
				SliceLoopCrossFadeLength: 0,
			}
		}

		// add layers to instrument
		instruments[i].Layers = Layers{
			Layer: layers,
		}

	}

	xpm.Program.Instruments = Instruments{
		Instrument: instruments,
	}
	return xpm
}

func (xpm *MPCVObject) Save(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	out, err := xml.MarshalIndent(xpm, "", "  ")
	if err != nil {
		return err
	}
	out = []byte(xml.Header + string(out))
	_, err = f.Write(out)
	return err
}

// NewXPMDrum creates a new MPC Drum program template.
// Drum programs are used for drum kits where each pad triggers a single sample
// at a fixed pitch (one sample per MIDI note/pad).
func NewXPMDrum() *MPCVObject {
	xpm := &MPCVObject{}
	xpm.Version = Version{
		FileVersion:        "2.1",
		Application:        "MPC-V",
		ApplicationVersion: "2.11.3.5",
		Platform:           "Linux",
	}
	xpm.Program = Program{
		Type:        "Drum",
		ProgramName: "",
		ProgramPads: ProgramPadsContent{
			Content: programPads,
		},
		CueBusEnable: "False",
		AudioRoute: AudioRoute{
			AudioRoute:              2,
			AudioRouteSubIndex:      0,
			AudioRouteChannelBitmap: 3,
			InsertsEnabled:          "True",
		},
		Send1:                      "0.000000",
		Send2:                      "0.000000",
		Send3:                      "0.000000",
		Send4:                      "0.000000",
		Volume:                     "0.707946",
		Mute:                       "False",
		Solo:                       "False",
		Pan:                        "0.500000",
		AutomationFilter:           1,
		Pitch:                      "0.000000",
		TuneCoarse:                 0,
		TuneFine:                   0,
		Mono:                       "False",
		ProgramPolyphony:           0,
		PortamentoTime:             "0.000000",
		PortamentoLegato:           "False",
		PortamentoQuantized:        "False",
		ProgramXfaderRoute:         0,
		KeygroupMasterTranspose:    "0.500000",
		KeygroupNumKeygroups:       0,
		KeygroupPitchBendRange:     "0.000000", // No pitch bend for drums
		KeygroupWheelToLfo:         "0.000000",
		KeygroupAftertouchToFilter: "0.000000",
	}

	// Initialize PadNoteMap for drum programs (pads 1-128)
	padNotes := make([]PadNote, 128)
	for i := 0; i < 128; i++ {
		padNotes[i] = PadNote{
			Number: fmt.Sprintf("%d", i+1),
			Note:   i, // MIDI note 0-127
		}
	}
	xpm.Program.PadNoteMap = &PadNoteMap{
		PadNote: padNotes,
	}

	// Initialize PadGroupMap for drum programs
	xpm.Program.PadGroupMap = &PadGroupMap{}

	// Initialize instruments for drum pads (128 pads)
	instruments := make([]Instrument, 128)
	for i := 0; i < 128; i++ {
		instruments[i] = Instrument{
			Number:       fmt.Sprintf("%d", i+1),
			CueBusEnable: "False",
			AudioRoute: AudioRoute{
				AudioRoute:              0,
				AudioRouteSubIndex:      0,
				AudioRouteChannelBitmap: 3,
				InsertsEnabled:          "True",
			},
			Send1:            "0.000000",
			Send2:            "0.000000",
			Send3:            "0.000000",
			Send4:            "0.000000",
			Volume:           "0.707946",
			Mute:             "False",
			Solo:             "False",
			Pan:              "0.500000",
			AutomationFilter: 1,
			TuneCoarse:       0,
			TuneFine:         0,
			Mono:             "False",
			Polyphony:        1, // Drums typically monophonic per pad
			FilterKeytrack:   "False",
			// For drums, no LowNote/HighNote - each pad is a single note
			// The note is defined in PadNoteMap
			IgnoreBaseNote:         "False", // Drums should use base note
			ZonePlay:               0,
			MuteGroup:              0,
			LfoPitch:               "0.000000",
			LfoCutoff:              "0.000000",
			LfoVolume:              "0.000000",
			LfoPan:                 "0.000000",
			TriggerMode:            0, // One-shot for drums
			FilterType:             0,
			Cutoff:                 "1.000000",
			Resonance:              "0.000000",
			FilterEnvAmt:           "0.500000",
			VelocityToFilter:       "0.000000",
			FilterAttack:           "0.000000",
			FilterHold:             "0.000000",
			FilterDecay:            "0.000000",
			FilterSustain:          "1.000000",
			FilterRelease:          "0.000000",
			FilterAttackCurve:      "0.500000",
			FilterDecayCurve:       "0.500000",
			FilterReleaseCurve:     "0.500000",
			VolumeHold:             "0.000000",
			VolumeAttack:           "0.000000",
			VolumeDecay:            "0.000000",
			VolumeSustain:          "1.000000",
			VolumeRelease:          "0.005000",
			VolumeAttackCurve:      "0.500000",
			VolumeDecayCurve:       "0.500000",
			VolumeReleaseCurve:     "0.500000",
			PitchAttack:            "0.000000",
			PitchHold:              "0.000000",
			PitchDecay:             "0.000000",
			PitchSustain:           "0.500000",
			PitchRelease:           "0.000000",
			PitchAttackCurve:       "0.500000",
			PitchDecayCurve:        "0.500000",
			PitchReleaseCurve:      "0.500000",
			PitchEnvAmount:         "0.500000",
			VelocityToPitch:        "0.000000",
			VelocityToVolumeAttack: "0.000000",
			VelocitySensitivity:    "0.500000",
			VelocityToPan:          "0.000000",
			LFO: LFO{
				Type:         "Sine",
				Rate:         "0.250000",
				Sync:         0,
				Reset:        "Free",
				PitchAmount:  "0.000000",
				CutoffAmount: "0.000000",
				VolumeAmount: "0.000000",
				PanAmount:    "0.000000",
				Delay:        "0.000000",
				FadeIn:       "0.000000",
				Attack:       "0.000000",
				Depth:        "1.000000",
				Phase:        "0.000000",
				Offset:       "0.000000",
			},
			WarpTempo:         "0.000000",
			BpmLock:           "False",
			WarpEnable:        "False",
			StretchPercentage: 100,
			Layers: Layers{
				Layer: []Layer{},
			},
		}
	}
	xpm.Program.Instruments = Instruments{
		Instrument: instruments,
	}

	return xpm
}
