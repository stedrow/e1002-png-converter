# E1002 PNG Converter

Fast, standalone CLI tool for converting PNG images to e-ink compliant dithered PNGs optimized for the **reTerminal E1002** display (Spectra E6 7.3").

Built in Go - a single compiled binary with zero runtime dependencies.

## Features

- **Fast**: Native compiled Go binary, instant startup
- **Standalone**: Single binary, no runtime dependencies
- **Cross-platform**: Builds for Linux, macOS, Windows
- **4-bit PNG**: Native 4-bit indexed color output (no external tools needed)
- **Device Profiles**: Pre-configured settings for popular e-ink displays
- **Multiple Dithering Algorithms**: Choose from Stucki, Floyd-Steinberg, or Atkinson
- **Brightness/Contrast Adjustment**: Fine-tune image tone before dithering for optimal e-ink display
- **Smart Resizing**: Auto-resize and center-crop to target dimensions while maintaining aspect ratio

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

# Use a device profile
./bin/e1002-convert input.png --device waveshare-7in5-v2

# List available device profiles
./bin/e1002-convert --list-devices

# Use Floyd-Steinberg dithering instead of Stucki
./bin/e1002-convert input.png --dither floyd-steinberg

# Use Atkinson dithering (lighter, more subtle)
./bin/e1002-convert input.png --dither atkinson

# Adjust brightness (darker images on e-ink)
./bin/e1002-convert input.png --brightness 20

# Adjust contrast (sharper appearance)
./bin/e1002-convert input.png --contrast 30

# Adjust both brightness and contrast
./bin/e1002-convert input.png --brightness 15 --contrast 25

# Combine all options for optimal e-ink output
./bin/e1002-convert input.png -d pimoroni-inky-impression --dither atkinson --brightness 10 --contrast 20 -o output.png

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
  -d, --device string      Device profile to use (e.g., reterminal-e1002, waveshare-7in5-v2)
      --list-devices       List all available device profiles
      --dither string      Dithering algorithm: stucki, floyd-steinberg, atkinson (default "stucki")
      --brightness int     Brightness adjustment: -100 (darker) to 100 (brighter) (default 0)
      --contrast int       Contrast adjustment: -100 (less contrast) to 100 (more contrast) (default 0)
      --no-resize          Disable automatic resizing to fit 800x480 display
      --max-width int      Maximum width for resizing (default 800)
      --max-height int     Maximum height for resizing (default 480)
  -h, --help              Help for e1002-convert
```

## Device Profiles

Built-in device profiles for popular e-ink displays:

- `reterminal-e1002` (alias: `e1002`) - reTerminal E1002 (Spectra E6 7.3"), 800x480, 7 colors
- `waveshare-7in5-v2` (alias: `waveshare-7.5`) - Waveshare 7.5inch V2, 800x480, 2 colors
- `waveshare-7in5-v3` - Waveshare 7.5inch V3, 800x480, 2 colors
- `waveshare-7in5-b` (alias: `waveshare-bwr`) - Waveshare 7.5inch B, 640x384, 3 colors
- `waveshare-5in65` - Waveshare 5.65inch, 600x448, 7 colors
- `pimoroni-inky-impression` (alias: `inky`) - Pimoroni Inky Impression 7.3", 800x480, 7 colors
- `kindle-paperwhite` - Kindle Paperwhite, 1236x1648, 16 grayscale levels

Use `--list-devices` to see all available profiles.

## Dithering Algorithms

Three error diffusion algorithms are available:

- **Stucki** (default) - Highest quality, spreads error over 12 neighboring pixels. Best for photos and complex images.
- **Floyd-Steinberg** - Classic algorithm, spreads error over 4 pixels. Good balance of speed and quality.
- **Atkinson** - Lighter dithering, preserves more highlights. Best for line art, comics, and high-contrast images.

## Brightness and Contrast Adjustments

Fine-tune image appearance before dithering to optimize for e-ink displays:

### How It Works

E-ink displays have fixed color palettes, but brightness/contrast adjustments affect which palette colors get selected during dithering:

- **Brightness** (-100 to 100): Shifts color selection toward lighter or darker palette colors
  - Positive values: More whites, yellows, lighter colors (brighter appearance)
  - Negative values: More blacks, dark blues, darker colors (darker appearance)

- **Contrast** (-100 to 100): Controls how extreme the color selection is
  - Positive values: More pure blacks/whites, sharper appearance
  - Negative values: More middle tones, softer appearance

### When to Use

- **Dark/muddy source images**: Try `--brightness 15 --contrast 20`
- **Washed-out images**: Try `--brightness -10 --contrast 25`
- **Text readability**: Try `--contrast 30` for sharper text
- **Photos**: Try `--brightness 10 --contrast 15` for better tone mapping

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
├── dither.go        # Error diffusion dithering algorithms (Stucki, Floyd-Steinberg, Atkinson)
├── image_utils.go   # Resize and crop utilities
├── png_utils.go     # PNG palette conversion
├── device.go        # Device profile management
├── devices.json     # Device profile definitions (embedded in binary)
├── go.mod           # Go module definition
├── Makefile         # Build commands
├── README.md        # This file
├── bin/             # Compiled binaries
├── input/           # Input images (for testing)
└── output/          # Output images (for testing)
```

## Algorithm

Uses a dual-palette error diffusion approach:

1. **Resize & Crop**: Image resized to fill target dimensions, then center-cropped
2. **Quantization**: Match pixels to actual E1002 display colors
3. **Dithering**: Apply selected error diffusion algorithm (Stucki, Floyd-Steinberg, or Atkinson)
4. **Output**: Write idealized palette colors to 4-bit indexed PNG

## License

MIT

## Contributing

This is the Go port of the Python e1002-png-converter. Both implementations produce identical output.
