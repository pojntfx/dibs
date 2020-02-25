package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config is the main configuration
type Config struct {
	Paths struct {
		Watch  string
		Ignore string
	}
	Commands struct {
		GenerateSources  string
		Build            string
		UnitTests        string
		IntegrationTests string
	}
}

func main() {
	var (
		configFilePath string
	)

	flag.StringVar(&configFilePath, "config-file", "dibs.yaml", "Config file to use")

	flag.Parse()

	configFileContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}

	if err := yaml.Unmarshal(configFileContent, &config); err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)
}
