package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	var (
		keepRunning bool
	)

	flag.BoolVar(&keepRunning, "keep-running", true, "Whether to keep the app running")

	flag.Parse()

	if keepRunning {
		for {
			time.Sleep(time.Second * 2)
			fmt.Println("Hello, world!")
		}
	}
}
