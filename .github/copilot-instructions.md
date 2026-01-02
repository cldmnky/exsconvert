# Copilot Instructions for exsconvert

## Repository Overview

**Purpose**: exsconvert is a Go CLI application that converts Logic Pro EXS24 sampler instrument files to Akai MPC keygroup format (XPM files). It parses binary EXS files and generates XML-based XPM files suitable for MPC hardware/software.

**Type**: Command-line tool  
**Language**: Go 1.18+  
**Size**: Small (~2,085 lines of Go code)  
**Frameworks**: Cobra (CLI), Ginkgo/Gomega (testing), klog (logging)

## Build and Validation Commands

### Prerequisites
- Go 1.18 or higher (tested with Go 1.24.11)
- No additional build tools required
- Dependencies are managed via go.mod

### Build Process (ALWAYS follow this order)
1. **Download dependencies** (first time or after go.mod changes):
   ```bash
   go mod download
   ```

2. **Build the binary**:
   ```bash
   go build .
   ```
   - This creates the `exsconvert` binary in the project root
   - Build time: ~2-5 seconds on clean build
   - Binary size: ~6MB

3. **Alternative build with custom output**:
   ```bash
   go build -v -o exsconvert .
   ```

### Testing (ALWAYS run before committing)
```bash
go test ./...
```
- Uses Ginkgo/Gomega BDD testing framework
- All tests should pass; test time: <1 second
- No separate Ginkgo CLI installation needed - tests run via `go test`
- Test files: `*_test.go` with `*_suite_test.go` for Ginkgo setup

**Test specific packages**:
```bash
go test -v ./pkg/exs/...
go test -v ./pkg/convert/...
go test -v ./pkg/xpm/...
```

**Generate test coverage**:
```bash
go test -coverprofile=coverage.out ./...
```

### Code Quality (ALWAYS run before committing)
1. **Format code** (ALWAYS run first):
   ```bash
   go fmt ./...
   ```
   - Must be run before committing any Go code changes
   - Will modify files in place

2. **Static analysis**:
   ```bash
   go vet ./...
   ```
   - Should produce no output if passing

### Clean Build
To test from a clean state:
```bash
go clean
rm -f exsconvert coverage.out
go build .
go test ./...
```

## Project Structure

```
.
├── main.go                    # Entry point - calls cmd.Execute()
├── cmd/                       # CLI commands (Cobra)
│   ├── root.go               # Root command setup, klog initialization
│   └── convert.go            # Main "convert" subcommand
├── pkg/
│   ├── exs/                  # EXS file parser
│   │   ├── decode.go         # Binary parsing logic for EXS24 files
│   │   ├── exs_sequence.go   # Sequence data structures
│   │   ├── utils.go          # Helper functions
│   │   ├── exs_test.go       # Tests
│   │   ├── exs_suite_test.go # Ginkgo test suite setup
│   │   └── testdata/         # Sample .exs files for testing
│   ├── convert/              # Conversion orchestration
│   │   ├── convert.go        # Interface definition
│   │   ├── xpm.go            # XPM conversion implementation
│   │   ├── convert_test.go   # Tests
│   │   └── convert_suite_test.go
│   └── xpm/                  # XPM file generation
│       ├── types.go          # XML data structures for MPC format
│       ├── xpm.go            # XML generation logic
│       ├── xpm_test.go       # Tests
│       ├── xpm_suite_test.go
│       └── testdata/         # Sample .xpm files
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
└── README.md                 # Minimal documentation
```

### Key Files and Their Purposes

**Entry Point**:
- `main.go` (12 lines) - Imports cmd package and calls Execute()

**CLI Layer** (`cmd/`):
- `root.go` - Defines root command `exs2mpc` (note: binary is named `exsconvert` but CLI command is `exs2mpc`)
- `convert.go` - Main conversion command with flags: `-p` (search path), `-o` (output path), `-l` (layers per instrument), `-s` (skip errors)

**Core Logic** (`pkg/`):
- `exs/decode.go` - Binary file parser for EXS24 format (chunks: header, zone, group, sample, options)
- `convert/xpm.go` - Orchestrates file finding and conversion
- `xpm/types.go` - Extensive XML struct definitions for MPC format (MPCVObject, Program, Instruments, Layers, etc.)

## Usage and Behavior

**Command structure**:
```bash
./exsconvert convert -p <search-path> -o <output-path> [-l layers] [-s]
```

**What the tool does**:
1. Recursively finds .exs files in the search path
2. Parses each EXS file to extract zones, groups, samples
3. Generates one or more XPM (XML) files in the output directory
4. Creates subdirectories per instrument

## Important Conventions and Notes

### Code Style
- Standard Go formatting (use `go fmt`)
- Uses k8s.io/klog for structured logging with verbosity levels
- No special linting configuration - standard Go conventions apply

### Testing Conventions
- Ginkgo BDD style tests with Gomega matchers
- Test suites defined in `*_suite_test.go` files
- Each package has its own test suite
- Test data stored in `testdata/` subdirectories
- Many tests are commented out (/* ... */) in the test files - this is intentional

### Binary Naming
- The built binary file is named `exsconvert` (this is what you build and run, for example: `./exsconvert --help`)
- The Cobra root command defined in `cmd/root.go` is named `exs2mpc`, so usage/help text and error messages will show `exs2mpc` as the command name
- This difference is intentional: you should always invoke the tool as `exsconvert`, but seeing `exs2mpc` in help output is expected and refers to the same tool

### Dependencies
- **Cobra**: CLI framework (github.com/spf13/cobra)
- **Ginkgo v2**: BDD testing (github.com/onsi/ginkgo/v2)
- **Gomega**: Test matchers (github.com/onsi/gomega)
- **klog**: Kubernetes-style logging (k8s.io/klog)

### File Handling
- Reads binary EXS files with endianness detection
- Writes XML files with specific MPC format requirements
- Error handling: can skip errors with `-s` flag (default: true)

## Common Pitfalls and Solutions

1. **Binary not found after build**: The binary `exsconvert` is in the project root, not in a `bin/` directory. Run it as `./exsconvert`.

2. **Tests must use `go test`, not `ginkgo` CLI**: While the project uses Ginkgo testing framework, tests are run through standard `go test` command. The ginkgo CLI is not required.

3. **Modified files after checkout**: Running `go fmt ./...` modifies source files. Always run it before checking git status.

4. **Test file naming**: Ginkgo tests use `_test.go` suffix and require a suite file with pattern `*_suite_test.go` containing `RunSpecs()`.

5. **Binary in git**: The `exsconvert` binary should NOT be committed. It's added to `.gitignore`.

## Validation Workflow

Before finalizing any code changes:
```bash
# 1. Format code
go fmt ./...

# 2. Static analysis
go vet ./...

# 3. Run tests
go test ./...

# 4. Build to verify compilation
go build .

# 5. (Optional) Test the binary
./exsconvert --help
```

## No CI/CD Configuration

This repository currently has no GitHub Actions or other CI/CD workflows. All validation must be done locally.

## Making Changes

When modifying code:
1. **Parser changes** (`pkg/exs/`): Be careful with binary parsing logic, chunk types, and endianness handling
2. **Conversion logic** (`pkg/convert/`): Ensure proper handling of file paths and error cases
3. **XML generation** (`pkg/xpm/`): Maintain XML structure compatibility with MPC format
4. **CLI changes** (`cmd/`): Update help text and flag descriptions appropriately

## Definition of Done

A code change is considered complete when:
- [ ] Code compiles without errors (`go build .`)
- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Static analysis passes (`go vet ./...`)
- [ ] Binary runs without errors (for CLI changes)
- [ ] Relevant documentation is updated
- [ ] No sensitive data or secrets committed

## Environment Notes

- **Local Development**: Standard Go development environment
- **GitHub Actions**: When Copilot runs in CI, it uses a containerized environment with Go pre-installed
- **No External Dependencies**: All builds are self-contained via go.mod

## Working with Copilot

- **Iterative Feedback**: Review PRs and use @copilot mentions in PR comments to request changes
- **Well-Scoped Issues**: Provide clear problem statements with specific acceptance criteria
- **Start Small**: Begin with focused tasks like bug fixes or feature additions to a single package

## Trust These Instructions

These instructions have been validated by:
- Building from clean state
- Running all tests
- Testing code formatting and static analysis
- Verifying build times and outputs
- Testing the built binary

ONLY search for additional information if you encounter behavior that contradicts these instructions or if you need details not covered here.
