package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup // instanciation de notre structure WaitGroup

func run(name string) {
	defer wg.Done()
	for i := 0; i < 3; i++ {
		time.Sleep(3 * time.Second)
		fmt.Println(name, " : ", i)
	}
}

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

func computeSumAsync2(array *[]int, resChan chan int, taskChan chan *[2]int) {
	for bounds := range taskChan {
		sum := 0
		for i := bounds[0]; i < bounds[1]; i++ {
			sum += (*array)[i]
		}
		resChan <- sum
	}
}

const ARRAY_SIZE = 10000000

func runResultChannel() {
	arr := make([]int, ARRAY_SIZE)

	for i := 0; i < ARRAY_SIZE; i++ {
		arr[i] = rand.Intn(10)
	}

	println("RESULT CHANNEL VERSION")
	println("Sync version ****")

	startTime := time.Now()
	println("Result: ", computeSum(&arr))
	endTime := time.Now()
	println("Time: ", endTime.Sub(startTime).Milliseconds())
	println()

	const WORKERS_COUNT = 4
	resChan := make(chan int, WORKERS_COUNT)

	wg.Add(WORKERS_COUNT)
	startTime = time.Now()
	for i := 0; i < WORKERS_COUNT; i++ {
		startIndex := (ARRAY_SIZE / WORKERS_COUNT) * i
		end := (ARRAY_SIZE / WORKERS_COUNT) * (i + 1)
		go computeSumAsync(&arr, startIndex, end, resChan)
	}
	wg.Wait()
	close(resChan)

	result := 0
	for partialRes := range resChan {
		result += partialRes
	}
	endTime = time.Now()

	println("Async version ****")
	println("Result: ", result)
	println("Time: ", endTime.Sub(startTime).Milliseconds())
}

func runTaskChannel() {
	arr := make([]int, ARRAY_SIZE)

	for i := 0; i < ARRAY_SIZE; i++ {
		arr[i] = rand.Intn(10)
	}

	println("TASK CHANNEL VERSION")
	println("Sync version ****")

	startTime := time.Now()
	println("Result: ", computeSum(&arr))
	endTime := time.Now()
	println("Time: ", endTime.Sub(startTime).Milliseconds())
	println()

	const WORKERS_COUNT = 4
	const TASK_COUNT = 50
	resChan := make(chan int, WORKERS_COUNT)
	taskChan := make(chan *[2]int, TASK_COUNT)

	startTime = time.Now()
	for i := 0; i < WORKERS_COUNT; i++ {
		go computeSumAsync2(&arr, resChan, taskChan)
	}

	// Put the tasks in the channel and then close the channel
	for i := 0; i < TASK_COUNT; i++ {
		startIndex := (ARRAY_SIZE / TASK_COUNT) * i
		end := (ARRAY_SIZE / TASK_COUNT) * (i + 1)
		taskChan <- &[2]int{startIndex, end}
	}
	close(taskChan)

	retrievedResCount := 0
	result := 0
	for partialRes := range resChan {
		result += partialRes
		retrievedResCount++
		if retrievedResCount == TASK_COUNT {
			break
		}
	}
	endTime = time.Now()

	println("Async version ****")
	println("Result: ", result)
	println("Time: ", endTime.Sub(startTime).Milliseconds())
}

func main() {
	runTaskChannel()
}
