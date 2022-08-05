package convert

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/cldmnky/exsconvert/pkg/exs"
	"github.com/cldmnky/exsconvert/pkg/xpm"
	"k8s.io/klog"
)

type XPM struct {
	SearchPath          string
	OutputPath          string
	LayersPerInstrument int
	SkipErrors          bool
}

func NewXPM(searchPath, outputPath string, layersPerInstrument int, skipErrors bool) *XPM {
	return &XPM{
		SearchPath:          searchPath,
		OutputPath:          outputPath,
		LayersPerInstrument: layersPerInstrument,
		SkipErrors:          skipErrors,
	}
}

func (x *XPM) Convert() error {
	exsFiles, err := x.findEXSFiles()
	if err != nil {
		return err
	}
	for _, exsFile := range exsFiles {
		exs, err := exs.NewFromFile(exsFile)
		if err != nil {
			return err
		}
		destPath := filepath.Join(x.OutputPath, exs.Name)
		err = os.MkdirAll(destPath, 0755)
		if err != nil {
			return err
		}
		err = x.toXPM(exs, destPath)
		if err != nil {
			return err
		}
		fmt.Printf("Converted %s\n", exs.Name)

	}
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

func (x *XPM) copySample(name, destPath string) (string, error) {
	samples := []string{}
	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == name {
			klog.V(2).Infof("found %s\n", path)
			samples = append(samples, path)
			return nil
		}
		return nil
	}
	err := filepath.WalkDir(x.SearchPath, walk)
	if err != nil {
		return "", err
	}
	if len(samples) == 0 {
		return "", fmt.Errorf("no sample found for %s", name)
	}
	if len(samples) > 1 {
		return "", fmt.Errorf("multiple samples found for %s", name)
	}

	var toUpperExt = func(fileName string) string {
		ext := filepath.Ext(fileName)
		fileName = fileName[:len(fileName)-len(ext)]
		return fmt.Sprintf("%s%s", fileName, strings.ToUpper(ext))
	}

	src := samples[0]
	dst := filepath.Join(destPath, filepath.Base(toUpperExt(src)))
	in, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return "", err
	}
	err = out.Close()
	if err != nil {
		return "", err
	}

	return toUpperExt(name), nil
}

func (x *XPM) toXPM(exsFile *exs.EXS, destPath string) error {
	keyGroup := xpm.NewXPMKeygroup()
	z := exsFile.GetZonesByKeyRanges(x.LayersPerInstrument)
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
	for _, group := range groups {
		klog.V(2).Infof("group: %s", group.Name)
	}
	if len(groups) > 1 {
		klog.Infof("group: %s", groups[0].Name)
		return fmt.Errorf("multiple groups found")
	}
	g := groups[0]

	keyGroup.Program.ProgramName = g.Name
	j := 0
	for _, zoneMap := range z {
		for _, zones := range zoneMap {
			keyGroup.Program.Instruments.Instrument[j].LowNote = int(zones[0].KeyLow)
			keyGroup.Program.Instruments.Instrument[j].HighNote = int(zones[0].KeyHigh)
			keyGroup.Program.Instruments.Instrument[j].Resonance = float32(g.Resonance)
			keyGroup.Program.Instruments.Instrument[j].VolumeAttack = float32(g.Attack1)
			keyGroup.Program.Instruments.Instrument[j].VolumeDecay = float32(g.Decay1)
			keyGroup.Program.Instruments.Instrument[j].VolumeSustain = float32(g.Sustain1)
			keyGroup.Program.Instruments.Instrument[j].VolumeRelease = float32(g.Release1)
			keyGroup.Program.Instruments.Instrument[j].Volume = float32(g.Volume)
			klog.V(2).Infof("Instrument: %s, LowNote: %d, HighNote: %d\n", keyGroup.Program.Instruments.Instrument[j].Number, keyGroup.Program.Instruments.Instrument[j].LowNote, keyGroup.Program.Instruments.Instrument[j].HighNote)
			for i, zone := range zones {
				sampleName := strings.TrimSpace(exsFile.Samples[zone.SampleIndex].FileName)
				xpmSampleName, err := x.copySample(sampleName, destPath)
				if err != nil {
					continue
				}
				// layers
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].VelStart = int(zone.VelLow)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].VelEnd = int(zone.VelHigh)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].SampleFile = xpmSampleName
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].SampleStart = int(zone.SampleStart)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].SampleEnd = int(zone.SampleEnd)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].Loop = zone.LoopOn
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].LoopStart = int(zone.LoopStart)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].LoopEnd = int(zone.LoopEnd)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].RootNote = int(zone.Key)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].LoopCrossfadeLength = int(zone.LoopCrossfade)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].Offset = int(zone.Offset)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].LoopTune = int(zone.LoopTune)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].TuneCoarse = int(zone.CoarseTuning)
				keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].TuneFine = int(zone.FineTuning)
				klog.V(2).Infof("  Layer: %d, VelStart: %d, VelEnd: %d, SampleFile: %s\n", i, keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].VelStart, keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].VelEnd, keyGroup.Program.Instruments.Instrument[j].Layers.Layer[i].SampleFile)
			}
			j++
		}
	}

	keyGroup.Program.KeygroupNumKeygroups = len(z)
	return keyGroup.Save(destPath + "/" + g.Name + ".xpm")
}