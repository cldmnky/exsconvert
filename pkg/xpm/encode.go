package xpm

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
)

var programPads = `{
	"ProgramPads": {
		"Universal": {
			"value0": true
		},
		"Type": {
			"value0": 1
		},
		"universalPad": 32512,
		"pads": {
			 "value0": 0,
			 "value1": 0,
			 "value2": 0,
			 "value3": 0,
			 "value4": 0,
			 "value5": 0,
			 "value6": 0,
			 "value7": 0,
			 "value8": 0,
			 "value9": 0,
			 "value10": 0,
			 "value11": 0,
			 "value12": 0,
			 "value13": 0,
			 "value14": 0,
			 "value15": 0,
			 "value16": 0,
			 "value17": 0,
			 "value18": 0,
			 "value19": 0,
			 "value20": 0,
			 "value21": 0,
			 "value22": 0,
			 "value23": 0,
			 "value24": 0,
			 "value25": 0,
			 "value26": 0,
			 "value27": 0,
			 "value28": 0,
			 "value29": 0,
			 "value30": 0,
			 "value31": 0,
			 "value32": 0,
			 "value33": 0,
			 "value34": 0,
			 "value35": 0,
			 "value36": 0,
			 "value37": 0,
			 "value38": 0,
			 "value39": 0,
			 "value40": 0,
			 "value41": 0,
			 "value42": 0,
			 "value43": 0,
			 "value44": 0,
			 "value45": 0,
			 "value46": 0,
			 "value47": 0,
			 "value48": 0,
			 "value49": 0,
			 "value50": 0,
			 "value51": 0,
			 "value52": 0,
			 "value53": 0,
			 "value54": 0,
			 "value55": 0,
			 "value56": 0,
			 "value57": 0,
			 "value58": 0,
			 "value59": 0,
			 "value60": 0,
			 "value61": 0,
			 "value62": 0,
			 "value63": 0,
			 "value64": 0,
			 "value65": 0,
			 "value66": 0,
			 "value67": 0,
			 "value68": 0,
			 "value69": 0,
			 "value70": 0,
			 "value71": 0,
			 "value72": 0,
			 "value73": 0,
			 "value74": 0,
			 "value75": 0,
			 "value76": 0,
			 "value77": 0,
			 "value78": 0,
			 "value79": 0,
			 "value80": 0,
			 "value81": 0,
			 "value82": 0,
			 "value83": 0,
			 "value84": 0,
			 "value85": 0,
			 "value86": 0,
			 "value87": 0,
			 "value88": 0,
			 "value89": 0,
			 "value90": 0,
			 "value91": 0,
			 "value92": 0,
			 "value93": 0,
			 "value94": 0,
			 "value95": 0,
			 "value96": 0,
			 "value97": 0,
			 "value98": 0,
			 "value99": 0,
			 "value100": 0,
			 "value101": 0,
			 "value102": 0,
			 "value103": 0,
			 "value104": 0,
			 "value105": 0,
			 "value106": 0,
			 "value107": 0,
			 "value108": 0,
			 "value109": 0,
			 "value110": 0,
			 "value111": 0,
			 "value112": 0,
			 "value113": 0,
			 "value114": 0,
			 "value115": 0,
			 "value116": 0,
			 "value117": 0,
			 "value118": 0,
			 "value119": 0,
			 "value120": 0,
			 "value121": 0,
			 "value122": 0,
			 "value123": 0,
			 "value124": 0,
			 "value125": 0,
			 "value126": 0,
			 "value127": 0
	},
		"UnusedPads": {
			"value0": 1
		}
	}
}`

func NewXPMKeygroup() *MPCVObject {
	xpm := &MPCVObject{}
	xpm.Version = Version{
		FileVersion:        "2.1",
		Application:        "MPC-V",
		ApplicationVersion: "2.9.1.2",
		Platform:           "Linux",
	}
	buff := new(bytes.Buffer)
	xml.EscapeText(buff, []byte(programPads))
	xpm.Program = Program{
		ProgramName:  "",
		ProgramPads:  buff.String(),
		CueBusEnable: false,
		AudioRoute: AudioRoute{
			AudioRoute:              2,
			AudioRouteSubIndex:      0,
			AudioRouteChannelBitmap: 3,
			InsertsEnabled:          true,
		},
		Send1:                      0,
		Send2:                      0,
		Send3:                      0,
		Send4:                      0,
		Volume:                     0.707946,
		Mute:                       false,
		Solo:                       false,
		Pan:                        0.5,
		AutomationFilter:           1,
		Pitch:                      0,
		TuneCoarse:                 0,
		TuneFine:                   0,
		Mono:                       false,
		ProgramPolyphony:           16,
		PortamentoTime:             0,
		PortamentoLegato:           false,
		PortamentoQuantized:        false,
		ProgramXfaderRoute:         0,
		KeygroupMasterTranspose:    0.5,
		KeygroupNumKeygroups:       1,
		KeygroupPitchBendRange:     0.5,
		KeygroupWheelToLfo:         1,
		KeygroupAftertouchToFilter: 0,
	}

	instruments := make([]Instrument, 128)
	for i := 0; i < 128; i++ {
		instruments[i] = Instrument{
			Number:       fmt.Sprintf("%d", i+1),
			CueBusEnable: false,
			AudioRoute: AudioRoute{
				AudioRoute:              0,
				AudioRouteSubIndex:      0,
				AudioRouteChannelBitmap: 3,
				InsertsEnabled:          true,
			},
			Send1:                    0,
			Send2:                    0,
			Send3:                    0,
			Send4:                    0,
			Volume:                   0.707946,
			Mute:                     false,
			Solo:                     false,
			Pan:                      0.5,
			AutomationFilter:         1,
			TuneCoarse:               0,
			TuneFine:                 0,
			Mono:                     false,
			FilterKeytrack:           0,
			LowNote:                  0,
			HighNote:                 127,
			IgnoreBaseNote:           false,
			ZonePlay:                 1,
			MuteGroup:                0,
			MuteTarget1:              0,
			MuteTarget2:              0,
			MuteTarget3:              0,
			MuteTarget4:              0,
			SimultTarget1:            0,
			SimultTarget2:            0,
			SimultTarget3:            0,
			SimultTarget4:            0,
			LfoPitch:                 0,
			LfoCutoff:                0,
			LfoVolume:                0.203125,
			LfoPan:                   0,
			OneShot:                  false,
			FilterType:               0,
			Cutoff:                   1,
			Resonance:                0,
			FilterEnvAmt:             0,
			AfterTouchToFilter:       0,
			VelocityToStart:          0,
			VelocityToFilterAttack:   0,
			VelocityToFilter:         0,
			VelocityToFilterEnvelope: 0,
			FilterAttack:             0,
			FilterDecay:              0.5,
			FilterSustain:            1,
			FilterRelease:            0,
			FilterHold:               0,
			FilterDecayType:          true,
			FilterADEnvelope:         true,
			VolumeHold:               0,
			VolumeDecayType:          true,
			VolumeADEnvelope:         true,
			VolumeAttack:             0,
			VolumeDecay:              0.5,
			VolumeSustain:            1,
			VolumeRelease:            0,
			VelocityToPitch:          0,
			VelocityToVolumeAttack:   0,
			VelocitySensitivity:      1,
			VelocityToPan:            0,
			LFO: LFO{
				Type:  "sine",
				Rate:  0.5,
				Sync:  0,
				Reset: false,
			},
		}

		// create Layers
		layers := make([]Layer, 4)
		for j := 0; j < 4; j++ {
			layers[j] = Layer{
				Number:                   fmt.Sprintf("%d", j+1),
				Active:                   true,
				Volume:                   1,
				Pan:                      0.5,
				TuneCoarse:               0,
				TuneFine:                 0,
				VelStart:                 0,
				VelEnd:                   127,
				SampleStart:              0,
				SampleEnd:                0,
				Loop:                     false,
				LoopStart:                0,
				LoopEnd:                  0,
				LoopCrossfadeLength:      0,
				LoopTune:                 0,
				Mute:                     false,
				RootNote:                 0,
				KeyTrack:                 false,
				SampleName:               "",
				SampleFile:               "",
				SliceIndex:               128,
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

	/* 	enc := xml.NewEncoder(os.Stdout)
	   	enc.Indent("  ", "    ")
	   	if err := enc.Encode(xpm); err != nil {
	   		fmt.Printf("error: %v\n", err)
	   	} */
	return xpm
}

func (xpm *MPCVObject) Save(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := xml.NewEncoder(f)
	enc.Indent("  ", "    ")
	if err := enc.Encode(xpm); err != nil {
		return err
	}
	return nil
}
