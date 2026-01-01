package main

import (
	"image/color"
)

// E1002 Display Configuration (Spectra E6 7.3")
const (
	E1002Width  = 800
	E1002Height = 480
)

// RGB color type for easier manipulation
type RGB struct {
	R, G, B uint8
}

// Output palette - idealized RGB values written to the output image
var Colors = []RGB{
	{0x00, 0x00, 0x00}, // Black
	{0xFF, 0xFF, 0xFF}, // White
	{0x00, 0x00, 0xFF}, // Blue
	{0x00, 0xFF, 0x00}, // Green
	{0xFF, 0x00, 0x00}, // Red
	{0xFF, 0xFF, 0x00}, // Yellow
}

// Palette colors - actual display colors for accurate matching
var PaletteColors = []RGB{
	{0x19, 0x1E, 0x21}, // Black (actual)
	{0xE8, 0xE8, 0xE8}, // White (actual)
	{0x21, 0x57, 0xBA}, // Blue (actual)
	{0x12, 0x5F, 0x20}, // Green (actual)
	{0xB2, 0x13, 0x18}, // Red (actual)
	{0xEF, 0xDE, 0x44}, // Yellow (actual)
}

// ToColor converts RGB to color.Color
func (rgb RGB) ToColor() color.Color {
	return color.RGBA{R: rgb.R, G: rgb.G, B: rgb.B, A: 255}
}

// Distance calculates Euclidean distance between two RGB colors
func (rgb RGB) Distance(other RGB) float64 {
	dr := float64(rgb.R) - float64(other.R)
	dg := float64(rgb.G) - float64(other.G)
	db := float64(rgb.B) - float64(other.B)
	return dr*dr + dg*dg + db*db
}

// FindClosestColorIndex finds the closest color in the quantization palette
func FindClosestColorIndex(r, g, b uint8) int {
	pixel := RGB{r, g, b}
	minDistance := float64(1e9)
	closestIndex := 0

	for i, paletteColor := range PaletteColors {
		distance := pixel.Distance(paletteColor)
		if distance < minDistance {
			minDistance = distance
			closestIndex = i
		}
	}

	return closestIndex
}
