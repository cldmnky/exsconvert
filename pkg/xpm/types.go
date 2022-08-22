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
	InsertsEnabled          string `xml:"InsertsEnabled"`
}

type Instruments struct {
	Text       string       `xml:",chardata"`
	Instrument []Instrument `xml:"Instrument"`
}

type LFO struct {
	Text  string `xml:",chardata"`
	Type  string `xml:"Type"`
	Rate  string `xml:"Rate"`
	Sync  int    `xml:"Sync"`
	Reset string `xml:"Reset"`
}

type Layers struct {
	Text  string  `xml:",chardata"`
	Layer []Layer `xml:"Layer"`
}

type Layer struct {
	Text                     string `xml:",chardata"`
	Number                   string `xml:"number,attr"`
	Active                   string `xml:"Active"`
	Volume                   string `xml:"Volume"`
	Pan                      string `xml:"Pan"`
	Pitch                    string `xml:"Pitch"`
	TuneCoarse               int    `xml:"TuneCoarse"`
	TuneFine                 int    `xml:"TuneFine"`
	VelStart                 int    `xml:"VelStart"`
	VelEnd                   int    `xml:"VelEnd"`
	SampleStart              int    `xml:"SampleStart"`
	SampleEnd                int    `xml:"SampleEnd"`
	Loop                     string `xml:"Loop"`
	LoopStart                int    `xml:"LoopStart"`
	LoopEnd                  int    `xml:"LoopEnd"`
	LoopCrossfadeLength      int    `xml:"LoopCrossfadeLength"`
	LoopTune                 int    `xml:"LoopTune"`
	Mute                     string `xml:"Mute"`
	RootNote                 int    `xml:"RootNote"`
	KeyTrack                 string `xml:"KeyTrack"`
	SampleName               string `xml:"SampleName"`
	SampleFile               string `xml:"SampleFile"`
	SliceIndex               int    `xml:"SliceIndex"`
	Direction                int    `xml:"Direction"`
	Offset                   int    `xml:"Offset"`
	SliceStart               int    `xml:"SliceStart"`
	SliceEnd                 int    `xml:"SliceEnd"`
	SliceLoopStart           int    `xml:"SliceLoopStart"`
	SliceLoop                int    `xml:"SliceLoop"`
	SliceLoopCrossFadeLength int    `xml:"SliceLoopCrossFadeLength"`
}

type Instrument struct {
	Text                     string     `xml:",chardata"`
	Number                   string     `xml:"number,attr"`
	CueBusEnable             string     `xml:"CueBusEnable"`
	AudioRoute               AudioRoute `xml:"AudioRoute"`
	Send1                    string     `xml:"Send1"`
	Send2                    string     `xml:"Send2"`
	Send3                    string     `xml:"Send3"`
	Send4                    string     `xml:"Send4"`
	Volume                   string     `xml:"Volume"`
	Mute                     string     `xml:"Mute"`
	Solo                     string     `xml:"Solo"`
	Pan                      string     `xml:"Pan"`
	AutomationFilter         int        `xml:"AutomationFilter"`
	TuneCoarse               int        `xml:"TuneCoarse"`
	TuneFine                 int        `xml:"TuneFine"`
	Mono                     string     `xml:"Mono"`
	Polyphony                int        `xml:"Polyphony"`
	FilterKeytrack           string     `xml:"FilterKeytrack"`
	LowNote                  int        `xml:"LowNote"`
	HighNote                 int        `xml:"HighNote"`
	IgnoreBaseNote           string     `xml:"IgnoreBaseNote"`
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
	LfoPitch                 string     `xml:"LfoPitch"`
	LfoCutoff                string     `xml:"LfoCutoff"`
	LfoVolume                string     `xml:"LfoVolume"`
	LfoPan                   string     `xml:"LfoPan"`
	OneShot                  string     `xml:"OneShot"`
	FilterType               int        `xml:"FilterType"`
	Cutoff                   string     `xml:"Cutoff"`
	Resonance                string     `xml:"Resonance"`
	FilterEnvAmt             string     `xml:"FilterEnvAmt"`
	AfterTouchToFilter       string     `xml:"AfterTouchToFilter"`
	VelocityToStart          string     `xml:"VelocityToStart"`
	VelocityToFilterAttack   string     `xml:"VelocityToFilterAttack"`
	VelocityToFilter         string     `xml:"VelocityToFilter"`
	VelocityToFilterEnvelope string     `xml:"VelocityToFilterEnvelope"`
	FilterAttack             string     `xml:"FilterAttack"`
	FilterDecay              string     `xml:"FilterDecay"`
	FilterSustain            string     `xml:"FilterSustain"`
	FilterRelease            string     `xml:"FilterRelease"`
	FilterHold               string     `xml:"FilterHold"`
	FilterDecayType          string     `xml:"FilterDecayType"`
	FilterADEnvelope         string     `xml:"FilterADEnvelope"`
	VolumeHold               string     `xml:"VolumeHold"`
	VolumeDecayType          string     `xml:"VolumeDecayType"`
	VolumeADEnvelope         string     `xml:"VolumeADEnvelope"`
	VolumeAttack             string     `xml:"VolumeAttack"`
	VolumeDecay              string     `xml:"VolumeDecay"`
	VolumeSustain            string     `xml:"VolumeSustain"`
	VolumeRelease            string     `xml:"VolumeRelease"`
	VelocityToPitch          string     `xml:"VelocityToPitch"`
	VelocityToVolumeAttack   string     `xml:"VelocityToVolumeAttack"`
	VelocitySensitivity      string     `xml:"VelocitySensitivity"`
	VelocityToPan            string     `xml:"VelocityToPan"`
	LFO                      LFO        `xml:"LFO"`
	WarpTempo                string     `xml:"WarpTempo"`
	BpmLock                  string     `xml:"BpmLock"`
	WarpEnable               string     `xml:"WarpEnable"`
	StretchPercentage        int        `xml:"StretchPercentage"`
	Layers                   Layers     `xml:"Layers"`
}

type ProgramPadsContent struct {
	Content string `xml:",innerxml"`
}

type Program struct {
	Text        string             `xml:",chardata"`
	Type        string             `xml:"type,attr"`
	ProgramName string             `xml:"ProgramName"`
	ProgramPads ProgramPadsContent `xml:"ProgramPads"`
	//ProgramPadsContent string     `xml:",innerxml"`
	CueBusEnable string     `xml:"CueBusEnable"`
	AudioRoute   AudioRoute `xml:"AudioRoute"`

	Send1                      string      `xml:"Send1"`
	Send2                      string      `xml:"Send2"`
	Send3                      string      `xml:"Send3"`
	Send4                      string      `xml:"Send4"`
	Volume                     string      `xml:"Volume"`
	Mute                       string      `xml:"Mute"`
	Solo                       string      `xml:"Solo"`
	Pan                        string      `xml:"Pan"`
	AutomationFilter           int         `xml:"AutomationFilter"`
	Pitch                      string      `xml:"Pitch"`
	TuneCoarse                 int         `xml:"TuneCoarse"`
	TuneFine                   int         `xml:"TuneFine"`
	Mono                       string      `xml:"Mono"`
	ProgramPolyphony           int         `xml:"Program_Polyphony"`
	PortamentoTime             string      `xml:"PortamentoTime"`
	PortamentoLegato           string      `xml:"PortamentoLegato"`
	PortamentoQuantized        string      `xml:"PortamentoQuantized"`
	ProgramXfaderRoute         int         `xml:"Program.Xfader.Route"`
	Instruments                Instruments `xml:"Instruments"`
	PadNoteMap                 PadNoteMap  `xml:"PadNoteMap"`
	PadGroupMap                PadGroupMap `xml:"PadGroupMap"`
	KeygroupMasterTranspose    string      `xml:"KeygroupMasterTranspose"`
	KeygroupNumKeygroups       int         `xml:"KeygroupNumKeygroups"`
	KeygroupPitchBendRange     string      `xml:"KeygroupPitchBendRange"`
	KeygroupWheelToLfo         string      `xml:"KeygroupWheelToLfo"`
	KeygroupAftertouchToFilter string      `xml:"KeygroupAftertouchToFilter"`
	QLinkAssignments           string      `xml:"QLinkAssignments"`
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
