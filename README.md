# exsconvert

Convert Logic Pro EXS24 sampler instruments to Akai MPC keygroup format (XPM files).

## Features

- Converts EXS24 instruments to MPC-compatible XPM keygroup programs
- Automatically copies and converts sample files (WAV format with uppercase extension)
- Preserves envelope parameters, filter settings, and sample mappings
- GUI and command-line interfaces available

## Requirements

- Go 1.18 or higher (for building from source)
- EXS24 files (.exs)
- Sample files (.wav) referenced by the EXS instruments

## Installation

### From Source

```bash
go build .
```

This creates the `exsconvert` binary in the project directory.

## Usage

### GUI Mode (Recommended)

```bash
./exsconvert gui
```

The GUI provides:
- File selection for EXS instruments
- Output directory selection
- Optional separate samples directory (if your WAV files are in a different location)
- Visual feedback on conversion progress
- XPM file viewer and file exploration

### Command Line Mode

```bash
./exsconvert convert -p /path/to/exs/files -o /path/to/output
```

#### Options

- `-p, --path` - Directory containing EXS files and samples (searches recursively)
- `-o, --output` - Output directory for converted XPM files
- `-l, --layers` - Layers per instrument (default: 1)
- `-s, --skip-errors` - Skip errors during conversion (default: true)

## Output Structure

The converter creates an MPC-compatible directory structure:

```
output/
├── InstrumentName1/
│   ├── InstrumentName1.xpm
│   ├── Sample1.WAV
│   ├── Sample2.WAV
│   └── Sample3.WAV
└── InstrumentName2/
    ├── InstrumentName2.xpm
    └── Sample4.WAV
```

**Important for MPC Compatibility:**
- Each XPM file and its sample files MUST be in the same directory
- Sample files use uppercase `.WAV` extension
- Samples should be 16-bit or 24-bit, 44.1 kHz WAV files

## Loading on MPC

1. Copy the output directory to your MPC's storage (SD card, USB drive, or internal drive)
2. In the MPC Browser, navigate to the location where you copied the files
3. Select the `.xpm` file and load it
4. The MPC will automatically load the WAV files from the same directory

## Sample File Handling

The converter automatically:
- Searches for WAV files referenced in the EXS instrument
- Searches recursively through the specified samples directory
- Copies found samples to the output directory
- Converts file extensions to uppercase (`.WAV`)
- Reports missing samples as warnings

**Tips:**
- Place your WAV files in the same directory as your EXS files, or
- Use the GUI to specify a separate samples directory, or
- Use the `-p` flag to point to a directory containing both EXS and WAV files

## Parameter Conversion

The converter properly scales EXS parameters to MPC format:

- **Envelope Times**: Logarithmic scaling for attack/decay/release
- **Envelope Levels**: Linear scaling for sustain
- **Filter Cutoff**: Linear scaling (0-127 → 0-1)
- **Filter Resonance**: Linear scaling (0-127 → 0-1)

## Troubleshooting

### Missing Samples

If samples are not found:
- Ensure WAV files are in the search path (`-p` directory or GUI samples directory)
- Check that filenames match exactly (case-sensitive on some systems)
- Verify files have `.wav` or `.WAV` extension

### MPC Won't Load XPM

- Verify WAV files are in the same directory as the XPM file
- Check that sample rate is 44.1 kHz
- Ensure bit depth is 16 or 24-bit

## Building

```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Build
go build .
```

## License

See LICENSE file.
