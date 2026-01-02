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
	outputPath   string
	noResize     bool
	maxWidth     int
	maxHeight    int
	deviceID     string
	ditherMethod string
	listDevs     bool
	brightness   int
	contrast     int
)

var rootCmd = &cobra.Command{
	Use:   "e1002-convert [input.png]",
	Short: "Convert PNG images to e-ink compliant dithered PNGs for E1002 display",
	Long: `E1002 PNG Converter - Converts PNG images to e-ink compliant dithered PNGs
optimized for the reTerminal E1002 display (Spectra E6 7.3").

Applies error diffusion dithering to convert standard PNG images to a 6-color palette
suitable for e-ink displays, producing high-quality results with minimal artifacts.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Allow 0 args when listing devices
		if listDevs {
			return nil
		}
		return cobra.ExactArgs(1)(cmd, args)
	},
	RunE: convertImage,
}

func init() {
	rootCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output PNG file path (default: input_dithered.png)")
	rootCmd.Flags().StringVarP(&deviceID, "device", "d", "", "Device profile to use (e.g., reterminal-e1002, waveshare-7in5-v2)")
	rootCmd.Flags().BoolVar(&listDevs, "list-devices", false, "List all available device profiles")
	rootCmd.Flags().StringVar(&ditherMethod, "dither", "stucki", "Dithering algorithm (stucki, floyd-steinberg, atkinson)")
	rootCmd.Flags().IntVar(&brightness, "brightness", 0, "Brightness adjustment: -100 (darker) to 100 (brighter)")
	rootCmd.Flags().IntVar(&contrast, "contrast", 0, "Contrast adjustment: -100 (less contrast) to 100 (more contrast)")
	rootCmd.Flags().BoolVar(&noResize, "no-resize", false, fmt.Sprintf("Disable automatic resizing to fit %dx%d display", E1002Width, E1002Height))
	rootCmd.Flags().IntVar(&maxWidth, "max-width", E1002Width, "Maximum width for resizing")
	rootCmd.Flags().IntVar(&maxHeight, "max-height", E1002Height, "Maximum height for resizing")
}

func convertImage(cmd *cobra.Command, args []string) error {
	// Handle --list-devices flag
	if listDevs {
		devices, err := ListDevices()
		if err != nil {
			return fmt.Errorf("failed to list devices: %w", err)
		}
		fmt.Println("Available device profiles:")
		for _, device := range devices {
			fmt.Println(device)
		}
		return nil
	}

	inputPath := args[0]

	// Load device profile if specified
	if deviceID != "" {
		device, err := GetDevice(deviceID)
		if err != nil {
			return fmt.Errorf("failed to load device profile: %w", err)
		}
		fmt.Printf("Using device profile: %s\n", device.Name)
		maxWidth = device.Width
		maxHeight = device.Height
	}

	// Validate dithering method
	method := DitherMethod(ditherMethod)
	if method != DitherStucki && method != DitherFloydSteinberg && method != DitherAtkinson {
		return fmt.Errorf("invalid dithering method: %s (must be stucki, floyd-steinberg, or atkinson)", ditherMethod)
	}

	// Validate brightness and contrast
	if brightness < -100 || brightness > 100 {
		return fmt.Errorf("brightness must be between -100 and 100, got %d", brightness)
	}
	if contrast < -100 || contrast > 100 {
		return fmt.Errorf("contrast must be between -100 and 100, got %d", contrast)
	}

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

	// Apply brightness/contrast adjustments if specified
	if brightness != 0 || contrast != 0 {
		fmt.Printf("Adjusting brightness: %+d, contrast: %+d\n", brightness, contrast)
		img = AdjustBrightnessContrast(img, brightness, contrast)
	}

	fmt.Printf("Applying %s dithering with E1002 palette...\n", ditherMethod)

	// Apply dithering
	start := time.Now()
	dithered := ApplyDithering(img, method)
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
