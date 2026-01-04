package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cobra"

	"github.com/cldmnky/exsconvert/pkg/convert"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch the graphical user interface",
	Long:  `Launch the Fyne-based GUI for converting EXS24 files to XPM format.`,
	Run: func(cmd *cobra.Command, args []string) {
		launchGUI()
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
}

type exsconvertGUI struct {
	window              fyne.Window
	exsFileEntry        *widget.Entry
	outputDirEntry      *widget.Entry
	samplesPathEntry    *widget.Entry
	layersEntry         *widget.Entry
	autoDetectCheck     *widget.Check
	skipErrorsCheck     *widget.Check
	statusLabel         *widget.Label
	resultList          *widget.List
	convertButton       *widget.Button
	selectedEXSFolder   string
	selectedOutputDir   string
	selectedSamplesPath string
	convertedFiles      []string
}

func launchGUI() {
	myApp := app.NewWithID("com.exsconvert.app")
	myWindow := myApp.NewWindow("EXS24 to XPM Converter")

	gui := &exsconvertGUI{
		window:            myWindow,
		convertedFiles:    make([]string, 0),
		selectedOutputDir: getDefaultOutputDir(),
	}

	gui.setupUI()

	myWindow.Resize(fyne.NewSize(700, 500))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

func (g *exsconvertGUI) setupUI() {
	// Header
	title := widget.NewLabelWithStyle("EXS24 to XPM Converter",
		fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Folder Selection Section
	g.exsFileEntry = widget.NewEntry()
	g.exsFileEntry.SetPlaceHolder("No folder selected...")
	g.exsFileEntry.Disable()

	selectFileButton := widget.NewButton("Select EXS Folder", func() {
		g.selectEXSFolder()
	})

	fileSection := container.NewBorder(nil, nil, nil, selectFileButton, g.exsFileEntry)

	// Output Directory Section
	g.outputDirEntry = widget.NewEntry()
	g.outputDirEntry.SetText(g.selectedOutputDir)
	g.outputDirEntry.Disable()

	selectOutputButton := widget.NewButton("Select Output Dir", func() {
		g.selectOutputDirectory()
	})

	outputSection := container.NewBorder(nil, nil, nil, selectOutputButton, g.outputDirEntry)

	// Samples Path Section (optional - for WAV files in different location)
	g.samplesPathEntry = widget.NewEntry()
	g.samplesPathEntry.SetPlaceHolder("(Optional) Defaults to EXS file directory")
	g.samplesPathEntry.Disable()

	selectSamplesButton := widget.NewButton("Select Samples Dir", func() {
		g.selectSamplesDirectory()
	})

	samplesSection := container.NewBorder(nil, nil, nil, selectSamplesButton, g.samplesPathEntry)

	// Options Section
	layersLabel := widget.NewLabel("Layers per Instrument:")
	g.layersEntry = widget.NewEntry()
	g.layersEntry.SetText("1")
	g.layersEntry.SetPlaceHolder("1")

	g.autoDetectCheck = widget.NewCheck("Auto-detect drum programs", nil)
	g.autoDetectCheck.SetChecked(true)

	g.skipErrorsCheck = widget.NewCheck("Skip errors during conversion", nil)
	g.skipErrorsCheck.SetChecked(true)

	optionsGrid := container.New(layout.NewFormLayout(),
		layersLabel, g.layersEntry,
		widget.NewLabel(""), g.autoDetectCheck,
		widget.NewLabel(""), g.skipErrorsCheck,
	)

	// Convert Button
	g.convertButton = widget.NewButton("Convert", func() {
		g.performConversion()
	})
	g.convertButton.Importance = widget.HighImportance
	g.convertButton.Disable()

	// Status Label
	g.statusLabel = widget.NewLabel("Ready to convert")
	g.statusLabel.Wrapping = fyne.TextWrapWord

	// Results List with click handler
	g.resultList = widget.NewList(
		func() int {
			return len(g.convertedFiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(g.convertedFiles[id])
		},
	)
	g.resultList.OnSelected = func(id widget.ListItemID) {
		g.showXPMFileOptions(id)
	}

	resultsCard := widget.NewCard("Converted Files", "Click a file to view or open", g.resultList)

	// Layout
	content := container.NewBorder(
		container.NewVBox(
			title,
			widget.NewSeparator(),
			widget.NewLabel("EXS24 Folder (will search recursively):"),
			fileSection,
			widget.NewLabel("Output Directory:"),
			outputSection,
			widget.NewLabel("Samples Directory (for WAV files):"),
			samplesSection,
			widget.NewLabel("Options:"),
			optionsGrid,
			g.convertButton,
			widget.NewSeparator(),
			g.statusLabel,
		),
		nil,
		nil,
		nil,
		resultsCard,
	)

	g.window.SetContent(content)
}

func (g *exsconvertGUI) selectEXSFolder() {
	dirDialog := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, g.window)
			return
		}
		if dir == nil {
			return
		}

		g.selectedEXSFolder = dir.Path()
		g.exsFileEntry.SetText(g.selectedEXSFolder)
		g.statusLabel.SetText(fmt.Sprintf("Selected folder: %s", filepath.Base(g.selectedEXSFolder)))
		g.updateConvertButton()
	}, g.window)

	dirDialog.Show()
}

func (g *exsconvertGUI) selectOutputDirectory() {
	dirDialog := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, g.window)
			return
		}
		if dir == nil {
			return
		}

		g.selectedOutputDir = dir.Path()
		g.outputDirEntry.SetText(g.selectedOutputDir)
		g.statusLabel.SetText(fmt.Sprintf("Output directory: %s", filepath.Base(g.selectedOutputDir)))
		g.updateConvertButton()
	}, g.window)

	dirDialog.Show()
}

func (g *exsconvertGUI) selectSamplesDirectory() {
	dirDialog := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, g.window)
			return
		}
		if dir == nil {
			return
		}

		g.selectedSamplesPath = dir.Path()
		g.samplesPathEntry.SetText(g.selectedSamplesPath)
		g.statusLabel.SetText(fmt.Sprintf("Samples directory: %s", filepath.Base(g.selectedSamplesPath)))
	}, g.window)

	dirDialog.Show()
}

func (g *exsconvertGUI) updateConvertButton() {
	if g.selectedEXSFolder != "" && g.selectedOutputDir != "" {
		g.convertButton.Enable()
	} else {
		g.convertButton.Disable()
	}
}

func (g *exsconvertGUI) performConversion() {
	if g.selectedEXSFolder == "" || g.selectedOutputDir == "" {
		dialog.ShowError(fmt.Errorf("Please select both EXS folder and output directory"), g.window)
		return
	}

	// Disable convert button during conversion
	g.convertButton.Disable()
	defer g.convertButton.Enable()

	// Parse layers
	layers := 1
	if g.layersEntry.Text != "" {
		fmt.Sscanf(g.layersEntry.Text, "%d", &layers)
	}

	// Update status
	g.statusLabel.SetText("Converting... Please wait.")
	g.convertedFiles = []string{}
	g.resultList.Refresh()

	// The search path for samples is the selected folder
	// (or the custom samples path if specified)
	searchPath := g.selectedSamplesPath
	if searchPath == "" {
		searchPath = g.selectedEXSFolder
	}

	// Create converter - it will search for .exs files recursively
	xpmConverter := convert.NewXPM(
		g.selectedEXSFolder,
		g.selectedOutputDir,
		layers,
		g.skipErrorsCheck.Checked,
		"", // Empty program type - will be auto-detected
	)

	// Set auto-detect mode
	xpmConverter.AutoDetectDrums = g.autoDetectCheck.Checked

	// Set search path for samples
	xpmConverter.SamplesSearchPath = searchPath

	// Perform conversion (will find all .exs files recursively)
	err := xpmConverter.Convert()
	if err != nil {
		g.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
		dialog.ShowError(fmt.Errorf("Conversion failed: %v", err), g.window)
		return
	}

	// Find generated XPM files
	g.findConvertedFiles()

	// Update status with information about sample copying
	statusMsg := fmt.Sprintf("✓ Conversion complete! Generated %d file(s)", len(g.convertedFiles))
	if searchPath != "" {
		statusMsg += fmt.Sprintf("\nSamples searched in: %s", searchPath)
	}
	g.statusLabel.SetText(statusMsg)

	if len(g.convertedFiles) > 0 {
		dialog.ShowInformation("Success",
			fmt.Sprintf("Successfully converted!\n\nGenerated %d XPM file(s) in:\n%s",
				len(g.convertedFiles), g.selectedOutputDir),
			g.window)
	}
}

func (g *exsconvertGUI) findConvertedFiles() {
	g.convertedFiles = []string{}

	// Walk the entire output directory to find XPM files
	filepath.Walk(g.selectedOutputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".xpm") {
			relPath, _ := filepath.Rel(g.selectedOutputDir, path)
			g.convertedFiles = append(g.convertedFiles, relPath)
		}
		return nil
	})

	g.resultList.Refresh()
}

func getDefaultOutputDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(homeDir, "Desktop")
}

// Custom file filter for EXS files
type exsFileFilter struct{}

func newEXSFileFilter() *exsFileFilter {
	return &exsFileFilter{}
}

func (e *exsFileFilter) Matches(uri fyne.URI) bool {
	return strings.HasSuffix(strings.ToLower(uri.Path()), ".exs")
}

func (g *exsconvertGUI) showXPMFileOptions(id widget.ListItemID) {
	if id < 0 || id >= len(g.convertedFiles) {
		return
	}

	filePath := filepath.Join(g.selectedOutputDir, g.convertedFiles[id])
	fileName := filepath.Base(filePath)

	// Create dialog with options
	viewButton := widget.NewButton("View Contents", func() {
		g.viewXPMFile(filePath)
	})
	viewButton.Importance = widget.HighImportance

	openButton := widget.NewButton("Open in System Editor", func() {
		g.openXPMFileInSystem(filePath)
	})

	openFolderButton := widget.NewButton("Show in Folder", func() {
		g.openFolderInSystem(filepath.Dir(filePath))
	})

	content := container.NewVBox(
		widget.NewLabelWithStyle(fileName, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewLabel("Choose an action:"),
		viewButton,
		openButton,
		openFolderButton,
	)

	dialog.ShowCustom("XPM File Options", "Close", content, g.window)
}

func (g *exsconvertGUI) viewXPMFile(filePath string) {
	// Read file contents
	content, err := os.ReadFile(filePath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to read file: %v", err), g.window)
		return
	}

	// Create scrollable text display
	textContent := widget.NewLabel(string(content))
	textContent.Wrapping = fyne.TextWrapWord

	scroll := container.NewScroll(textContent)
	scroll.SetMinSize(fyne.NewSize(600, 400))

	// Show in dialog
	fileName := filepath.Base(filePath)
	d := dialog.NewCustom("XPM File: "+fileName, "Close", scroll, g.window)
	d.Resize(fyne.NewSize(700, 500))
	d.Show()
}

func (g *exsconvertGUI) openXPMFileInSystem(filePath string) {
	var cmd *exec.Cmd

	switch {
	case fileExists("/usr/bin/open"): // macOS
		cmd = exec.Command("open", filePath)
	case fileExists("/usr/bin/xdg-open"): // Linux
		cmd = exec.Command("xdg-open", filePath)
	default: // Windows
		cmd = exec.Command("cmd", "/c", "start", "", filePath)
	}

	err := cmd.Start()
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to open file: %v", err), g.window)
	}
}

func (g *exsconvertGUI) openFolderInSystem(folderPath string) {
	var cmd *exec.Cmd

	switch {
	case fileExists("/usr/bin/open"): // macOS
		cmd = exec.Command("open", folderPath)
	case fileExists("/usr/bin/xdg-open"): // Linux
		cmd = exec.Command("xdg-open", folderPath)
	default: // Windows
		cmd = exec.Command("explorer", folderPath)
	}

	err := cmd.Start()
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to open folder: %v", err), g.window)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
