package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func loadImg(filePath string) *image.NRGBA {
	imgFile, err := os.Open(filePath)
	defer imgFile.Close()
	if err != nil {
		fmt.Println("Cannot read file:", err)
		return nil
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Cannot decode file:", err)
		return nil
	}
	return img.(*image.NRGBA)
}

func saveImg(filePath string, img *image.NRGBA) {
	imgFile, err := os.Create(filePath)
	defer imgFile.Close()
	if err != nil {
		fmt.Println("Cannot create file:", err)
		return
	}
	png.Encode(imgFile, img.SubImage(img.Rect))
}
