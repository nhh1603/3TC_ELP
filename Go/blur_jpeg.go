package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

func gaussianBlur(img image.Image, radius int) *image.RGBA {
	// Create a new RGBA image with the same bounds as the original
	bounds := img.Bounds()
	destImg := image.NewRGBA(bounds)
	kernel := generateKernel(radius)

	// Iterate over each pixel in the image
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			var r, g, b float64

			// Iterate over each surrounding pixel within the specified radius
			for i := -radius; i <= radius; i++ {
				for j := -radius; j <= radius; j++ {
					// Get the color of a surrounding pixel
					xn := x + i
					yn := y + j
					if xn < bounds.Min.X || xn >= bounds.Max.X || yn < bounds.Min.Y || yn >= bounds.Max.Y {
						continue
					}
					neighborPixelColor := img.At(xn, yn)
					sr, sg, sb, _ := neighborPixelColor.RGBA()

					// Apply the Gaussian blur weight to the surrounding pixel color
					blurWeight := kernel[i+radius][j+radius]
					r += float64(sr>>8) * blurWeight
					g += float64(sg>>8) * blurWeight
					b += float64(sb>>8) * blurWeight
				}
			}

			// Set the color of the destination pixel
			destImg.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}

	return destImg
}

// Generate a normalized gaussian kernel with given radius
func generateKernel(radius int) [][]float64 {
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
	return 1 / (2 * math.Pi * float64(sqrSigma)) * math.Exp(-(float64(x*x+y*y))/(2*float64(sqrSigma)))
}

func main() {
	fmt.Println("Start")

	var img image.Image = loadImage("b&w_large.png")

	// Apply the Gaussian blur
	blurredImg := gaussianBlur(img, 1)
	saveImage("b&w_large_blurred.png", blurredImg, PNG)
	fmt.Println("Blurred image generated")
}
