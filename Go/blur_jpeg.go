package main

import (
	"fmt"
	"os"
	"image"
	"image/jpeg"
	"image/color"
	"math"
)

func gaussianBlur(img image.Image, radius int) *image.RGBA {
	// Create a new RGBA image with the same bounds as the original
	bounds := img.Bounds()
	dest := image.NewRGBA(bounds)

	// Iterate over each pixel in the image
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			var r, g, b, a float64

			// Iterate over each surrounding pixel within the specified radius
			for i := -radius; i <= radius; i++ {
				for j := -radius; j <= radius; j++ {
					// Get the color of the surrounding pixel
					xn := x + i
					yn := y + j
					if xn < bounds.Min.X || xn >= bounds.Max.X || yn < bounds.Min.Y || yn >= bounds.Max.Y {
						continue
					}
					surroundingPixelColor := img.At(xn, yn)
					sr, sg, sb, sa := surroundingPixelColor.RGBA()

					// Apply the Gaussian blur weight to the surrounding pixel color
					r += float64(sr>>8) * gaussianWeight(i, j, radius)
					g += float64(sg>>8) * gaussianWeight(i, j, radius)
					b += float64(sb>>8) * gaussianWeight(i, j, radius)
					a += float64(sa>>8) * gaussianWeight(i, j, radius)
				}
			}

			// Set the color of the destination pixel
			dest.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	fmt.Println("Hi")

	return dest
}

func gaussianWeight(x, y, radius int) float64 {
	// Gaussian blur weight function
	return 1 / (2 * math.Pi * float64(radius)) * math.Exp(-(float64(x*x+y*y))/(2*float64(radius)))
}

func main() {
	fmt.Println("Start")
	file, _ := os.Open("original.jpg")
	//fmt.Println("1")
	img, _, _ := image.Decode(file)
	file.Close()
	fmt.Println("1")

	// Apply the Gaussian blur
	blurredImg := gaussianBlur(img, 10)

	// Save the blurred image to a new file
	file, _ = os.Create("blurred.jpg")
	jpeg.Encode(file, blurredImg, &jpeg.Options{Quality: 90})
}