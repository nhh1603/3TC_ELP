package main

import (
	"fmt"
	"image"
	"time"
)

const TASKS_COUNT = 50
const WORKERS_COUNT = 8

func worker(fromImg image.Image, toImg *image.RGBA, taskChan chan []int, resChan chan int) {
	for t := range taskChan {
		convolveRange(fromImg, toImg, t[0], t[1], generateKernel(5))
		resChan <- 1
	}
}

func generateImage(fromImg image.Image, toImg *image.RGBA) {
	totalPixels := fromImg.Bounds().Max.X*fromImg.Bounds().Max.Y
	taskPerWorker := totalPixels / TASKS_COUNT

	taskChan := make(chan []int, TASKS_COUNT)
	resChan := make(chan int, TASKS_COUNT)

	for i := 0; i < WORKERS_COUNT; i++ {
		go worker(fromImg, toImg, taskChan, resChan)
	}

	for j := 0; j < TASKS_COUNT; j++ {
		taskChan <- []int{j*taskPerWorker, (j+1)*taskPerWorker-1}
	}
	close(taskChan)

	// Assure that all goroutines have completed
	for a := 0; a < TASKS_COUNT; a++ {
		<-resChan
	}
	
	saveImage("assets/original_blurred.png", toImg, JPEG)
}

func main() {
	fmt.Println("Start")

	var img image.Image = loadImage("assets/original.jpg")
	destImg := image.NewRGBA(img.Bounds())

	start := time.Now()
	generateImage(img, destImg)
	end := time.Now()
	fmt.Println("Blurred image generated")
	fmt.Println("Time: ", end.Sub(start).Milliseconds())
}