package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	outputPath string
	noResize   bool
	maxWidth   int
	maxHeight  int
)

var rootCmd = &cobra.Command{
	Use:   "e1002-convert [input.png]",
	Short: "Convert PNG images to e-ink compliant dithered PNGs for E1002 display",
	Long: `E1002 PNG Converter - Converts PNG images to e-ink compliant dithered PNGs
optimized for the reTerminal E1002 display (Spectra E6 7.3").

Applies Stucki dithering to convert standard PNG images to a 6-color palette
suitable for e-ink displays, producing high-quality results with minimal artifacts.`,
	Args: cobra.ExactArgs(1),
	RunE: convertImage,
}

func init() {
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output PNG file path (default: input_dithered.png)")
	rootCmd.Flags().BoolVar(&noResize, "no-resize", false, fmt.Sprintf("Disable automatic resizing to fit %dx%d display", E1002Width, E1002Height))
	rootCmd.Flags().IntVar(&maxWidth, "max-width", E1002Width, "Maximum width for resizing")
	rootCmd.Flags().IntVar(&maxHeight, "max-height", E1002Height, "Maximum height for resizing")
}

func convertImage(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Validate input file
	if !strings.HasSuffix(strings.ToLower(inputPath), ".png") {
		return fmt.Errorf("input file must be a PNG: %s", inputPath)
	}

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", inputPath)
	}

	// Generate output path if not provided
	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		base := strings.TrimSuffix(inputPath, ext)
		outputPath = base + "_dithered.png"
	}

	fmt.Printf("Loading image: %s\n", inputPath)

	// Load image
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode PNG: %w", err)
	}

	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	fmt.Printf("Original size: %dx%d\n", originalWidth, originalHeight)

	// Resize and crop if requested
	if !noResize {
		if originalWidth != maxWidth || originalHeight != maxHeight {
			// Calculate resize dimensions
			newWidth, newHeight := CalculateResizeDimensions(
				originalWidth, originalHeight,
				maxWidth, maxHeight,
			)
			fmt.Printf("Resizing to: %dx%d (maintaining aspect ratio)\n", newWidth, newHeight)

			// Resize and crop
			if newWidth != maxWidth || newHeight != maxHeight {
				fmt.Printf("Center cropping to: %dx%d\n", maxWidth, maxHeight)
			}
			img = ResizeAndCrop(img, maxWidth, maxHeight)
		} else {
			fmt.Printf("Image already %dx%d, no resize needed\n", maxWidth, maxHeight)
		}
	}

	bounds = img.Bounds()
	fmt.Printf("Processing size: %dx%d\n", bounds.Dx(), bounds.Dy())
	fmt.Println("Applying Stucki dithering with E1002 palette...")

	// Apply dithering
	start := time.Now()
	dithered := ApplyStuckiDithering(img)
	elapsed := time.Since(start)
	fmt.Printf("Dithering completed in %.1f seconds\n", elapsed.Seconds())

	// Convert to palette
	fmt.Println("Converting to indexed color palette...")
	paletted := ConvertToPalette(dithered)

	// Save
	fmt.Printf("Saving to: %s\n", outputPath)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Use best compression and let the encoder optimize bit depth
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	if err := encoder.Encode(outFile, paletted); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}
	outFile.Close()

	fmt.Println("✓ Conversion complete!")

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
