package main

import "fmt"

func main() {
	fmt.Println("Start")

	result := cfshttplogging.SampleFuncFromConnor()
	fmt.Printf("result: %v\n", result)

	fmt.Println("End")
}
