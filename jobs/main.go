package main

import (
	"fmt"
	"log"

	endpointupdater "github.com/fantasticatif/health_monitor/jobs/endpoint_updater"
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

	if len(args.args) == 0 {
		log.Fatal("action argument is missing, allowed values are endpoint_updater")
	}
	if args.args[0] == "endpoint_updater" {
		endpointupdater.Run()
	} else {
		log.Fatal("invalid action argument is provided, allowed values are endpoint_updater")
	}
}
