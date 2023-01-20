package main

import (
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup // instanciation de notre structure WaitGroup

func computeSum(array *[]int) int {
	sum := 0
	for _, v := range *array {
		sum += v
	}
	return sum
}

func computeSumAsync(array *[]int, start, end int, resChan chan int) {
	defer wg.Done()
	sum := 0
	for i := start; i < end; i++ {
		sum += (*array)[i]
	}
	resChan <- sum
}

func main2() {

	const ARRAY_SIZE = 100000000
	arr := make([]int, ARRAY_SIZE)

	for i := 0; i < ARRAY_SIZE; i++ {
		arr[i] = rand.Intn(10)
	}

	println("Sync version ****")

	startTime := time.Now()
	println("Result: ", computeSum(&arr))
	endTime := time.Now()
	println("Time: ", endTime.Sub(startTime).Milliseconds())
	println()

	const WORKERS_COUNT = 8
	resChan := make(chan int, WORKERS_COUNT)

	wg.Add(WORKERS_COUNT)
	startTime = time.Now()
	for i := 0; i < WORKERS_COUNT; i++ {
		startIndex := (ARRAY_SIZE / WORKERS_COUNT) * i
		end := (ARRAY_SIZE / WORKERS_COUNT) * (i + 1)
		go computeSumAsync(&arr, startIndex, end, resChan)
	}
	wg.Wait()
	endTime = time.Now()
	close(resChan)

	result := 0
	for partialRes := range resChan {
		result += partialRes
	}

	println("Async version ****")
	println("Result: ", result)
	println("Time: ", endTime.Sub(startTime).Milliseconds())

}
