package main

import (
	"fmt"
	"os"
	"strings"
)

type argsHandler struct {
	excecutionPath string
	args           []string
	options        map[string]string
}

func getArgsHandler() argsHandler {
	args := os.Args
	fmt.Printf("args: %s", args)
	h := argsHandler{}
	options := make(map[string]string)
	commandArgs := []string{}
	for indx, val := range args {
		if indx == 0 {
			h.excecutionPath = val
		} else {
			if strings.HasPrefix(val, "--") {
				val_trimmed, _ := strings.CutPrefix(val, "--")
				position := strings.Index(val_trimmed, "=")
				if position < 0 {
					options[val_trimmed] = "true"
				} else {
					key := val_trimmed[:position]
					val := val_trimmed[position+1:]
					options[key] = val
				}
			} else {
				commandArgs = append(commandArgs, val)
			}

		}
	}
	h.options = options
	h.args = commandArgs
	return h
}
