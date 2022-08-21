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
	converter           convert.Convert
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert exs files to mpc keygroups",
	Long:  `Make it so`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// output xml
		//
		xpmConverter := convert.NewXPM(searchPath, outputPath, layersPerInstrument, skipErrors)
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
	convertCmd.Flags().StringVarP(&searchPath, "search-path", "p", "", "search path for exs files and samples")
	convertCmd.Flags().StringVarP(&outputPath, "output-path", "o", "", "output path for key group files")
	convertCmd.Flags().IntVarP(&layersPerInstrument, "layers-per-instrument", "l", 4, "number of layers per instrument")
	convertCmd.Flags().BoolVarP(&skipErrors, "skip-errors", "s", true, "skip errors")
}
