/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/cldmnky/exsconvert/pkg/convert"
	"github.com/spf13/cobra"
)

var (
	searchPath          string
	outputPath          string
	layersPerInstrument int
	skipErrors          bool
	programType         string
	autoDetect          bool
	samplesPath         string
	converter           convert.Convert
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert exs files to mpc keygroups",
	Long:  `Make it so`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// output xml
		//
		xpmConverter := convert.NewXPM(searchPath, outputPath, layersPerInstrument, skipErrors, programType)

		// Set auto-detect if enabled
		if autoDetect {
			xpmConverter.AutoDetectDrums = true
			xpmConverter.ProgramType = "" // Clear program type to enable auto-detect
		}

		// Set custom samples path if provided
		if samplesPath != "" {
			xpmConverter.SamplesSearchPath = samplesPath
		}

		converter = xpmConverter
		err := converter.Convert()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVarP(&searchPath, "search-path", "p", "", "search path for exs files (will search recursively)")
	convertCmd.Flags().StringVarP(&outputPath, "output-path", "o", "", "output path for XPM files")
	convertCmd.Flags().StringVarP(&samplesPath, "samples-path", "w", "", "search path for WAV samples (defaults to search-path)")
	convertCmd.Flags().IntVarP(&layersPerInstrument, "layers-per-instrument", "l", 4, "number of layers per instrument")
	convertCmd.Flags().BoolVarP(&skipErrors, "skip-errors", "s", true, "skip errors")
	convertCmd.Flags().StringVarP(&programType, "program-type", "t", "", "program type: Keygroup or Drum (leave empty to auto-detect)")
	convertCmd.Flags().BoolVarP(&autoDetect, "auto-detect", "a", false, "auto-detect drum programs (overrides -t)")
}
