package main

import (
	"flag"
	"github.com/pojntfx/dibs/v2/pkg/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// Config is a dibs configuration
type Config struct {
	Paths struct {
		Watch   string `yaml:"watch"`
		Include string `yaml:"include"`
	}
	Commands struct {
		GenerateSources  string `yaml:"generateSources"`
		Build            string `yaml:"build"`
		UnitTests        string `yaml:"unitTests"`
		IntegrationTests string `yaml:"integrationTests"`
		Start            string `yaml:"start"`
	}
}

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config-file", "dibs.yaml", "The config file to use")
	flag.Parse()

	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatal(err)
	}

	stdoutChan, stderrChan := make(chan string), make(chan string)

	commandFlow := utils.NewCommandFlow([]string{
		config.Commands.GenerateSources,
		config.Commands.Build,
		config.Commands.UnitTests,
		config.Commands.IntegrationTests,
		config.Commands.Start,
	}, stdoutChan, stderrChan)

	if err := commandFlow.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case stdout := <-stdoutChan:
				log.Println("STDOUT", stdout)
			case stderr := <-stderrChan:
				log.Println("STDERR", stderr)
			}
		}
	}()

	eventChan := make(chan string)

	pathWatcher := utils.NewPathWatcher(config.Paths.Watch, config.Paths.Include, eventChan)

	go func() {
		for {
			select {
			case <-eventChan:
				if err := commandFlow.Restart(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	defer func() {
		if err := commandFlow.Stop(); err != nil {
			log.Fatal(err)
		}
	}()
	if err := pathWatcher.Start(); err != nil {
		log.Fatal(err)
	}
}
