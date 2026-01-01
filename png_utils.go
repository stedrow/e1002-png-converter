package main

import (
	"image"
	"image/color"
)

// ConvertToPalette converts an RGBA image to an indexed palette image with minimal palette
func ConvertToPalette(img *image.RGBA) *image.Paletted {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create minimal color palette with only our 6 colors
	// This allows PNG encoder to use 4-bit depth (2^4 = 16 colors, we use 6)
	palette := make(color.Palette, len(Colors))
	for i, c := range Colors {
		palette[i] = c.ToColor()
	}

	// Create paletted image
	paletted := image.NewPaletted(bounds, palette)

	// Map pixels to palette indices
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()

			// Find matching color in our palette
			pixel := RGB{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}

			// Find exact match
			var colorIndex uint8
			for i, c := range Colors {
				if c.R == pixel.R && c.G == pixel.G && c.B == pixel.B {
					colorIndex = uint8(i)
					break
				}
			}

			paletted.SetColorIndex(x+bounds.Min.X, y+bounds.Min.Y, colorIndex)
		}
	}

	return paletted
}
