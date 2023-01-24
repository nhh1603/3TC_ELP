package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"
)

// Computes the convolution of the pixels in the range [startPix, endPix] with the given kernel
func convolveRange(fromImg image.Image, toImg *image.RGBA, startPix, endPix int, kernel [][]float64) {
	fromBounds := fromImg.Bounds()

	if endPix >= fromBounds.Max.X * fromBounds.Max.Y {
		// panic("End row is out of image's bounds")
		endPix = fromBounds.Max.X * fromBounds.Max.Y - 1
	} else if startPix < 0 {
		// panic("Start row is out of image's bounds")
		startPix = 0
	}

	radius := (len(kernel) - 1) / 2

	// Iterate over each pixel from the starting to the ending row of pixels
	for i := startPix; i <= endPix; i++ {
		currentPix := [2]int {i / fromBounds.Max.X, i % fromBounds.Max.X}
		var r, g, b float64
		// Iterate over each surrounding pixel within the specified radius
		for i := -radius; i <= radius; i++ {
			for j := -radius; j <= radius; j++ {
				neighborX := currentPix[1] + i
				neighborY := currentPix[0] + j
				if neighborX < fromBounds.Min.X || neighborX >= fromBounds.Max.X ||
					neighborY < fromBounds.Min.Y || neighborY >= fromBounds.Max.Y {
					// Ignore neighbor pixels out of range
					continue
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
		toImg.Set(currentPix[1], currentPix[0], color.RGBA{uint8(r), uint8(g), uint8(b), 255})
	}
}

// Computes the convolution of the pixels in the range [startRow, endRow] with the given kernel
func convolveRegion(fromImg image.Image, toImg *image.RGBA, region image.Rectangle, kernel [][]float64) {
	fromBounds := fromImg.Bounds()

	if region.Max.Y > fromBounds.Max.Y {
		panic("End row is out of image's bounds")
	} else if region.Min.Y < fromBounds.Min.Y {
		panic("Start row is out of image's bounds")
	}

	radius := (len(kernel) - 1) / 2

	// Iterate over each pixel from the starting to the ending row of pixels
	for x := region.Min.X; x < region.Max.X; x++ {
		for y := region.Min.Y; y < region.Max.Y; y++ {
			var r, g, b float64

			// Iterate over each surrounding pixel within the specified radius
			for i := -radius; i <= radius; i++ {
				for j := -radius; j <= radius; j++ {
					neighborX := x + i
					neighborY := y + j
					if neighborX < fromBounds.Min.X || neighborX >= fromBounds.Max.X ||
						neighborY < fromBounds.Min.Y || neighborY >= fromBounds.Max.Y {
						// Ignore neighbor pixels out of range
						continue
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
			toImg.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
		}
	}
}

// Computes the convolution of the pixels in the range [startRow, endRow] with the given kernel
func getConvolvedWithin(img image.Image, startRow, endRow int, kernel [][]float64) [][]color.RGBA {
	bounds := img.Bounds()

	if endRow >= bounds.Max.Y {
		panic("End row is out of image's bounds")
	} else if startRow < bounds.Min.Y {
		panic("Start row is out of image's bounds")
	}

	resultPixels := make([][]color.RGBA, endRow-startRow+1)
	radius := (len(kernel) - 1) / 2

	// Iterate over each pixel from the starting to the ending row of pixels
	for y := startRow; y <= endRow; y++ {
		if resultPixels[y] == nil {
			resultPixels[y] = make([]color.RGBA, bounds.Max.X-bounds.Min.X)
		}
		pixelRow := resultPixels[y]
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var r, g, b float64

			// Iterate over each surrounding pixel within the specified radius
			for i := -radius; i <= radius; i++ {
				for j := -radius; j <= radius; j++ {
					neighborX := x + i
					neighborY := y + j
					if neighborX < bounds.Min.X || neighborX >= bounds.Max.X ||
						neighborY < bounds.Min.Y || neighborY >= bounds.Max.Y {
						// Ignore neighbor pixels out of range
						continue
					}
					// Get the color of a surrounding pixel
					neighborR, neighborG, neighborB, _ := img.At(neighborX, neighborY).RGBA()

					// Apply the Gaussian blur weight to the surrounding pixel color
					weight := kernel[i+radius][j+radius]
					r += float64(neighborR>>8) * weight
					g += float64(neighborG>>8) * weight
					b += float64(neighborB>>8) * weight
				}
			}
			// Set the color of the destination pixel
			pixelRow[x] = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
		}
	}

	return resultPixels
}

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
					neighborX := x + i
					neighborY := y + j
					if neighborX < bounds.Min.X || neighborX >= bounds.Max.X ||
						neighborY < bounds.Min.Y || neighborY >= bounds.Max.Y {
						// Ignore neighbor pixels out of range
						continue
					}
					neighborPixelColor := img.At(neighborX, neighborY)
					neighborR, neighborG, neighborB, _ := neighborPixelColor.RGBA()

					// Apply the Gaussian blur weight to the surrounding pixel color
					blurWeight := kernel[i+radius][j+radius]
					r += float64(neighborR>>8) * blurWeight
					g += float64(neighborG>>8) * blurWeight
					b += float64(neighborB>>8) * blurWeight
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
	return 1 / (2 * math.Pi * sqrSigma) * math.Exp(-(float64(x*x+y*y))/(2*sqrSigma))
}

// Function main to test  
func main2() {
	fmt.Println("Start")

	var img image.Image = loadImage("assets/original.jpg")
	destImg := image.NewRGBA(img.Bounds())
	// Apply the Gaussian blur
	// convolveRegion(img, destImg, image.Rect(50, 100, 200, 190), generateKernel(5))
	start := time.Now()
	convolveRange(img, destImg, 10, 4000, generateKernel(5))
	end := time.Now()
	saveImage("assets/original_blurred.png", destImg, JPEG)
	fmt.Println("Blurred image generated")
	fmt.Println("Time: ", end.Sub(start).Milliseconds())
}
