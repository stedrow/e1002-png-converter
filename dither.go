package main

import (
	"image"
	"image/color"
	"math"
)

// StuckiMatrix defines the error distribution pattern for Stucki dithering
var StuckiMatrix = []struct {
	dx, dy int
	weight float64
}{
	{1, 0, 8.0 / 42.0},   // Right
	{2, 0, 4.0 / 42.0},   // Right + 1
	{-2, 1, 2.0 / 42.0},  // Bottom-left-left
	{-1, 1, 4.0 / 42.0},  // Bottom-left
	{0, 1, 8.0 / 42.0},   // Bottom
	{1, 1, 4.0 / 42.0},   // Bottom-right
	{2, 1, 2.0 / 42.0},   // Bottom-right-right
	{-2, 2, 1.0 / 42.0},  // Bottom2-left-left
	{-1, 2, 2.0 / 42.0},  // Bottom2-left
	{0, 2, 4.0 / 42.0},   // Bottom2
	{1, 2, 2.0 / 42.0},   // Bottom2-right
	{2, 2, 1.0 / 42.0},   // Bottom2-right-right
}

// ApplyStuckiDithering applies the Stucki dithering algorithm to an image
func ApplyStuckiDithering(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create working image with float64 values for error propagation
	pixels := make([][]struct{ r, g, b float64 }, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]struct{ r, g, b float64 }, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			pixels[y][x].r = float64(r >> 8)
			pixels[y][x].g = float64(g >> 8)
			pixels[y][x].b = float64(b >> 8)
		}
	}

	// Create output image
	output := image.NewRGBA(bounds)

	// Apply Stucki dithering
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldR := pixels[y][x].r
			oldG := pixels[y][x].g
			oldB := pixels[y][x].b

			// Find closest color using quantization palette
			closestIdx := FindClosestColorIndex(
				uint8(clamp(oldR)),
				uint8(clamp(oldG)),
				uint8(clamp(oldB)),
			)

			// Get the output color
			newColor := Colors[closestIdx]
			newR := float64(newColor.R)
			newG := float64(newColor.G)
			newB := float64(newColor.B)

			// Set the new color in output
			output.SetRGBA(x+bounds.Min.X, y+bounds.Min.Y, color.RGBA{
				R: newColor.R,
				G: newColor.G,
				B: newColor.B,
				A: 255,
			})

			// Calculate error
			errorR := oldR - newR
			errorG := oldG - newG
			errorB := oldB - newB

			// Calculate total error
			totalError := math.Abs(errorR) + math.Abs(errorG) + math.Abs(errorB)

			// Skip error diffusion for near-perfect matches
			if totalError < 8 {
				continue
			}

			// Check if near-white or near-black
			isNearWhite := oldR > 245 && oldG > 245 && oldB > 245
			isNearBlack := oldR < 10 && oldG < 10 && oldB < 10
			quantizedIsWhite := newR > 245 && newG > 245 && newB > 245
			quantizedIsBlack := newR < 10 && newG < 10 && newB < 10

			if (isNearWhite && quantizedIsWhite) || (isNearBlack && quantizedIsBlack) {
				continue
			}

			// Distribute error to neighboring pixels
			for _, offset := range StuckiMatrix {
				newX := x + offset.dx
				newY := y + offset.dy

				if newX >= 0 && newX < width && newY >= 0 && newY < height {
					pixels[newY][newX].r = clamp(pixels[newY][newX].r + errorR*offset.weight)
					pixels[newY][newX].g = clamp(pixels[newY][newX].g + errorG*offset.weight)
					pixels[newY][newX].b = clamp(pixels[newY][newX].b + errorB*offset.weight)
				}
			}
		}
	}

	return output
}

// clamp ensures a float64 value is within 0-255 range
func clamp(val float64) float64 {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return val
}
