package main

import (
	"fmt"
)

func main() {
	fmt.Println("job begin")

	args := getArgsHandler()

	fmt.Println("--- args")
	for index, val := range args.args {
		fmt.Printf("Index: %d , val: %s \n", index, val)
	}

	fmt.Println("--- options")
	for index, val := range args.options {
		fmt.Printf("Key: %s , val: %s \n", index, val)
	}
}
