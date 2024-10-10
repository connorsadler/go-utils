package main

import (
	"fmt"

	"github.com/connorsadler/go-utils/cfshttplogging"
)

func main() {
	fmt.Println("Start")

	result := cfshttplogging.SampleFuncFromConnor()
	fmt.Printf("result: %v\n", result)

	fmt.Println("End")
}
