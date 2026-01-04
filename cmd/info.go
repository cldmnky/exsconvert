package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cldmnky/exsconvert/pkg/exs"
	"github.com/spf13/cobra"
)

var (
	infoVerbose bool
	infoShowAll bool
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info [exs-file]",
	Short: "Display information about an EXS file",
	Long: `Display detailed information about an EXS24 instrument file for troubleshooting.

This command parses the EXS file and displays:
- Basic file information (name, size, endianness)
- Zone information (key ranges, velocity ranges, root notes)
- Group information
- Sample information
- Global parameters

Examples:
  exsconvert info myinstrument.exs
  exsconvert info --verbose myinstrument.exs  # Verbose output with all details
  exsconvert info -a myinstrument.exs         # Show all zones (not just summary)`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().BoolVar(&infoVerbose, "verbose", false, "Show verbose output with detailed information")
	infoCmd.Flags().BoolVarP(&infoShowAll, "all", "a", false, "Show all zones/groups/samples (not just summary)")
}

func runInfo(cmd *cobra.Command, args []string) error {
	exsPath := args[0]

	// Check if file exists
	if _, err := os.Stat(exsPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", exsPath)
	}

	// Parse the EXS file
	fmt.Printf("Parsing EXS file: %s\n", exsPath)
	exsFile, err := exs.NewFromFile(exsPath)
	if err != nil {
		return fmt.Errorf("failed to parse EXS file: %w", err)
	}

	// Display basic information
	printBasicInfo(exsFile, exsPath)
	fmt.Println()

	// Display zones information
	printZonesInfo(exsFile)
	fmt.Println()

	// Display groups information
	printGroupsInfo(exsFile)
	fmt.Println()

	// Display samples information
	printSamplesInfo(exsFile)
	fmt.Println()

	// Display parameters if verbose
	if infoVerbose && exsFile.Params != nil {
		printParamsInfo(exsFile)
		fmt.Println()
	}

	// Display sequences if any
	if len(exsFile.Sequences) > 0 {
		printSequencesInfo(exsFile)
		fmt.Println()
	}

	return nil
}

func printBasicInfo(exsFile *exs.EXS, path string) {
	fileInfo, _ := os.Stat(path)
	fmt.Println("═══ Basic Information ═══")
	fmt.Printf("  Name:          %s\n", exsFile.Name)
	fmt.Printf("  File:          %s\n", filepath.Base(path))
	fmt.Printf("  Size:          %d bytes (%.2f KB)\n", fileInfo.Size(), float64(fileInfo.Size())/1024.0)
	fmt.Printf("  Endianness:    %s\n", getEndianness(exsFile.BigEndian))
	fmt.Printf("  Size Expanded: %v\n", exsFile.IsSizeExpanded)
	fmt.Printf("  Zones:         %d\n", len(exsFile.Zones))
	fmt.Printf("  Groups:        %d\n", len(exsFile.Groups))
	fmt.Printf("  Samples:       %d\n", len(exsFile.Samples))
}

func printZonesInfo(exsFile *exs.EXS) {
	fmt.Println("═══ Zones Information ═══")
	fmt.Printf("  Total Zones: %d\n", len(exsFile.Zones))

	if len(exsFile.Zones) == 0 {
		fmt.Println("  No zones found")
		return
	}

	// Show summary or all zones based on flag
	limit := 10
	if infoShowAll {
		limit = len(exsFile.Zones)
	}

	fmt.Println()
	fmt.Println("  Key Range Analysis:")
	minKey := int8(127)
	maxKey := int8(0)
	for _, zone := range exsFile.Zones {
		if zone.KeyLow < minKey {
			minKey = zone.KeyLow
		}
		if zone.KeyHigh > maxKey {
			maxKey = zone.KeyHigh
		}
	}
	fmt.Printf("    Overall key range: %d (%s) to %d (%s)\n",
		minKey, midiNoteName(int(minKey)),
		maxKey, midiNoteName(int(maxKey)))

	fmt.Println()
	if infoVerbose {
		fmt.Printf("  Showing %d of %d zones:\n", min(limit, len(exsFile.Zones)), len(exsFile.Zones))
		fmt.Println("  ┌──────────────────────────────────────────────────────────────────────────────┐")
		fmt.Println("  │ # │ Name                  │ Key Range    │ Vel    │ Root│ Sample           │")
		fmt.Println("  ├──────────────────────────────────────────────────────────────────────────────┤")

		for i := 0; i < min(limit, len(exsFile.Zones)); i++ {
			zone := exsFile.Zones[i]
			sampleName := "N/A"
			if int(zone.SampleIndex) < len(exsFile.Samples) {
				sampleName = truncateString(strings.TrimSpace(exsFile.Samples[zone.SampleIndex].FileName), 16)
			}

			fmt.Printf("  │%3d│ %-21s │ %3d-%-3d %4s │ %3d-%-3d│  %3d│ %-16s │\n",
				i+1,
				truncateString(zone.Name, 21),
				zone.KeyLow, zone.KeyHigh, midiNoteName(int(zone.KeyLow)),
				zone.VelLow, zone.VelHigh,
				zone.Key,
				sampleName)
		}
		fmt.Println("  └──────────────────────────────────────────────────────────────────────────────┘")
	} else {
		// Show first few zones in compact format
		fmt.Printf("  First %d zones:\n", min(limit, len(exsFile.Zones)))
		for i := 0; i < min(limit, len(exsFile.Zones)); i++ {
			zone := exsFile.Zones[i]
			fmt.Printf("    %2d. %-25s  Keys: %3d-%-3d (%s to %s)  Vel: %3d-%-3d  RootKey: %3d\n",
				i+1,
				truncateString(zone.Name, 25),
				zone.KeyLow, zone.KeyHigh,
				midiNoteName(int(zone.KeyLow)), midiNoteName(int(zone.KeyHigh)),
				zone.VelLow, zone.VelHigh,
				zone.Key)
		}
	}

	if !infoShowAll && len(exsFile.Zones) > limit {
		fmt.Printf("\n  ... and %d more zones (use -a to show all)\n", len(exsFile.Zones)-limit)
	}
}

func printGroupsInfo(exsFile *exs.EXS) {
	fmt.Println("═══ Groups Information ═══")
	fmt.Printf("  Total Groups: %d\n", len(exsFile.Groups))

	if len(exsFile.Groups) == 0 {
		fmt.Println("  No groups found")
		return
	}

	limit := 10
	if infoShowAll {
		limit = len(exsFile.Groups)
	}

	if infoVerbose {
		fmt.Println()
		fmt.Printf("  Showing %d of %d groups:\n", min(limit, len(exsFile.Groups)), len(exsFile.Groups))
		for i := 0; i < min(limit, len(exsFile.Groups)); i++ {
			group := exsFile.Groups[i]
			fmt.Printf("\n  Group %d:\n", i+1)
			fmt.Printf("    ID:           %d\n", group.ID)
			fmt.Printf("    Name:         %s\n", group.Name)
			fmt.Printf("    Volume:       %d\n", group.Volume)
			fmt.Printf("    Pan:          %d\n", group.Pan)
			fmt.Printf("    Key Range:    %d-%d\n", group.KeyLow, group.KeyHigh)
			fmt.Printf("    Vel Range:    %d-%d\n", group.VelLow, group.VelHigh)
			fmt.Printf("    Polyphony:    %d\n", group.Polyphony)
			fmt.Printf("    SelectGroup:  %d\n", group.SelectGroup)
			fmt.Printf("    SelectNumber: %d\n", group.SelectNumber)
			if group.SelectGroup >= 0 {
				fmt.Printf("    ⚡ Round Robin Enabled\n")
			}
		}
	} else {
		fmt.Println()
		for i := 0; i < min(limit, len(exsFile.Groups)); i++ {
			group := exsFile.Groups[i]
			rrIndicator := ""
			if group.SelectGroup >= 0 {
				rrIndicator = " [RR]"
			}
			fmt.Printf("    %2d. %-25s  ID: %3d  Vol: %4d  Keys: %3d-%-3d  SelectGrp: %3d%s\n",
				i+1,
				truncateString(group.Name, 25),
				group.ID,
				group.Volume,
				group.KeyLow, group.KeyHigh,
				group.SelectGroup,
				rrIndicator)
		}
	}

	if !infoShowAll && len(exsFile.Groups) > limit {
		fmt.Printf("\n  ... and %d more groups (use -a to show all)\n", len(exsFile.Groups)-limit)
	}
}

func printSamplesInfo(exsFile *exs.EXS) {
	fmt.Println("═══ Samples Information ═══")
	fmt.Printf("  Total Samples: %d\n", len(exsFile.Samples))

	if len(exsFile.Samples) == 0 {
		fmt.Println("  No samples found")
		return
	}

	limit := 10
	if infoShowAll {
		limit = len(exsFile.Samples)
	}

	fmt.Println()
	if infoVerbose {
		fmt.Printf("  Showing %d of %d samples:\n", min(limit, len(exsFile.Samples)), len(exsFile.Samples))
		for i := 0; i < min(limit, len(exsFile.Samples)); i++ {
			sample := exsFile.Samples[i]
			fmt.Printf("\n  Sample %d:\n", i+1)
			fmt.Printf("    Name:     %s\n", sample.Name)
			fmt.Printf("    FileName: %s\n", strings.TrimSpace(sample.FileName))
			if infoVerbose && sample.Path != "" {
				fmt.Printf("    Path:     %s\n", sample.Path)
			}
		}
	} else {
		fmt.Printf("  First %d samples:\n", min(limit, len(exsFile.Samples)))
		for i := 0; i < min(limit, len(exsFile.Samples)); i++ {
			sample := exsFile.Samples[i]
			fmt.Printf("    %3d. %s\n", i+1, strings.TrimSpace(sample.FileName))
		}
	}

	if !infoShowAll && len(exsFile.Samples) > limit {
		fmt.Printf("\n  ... and %d more samples (use -a to show all)\n", len(exsFile.Samples)-limit)
	}
}

func printParamsInfo(exsFile *exs.EXS) {
	params := exsFile.Params
	fmt.Println("═══ Global Parameters ═══")
	fmt.Printf("  Output Volume:    %d dB\n", params.OutputVolume)
	fmt.Printf("  Pitch Bend Up:    %d cents\n", params.PitchBendUp)
	fmt.Printf("  Pitch Bend Down:  %d cents\n", params.PitchBendDown)
	fmt.Printf("  Mono Mode:        %d\n", params.MonoMode)
	fmt.Printf("  Voices:           %d\n", params.Voices)
	fmt.Printf("  Unison:           %v\n", params.Unison)
	fmt.Println()
	fmt.Println("  Filter:")
	fmt.Printf("    Filter On:        %v\n", params.FilterOn)
	fmt.Printf("    Filter Type:      %d\n", params.FilterType)
	fmt.Printf("    Filter Cutoff:    %d\n", params.FilterCutoff)
	fmt.Printf("    Filter Resonance: %d\n", params.FilterResonance)
	fmt.Printf("    Filter Drive:     %d\n", params.FilterDrive)
	fmt.Println()
	fmt.Println("  Envelopes:")
	fmt.Printf("    Env1 (Filter): A=%d D=%d S=%d R=%d\n",
		params.Env1Attack, params.Env1Decay, params.Env1Sustain, params.Env1Release)
	fmt.Printf("    Env2 (Volume): A=%d D=%d S=%d R=%d\n",
		params.Env2Attack, params.Env2Decay, params.Env2Sustain, params.Env2Release)
}

func printSequencesInfo(exsFile *exs.EXS) {
	fmt.Println("═══ Round Robin Sequences ═══")
	fmt.Printf("  Total Sequences: %d\n", len(exsFile.Sequences))
	fmt.Println()
	for i, seq := range exsFile.Sequences {
		fmt.Printf("  Sequence %d: %v\n", i+1, seq)
	}
}

// Helper functions
func getEndianness(bigEndian bool) string {
	if bigEndian {
		return "Big Endian"
	}
	return "Little Endian"
}

func midiNoteName(note int) string {
	if note < 0 || note > 127 {
		return "???"
	}
	noteNames := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	octave := (note / 12) - 2
	noteName := noteNames[note%12]
	return fmt.Sprintf("%s%d", noteName, octave)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
