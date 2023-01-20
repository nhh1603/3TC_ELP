package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func loadImage(filePath string) image.Image {
	imgFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Cannot read file:", err)
		os.Exit(-1)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Cannot decode file:", err)
		os.Exit(-1)
	}

	return img
}

type ImageFormat int8

const (
	PNG ImageFormat = iota
	JPEG
)

func saveImage(filePath string, img image.Image, format ImageFormat) {
	imgFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Cannot create file:", err)
		os.Exit(-1)
	}
	defer imgFile.Close()

	switch format {
	case PNG:
		png.Encode(imgFile, img)
	case JPEG:
		jpeg.Encode(imgFile, img, &jpeg.Options{Quality: 90})
	}
}
