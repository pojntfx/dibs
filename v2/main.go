package main

import (
	"flag"
	"github.com/pojntfx/dibs/v2/pkg/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
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
	Docker struct {
		Build dockerConfig
	}
}

type dockerConfig struct {
	File    string `yaml:"file"`
	Context string `yaml:"context"`
	Tag     string `yaml:"tag"`
}

func runCommandWithLog(execLine string, stdoutChan, stderrChan chan string) {
	command := utils.NewManageableCommand(execLine, stdoutChan, stderrChan)

	if err := command.Start(); err != nil {
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

	if err := command.Wait(); err != nil {
		if err.Error() == "exit status 2" { // -help
			return
		}
		log.Fatal(err)
	}
}

func main() {
	var (
		configFilePath   string
		dev              bool
		generateSources  bool
		build            bool
		buildImage       bool
		unitTests        bool
		integrationTests bool
	)
	flag.StringVar(&configFilePath, "configFile", "dibs.yaml", "The config file to use")
	flag.BoolVar(&dev, "dev", false, "Run the development mode for the project")
	flag.BoolVar(&generateSources, "generateSources", false, "Generate the sources for the project")
	flag.BoolVar(&build, "build", false, "Build the project")
	flag.BoolVar(&buildImage, "buildImage", false, "Build the Docker image for the project")
	flag.BoolVar(&unitTests, "unitTests", false, "Run the unit tests to the project")
	flag.BoolVar(&integrationTests, "integrationTests", false, "Run the integration tests to the project")
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

	if dev {
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

	if generateSources {
		runCommandWithLog(config.Commands.GenerateSources, stdoutChan, stderrChan)
	}

	if build {
		runCommandWithLog(config.Commands.Build, stdoutChan, stderrChan)
	}

	if buildImage {
		d := utils.NewDockerManager(stdoutChan, stderrChan)

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

		if err := d.Build(filepath.Join(configFilePath, "..", config.Docker.Build.File), filepath.Join(configFilePath, "..", config.Docker.Build.Context), config.Docker.Build.Tag); err != nil {
			log.Fatal(err)
		}
	}

	if unitTests {
		runCommandWithLog(config.Commands.UnitTests, stdoutChan, stderrChan)
	}

	if integrationTests {
		runCommandWithLog(config.Commands.IntegrationTests, stdoutChan, stderrChan)
	}
}
