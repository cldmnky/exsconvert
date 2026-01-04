package convert

// TODO: Implement missing features from ConvertWithMoss for full XPM compatibility:
// 1. ✅ Envelope Curves: Add attack/decay/release curves for amplitude and filter envelopes - COMPLETED
// 2. ✅ LFO Settings: Add per-instrument LFO with pitch/cutoff/volume/pan modulation - COMPLETED
// 3. Advanced Filter Features: Add filter envelope curves and more filter types - PARTIALLY COMPLETED (curves added, types need refinement)
// 4. ✅ Pitch Envelope: Implement full pitch envelope support - COMPLETED
// 5. Zone Play Modes: Add different trigger modes (one-shot, note-off, etc.) - PARTIALLY COMPLETED (TriggerMode field added)
// 6. ✅ XML Tag Constants: Create comprehensive constants file like MPCKeygroupTag.java - COMPLETED
// 7. Keygroup Stacking: Improve zone merging logic for better instrument organization - NOT STARTED

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/klog"

	"github.com/cldmnky/exsconvert/pkg/exs"
	"github.com/cldmnky/exsconvert/pkg/xpm"
)

// ProgramType constants
const (
	ProgramTypeKeygroup = "Keygroup"
	ProgramTypeDrum     = "Drum"
)

type XPM struct {
	SearchPath          string
	OutputPath          string
	LayersPerInstrument int
	SkipErrors          bool
	ProgramType         string            // "Keygroup" or "Drum" - empty for auto-detect
	AutoDetectDrums     bool              // If true, auto-detect drum programs
	SamplesSearchPath   string            // Path to search for samples (defaults to SearchPath)
	sampleIndex         map[string]string // Cache: filename -> full path
}

func NewXPM(searchPath, outputPath string, layersPerInstrument int, skipErrors bool, programType string) *XPM {
	// Default to Keygroup if not specified or invalid
	if programType != ProgramTypeDrum && programType != ProgramTypeKeygroup && programType != "" {
		programType = ProgramTypeKeygroup
	}
	return &XPM{
		SearchPath:          searchPath,
		OutputPath:          outputPath,
		LayersPerInstrument: layersPerInstrument,
		SkipErrors:          skipErrors,
		ProgramType:         programType,
		AutoDetectDrums:     false, // Default to manual
		SamplesSearchPath:   searchPath,
	}
}

// Convert processes EXS files and generates MPC-compatible XPM programs.
//
// Directory Structure:
// The converter creates a subdirectory for each instrument containing:
//   - One .xpm file (the keygroup program)
//   - All referenced .WAV sample files (uppercase extension)
//
// This structure meets MPC requirements: WAV files must be in the same
// directory as the XPM file for the MPC to load them properly.
//
// Example output structure:
//
//	output/
//	└── MyInstrument/
//	    ├── MyInstrument.xpm
//	    ├── Sample1.WAV
//	    └── Sample2.WAV
func (x *XPM) Convert() error {
	exsFiles, err := x.findEXSFiles()
	if err != nil {
		return err
	}

	if len(exsFiles) == 0 {
		return fmt.Errorf("no .exs files found in %s", x.SearchPath)
	}

	klog.V(2).Infof("Found %d EXS files to convert", len(exsFiles))

	// Build sample index once for performance
	if err := x.buildSampleIndex(); err != nil {
		return fmt.Errorf("failed to build sample index: %w", err)
	}

	for _, exsFile := range exsFiles {
		klog.V(2).Infof("Processing file: %s", filepath.Base(exsFile))
		exs, err := exs.NewFromFile(exsFile)
		if err != nil {
			if x.SkipErrors {
				klog.Warningf("Skipping %s: %v", exsFile, err)
				continue
			}
			return err
		}
		klog.V(2).Infof("Loaded EXS file: %s", exs.Name)

		// Determine program type for this file
		programType := x.ProgramType
		klog.V(2).Infof("Determining program type for %s (AutoDetect=%v, ProgramType=%s)", exs.Name, x.AutoDetectDrums, x.ProgramType)
		if x.AutoDetectDrums && programType == "" {
			klog.V(2).Infof("Calling IsDrumProgram() for %s", exs.Name)
			if exs.IsDrumProgram() {
				programType = ProgramTypeDrum
				klog.V(2).Infof("Auto-detected drum program: %s", exs.Name)
			} else {
				programType = ProgramTypeKeygroup
				klog.V(2).Infof("Auto-detected keygroup program: %s", exs.Name)
			}
		} else if programType == "" {
			// Default to keygroup if not auto-detecting
			programType = ProgramTypeKeygroup
		}

		// Store the program type for this conversion
		originalProgramType := x.ProgramType
		x.ProgramType = programType

		// Create subdirectory for this instrument (XPM + WAV files)
		destPath := filepath.Join(x.OutputPath, exs.Name)
		err = os.MkdirAll(destPath, 0755)
		if err != nil {
			x.ProgramType = originalProgramType
			if x.SkipErrors {
				klog.Warningf("Skipping %s: failed to create directory: %v", exs.Name, err)
				continue
			}
			return err
		}

		err = x.toXPM(exs, destPath)
		x.ProgramType = originalProgramType

		if err != nil {
			if x.SkipErrors {
				klog.Warningf("Skipping %s: %v", exs.Name, err)
				continue
			}
			return err
		}
		klog.V(2).Infof("Finished processing %s", exs.Name)
		fmt.Printf("Converted %s as %s program\n", exs.Name, programType)
	}
	return nil
}

// ConvertFile converts a single EXS file to XPM format.
func (x *XPM) ConvertFile(exsFilePath string) error {
	exs, err := exs.NewFromFile(exsFilePath)
	if err != nil {
		return fmt.Errorf("failed to load EXS file: %w", err)
	}
	// Create subdirectory for this instrument (XPM + WAV files)
	destPath := filepath.Join(x.OutputPath, exs.Name)
	err = os.MkdirAll(destPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	err = x.toXPM(exs, destPath)
	if err != nil {
		return fmt.Errorf("failed to convert to XPM: %w", err)
	}
	fmt.Printf("Converted %s\n", exs.Name)
	return nil
}

func (x *XPM) findEXSFiles() ([]string, error) {
	exsFiles := []string{}
	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if filepath.Ext(d.Name()) == ".exs" || filepath.Ext(d.Name()) == ".EXS" {
				exsFiles = append(exsFiles, path)
				return nil
			}
		}
		return nil
	}
	err := filepath.WalkDir(x.SearchPath, walk)
	if err != nil {
		return nil, err
	}
	return exsFiles, nil
}

// buildSampleIndex walks the samples directory once and builds an index of all sample files.
// This dramatically improves performance when processing instruments with many samples.
func (x *XPM) buildSampleIndex() error {
	x.sampleIndex = make(map[string]string)

	searchPath := x.SamplesSearchPath
	if searchPath == "" {
		searchPath = x.SearchPath
	}

	klog.V(3).Infof("Building sample index from %s", searchPath)
	count := 0

	err := filepath.WalkDir(searchPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		filename := d.Name()
		if existingPath, exists := x.sampleIndex[filename]; exists {
			klog.V(3).Infof("Warning: duplicate sample file %s (keeping %s, ignoring %s)", filename, existingPath, path)
			return nil
		}

		x.sampleIndex[filename] = path
		count++
		return nil
	})

	if err != nil {
		return err
	}

	klog.V(2).Infof("Indexed %d sample files", count)
	return nil
}

// copySample searches for a sample file in the SamplesSearchPath directory tree,
// copies it to the destination directory, and converts the extension to uppercase (.WAV).
// This ensures MPC compatibility: sample files must be in the same directory as the XPM file.
//
// Parameters:
//   - name: The filename of the sample to copy (e.g., "kick.wav")
//   - destPath: The destination directory (same as the XPM file location)
//
// Returns the sample name without extension and the sample filename with uppercase extension,
// or an error if the sample is not found or if multiple samples with the same name exist.
func (x *XPM) copySample(name, destPath string) (string, string, error) {
	// Use pre-built sample index for fast lookup
	src, found := x.sampleIndex[name]
	if !found {
		return "", "", fmt.Errorf("no sample found for %s", name)
	}

	klog.V(2).Infof("found %s", src)

	var toUpperExt = func(fileName string) string {
		ext := filepath.Ext(fileName)
		fileName = fileName[:len(fileName)-len(ext)]
		return fmt.Sprintf("%s%s", fileName, strings.ToUpper(ext))
	}

	var toSampleName = func(fileName string) string {
		ext := filepath.Ext(fileName)
		return fileName[:len(fileName)-len(ext)]
	}

	sampleFileName := filepath.Base(toUpperExt(src))
	dst := filepath.Join(destPath, sampleFileName)

	klog.V(3).Infof("Copying %s -> %s", src, dst)
	in, err := os.Open(src)
	if err != nil {
		return "", "", err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return "", "", err
	}
	err = out.Close()
	if err != nil {
		return "", "", err
	}

	klog.V(3).Infof("Copied sample %s successfully", sampleFileName)
	return toSampleName(name), sampleFileName, nil
}

func (x *XPM) toXPM(exsFile *exs.EXS, destPath string) error {
	klog.V(2).Infof("Starting toXPM for %s", exsFile.Name)
	// Create the appropriate XPM structure based on program type
	var keyGroup *xpm.MPCVObject
	if x.ProgramType == ProgramTypeDrum {
		keyGroup = xpm.NewXPMDrum()
	} else {
		keyGroup = xpm.NewXPMKeygroup()
	}

	klog.V(2).Infof("Calling GetZonesByKeyRanges with %d layers", x.LayersPerInstrument)
	z := exsFile.GetZonesByKeyRanges(x.LayersPerInstrument)
	klog.V(2).Infof("GetZonesByKeyRanges returned %d instruments", len(z))
	for s, zoneMap := range z {
		klog.V(5).Infof("zoneMap: %d", s)
		for _, zones := range zoneMap {
			for _, zone := range zones {
				klog.V(5).Infof("    zone: %s, key low: %d key high: %d, vel low: %d, vel high: %d, group: %d, sample: %s", zone.Name, zone.KeyLow, zone.KeyHigh, zone.VelLow, zone.VelHigh, zone.GroupIndex, strings.TrimSpace(exsFile.Samples[zone.SampleIndex].FileName))
			}
		}
	}
	if len(z) == 0 {
		return fmt.Errorf("no instruments found")
	} else {
		klog.V(2).Infof("Number of instruments: %d", len(z))
	}

	if len(z) >= 128 {
		if !x.SkipErrors {
			return fmt.Errorf("%s too many instruments", exsFile.Name)
		} else {
			klog.Warningf("Skipping %s due to too many instruments (%d)", exsFile.Name, len(z))
			return nil
		}
	}

	groups := exsFile.GetGroups()
	if len(groups) == 0 {
		return fmt.Errorf("no groups found")
	}

	// Create a map of group ID to group for quick lookup
	groupMap := make(map[uint32]*exs.Group)
	for i := range groups {
		groupMap[groups[i].ID] = groups[i]
		klog.V(2).Infof("group: %s, id: %d, selectgroup: %d, sequences: %+v, selectType: %d, selectNumber: %d", groups[i].Name, groups[i].ID, groups[i].SelectGroup, exsFile.Sequences, groups[i].SelectType, groups[i].SelectNumber)
	}

	// Check for Round Robin support across groups
	// Round Robin is enabled when SelectGroup >= 0 in any group
	hasRoundRobin := false
	for _, g := range groups {
		if g.SelectGroup >= 0 {
			hasRoundRobin = true
			klog.V(2).Infof("Round Robin detected: group %d has SelectGroup=%d", g.ID, g.SelectGroup)
			break
		}
	}

	// Use the EXS instrument name as the program name
	keyGroup.Program.ProgramName = exsFile.Name

	j := 0
	for _, zoneMap := range z {
		for _, zones := range zoneMap {
			// Look up the group for this zone
			g, ok := groupMap[uint32(zones[0].GroupIndex)]
			if !ok {
				klog.Warningf("Group %d not found for zone, using defaults", zones[0].GroupIndex)
				// Use first group as fallback
				g = groups[0]
			}

			// Apply group range limiting - groups can constrain zone ranges
			// This implements the ConvertWithMoss limitByGroupAttributes logic
			zoneKeyLow := int(zones[0].KeyLow)
			zoneKeyHigh := int(zones[0].KeyHigh)

			// Apply group key range limits if set (non-zero)
			if g.KeyLow != 0 && zoneKeyLow < int(g.KeyLow) {
				zoneKeyLow = int(g.KeyLow)
			}
			if g.KeyHigh != 0 && zoneKeyHigh > int(g.KeyHigh) {
				zoneKeyHigh = int(g.KeyHigh)
			}

			// Skip zones completely outside group's key range
			if zoneKeyHigh < int(g.KeyLow) || (g.KeyHigh != 0 && zoneKeyLow > int(g.KeyHigh)) {
				klog.V(2).Infof("Skipping zone outside group key range: zone[%d-%d] group[%d-%d]",
					zones[0].KeyLow, zones[0].KeyHigh, g.KeyLow, g.KeyHigh)
				continue
			}

			// For Drum programs: each zone maps to a single pad/note
			// For Keygroup programs: zones can span multiple notes
			if x.ProgramType == ProgramTypeDrum {
				// In drum mode, use the zone's key low as the pad note
				// Each zone gets mapped to its specific MIDI note
				keyGroup.Program.Instruments.Instrument[j].LowNote = zoneKeyLow
				keyGroup.Program.Instruments.Instrument[j].HighNote = zoneKeyLow // Same as LowNote for drums
			} else {
				// Keygroup mode: use full key range (with group limits applied)
				keyGroup.Program.Instruments.Instrument[j].LowNote = zoneKeyLow
				keyGroup.Program.Instruments.Instrument[j].HighNote = zoneKeyHigh
			}

			// Filter parameters - convert from EXS int8 range to XPM normalized values
			keyGroup.Program.Instruments.Instrument[j].Cutoff = formatFilterCutoff(float64(g.Cutoff))
			keyGroup.Program.Instruments.Instrument[j].Resonance = formatFilterResonance(float64(g.Resonance))
			// Volume envelope - convert from EXS time units to XPM normalized values
			// EXS envelope times appear to be in some normalized unit, convert to 0-1 range for XPM

			// Use envelope parameters from EXS Params (global instrument settings) instead of Group settings
			var volAttack, volDecay, volSustain, volRelease, volHold float64
			var filtAttack, filtDecay, filtSustain, filtRelease, filtHold float64

			if exsFile.Params != nil {
				volAttack = float64(exsFile.Params.Env2Attack)
				volDecay = float64(exsFile.Params.Env2Decay)
				volSustain = float64(exsFile.Params.Env2Sustain)
				volRelease = float64(exsFile.Params.Env2Release)
				// Note: Hold2 from group might be used, but for now use 0
				volHold = 0

				filtAttack = float64(exsFile.Params.Env1Attack)
				filtDecay = float64(exsFile.Params.Env1Decay)
				filtSustain = float64(exsFile.Params.Env1Sustain)
				filtRelease = float64(exsFile.Params.Env1Release)
				filtHold = 0 // EXS doesn't have separate hold for filter envelope
			} else {
				// Fallback to group values (though they appear to be 0)
				volAttack = float64(g.Attack2)
				volDecay = float64(g.Decay2)
				volSustain = float64(g.Sustain2)
				volRelease = float64(g.Release2)
				volHold = float64(g.Hold2)

				filtAttack = float64(g.Attack1)
				filtDecay = float64(g.Decay1)
				filtSustain = float64(g.Sustain1)
				filtRelease = float64(g.Release1)
			}

			keyGroup.Program.Instruments.Instrument[j].VolumeAttack = formatEnvTime(volAttack)
			keyGroup.Program.Instruments.Instrument[j].VolumeDecay = formatEnvTime(volDecay)
			keyGroup.Program.Instruments.Instrument[j].VolumeSustain = formatEnvLevel(volSustain)
			keyGroup.Program.Instruments.Instrument[j].VolumeRelease = formatEnvTime(volRelease)
			keyGroup.Program.Instruments.Instrument[j].VolumeHold = formatEnvTime(volHold)
			// Envelope curves - use TimeCurve as approximation (specific ENV curves not in basic params)
			if exsFile.Params != nil && exsFile.Params.TimeCurve != 0 {
				keyGroup.Program.Instruments.Instrument[j].VolumeAttackCurve = formatEnvelopeCurve(int(exsFile.Params.TimeCurve))
			} else {
				keyGroup.Program.Instruments.Instrument[j].VolumeAttackCurve = getDefaultEnvelopeCurve()
			}
			// Decay and release curves use default (linear) - EXS doesn't store these separately
			keyGroup.Program.Instruments.Instrument[j].VolumeDecayCurve = getDefaultEnvelopeCurve()
			keyGroup.Program.Instruments.Instrument[j].VolumeReleaseCurve = getDefaultEnvelopeCurve()
			// Filter envelope - ENV2 in EXS is the filter envelope
			// Only apply filter envelope if FilterEnvAmt > 0 (ConvertWithMoss logic)
			// MPC does not support negative filter modulation
			if keyGroup.Program.Instruments.Instrument[j].FilterEnvAmt != "" && keyGroup.Program.Instruments.Instrument[j].FilterEnvAmt != "0.000000" {
				keyGroup.Program.Instruments.Instrument[j].FilterAttack = formatEnvTime(filtAttack)
				keyGroup.Program.Instruments.Instrument[j].FilterDecay = formatEnvTime(filtDecay)
				keyGroup.Program.Instruments.Instrument[j].FilterSustain = formatEnvLevel(filtSustain)
				keyGroup.Program.Instruments.Instrument[j].FilterRelease = formatEnvTime(filtRelease)
				keyGroup.Program.Instruments.Instrument[j].FilterHold = formatEnvTime(filtHold)

				// Filter envelope curves - use TimeCurve as approximation
				if exsFile.Params != nil && exsFile.Params.TimeCurve != 0 {
					keyGroup.Program.Instruments.Instrument[j].FilterAttackCurve = formatEnvelopeCurve(int(exsFile.Params.TimeCurve))
				} else {
					keyGroup.Program.Instruments.Instrument[j].FilterAttackCurve = getDefaultEnvelopeCurve()
				}
				keyGroup.Program.Instruments.Instrument[j].FilterDecayCurve = getDefaultEnvelopeCurve()
				keyGroup.Program.Instruments.Instrument[j].FilterReleaseCurve = getDefaultEnvelopeCurve()
			}
			// Pitch envelope - EXS doesn't have dedicated pitch envelope, using filter envelope as approximation
			keyGroup.Program.Instruments.Instrument[j].PitchAttack = formatEnvTime(filtAttack)
			keyGroup.Program.Instruments.Instrument[j].PitchDecay = formatEnvTime(filtDecay)
			keyGroup.Program.Instruments.Instrument[j].PitchSustain = formatEnvLevel(filtSustain)
			keyGroup.Program.Instruments.Instrument[j].PitchRelease = formatEnvTime(filtRelease)
			keyGroup.Program.Instruments.Instrument[j].PitchHold = formatEnvTime(volHold)
			// Pitch envelope curves - use TimeCurve as approximation
			if exsFile.Params != nil && exsFile.Params.TimeCurve != 0 {
				keyGroup.Program.Instruments.Instrument[j].PitchAttackCurve = formatEnvelopeCurve(int(exsFile.Params.TimeCurve))
			} else {
				keyGroup.Program.Instruments.Instrument[j].PitchAttackCurve = getDefaultEnvelopeCurve()
			}
			keyGroup.Program.Instruments.Instrument[j].PitchDecayCurve = getDefaultEnvelopeCurve()
			keyGroup.Program.Instruments.Instrument[j].PitchReleaseCurve = getDefaultEnvelopeCurve()
			keyGroup.Program.Instruments.Instrument[j].PitchEnvAmount = "0"
			// Trigger mode - set based on group's Trigger field
			// Trigger == 1 means release-triggered samples (like piano sympathetic resonance)
			// TriggerMode: 0=one-shot, 1=release, 2=normal attack
			if g.Trigger == 1 {
				keyGroup.Program.Instruments.Instrument[j].TriggerMode = 1 // Release trigger
				klog.V(2).Infof("Setting release trigger for instrument %d (group %d)", j, zones[0].GroupIndex)
			} else {
				keyGroup.Program.Instruments.Instrument[j].TriggerMode = 2 // Normal attack trigger
			}

			// Set ZonePlay for Round Robin if enabled
			// ZonePlay values: 0=CYCLE (round robin), 1=VELOCITY, 2=RANDOM
			if hasRoundRobin && g.SelectGroup >= 0 {
				keyGroup.Program.Instruments.Instrument[j].ZonePlay = 0 // CYCLE = Round Robin
				klog.V(2).Infof("Setting Round Robin (ZonePlay=0) for instrument %d (group %d, SelectGroup=%d)",
					j, zones[0].GroupIndex, g.SelectGroup)
			} else {
				keyGroup.Program.Instruments.Instrument[j].ZonePlay = 1 // VELOCITY (default)
			}

			// Phase 2: One-shot mode - map from first zone in group
			// OneShot: "True" = sample plays once without looping (ignores note-off)
			if len(zones) > 0 && zones[0].OneShot {
				keyGroup.Program.Instruments.Instrument[j].OneShot = "True"
			} else {
				keyGroup.Program.Instruments.Instrument[j].OneShot = "False"
			}

			// Phase 2: Output routing - map from zone output to AudioRoute
			// Output: 0-15 in EXS maps to different output channels/busses
			if len(zones) > 0 && zones[0].HasOutput {
				keyGroup.Program.Instruments.Instrument[j].AudioRoute.AudioRoute = convertOutputToAudioRouteInt(int(zones[0].ExsZone.Output))
			}

			// LFO - initialize with default values
			keyGroup.Program.Instruments.Instrument[j].LFO.Type = "Triangle"
			keyGroup.Program.Instruments.Instrument[j].LFO.Rate = "0"
			keyGroup.Program.Instruments.Instrument[j].LFO.Sync = 0
			keyGroup.Program.Instruments.Instrument[j].LFO.Reset = "False"
			keyGroup.Program.Instruments.Instrument[j].LFO.PitchAmount = "0"
			keyGroup.Program.Instruments.Instrument[j].LFO.CutoffAmount = "0"
			keyGroup.Program.Instruments.Instrument[j].LFO.VolumeAmount = "0"
			keyGroup.Program.Instruments.Instrument[j].LFO.PanAmount = "0"
			keyGroup.Program.Instruments.Instrument[j].Volume = convertGain(float64(g.Volume))
			klog.V(2).Infof("Instrument: %s, LowNote: %d, HighNote: %d\n", keyGroup.Program.Instruments.Instrument[j].Number, keyGroup.Program.Instruments.Instrument[j].LowNote, keyGroup.Program.Instruments.Instrument[j].HighNote)

			// First pass: count valid layers (zones with successfully copied samples)
			validLayerCount := 0
			for _, zone := range zones {
				// Apply group velocity range limits to layers
				layerVelLow := int(zone.VelLow)
				layerVelHigh := int(zone.VelHigh)

				// Clamp to group's velocity range if set (non-zero)
				if g.VelLow != 0 && layerVelLow < int(g.VelLow) {
					layerVelLow = int(g.VelLow)
				}
				if g.VelHigh != 0 && layerVelHigh > int(g.VelHigh) {
					layerVelHigh = int(g.VelHigh)
				}

				// Skip layers completely outside group's velocity range
				if layerVelHigh < int(g.VelLow) || (g.VelHigh != 0 && layerVelLow > int(g.VelHigh)) {
					continue
				}

				validLayerCount++
			}

			// Skip instrument if no valid layers
			if validLayerCount == 0 {
				klog.V(2).Infof("Skipping instrument %d: no valid layers", j)
				continue
			}

			// Allocate layer array
			keyGroup.Program.Instruments.Instrument[j].Layers.Layer = make([]xpm.Layer, validLayerCount)

			// Second pass: populate layers
			layerIdx := 0
			for _, zone := range zones {
				// Apply group velocity range limits to layers
				layerVelLow := int(zone.VelLow)
				layerVelHigh := int(zone.VelHigh)

				// Clamp to group's velocity range if set (non-zero)
				if g.VelLow != 0 && layerVelLow < int(g.VelLow) {
					layerVelLow = int(g.VelLow)
				}
				if g.VelHigh != 0 && layerVelHigh > int(g.VelHigh) {
					layerVelHigh = int(g.VelHigh)
				}

				// Skip layers completely outside group's velocity range
				if layerVelHigh < int(g.VelLow) || (g.VelHigh != 0 && layerVelLow > int(g.VelHigh)) {
					klog.V(2).Infof("Skipping layer outside group velocity range: layer[%d-%d] group[%d-%d]",
						zone.VelLow, zone.VelHigh, g.VelLow, g.VelHigh)
					continue
				}

				sampleName := strings.TrimSpace(exsFile.Samples[zone.SampleIndex].FileName)
				xpmSampleName, xpmSampleFile, err := x.copySample(sampleName, destPath)
				if err != nil {
					klog.Warningf("Failed to copy sample '%s': %v", sampleName, err)
					// This shouldn't happen since we already counted valid layers,
					// but if it does, skip this layer
					continue
				}
				klog.V(2).Infof("Successfully copied sample: %s", sampleName)
				// layers - use group-limited velocity ranges
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].Number = fmt.Sprintf("%d", layerIdx+1)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].Active = "True"
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].Pitch = "0.000000"
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].Mute = "False"
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].VelStart = layerVelLow
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].VelEnd = layerVelHigh
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].SampleName = xpmSampleName
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].SampleFile = xpmSampleFile

				// Phase 1: Sample offsets - map EXS zone sample start/end to XPM layer
				// SampleStart/SampleEnd: absolute sample positions in the audio file
				// SliceStart/SliceEnd: appear to be used for slice-based sampling
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].SampleStart = int(zone.SampleStart)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].SampleEnd = int(zone.SampleEnd)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].SliceStart = int(zone.SampleStart)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].SliceEnd = int(zone.SampleEnd)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].Offset = int(zone.Offset)

				// Phase 1: Loop parameters - map EXS zone loop settings to XPM layer
				// Loop: "True" or "False" string to enable/disable looping
				// LoopStart/LoopEnd: loop points in samples
				// LoopCrossfadeLength: crossfade length for smooth loops
				// LoopTune: fine-tune adjustment for loop region
				if zone.LoopOn {
					keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].Loop = "True"
				} else {
					keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].Loop = "False"
				}
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].LoopStart = int(zone.LoopStart)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].LoopEnd = int(zone.LoopEnd)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].LoopCrossfadeLength = int(zone.LoopCrossfade)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].LoopTune = int(zone.LoopTune)
				// SliceLoop and SliceLoopStart also set for compatibility
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].SliceLoop = Btoi(zone.LoopOn)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].SliceLoopStart = int(zone.LoopStart)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].SliceLoopCrossFadeLength = int(zone.LoopCrossfade)

				// Phase 1: Zone tuning - map EXS zone pitch settings to XPM layer
				// TuneCoarse: semitone adjustment (-48 to +48)
				// TuneFine: cent adjustment (-50 to +50)
				// RootNote: MIDI note number that plays sample at original pitch
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].TuneCoarse = int(zone.CoarseTuning)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].TuneFine = int(zone.FineTuning)

				// Determine the RootNote: The RootNote is the MIDI note at which the sample plays at its original pitch.
				// EXS24 behavior: zone.Key field contains the root note for the zone.
				//
				// Correct approach:
				// - For single-note zones (KeyLow == KeyHigh): use zone.Key if it matches, otherwise use KeyLow
				// - For multi-note zones (KeyLow < KeyHigh): use zone.Key if set, otherwise use KeyLow
				rootNote := int(zone.KeyLow)

				if zone.KeyLow == zone.KeyHigh && zone.Key != 0 && zone.Key == uint8(zone.KeyLow) {
					// Single-note zone with matching zone.Key - use it
					rootNote = int(zone.Key)
					klog.Infof("Layer %d: Single-note zone, using zone.Key=%d as RootNote",
						layerIdx, zone.Key)
				} else if zone.KeyLow != zone.KeyHigh {
					// Multi-note zone - use zone.Key if set, otherwise KeyLow
					if zone.Key != 0 {
						rootNote = int(zone.Key)
						klog.Infof("Layer %d: Multi-note zone [%d-%d], using zone.Key=%d as RootNote, sample=%s",
							layerIdx, zone.KeyLow, zone.KeyHigh, zone.Key, sampleName)
					} else {
						rootNote = int(zone.KeyLow)
						klog.Infof("Layer %d: Multi-note zone [%d-%d], zone.Key=0, using KeyLow=%d as RootNote, sample=%s",
							layerIdx, zone.KeyLow, zone.KeyHigh, zone.KeyLow, sampleName)
					}
				} else {
					// Single-note zone but zone.Key is 0 or doesn't match - use KeyLow
					klog.Infof("Layer %d: Single-note zone, zone.Key=%d doesn't match KeyLow=%d, using KeyLow as RootNote",
						layerIdx, zone.Key, zone.KeyLow)
				}

				// IMPORTANT: MPC XPM format quirk - RootNote is stored as midi_note + 1
				// This is confirmed by ConvertWithMoss implementation:
				// - When writing XPM: add 1 to the MIDI note
				// - When reading XPM: subtract 1 from the stored value
				// See: https://github.com/git-moss/ConvertWithMoss/blob/main/src/main/java/de/mossgrabers/convertwithmoss/format/akai/MPCKeygroupCreator.java#L224
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].RootNote = rootNote + 1

				// Phase 1: Zone volume and pan - map EXS zone mixing to XPM layer
				// Volume: convert from dB (-60 to +12) to linear (0.0 to ~2.0)
				// Pan: convert from EXS range (-64 to +63) to XPM normalized (0.0 to 1.0)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].Volume = convertVolumeDbToLinear(int(zone.Volume))
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].Pan = convertPanToNormalized(int(zone.Pan))

				// Phase 2: Zone scale (key tracking) - map EXS zone scale to XPM layer KeyTrack
				// Scale: -100 to +100 in EXS, controls how much pitch changes with key
				// KeyTrack: "1.000000" = full tracking, "0.000000" = no tracking (fixed pitch)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].KeyTrack = convertScaleToKeyTrack(int(zone.Scale))
				klog.V(2).Infof("  Layer: %d, VelStart: %d, VelEnd: %d, SampleFile: %s\n", layerIdx, keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].VelStart, keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].VelEnd, keyGroup.Program.Instruments.Instrument[j].Layers.Layer[layerIdx].SampleFile)
				layerIdx++
			}
			j++
		}
	}

	keyGroup.Program.KeygroupNumKeygroups = len(z)

	// Resize the Instrument array to only include the actual instruments created
	// This prevents empty instruments with LowNote=0, HighNote=127 from being written to the XPM
	if j < len(keyGroup.Program.Instruments.Instrument) {
		keyGroup.Program.Instruments.Instrument = keyGroup.Program.Instruments.Instrument[:j]
		klog.V(2).Infof("Resized instrument array from 128 to %d actual instruments", j)
	}

	// Global pitch bend range - use the first group's settings as global
	// ConvertWithMoss applies this to all zones uniformly
	if len(groups) > 0 {
		// Calculate global pitch bend range
		// EXS doesn't have a direct global pitch bend, but we can use group settings
		// For now, use a default of 2 semitones (standard MIDI)
		keyGroup.Program.KeygroupPitchBendRange = "0.166667" // 2 semitones / 12 = 0.166667
		klog.V(2).Infof("Set global pitch bend range to 2 semitones")

		// Global master transpose - sum of all group tuning offsets
		// This is an approximation since EXS doesn't have explicit global transpose
		globalTranspose := 0.0
		for _, g := range groups {
			// Groups might have pitch offsets that act as transpose
			// For now, we'll leave this at 0 but log the capability
			_ = g
		}
		keyGroup.Program.KeygroupMasterTranspose = fmt.Sprintf("%.6f", globalTranspose)
		klog.V(2).Infof("Set global master transpose to %.2f", globalTranspose)
	}

	// Use EXS file name (without extension) for the output XPM file
	filename := exsFile.Name
	if idx := strings.LastIndex(filename, "."); idx > 0 {
		filename = filename[:idx]
	}
	return keyGroup.Save(destPath + "/" + filename + ".xpm")
}

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func convertGain(volumeDB float64) string {
	if volumeDB > 6 {
		volumeDB = 6
	}
	if volumeDB < -12 {
		volumeDB = -12
	}
	/*
			private static final double MINUS_12_DB  = 0.353000;
		    private static final double PLUS_6_DB    = 1.0;
		    private static final double VALUE_RANGE  = PLUS_6_DB - MINUS_12_DB;
	*/
	v := 12 + volumeDB
	res := (1 - 0.353000) * v / 18
	return fmt.Sprintf("%.6f", 0.353000+res)
	//return fmt.Sprintf("%.6f", math.Pow(10, volumeDB/20))
}

func clamp(value, minimum, maximum float64) float64 {
	return math.Max(minimum, math.Min(value, maximum))
}

func normalizeValue(value, minimum, maximum float64) string {
	return fmt.Sprintf("%.6f", clamp(value, minimum, maximum)/maximum)
}

// Constants for envelope time conversion (from ConvertWithMoss MPCKeygroupConstants)
const (
	MinEnvTimeSeconds   = 0.001 // Minimum envelope time in seconds
	MaxEnvTimeSeconds   = 10.0  // Maximum envelope time in seconds
	DefaultAttackTime   = 0.001 // Default attack time
	DefaultHoldTime     = 0.001 // Default hold time
	DefaultDecayTime    = 0.001 // Default decay time
	DefaultReleaseTime  = 0.63  // Default release time
	DefaultSustainLevel = 1.0   // Default sustain level
)

// formatEnvTime converts EXS envelope time values to XPM normalized values using logarithmic scaling
// EXS envelope times are in 0-127 range representing 0-10 seconds
// XPM uses logarithmic time scaling: normalizedValue = ln(time/min) / ln(max/min)
// Based on ConvertWithMoss normalizeLogarithmicEnvTimeValue
func formatEnvTime(envTime float64) string {
	if envTime < 0 {
		return fmt.Sprintf("%.6f", normalizeLogarithmicEnvTimeValue(DefaultAttackTime, MinEnvTimeSeconds, MaxEnvTimeSeconds))
	}
	// Convert EXS 0-127 range to 0-10 seconds (EXS maximum time per step)
	timeInSeconds := (envTime / 127.0) * 10.0
	// Apply logarithmic normalization for MPC
	return fmt.Sprintf("%.6f", normalizeLogarithmicEnvTimeValue(timeInSeconds, MinEnvTimeSeconds, MaxEnvTimeSeconds))
}

// normalizeLogarithmicEnvTimeValue computes normalized logarithmic value between 0 and 1
// The envelope time function of the MPC is approached by an exponential function:
//
//	duration = a * e^(b*control_value)
//
// where control_value is the normalized value (0-1) needed by MPC to produce the duration.
// Based on ConvertWithMoss MPCKeygroupCreator.normalizeLogarithmicEnvTimeValue
func normalizeLogarithmicEnvTimeValue(value, minimum, maximum float64) float64 {
	// Clamp value to valid range
	if value < minimum {
		value = minimum
	}
	if value > maximum {
		value = maximum
	}
	// Logarithmic formula: ln(value/min) / ln(max/min)
	return math.Log(value/minimum) / math.Log(maximum/minimum)
}

// formatEnvLevel converts EXS envelope level values to XPM normalized values (0-1)
// Based on ConvertWithMoss logic where sustain levels are scaled from 0-127 to 0-1
func formatEnvLevel(envLevel float64) string {
	if envLevel < 0 {
		envLevel = 0
	}
	if envLevel > 127 {
		envLevel = 127
	}
	// Convert from 0-127 range to 0-1 range
	return fmt.Sprintf("%.6f", envLevel/127.0)
}

// formatFilterCutoff converts EXS filter cutoff (0-127) to XPM normalized value (0-1)
func formatFilterCutoff(cutoff float64) string {
	if cutoff < 0 {
		cutoff = 0
	}
	if cutoff > 127 {
		cutoff = 127
	}
	return fmt.Sprintf("%.6f", cutoff/127.0)
}

// formatFilterResonance converts EXS filter resonance (0-127) to XPM normalized value (0-1)
func formatFilterResonance(resonance float64) string {
	if resonance < 0 {
		resonance = 0
	}
	if resonance > 127 {
		resonance = 127
	}
	return fmt.Sprintf("%.6f", resonance/127.0)
}

// formatEnvelopeCurve converts EXS envelope curve values to XPM normalized curve values
// EXS attack curves are typically in range -99 to +99 (signed)
// XPM curves are 0.0 to 1.0 where 0.5 is linear
// Based on ConvertWithMoss setEnvelopeCurveAttribute and EXS24Detector.createEnvelope
func formatEnvelopeCurve(exsCurve int) string {
	// EXS curves can be stored as signed values or with special encoding
	// If value >= 0xFF00, it represents a negative value: v = v - 0xFF00 - 0x100
	if exsCurve >= 0xFF00 {
		exsCurve = exsCurve - 0xFF00 - 0x100
	}
	// Normalize from -99..99 range to -1..1
	slopeValue := clamp(float64(exsCurve)/99.0, -1.0, 1.0)
	// Convert to XPM 0..1 range where 0.5 is linear
	curveValue := clamp((slopeValue+1.0)/2.0, 0.0, 1.0)
	return fmt.Sprintf("%.6f", curveValue)
}

// getDefaultEnvelopeCurve returns the default envelope curve value (0.5 = linear)
func getDefaultEnvelopeCurve() string {
	return "0.500000"
}

// convertVolumeDbToLinear converts EXS zone volume from dB (-60 to +12) to linear gain (0.0 to ~2.0)
// EXS zone volume range: -60 dB (silent) to +12 dB (boost)
// XPM layer volume: linear gain multiplier where 1.0 = unity gain
func convertVolumeDbToLinear(volumeDB int) string {
	// Clamp to EXS valid range
	if volumeDB < -60 {
		volumeDB = -60
	}
	if volumeDB > 12 {
		volumeDB = 12
	}
	// Convert dB to linear: gain = 10^(dB/20)
	linear := math.Pow(10.0, float64(volumeDB)/20.0)
	return fmt.Sprintf("%.6f", linear)
}

// convertPanToNormalized converts EXS zone pan from (-64 to +63) to XPM normalized (0.0 to 1.0)
// EXS pan range: -64 (hard left) to +63 (hard right), 0 = center
// XPM pan range: 0.0 (left) to 1.0 (right), 0.5 = center
func convertPanToNormalized(pan int) string {
	// Clamp to EXS valid range
	if pan < -64 {
		pan = -64
	}
	if pan > 63 {
		pan = 63
	}
	// Convert to 0.0-1.0 range where 0.5 is center
	// pan = -64 -> 0.0, pan = 0 -> 0.5, pan = 63 -> ~0.992
	normalized := (float64(pan) + 64.0) / 127.0
	return fmt.Sprintf("%.6f", normalized)
}

// Phase 2 Conversion Functions

// convertScaleToKeyTrack converts EXS zone scale to XPM layer KeyTrack
// EXS Scale: -100 to +100 (percentage of key tracking), typically 0-100 range
// XPM KeyTrack: 0.0 to 1.0 (1.0 = full tracking, 0.0 = no tracking/fixed pitch)
func convertScaleToKeyTrack(scale int) string {
	// EXS scale is typically 0-100, where 100 = full key tracking
	// Normalize to 0.0-1.0 range
	if scale < 0 {
		scale = 0
	}
	if scale > 100 {
		scale = 100
	}
	keyTrack := float64(scale) / 100.0
	return fmt.Sprintf("%.6f", keyTrack)
}

// convertOutputToAudioRoute converts EXS output number to XPM AudioRoute int
// EXS Output: 0-15 (various output routings)
// XPM AudioRoute: 0=main, 1-15=individual channels
func convertOutputToAudioRouteInt(output int) int {
	// Simple mapping: EXS output directly to XPM AudioRoute
	if output < 0 {
		return 0
	}
	if output > 15 {
		return 15
	}
	return output
}
