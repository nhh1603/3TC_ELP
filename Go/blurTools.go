package main

import (
	"image"
	"image/color"
	"math"
)

// Computes the convolution of the pixels in the range [startPix, endPix] with the given kernel
func convolveRange(fromImg image.Image, toImg *image.RGBA, startPix, endPix int, kernel [][]float64) {
	fromBounds := fromImg.Bounds()
	maxBounds := fromBounds.Max

	if endPix >= maxBounds.X*maxBounds.Y {
		panic("End pixel is out of image's bounds")
	} else if startPix < 0 {
		panic("Start pixel is out of image's bounds")
	}

	radius := (len(kernel) - 1) / 2

	// Iterate over each pixel from startPix to endPix, row after row
	for pixIndex := startPix; pixIndex <= endPix; pixIndex++ {
		currentPix := image.Point{X: pixIndex % maxBounds.X, Y: pixIndex / maxBounds.X}
		var r, g, b float64
		// Iterate over each surrounding pixel within the specified radius
		for i := -radius; i <= radius; i++ {
			for j := -radius; j <= radius; j++ {
				neighborX := currentPix.X + i
				neighborY := currentPix.Y + j

				// Mirror the pixel against the edge of image in case it is out of range
				if neighborX < 0 {
					neighborX = -neighborX
				} else if neighborX >= fromBounds.Max.X {
					neighborX = 2*(fromBounds.Max.X-1) - neighborX
				}
				if neighborY < 0 {
					neighborY = -neighborY
				} else if neighborY >= fromBounds.Max.Y {
					neighborY = 2*(fromBounds.Max.Y-1) - neighborY
				}
				// Get the color of a surrounding pixel
				neighborR, neighborG, neighborB, _ := fromImg.At(neighborX, neighborY).RGBA()
				// Apply the Gaussian blur weight to the surrounding pixel color
				weight := kernel[i+radius][j+radius]
				r += float64(neighborR>>8) * weight
				g += float64(neighborG>>8) * weight
				b += float64(neighborB>>8) * weight
			}
		}
		// Set the color of the destination pixel
		toImg.Set(currentPix.X, currentPix.Y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
	}
}

// Generate a normalized gaussian kernel with given radius
func generateGaussKernel(radius int) [][]float64 {
	size := radius*2 + 1
	kernel := make([][]float64, size)
	sqrSigma := math.Max(float64(radius*radius)/4, 1)
	sum := 0.0

	for i := -radius; i <= radius; i++ {
		kernel[i+radius] = make([]float64, size)
		for j := -radius; j <= radius; j++ {
			weight := gaussianWeight(i, j, sqrSigma)
			kernel[i+radius][j+radius] = weight
			sum += weight
		}
	}

	// Normalize the kernel
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			kernel[i][j] /= sum
		}
	}
	return kernel
}

// 2D Gaussian function
func gaussianWeight(x int, y int, sqrSigma float64) float64 {
	return 1 / (2 * math.Pi * sqrSigma) * math.Exp(-(float64(x*x+y*y))/(2*sqrSigma))
}
