# E1002 PNG Converter

Fast, standalone CLI tool for converting PNG images to e-ink compliant dithered PNGs optimized for the **reTerminal E1002** display (Spectra E6 7.3").

Built in Go - a single compiled binary with zero runtime dependencies.

## Features

- **Fast**: Native compiled Go binary, instant startup
- **Standalone**: Single binary, no runtime dependencies
- **Cross-platform**: Builds for Linux, macOS, Windows
- **4-bit PNG**: Native 4-bit indexed color output (no external tools needed)
- **Smart Resizing**: Auto-resize and center-crop to 800x480 while maintaining aspect ratio
- **Stucki Dithering**: High-quality error diffusion with minimal artifacts

## Installation

### From Source

```bash
# Install Go 1.21+ from https://go.dev/dl/

# Clone and build
git clone https://github.com/stedrow/e1002-png-converter.git
cd e1002-png-converter
make build

# Optional: Install to system
make install
```

### Download Binary

Download pre-built binaries from the releases page (coming soon).

## Usage

### Basic Usage

```bash
# Convert an image (auto-resizes and crops to 800x480)
./bin/e1002-convert input.png

# Specify output path
./bin/e1002-convert input.png -o output.png

# Disable automatic resizing
./bin/e1002-convert input.png --no-resize

# Custom max dimensions
./bin/e1002-convert input.png --max-width 1600 --max-height 1200
```

### Using Make

```bash
# Build
make build

# Convert an image
make run INPUT=input/yourimage.png

# Run test
make test

# Install system-wide
make install
```

## Command-Line Options

```
Usage:
  e1002-convert [input.png] [flags]

Flags:
  -o, --output string      Output PNG file path (default: input_dithered.png)
      --no-resize          Disable automatic resizing to fit 800x480 display
      --max-width int      Maximum width for resizing (default 800)
      --max-height int     Maximum height for resizing (default 480)
  -h, --help              Help for e1002-convert
```

## Building

### Build for Current Platform

```bash
make build
```

### Cross-Compilation

```bash
# Linux
make build-linux

# Windows
make build-windows

# macOS (Intel + Apple Silicon)
make build-darwin

# All platforms
make build-all
```

## Dependencies

### Build Dependencies

- Go 1.21+
- `github.com/nfnt/resize` - Image resizing
- `github.com/spf13/cobra` - CLI framework

### Runtime Dependencies

**None!** The compiled binary has zero runtime dependencies.

## Project Structure

```
e1002-png-converter/
├── main.go          # CLI entry point
├── colors.go        # E1002 color palette definitions
├── dither.go        # Stucki dithering algorithm
├── image_utils.go   # Resize and crop utilities
├── png_utils.go     # PNG palette conversion
├── go.mod           # Go module definition
├── Makefile         # Build commands
├── README.md        # This file
├── bin/             # Compiled binaries
├── input/           # Input images (for testing)
└── output/          # Output images (for testing)
```

## Algorithm

Uses a dual-palette Stucki dithering approach:

1. **Resize & Crop**: Image resized to fill 800x480, then center-cropped
2. **Quantization**: Match pixels to actual E1002 display colors
3. **Dithering**: Apply Stucki error diffusion
4. **Output**: Write idealized palette colors to 4-bit indexed PNG

## License

MIT

## Contributing

This is the Go port of the Python e1002-png-converter. Both implementations produce identical output.
