package main

import (
	"fmt"
	"image"
	"sync"
	"time"
)

const TASKS_COUNT = 50
const WORKERS_COUNT = 8

var waitGroup sync.WaitGroup

func worker(fromImg image.Image, toImg *image.RGBA, taskChan chan *[2]int, kernel [][]float64) {
	for t := range taskChan {
		convolveRange(fromImg, toImg, t[0], t[1], kernel)
	}
	waitGroup.Done()
}

func generateBlurImage(fromImg image.Image) *image.RGBA {
	fromBounds := fromImg.Bounds()
	totalPixels := fromBounds.Max.X * fromBounds.Max.Y
	destImg := image.NewRGBA(fromBounds)
	pixelsPerWorker := totalPixels / TASKS_COUNT
	kernel := generateGaussKernel(5)

	taskChan := make(chan *[2]int, WORKERS_COUNT)

	waitGroup.Add(WORKERS_COUNT)
	for i := 0; i < WORKERS_COUNT; i++ {
		go worker(fromImg, destImg, taskChan, kernel)
	}

	for i := 0; i < TASKS_COUNT; i++ {
		startPixel := i * pixelsPerWorker
		endPixel := (i+1)*pixelsPerWorker - 1
		taskChan <- &[2]int{startPixel, endPixel}
	}
	// The last task may have at most TASKS_COUNT additional pixels compared to other tasks,
	// but considering this const is relatively small, it's fine to let the last goroutine finish up.
	taskChan <- &[2]int{TASKS_COUNT * pixelsPerWorker, totalPixels - 1}
	close(taskChan)

	// Assure that all goroutines have completed
	waitGroup.Wait()

	return destImg
}

func main() {
	fmt.Println("Start")

	var img image.Image = loadImage("assets/original.jpg")

	start := time.Now()
	blurredImg := generateBlurImage(img)
	end := time.Now()
	saveImage("assets/original_blurred.jpg", blurredImg, JPEG)
	fmt.Println("Blurred image generated")
	fmt.Println("Time: ", end.Sub(start).Milliseconds())
}
