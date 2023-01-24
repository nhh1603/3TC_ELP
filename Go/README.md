This is a project that can apply Gaussian Blur to a given image, written in Golang.

File imageTools.go contains all the necessary functions to work with images
File blurTools.go contains the blur algorithm
File main.go contains the main program and goroutines implementation to execute the work in parallel, in order to maximize the execution time and the usage of CPU's cores.

In order to test this project, please upload an image to the assets folder, change its directory address in main, then run the following command: go run main.go blurTools.go imageTools.go