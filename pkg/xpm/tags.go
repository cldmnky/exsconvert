package xpm

// XPMTag contains XML tag constants for MPC XPM format
// Based on ConvertWithMoss MPCKeygroupTag.java
type XPMTag struct{}

// Root level elements
const (
	Root               = "MPCVObject"
	RootProgram        = "Program"
	RootVersion        = "Version"
	VersionFileVersion = "File_Version"
	VersionPlatform    = "Platform"
)

// Program level elements
const (
	ProgramName               = "ProgramName"
	ProgramInstruments        = "Instruments"
	ProgramPads               = "ProgramPads"
	ProgramPadNoteMap         = "PadNoteMap"
	ProgramNumKeygroups       = "KeygroupNumKeygroups"
	ProgramPitchBendRange     = "KeygroupPitchBendRange"
	ProgramWheelToLfo         = "KeygroupWheelToLfo"
	ProgramAftertouchToFilter = "KeygroupAftertouchToFilter"
)

// Instrument level elements
const (
	InstrumentsInstrument         = "Instrument"
	InstrumentLowNote             = "LowNote"
	InstrumentHighNote            = "HighNote"
	InstrumentIgnoreBaseNote      = "IgnoreBaseNote"
	InstrumentZonePlay            = "ZonePlay"
	InstrumentTriggerMode         = "TriggerMode"
	InstrumentOneShot             = "OneShot"
	InstrumentLayers              = "Layers"
	InstrumentFilterType          = "FilterType"
	InstrumentFilterCutoff        = "Cutoff"
	InstrumentFilterResonance     = "Resonance"
	InstrumentFilterEnvAmt        = "FilterEnvAmt"
	InstrumentVelocityToFilter    = "VelocityToFilter"
	InstrumentFilterAttack        = "FilterAttack"
	InstrumentFilterHold          = "FilterHold"
	InstrumentFilterDecay         = "FilterDecay"
	InstrumentFilterSustain       = "FilterSustain"
	InstrumentFilterRelease       = "FilterRelease"
	InstrumentFilterAttackCurve   = "FilterAttackCurve"
	InstrumentFilterDecayCurve    = "FilterDecayCurve"
	InstrumentFilterReleaseCurve  = "FilterReleaseCurve"
	InstrumentVelocityToAmpAmount = "VelocitySensitivity"
	InstrumentVolumeAttack        = "VolumeAttack"
	InstrumentVolumeHold          = "VolumeHold"
	InstrumentVolumeDecay         = "VolumeDecay"
	InstrumentVolumeSustain       = "VolumeSustain"
	InstrumentVolumeRelease       = "VolumeRelease"
	InstrumentVolumeAttackCurve   = "VolumeAttackCurve"
	InstrumentVolumeDecayCurve    = "VolumeDecayCurve"
	InstrumentVolumeReleaseCurve  = "VolumeReleaseCurve"
	InstrumentPitchAttack         = "PitchAttack"
	InstrumentPitchHold           = "PitchHold"
	InstrumentPitchDecay          = "PitchDecay"
	InstrumentPitchSustain        = "PitchSustain"
	InstrumentPitchRelease        = "PitchRelease"
	InstrumentPitchAttackCurve    = "PitchAttackCurve"
	InstrumentPitchDecayCurve     = "PitchDecayCurve"
	InstrumentPitchReleaseCurve   = "PitchReleaseCurve"
	InstrumentPitchEnvAmount      = "PitchEnvAmount"
)

// Layer level elements
const (
	LayersLayer             = "Layer"
	LayerSampleName         = "SampleName"
	LayerActive             = "Active"
	LayerVolume             = "Volume"
	LayerPan                = "Pan"
	LayerPitch              = "Pitch"
	LayerCoarseTune         = "TuneCoarse"
	LayerFineTune           = "TuneFine"
	LayerVelStart           = "VelStart"
	LayerVelEnd             = "VelEnd"
	LayerSampleStart        = "SampleStart"
	LayerSampleEnd          = "SampleEnd"
	LayerLoopStart          = "LoopStart"
	LayerLoopEnd            = "LoopEnd"
	LayerLoopCrossfade      = "LoopCrossfadeLength"
	LayerLoopTune           = "LoopTune"
	LayerRootNote           = "RootNote"
	LayerKeyTrack           = "KeyTrack"
	LayerSampleFile         = "SampleFile"
	LayerSliceIndex         = "SliceIndex"
	LayerDirection          = "Direction"
	LayerOffset             = "Offset"
	LayerSliceStart         = "SliceStart"
	LayerSliceEnd           = "SliceEnd"
	LayerSliceLoop          = "SliceLoop"
	LayerSliceLoopStart     = "SliceLoopStart"
	LayerSliceLoopCrossfade = "SliceLoopCrossFadeLength"
	LayerSliceTailPosition  = "SliceTailPosition"
	LayerSliceTailLength    = "SliceTailLength"
	LayerPitchRandom        = "PitchRandom"
	LayerVolumeRandom       = "VolumeRandom"
	LayerPanRandom          = "PanRandom"
	LayerOffsetRandom       = "OffsetRandom"
)

// Attributes
const (
	ProgramType      = "type"
	InstrumentNumber = "number"
	LayerNumber      = "number"
)

// Values
const (
	TypeKeygroup = "Keygroup"
	TypeDrum     = "Drum"
	True         = "True"
	False        = "False"
)
