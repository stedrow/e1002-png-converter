package main

import (
	"image"
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
