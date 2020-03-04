package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
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
			unix, err := exec.Command("uname", "-a").CombinedOutput()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Hello, world from UNIX \"%v\"!\n", strings.Split(string(unix), "\n")[0])
			time.Sleep(time.Second * 2)
		}
	}
}
