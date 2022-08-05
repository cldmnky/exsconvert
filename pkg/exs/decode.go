package exs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"k8s.io/klog"
)

const (
	headerChunk  = uint32(0x00)
	zoneChunk    = uint32(0x01)
	groupChunk   = uint32(0x02)
	sampleChunk  = uint32(0x03)
	optionsChunk = uint32(0x04)
)

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
}

type Zone struct {
	ExsZone
	Name            string
	Pitch           bool
	OneShot         bool
	Reverse         bool
	VelocityRangeOn bool
	LoopOn          bool
	LoopEqualPower  bool
}

type Group struct {
	ExsGroup
	Name  string
	Decay bool
}

type Sample struct {
	ExsSample
	Name     string
	FileName string
	Path     string
}

type ExsHeader struct {
	_     [4]byte
	Size  uint32
	_     [8]byte
	Magic [4]byte
}

type ExsChunkHeader struct {
	Signature uint32
	Size      uint32
	Magic     [4]byte
}

func NewFromFile(fileName string) (*EXS, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	// basename of filename
	name := strings.TrimSuffix(path.Base(fileName), ".exs")
	return NewFromReader(bytes.NewReader(file), name)
}

// NewFromReader creates a new EXS from a reader.
func NewFromReader(r *bytes.Reader, name string) (*EXS, error) {
	exs := &EXS{
		Name: name,
	}
	header, err := exs.readHeader(r)
	if err != nil {
		return nil, err
	}
	if string(header.Magic[:]) != "SOBT" &&
		string(header.Magic[:]) != "SOBJ" &&
		string(header.Magic[:]) != "TBOS" &&
		string(header.Magic[:]) != "JBOS" {
		return nil, errors.New("not an exs file")
	}
	klog.V(5).Infof("Magic: %s", header.Magic)
	if bytes.Equal(header.Magic[:], []byte("SOBT")) {
		exs.BigEndian = true
	}
	if bytes.Equal(header.Magic[:], []byte("SOBJ")) {
		exs.BigEndian = true
	}
	header, err = exs.readHeader(r)
	if err != nil {
		return nil, err
	}

	// determine if the file is size expanded
	// by checking the size of the header
	if header.Size > 0x8000 {
		klog.V(5).Infof("Size expanded file")
		exs.IsSizeExpanded = true
	}
	exs.Size, err = exs.readSize(r)
	if err != nil {
		return nil, err
	}
	klog.V(5).Infof("Size: %d", exs.Size)

	err = exs.readChunks(r)
	if err != nil {
		return nil, err
	}
	return exs, nil
}

// readSize reads the size of the exs file.
func (exs *EXS) readSize(r *bytes.Reader) (int, error) {
	// get the size of the reader
	size, err := r.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}
	return int(size), nil
}

// readHeader reads the header of the exs file.
func (exs *EXS) readHeader(r *bytes.Reader) (*ExsHeader, error) {
	var header ExsHeader
	if !exs.BigEndian {
		err := binary.Read(r, binary.LittleEndian, &header)
		if err != nil {
			return nil, err
		}
	} else {
		err := binary.Read(r, binary.BigEndian, &header)
		if err != nil {
			return nil, err
		}
	}
	return &header, nil
}

// getZonesByKeyRange returns the zones that are in the given key range.
func (exs *EXS) GetZonesByKeyRanges(zonesPerRange int) []map[string][]*Zone {
	var zones map[string][]*Zone

	for _, zone := range exs.Zones {
		hasSamples := false

		if zone.SampleIndex == -1 {
			continue
		}

		if zones == nil {
			zones = make(map[string][]*Zone)
		}
		// check if the zone has samples
		for s, _ := range exs.Samples {
			if int(zone.SampleIndex) == s {
				hasSamples = true
				break
			}
		}

		if hasSamples {
			z := fmt.Sprintf("%d-%d", zone.KeyLow, zone.KeyHigh)
			zones[z] = append(zones[z], zone)
		}
	}
	// sort zone map by keylow
	for _, z := range zones {
		sort.Slice(z, func(i, j int) bool {
			return z[i].VelLow < z[j].VelLow
		})
	}

	var ranges []map[string][]*Zone
	for s, zone := range zones {
		for i := 0; i < len(zone); i += zonesPerRange {
			end := i + zonesPerRange
			if end > len(zone) {
				end = len(zone)
			}
			ranges = append(ranges, map[string][]*Zone{s: zone[i:end]})
		}
	}
	return ranges
}

// GetGroups returns the groups in the exs file.
func (exs *EXS) GetGroups() []*Group {
	groups := make([]*Group, 0)

	for _, group := range exs.Groups {
		hasZones := false
		klog.V(2).Infof("Group: %s, ID: %d, Select number: %d, Select group: %d", group.Name, group.ID, group.SelectNumber, group.SelectGroup)
		for _, zone := range exs.Zones {
			klog.V(5).Infof("Zone: %s, Sample: %d, Group Index: %d", zone.Name, zone.SampleIndex, zone.GroupIndex)
			if zone.GroupIndex == int32(group.ID) {
				klog.V(2).Info("-----> Zone has group")
				if zone.SampleIndex != -1 && len(exs.Samples) > int(zone.SampleIndex) {
					hasZones = true
					break
				}
			}
		}
		if hasZones {
			groups = append(groups, group)
		}

	}

	if len(groups) == 0 {
		return []*Group{
			{
				Name: exs.Name,
				ExsGroup: ExsGroup{
					ID:           0,
					SelectNumber: 0,
					SelectGroup:  -1,
				},
			},
		}
	}
	klog.V(5).Infof("-----> Groups: %v", groups)
	return groups
}

// readChunks
func (exs *EXS) readChunks(r *bytes.Reader) error {
	i := 0
	// read until end of reader
	for i+84 < exs.Size {
		r.Seek(int64(i), io.SeekStart)
		header, err := exs.readChunkHeader(r)
		if err != nil {
			return err
		}
		chunkType := header.Signature & 0x0F000000 >> 24
		switch chunkType {
		case headerChunk:
			klog.V(5).Infof("Chunk: %d (exs header chunk), size: %d", chunkType, header.Size)
		case zoneChunk:
			if header.Size < 110 {
				return errors.New("invalid zone chunk size")
			}
			zone, err := exs.readZone(r, i)
			if err != nil {
				return err
			}
			exs.Zones = append(exs.Zones, zone)
			klog.V(2).Infof("Zone: %s, size: %d, keyLow: %d, keyHigh: %d, sample: %d", zone.Name, header.Size, zone.KeyLow, zone.KeyHigh, zone.SampleIndex)
		case groupChunk:
			klog.V(5).Infof("Exs chunk type: %d (group), size: %d",
				chunkType,
				header.Size)
			group, err := exs.readGroup(r, i)
			if err != nil {
				return err
			}
			exs.Groups = append(exs.Groups, group)
		case sampleChunk:
			klog.V(5).Infof("Exs chunk type: %d (sample), size: %d", chunkType, header.Size)
			if header.Size != 336 && header.Size != 592 && header.Size != 600 {
				return errors.New("invalid sample chunk size")
			}
			sample, err := exs.readSample(r, i)
			if err != nil {
				return err
			}
			exs.Samples = append(exs.Samples, sample)
		case optionsChunk:
			klog.V(5).Infof("Exs chunk type: %d (options), size: %d", chunkType, header.Size)
			exsParams, err := exs.readParams(r, i)
			if err != nil {
				return err
			}
			params := NewParamsFromExsParams(exsParams)
			exs.Params = params
		default:
			klog.V(5).Infof("Exs chunk type: %d (unknown)", chunkType)
		}
		i = i + int(header.Size) + 84
	}
	err := exs.ReadSequences()
	if err != nil {
		return err
	}
	err = exs.ConvertSeqNumbers()
	if err != nil {
		return err
	}

	klog.V(2).Infof("Exs %s contains %d groups, %d zones, %d samples", exs.Name, len(exs.Groups), len(exs.Zones), len(exs.Samples))

	return nil
}

// readChunkHeader reads the header of the exs file.
func (exs *EXS) readChunkHeader(r *bytes.Reader) (*ExsChunkHeader, error) {
	var header ExsChunkHeader
	if !exs.BigEndian {
		err := binary.Read(r, binary.LittleEndian, &header)
		if err != nil {
			return nil, err
		}
	} else {
		err := binary.Read(r, binary.BigEndian, &header)
		if err != nil {
			return nil, err
		}
	}
	return &header, nil
}

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
	_             [43]byte
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

func (exs *EXS) readZone(reader *bytes.Reader, pos int) (*Zone, error) {
	_, err := reader.Seek(int64(pos), io.SeekStart)
	if err != nil {
		return nil, err
	}
	var exsZone ExsZone
	if exs.BigEndian {
		err := binary.Read(reader, binary.BigEndian, &exsZone)
		if err != nil {
			return nil, err
		}
	} else {
		err := binary.Read(reader, binary.LittleEndian, &exsZone)
		if err != nil {
			return nil, err
		}
	}
	zone := &Zone{
		ExsZone:         exsZone,
		Name:            getString64(exsZone.Name),
		Pitch:           exsZone.Opts&0x02 == 0,
		OneShot:         exsZone.Opts&0x01 != 0,
		Reverse:         exsZone.Opts&0x04 != 0,
		VelocityRangeOn: exsZone.Opts&0x08 != 0,
		LoopOn:          exsZone.LoopOpts&0x01 != 0,
		LoopEqualPower:  exsZone.LoopOpts&0x02 != 0,
	}
	return zone, nil
}

type ExsGroup struct {
	_         [8]byte
	ID        uint32 // 8
	_         [8]byte
	Name      [64]byte // 20
	Volume    int8     // 84
	Pan       int8     // 85
	Polyphony int8     // 86
	Decay     int8     // 87
	Exclusive int8     // 88
	VelLow    uint8    // 89
	VelHigh   uint8    // 90
	_         [9]byte
	DecayTime uint32 // 100
	_         [20]byte
	Cutoff    int8 // 125
	_         [1]byte
	Resonance int8 // 127
	_         [12]byte
	// ADSR 2
	Attack2      int32 // 140
	Decay2       int32 // 144
	Sustain2     int32 // 148
	Release2     int32 // 152
	_            [9]byte
	SelectGroup  int32 // 164
	SelectType   uint8 // 168
	SelectNumber uint8 // 169
	SelectHigh   uint8 // 170
	SelectLow    uint8 // 171
	_            [12]byte
	// ADSR 1
	Attack1  int32 // 184
	Decay1   int32 // 188
	Sustain1 int32 // 192
	Release1 int32 // 196
}

func (exs *EXS) readGroup(reader *bytes.Reader, pos int) (*Group, error) {
	_, err := reader.Seek(int64(pos), io.SeekStart)
	if err != nil {
		return nil, err
	}
	var exsGroup ExsGroup
	if exs.BigEndian {
		err := binary.Read(reader, binary.BigEndian, &exsGroup)
		if err != nil {
			return nil, err
		}
	} else {
		err := binary.Read(reader, binary.LittleEndian, &exsGroup)
		if err != nil {
			return nil, err
		}
	}
	group := &Group{
		ExsGroup: exsGroup,
		Name:     getString64(exsGroup.Name),
		Decay:    exsGroup.Decay&0x40 != 0,
	}
	klog.V(5).Infof("Group: name: %s", string(group.Name[:]))
	return group, nil
}

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

// readSampole reads the sample data.
func (exs *EXS) readSample(reader *bytes.Reader, pos int) (*Sample, error) {
	_, err := reader.Seek(int64(pos), io.SeekStart)
	if err != nil {
		return nil, err
	}
	var sample ExsSample
	if exs.BigEndian {
		err := binary.Read(reader, binary.BigEndian, &sample)
		if err != nil {
			return nil, err
		}
	} else {
		err := binary.Read(reader, binary.LittleEndian, &sample)
		if err != nil {
			return nil, err
		}
	}
	klog.V(5).Infof("Sample: name: %s, filename: %s, path: %s", string(sample.Name[:]), string(sample.FileName[:]), string(sample.Path[:]))
	return &Sample{
		ExsSample: sample,
		Name:      getString64(sample.Name),
		FileName:  getString256(sample.FileName),
		Path:      getString256(sample.Path),
	}, nil
}

type ExsParams struct {
	_      [8]byte
	ID     uint32 // 8
	_      [8]byte
	Name   [64]byte // 20
	_      [4]byte
	Keys   [100]uint8 // 88
	Values [100]int16
}

func (exs *EXS) readParams(reader *bytes.Reader, pos int) (*ExsParams, error) {
	_, err := reader.Seek(int64(pos), io.SeekStart)
	if err != nil {
		return nil, err
	}
	var params ExsParams
	if exs.BigEndian {
		err := binary.Read(reader, binary.BigEndian, &params)
		if err != nil {
			return nil, err
		}
	} else {
		err := binary.Read(reader, binary.LittleEndian, &params)
		if err != nil {
			return nil, err
		}
	}
	klog.V(5).Infof("Params name: %s,  %+v", string(params.Name[:]), params)
	return &params, nil
}

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
	// Time
	TimeVia   int16 // [0,100]cent default:0
	TimeCurve int16 // [-99,99] default:0
	// Envelope 2
	Env2Attack       int16 // [0,127] for (0~10000ms), <=env1_attack_via_vel @see: ym_exs_time_to_second()
	Env2AttackViaVel int16 // [0,127] for (0~10000ms), >=env1_attack @see: ym_exs_time_to_second()
	Env2Decay        int16 // [0,127] for (0~10000ms), @see: ym_exs_time_to_second()
	Env2Sustain      int16 // [0,127] for (0,100%), liner mapping, default:127
	Env2Release      int16 // [0,127] for (0~10000ms), @see: ym_exs_time_to_second()
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

func NewParamsFromExsParams(exsParams *ExsParams) *Params {
	params := &Params{}
	i := 0
	for i < 100 {
		key := exsParams.Keys[i]
		value := exsParams.Values[i]
		switch key {
		case 7:
			params.OutputVolume = value
		case 8:
			params.KeyScale = value
		case 3:
			params.PitchBendUp = value
		case 4:
			params.PitchBendDown = value
		case 10:
			params.MonoMode = value
		case 5:
			params.Voices = value
		case 171:
			params.Unison = value != 0 // 0:off 1:on
		case 45:
			params.Transpose = value
		case 14:
			params.CoarseTune = value
		case 15:
			params.FineTune = value
		case 20:
			params.GlideTime = value
		case 44:
			params.FilterOn = value != 0 // 0:off 1:on
		case 46:
			params.FilterViaKey = value
		case 60:
			params.Lfo1DecayDelay = value
		case 61:
			params.Lfo1Rate = value
		case 62:
			params.Lfo1Waveform = value
		case 63:
			params.Lfo2Rate = value
		case 64:
			params.Lfo2Waveform = value
		case 72:
			params.Pitcher = value
		case 73:
			params.PitcherViaVel = value
		case 75:
			params.FilterDrive = value
		case 76:
			params.Env1Attack = value
		case 77:
			params.Env1AttackViaVel = value
		case 78:
			params.Env1Decay = value
		case 79:
			params.Env1Sustain = value
		case 80:
			params.Env1Release = value
		case 81:
			params.Env2Sustain = value
		case 82:
			params.Env2Attack = value
		case 83:
			params.Env2AttackViaVel = value
		case 84:
			params.Env2Decay = value
		case 85:
			params.Env2Release = value
		case 89:
			params.LevelViaVel = value
		case 90:
			params.LevelFixed = value
		case 91:
			params.TimeCurve = value
		case 92:
			params.TimeVia = value
		case 95:
			params.VelocityOffset = value
		case 97:
			params.VelocityXFade = value
		case 98:
			params.RandomDetune = value
		case 163:
			params.SampleSelectRandom = value
		case 164:
			params.VelocityRandom = value
		case 165:
			params.VelocityXFadeType = value
		case 166:
			params.CoarseTune = value
		case 167:
			params.Lfo3Rate = value
		case 170:
			params.FilterFat = value != 0 // 0:off 1:on
		case 172:
			params.HoldVia = value
		case 173:
			params.Destination[0] = value
		case 174:
			params.Source[0] = value
		case 175:
			params.Via[0] = value
		case 176:
			params.Amount[0] = value
		case 177:
			params.AmountVia[0] = value
		case 178:
			params.InvertVia[0] = value != 0
		case 179:
			params.Destination[1] = value
		case 180:
			params.Source[1] = value
		case 181:
			params.Via[1] = value
		case 182:
			params.Amount[1] = value
		case 183:
			params.AmountVia[1] = value
		case 184:
			params.InvertVia[1] = value != 0
		case 185:
			params.Destination[2] = value
		case 186:
			params.Source[2] = value
		case 187:
			params.Via[2] = value
		case 188:
			params.Amount[2] = value
		case 189:
			params.AmountVia[2] = value
		case 190:
			params.InvertVia[2] = value != 0
		case 191:
			params.Destination[3] = value
		case 192:
			params.Source[3] = value
		case 193:
			params.Via[3] = value
		case 194:
			params.Amount[3] = value
		case 195:
			params.AmountVia[3] = value
		case 196:
			params.InvertVia[3] = value != 0
		case 197:
			params.Destination[4] = value
		case 198:
			params.Source[4] = value
		case 199:
			params.Via[4] = value
		case 200:
			params.Amount[4] = value
		case 201:
			params.AmountVia[4] = value
		case 202:
			params.InvertVia[4] = value != 0
		case 203:
			params.Destination[5] = value
		case 204:
			params.Source[5] = value
		case 205:
			params.Via[5] = value
		case 206:
			params.Amount[5] = value
		case 207:
			params.AmountVia[5] = value
		case 208:
			params.InvertVia[5] = value != 0
		case 209:
			params.Destination[6] = value
		case 210:
			params.Source[6] = value
		case 211:
			params.Via[6] = value
		case 212:
			params.Amount[6] = value
		case 213:
			params.AmountVia[6] = value
		case 214:
			params.InvertVia[6] = value != 0
		case 215:
			params.Destination[7] = value
		case 216:
			params.Source[7] = value
		case 217:
			params.Via[7] = value
		case 218:
			params.Amount[7] = value
		case 219:
			params.AmountVia[7] = value
		case 220:
			params.InvertVia[7] = value != 0
		case 221:
			params.Destination[8] = value
		case 222:
			params.Source[8] = value
		case 223:
			params.Via[8] = value
		case 224:
			params.Amount[8] = value
		case 225:
			params.AmountVia[8] = value
		case 226:
			params.InvertVia[8] = value != 0
		case 227:
			params.Destination[9] = value
		case 228:
			params.Source[9] = value
		case 229:
			params.Via[9] = value
		case 230:
			params.Amount[9] = value
		case 231:
			params.AmountVia[9] = value
		case 232:
			params.InvertVia[9] = value != 0
		case 233:
			params.Invert[0] = value != 0
		case 234:
			params.Invert[1] = value != 0
		case 235:
			params.Invert[2] = value != 0
		case 236:
			params.Invert[3] = value != 0
		case 237:
			params.Invert[4] = value != 0
		case 238:
			params.Invert[5] = value != 0
		case 239:
			params.Invert[6] = value != 0
		case 240:
			params.Invert[7] = value != 0
		case 241:
			params.Invert[8] = value != 0
		case 242:
			params.Invert[9] = value != 0
		case 243:
			params.FilterType = value
		case 244:
			params.Bypass[0] = value != 0
		case 245:
			params.Bypass[1] = value != 0
		case 246:
			params.Bypass[2] = value != 0
		case 247:
			params.Bypass[3] = value != 0
		case 248:
			params.Bypass[4] = value != 0
		case 249:
			params.Bypass[5] = value != 0
		case 250:
			params.Bypass[6] = value != 0
		case 251:
			params.Bypass[7] = value != 0
		case 252:
			params.Bypass[8] = value != 0
		case 253:
			params.Bypass[9] = value != 0

		default:
			klog.V(5).Infof("unknown parameter %d", i)
		}

		i++
	}
	//klog.Infof("params: %+v", params)
	return params
}
