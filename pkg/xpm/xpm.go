package xpm

import (
	"encoding/xml"
)

type MPCVObject struct {
	XMLName xml.Name `xml:"MPCVObject"`
	Text    string   `xml:",chardata"`
	Version Version  `xml:"Version"`
	Program Program  `xml:"Program"`
}

type Version struct {
	Text               string `xml:",chardata"`
	FileVersion        string `xml:"File_Version"`
	Application        string `xml:"Application"`
	ApplicationVersion string `xml:"Application_Version"`
	Platform           string `xml:"Platform"`
}

type AudioRoute struct {
	Text                    string `xml:",chardata"`
	AudioRoute              int    `xml:"AudioRoute"`
	AudioRouteSubIndex      int    `xml:"AudioRouteSubIndex"`
	AudioRouteChannelBitmap int    `xml:"AudioRouteChannelBitmap"`
	InsertsEnabled          bool   `xml:"InsertsEnabled"`
}

type Instruments struct {
	Text       string       `xml:",chardata"`
	Instrument []Instrument `xml:"Instrument"`
}

type LFO struct {
	Text  string  `xml:",chardata"`
	Type  string  `xml:"Type"`
	Rate  float32 `xml:"Rate"`
	Sync  int     `xml:"Sync"`
	Reset bool    `xml:"Reset"`
}

type Layers struct {
	Text  string  `xml:",chardata"`
	Layer []Layer `xml:"Layer"`
}

type Layer struct {
	Text                     string  `xml:",chardata"`
	Number                   string  `xml:"number,attr"`
	Active                   bool    `xml:"Active"`
	Volume                   float32 `xml:"Volume"`
	Pan                      float32 `xml:"Pan"`
	Pitch                    float32 `xml:"Pitch"`
	TuneCoarse               int     `xml:"TuneCoarse"`
	TuneFine                 int     `xml:"TuneFine"`
	VelStart                 int     `xml:"VelStart"`
	VelEnd                   int     `xml:"VelEnd"`
	SampleStart              int     `xml:"SampleStart"`
	SampleEnd                int     `xml:"SampleEnd"`
	Loop                     bool    `xml:"Loop"`
	LoopStart                int     `xml:"LoopStart"`
	LoopEnd                  int     `xml:"LoopEnd"`
	LoopCrossfadeLength      int     `xml:"LoopCrossfadeLength"`
	LoopTune                 int     `xml:"LoopTune"`
	Mute                     bool    `xml:"Mute"`
	RootNote                 int     `xml:"RootNote"`
	KeyTrack                 bool    `xml:"KeyTrack"`
	SampleName               string  `xml:"SampleName"`
	SampleFile               string  `xml:"SampleFile"`
	SliceIndex               int     `xml:"SliceIndex"`
	Direction                int     `xml:"Direction"`
	Offset                   int     `xml:"Offset"`
	SliceStart               int     `xml:"SliceStart"`
	SliceEnd                 int     `xml:"SliceEnd"`
	SliceLoopStart           int     `xml:"SliceLoopStart"`
	SliceLoop                int     `xml:"SliceLoop"`
	SliceLoopCrossFadeLength int     `xml:"SliceLoopCrossFadeLength"`
}

type Instrument struct {
	Text                     string     `xml:",chardata"`
	Number                   string     `xml:"number,attr"`
	CueBusEnable             bool       `xml:"CueBusEnable"`
	AudioRoute               AudioRoute `xml:"AudioRoute"`
	Send1                    float32    `xml:"Send1"`
	Send2                    float32    `xml:"Send2"`
	Send3                    float32    `xml:"Send3"`
	Send4                    float32    `xml:"Send4"`
	Volume                   float32    `xml:"Volume"`
	Mute                     bool       `xml:"Mute"`
	Solo                     bool       `xml:"Solo"`
	Pan                      float32    `xml:"Pan"`
	AutomationFilter         int        `xml:"AutomationFilter"`
	TuneCoarse               int        `xml:"TuneCoarse"`
	TuneFine                 int        `xml:"TuneFine"`
	Mono                     bool       `xml:"Mono"`
	Polyphony                int        `xml:"Polyphony"`
	FilterKeytrack           float32    `xml:"FilterKeytrack"`
	LowNote                  int        `xml:"LowNote"`
	HighNote                 int        `xml:"HighNote"`
	IgnoreBaseNote           bool       `xml:"IgnoreBaseNote"`
	ZonePlay                 int        `xml:"ZonePlay"`
	MuteGroup                int        `xml:"MuteGroup"`
	MuteTarget1              int        `xml:"MuteTarget1"`
	MuteTarget2              int        `xml:"MuteTarget2"`
	MuteTarget3              int        `xml:"MuteTarget3"`
	MuteTarget4              int        `xml:"MuteTarget4"`
	SimultTarget1            int        `xml:"SimultTarget1"`
	SimultTarget2            int        `xml:"SimultTarget2"`
	SimultTarget3            int        `xml:"SimultTarget3"`
	SimultTarget4            int        `xml:"SimultTarget4"`
	LfoPitch                 float32    `xml:"LfoPitch"`
	LfoCutoff                float32    `xml:"LfoCutoff"`
	LfoVolume                float32    `xml:"LfoVolume"`
	LfoPan                   float32    `xml:"LfoPan"`
	OneShot                  bool       `xml:"OneShot"`
	FilterType               int        `xml:"FilterType"`
	Cutoff                   float32    `xml:"Cutoff"`
	Resonance                float32    `xml:"Resonance"`
	FilterEnvAmt             float32    `xml:"FilterEnvAmt"`
	AfterTouchToFilter       float32    `xml:"AfterTouchToFilter"`
	VelocityToStart          float32    `xml:"VelocityToStart"`
	VelocityToFilterAttack   float32    `xml:"VelocityToFilterAttack"`
	VelocityToFilter         float32    `xml:"VelocityToFilter"`
	VelocityToFilterEnvelope float32    `xml:"VelocityToFilterEnvelope"`
	FilterAttack             float32    `xml:"FilterAttack"`
	FilterDecay              float32    `xml:"FilterDecay"`
	FilterSustain            float32    `xml:"FilterSustain"`
	FilterRelease            float32    `xml:"FilterRelease"`
	FilterHold               float32    `xml:"FilterHold"`
	FilterDecayType          bool       `xml:"FilterDecayType"`
	FilterADEnvelope         bool       `xml:"FilterADEnvelope"`
	VolumeHold               float32    `xml:"VolumeHold"`
	VolumeDecayType          bool       `xml:"VolumeDecayType"`
	VolumeADEnvelope         bool       `xml:"VolumeADEnvelope"`
	VolumeAttack             float32    `xml:"VolumeAttack"`
	VolumeDecay              float32    `xml:"VolumeDecay"`
	VolumeSustain            float32    `xml:"VolumeSustain"`
	VolumeRelease            float32    `xml:"VolumeRelease"`
	VelocityToPitch          float32    `xml:"VelocityToPitch"`
	VelocityToVolumeAttack   float32    `xml:"VelocityToVolumeAttack"`
	VelocitySensitivity      float32    `xml:"VelocitySensitivity"`
	VelocityToPan            float32    `xml:"VelocityToPan"`
	LFO                      LFO        `xml:"LFO"`
	Layers                   Layers     `xml:"Layers"`
}

type Program struct {
	Text         string     `xml:",chardata"`
	Type         string     `xml:"type,attr"`
	ProgramName  string     `xml:"ProgramName"`
	ProgramPads  string     `xml:"ProgramPads"`
	CueBusEnable bool       `xml:"CueBusEnable"`
	AudioRoute   AudioRoute `xml:"AudioRoute"`

	Send1                      float32     `xml:"Send1"`
	Send2                      float32     `xml:"Send2"`
	Send3                      float32     `xml:"Send3"`
	Send4                      float32     `xml:"Send4"`
	Volume                     float32     `xml:"Volume"`
	Mute                       bool        `xml:"Mute"`
	Solo                       bool        `xml:"Solo"`
	Pan                        float32     `xml:"Pan"`
	AutomationFilter           int         `xml:"AutomationFilter"`
	Pitch                      float32     `xml:"Pitch"`
	TuneCoarse                 int         `xml:"TuneCoarse"`
	TuneFine                   int         `xml:"TuneFine"`
	Mono                       bool        `xml:"Mono"`
	ProgramPolyphony           int         `xml:"Program_Polyphony"`
	PortamentoTime             float32     `xml:"PortamentoTime"`
	PortamentoLegato           bool        `xml:"PortamentoLegato"`
	PortamentoQuantized        bool        `xml:"PortamentoQuantized"`
	ProgramXfaderRoute         int         `xml:"Program.Xfader.Route"`
	Instruments                Instruments `xml:"Instruments"`
	PadNoteMap                 PadNoteMap  `xml:"PadNoteMap"`
	PadGroupMap                PadGroupMap `xml:"PadGroupMap"`
	KeygroupMasterTranspose    float32     `xml:"KeygroupMasterTranspose"`
	KeygroupNumKeygroups       int         `xml:"KeygroupNumKeygroups"`
	KeygroupPitchBendRange     float32     `xml:"KeygroupPitchBendRange"`
	KeygroupWheelToLfo         float32     `xml:"KeygroupWheelToLfo"`
	KeygroupAftertouchToFilter float32     `xml:"KeygroupAftertouchToFilter"`
	QLinkAssignments           float32     `xml:"QLinkAssignments"`
}

type PadNoteMap struct {
	Text    string    `xml:",chardata"`
	PadNote []PadNote `xml:"PadNote"`
}

type PadNote struct {
	Text   string `xml:",chardata"`
	Number string `xml:"number,attr"`
	Note   int    `xml:"Note"`
}

type PadGroupMap struct {
	Text     string     `xml:",chardata"`
	PadGroup []PadGroup `xml:"PadGroup"`
}

type PadGroup struct {
	Text   string `xml:",chardata"`
	Number string `xml:"number,attr"`
	Group  int    `xml:"Group"`
}
