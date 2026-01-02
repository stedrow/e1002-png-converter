package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/nfnt/resize"
)

// CalculateResizeDimensions calculates new dimensions to fill target size while maintaining aspect ratio
func CalculateResizeDimensions(originalWidth, originalHeight, targetWidth, targetHeight int) (int, int) {
	// Calculate scale factor to ensure both dimensions are >= target
	scaleWidth := float64(targetWidth) / float64(originalWidth)
	scaleHeight := float64(targetHeight) / float64(originalHeight)
	scale := scaleWidth
	if scaleHeight > scale {
		scale = scaleHeight
	}

	newWidth := int(float64(originalWidth) * scale)
	newHeight := int(float64(originalHeight) * scale)

	return newWidth, newHeight
}

// CenterCrop crops an image to target dimensions from the center
func CenterCrop(img image.Image, targetWidth, targetHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate crop box
	left := (width - targetWidth) / 2
	top := (height - targetHeight) / 2

	// Create cropped image
	cropped := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.Draw(cropped, cropped.Bounds(),
		img, image.Point{X: bounds.Min.X + left, Y: bounds.Min.Y + top},
		draw.Src)

	return cropped
}

// ResizeAndCrop resizes an image to fill the target dimensions, then center crops
func ResizeAndCrop(img image.Image, targetWidth, targetHeight int) image.Image {
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// Check if resize needed
	if originalWidth == targetWidth && originalHeight == targetHeight {
		return img
	}

	// Calculate resize dimensions
	newWidth, newHeight := CalculateResizeDimensions(
		originalWidth, originalHeight,
		targetWidth, targetHeight,
	)

	// Resize using Lanczos3 (high quality)
	resized := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	// Center crop if needed
	if newWidth != targetWidth || newHeight != targetHeight {
		return CenterCrop(resized, targetWidth, targetHeight)
	}

	return resized
}

// AdjustBrightnessContrast applies brightness and contrast adjustments to an image
// brightness: -100 to 100 (adds/subtracts from RGB values)
// contrast: -100 to 100 (scales deviation from middle gray)
func AdjustBrightnessContrast(img image.Image, brightness, contrast int) image.Image {
	// Skip if no adjustments needed
	if brightness == 0 && contrast == 0 {
		return img
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	adjusted := image.NewRGBA(bounds)

	// Convert contrast to a scaling factor
	// contrast = 0 means no change (factor = 1.0)
	// contrast = 100 means maximum increase (factor = 2.0)
	// contrast = -100 means maximum decrease (factor = 0.0)
	contrastFactor := (100.0 + float64(contrast)) / 100.0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()

			// Convert from 16-bit to 8-bit
			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)

			// Apply contrast adjustment (scale deviation from middle gray)
			if contrast != 0 {
				r8 = ((r8 - 128) * contrastFactor) + 128
				g8 = ((g8 - 128) * contrastFactor) + 128
				b8 = ((b8 - 128) * contrastFactor) + 128
			}

			// Apply brightness adjustment
			if brightness != 0 {
				r8 += float64(brightness)
				g8 += float64(brightness)
				b8 += float64(brightness)
			}

			// Clamp to valid range [0, 255]
			r8 = clampFloat(r8, 0, 255)
			g8 = clampFloat(g8, 0, 255)
			b8 = clampFloat(b8, 0, 255)

			adjusted.Set(x+bounds.Min.X, y+bounds.Min.Y, color.RGBA{
				R: uint8(r8),
				G: uint8(g8),
				B: uint8(b8),
				A: uint8(a >> 8),
			})
		}
	}

	return adjusted
}

// clampFloat clamps a float64 value to the given range
func clampFloat(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
