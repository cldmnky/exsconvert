package exs

// ============================================================================
// Public Types
// ============================================================================

// EXS represents a parsed Logic Pro EXS24 sampler instrument file.
type EXS struct {
	Name           string
	BigEndian      bool
	IsSizeExpanded bool
	Size           int
	Zones          []*Zone
	Groups         []*Group
	Samples        []*Sample
	Params         *Params
	Sequences      [][]int32
	Instrument     *ExsInstrument
}

// Zone represents a zone in the EXS24 file.
type Zone struct {
	ExsZone
	Name            string
	Pitch           bool
	OneShot         bool
	Reverse         bool
	VelocityRangeOn bool
	LoopOn          bool
	LoopEqualPower  bool
	LoopEndRelease  bool
	PlayMode        uint8
	HasOutput       bool
}

// Group represents a group in the EXS24 file.
type Group struct {
	ExsGroup
	Name  string
	Decay bool
}

// Sample represents a sample in the EXS24 file.
type Sample struct {
	ExsSample
	Name     string
	FileName string
	Path     string
}

// Params represents the parsed parameters from ExsParams.
type Params struct {
	// Global
	OutputVolume  int16 ///< [-60,0]dB default:-6
	KeyScale      int16 ///< [-24,24]dB default:0
	PitchBendUp   int16 ///< [-12,12]cent default:0
	PitchBendDown int16 ///< [-12,12]cent default:0
	MonoMode      int16 ///< 0:off 1:on 2: legato
	Voices        int16 ///< [1,64] default:16
	Unison        bool
	// Start via Vel
	Transpose     int16
	CoarseTune    int16
	FineTune      int16
	GlideTime     int16
	Pitcher       int16
	PitcherViaVel int16
	// Pitch Mod. Wheel
	FilterOn   bool
	FilterType int16 // 5:HP12 2:LP24 1:LP18 0:LP12 3:LP6 4:BP12  (default:0)
	FilterFat  bool
	// Filter Mod. Wheel
	FilterCutoff    int16 // [0,1000] for (0~100%) default:1000
	FilterResonance int16 // [0,1000] for (0~100%) default:0
	FilterDrive     int16 // [0,1000] for (0~100%) default:100
	FilterViaKey    int16 // [0,1000] for (0~100%) default:0
	// Filter ADSR via Vel
	LevelFixed  int16 // [-96,0] for (-48dB~0dB) default:0   <=level_via_vel
	LevelViaVel int16 // [-96,0] for (-48dB~0dB) default:0   >=level_fixed;
	// Tremolo
	Lfo1DecayDelay int16 // [-9999,9999]ms negative:decay positive:delay
	Lfo1Rate       int16 // [0,127] default:98(4.8Hz)
	Lfo1Waveform   int16 // [0,6]
	Lfo2Waveform   int16 // [0,6]
	Lfo2Rate       int16 // [0,127] default:34(DC)
	Lfo3Rate       int16 // [0,127] default:98(4.8Hz)
	// Envelope
	Env1Attack       int16 // [0,127] for (0~10000ms), <=env1_attack_via_vel @see: ym_exs_time_to_second()
	Env1AttackViaVel int16 // [0,127] for (0~10000ms), >=env1_attack @see: ym_exs_time_to_second()
	Env1Decay        int16 // [0,127] for (0~10000ms), @see: ym_exs_time_to_second
	Env1Sustain      int16 // [0,127] for (0,100%), liner mapping, default:0
	Env1Release      int16 // [0,127] for (0~10000ms), @see: ym_exs_time_to_second
	Env1AttackCurve  int16 // [-99,99] envelope attack curve, signed value
	// Time
	TimeVia   int16 // [0,100]cent default:0
	TimeCurve int16 // [-99,99] default:0
	// Envelope 2
	Env2Attack       int16 // [0,127] for (0~10000ms), <=env1_attack_via_vel @see: ym_exs_time_to_second()
	Env2AttackViaVel int16 // [0,127] for (0~10000ms), >=env1_attack @see: ym_exs_time_to_second()
	Env2Decay        int16 // [0,127] for (0~10000ms), @see: ym_exs_time_to_second()
	Env2Sustain      int16 // [0,127] for (0,100%), liner mapping, default:127
	Env2Release      int16 // [0,127] for (0~10000ms), @see: ym_exs_time_to_second()
	Env2AttackCurve  int16 // [-99,99] envelope attack curve, signed value
	// Velocity
	VelocityOffset     int16 // [-127,127] default:0
	VelocityRandom     int16 // [0,127] default:0
	VelocityXFade      int16 // [0,127] default:0
	VelocityXFadeType  int16 //  0:dB lin 1:liner 2:Eq.Pow
	CoarseTuneRemote   int16 // [-1,127] default:-1(OFF)
	HoldVia            int16 // [-17,120	] default:64  @see: ym_exs_src_via_t
	SampleSelectRandom int16 // [0,127] default:0
	RandomDetune       int16 // [0,50]cent default:0
	// Modulator
	Destination [10]int16 // [0,30] @see: ym_exs_dest_t
	Source      [10]int16 // [-17,120]  @see: ym_exs_src_via_t
	Via         [10]int16 // [-17,120]  @see: ym_exs_src_via_t
	Amount      [10]int16 // [-1000,1000] for (-100%~100%)  <=amount_via
	AmountVia   [10]int16 // [-1000,1000] for (-100%~100%)  >=amount
	Invert      [10]bool
	InvertVia   [10]bool
	Bypass      [10]bool
}

// ============================================================================
// Binary Format Types (internal representations)
// ============================================================================

// ExsHeader represents the file header structure.
type ExsHeader struct {
	_     [4]byte
	Size  uint32
	_     [8]byte
	Magic [4]byte
}

// ExsChunkHeader represents a chunk header in the EXS file.
type ExsChunkHeader struct {
	Signature uint32
	Size      uint32
	Magic     [4]byte
}

// ExsZone represents the binary structure of a zone.
type ExsZone struct {
	_             [8]byte
	ID            uint32 // 8
	_             [8]byte
	Name          [64]byte // 20
	Opts          uint8    // 84
	Key           uint8    // 85
	FineTuning    int8     // 86
	Pan           int8     // 87
	Volume        int8     // 88
	Scale         int8     // 89
	KeyLow        int8     // 90
	KeyHigh       int8     // 91
	_             [1]byte
	VelLow        int8 // 93
	VelHigh       int8 // 94
	_             [1]byte
	SampleStart   int32  // 96
	SampleEnd     int32  // 100
	LoopStart     int32  // 104
	LoopEnd       int32  // 108
	LoopCrossfade int32  // 112
	LoopTune      int8   // 116
	LoopOpts      uint32 // 117
	PlayMode      uint8  // 121
	_             [42]byte
	CoarseTuning  int8 // 164
	_             [1]byte
	Output        int8 // 166
	_             [5]byte
	GroupIndex    int32 // 172
	SampleIndex   int32 // 176
	_             [8]byte
	SampleFade    int32 // 188
	Offset        int32 // 192
}

// ExsGroup represents the binary structure of a group.
type ExsGroup struct {
	_         [8]byte
	ID        uint32 // 8
	_         [8]byte
	Name      [64]byte // 20
	Volume    int8     // 84
	Pan       int8     // 85
	Polyphony int8     // 86
	Decay     uint8    // 87
	Exclusive int8     // 88
	VelLow    uint8    // 89
	VelHigh   uint8    // 90
	_         [9]byte  //91
	DecayTime uint32   // 100
	_         [21]byte // 104
	Cutoff    int8     // 125
	_         [1]byte
	Resonance int8 // 127
	_         [12]byte
	// ADSR 2
	Attack2      int32 // 140
	Decay2       int32 // 144
	Sustain2     int32 // 148
	Release2     int32 // 152
	_            [1]byte
	Trigger      int8 // 157
	Output       int8 // 158
	_            [5]byte
	SelectGroup  int32 // 164  exsSequence
	SelectType   uint8 // 168 < 0:-- 1:Note 2:Group 3:Control 4:Bend 5:Midi Channel 6:Articulation ID 7:Tempo
	SelectNumber uint8 // 169
	SelectHigh   uint8 // 170
	SelectLow    uint8 // 171
	KeyLow       uint8 // 172
	KeyHigh      uint8 // 173
	_            [6]byte
	Hold2        int32 // 180
	// ADSR 1
	Attack1  int32 // 184
	Decay1   int32 // 188
	Sustain1 int32 // 192
	Release1 int32 // 196
}

// ExsSample represents the binary structure of a sample.
type ExsSample struct {
	_        [8]byte
	ID       uint32 // 8
	_        [8]byte
	Name     [64]byte // 20
	_        [4]byte
	Length   int32 // 88
	Rate     int32 // 92
	BitDepth uint8 // 96
	_        [15]byte
	Type     int32 // 112
	_        [48]byte
	Path     [256]byte // 164
	FileName [256]byte // 420
}

// ExsParams represents the binary structure of parameters.
type ExsParams struct {
	_      [8]byte
	ID     uint32 // 8
	_      [8]byte
	Name   [64]byte // 20
	_      [4]byte
	Keys   [100]uint8 // 88
	Values [100]int16
}

// ExsInstrument represents the binary structure of instrument metadata.
type ExsInstrument struct {
	_          [4]byte
	NumZones   uint32
	NumGroups  uint32
	NumSamples uint32
}
