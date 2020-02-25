package main

import (
	"bufio"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config is the main configuration
type Config struct {
	Paths struct {
		Watch  string `yaml:"watch"`
		Ignore string `yaml:"ignore"`
	}
	Commands struct {
		GenerateSources  string `yaml:"generateSources"`
		Build            string `yaml:"build"`
		UnitTests        string `yaml:"unitTests"`
		IntegrationTests string `yaml:"integrationTests"`
	}
}

// Read reads from a reader till EOF
func Read(reader io.Reader, outChan chan string) {
	bufStdout := bufio.NewReader(reader)

	for {
		line, _, err := bufStdout.ReadLine()
		if err != nil {
			return
		}

		outChan <- string(line)
	}
}

// RunCommand runs a command and sends it's outputs into channels
func RunCommand(cmd *exec.Cmd, stdoutChan, stderrChan chan string, errChan chan error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errChan <- err

		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		errChan <- err

		return
	}

	go Read(stdout, stdoutChan)
	go Read(stderr, stderrChan)

	if err := cmd.Run(); err != nil {
		errChan <- err
	}
}

// GetCommandWrappedInSh returns a command wrapped in a call to `sh`
func GetCommandWrappedInSh(args []string) *exec.Cmd {
	wrappedArgs := append([]string{"sh", "-c"}, strings.Join(args, " "))

	return exec.Command(wrappedArgs[0], wrappedArgs[1:]...)
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

	generateSourcesCommands, buildCommands, unitTestsCommands, integrationTestsCommands :=
		strings.Fields(config.Commands.GenerateSources),
		strings.Fields(config.Commands.Build),
		strings.Fields(config.Commands.UnitTests),
		strings.Fields(config.Commands.IntegrationTests)
	generateSourcesCommand, buildCommand, unitTestsCommand, integrationTestsCommand :=
		GetCommandWrappedInSh(generateSourcesCommands),
		GetCommandWrappedInSh(buildCommands),
		GetCommandWrappedInSh(unitTestsCommands),
		GetCommandWrappedInSh(integrationTestsCommands)
	stdoutChan, stderrChan, errChan := make(chan string), make(chan string), make(chan error)

	done := false
	go func(done *bool) {
		log.Println("Generating sources")
		RunCommand(generateSourcesCommand, stdoutChan, stderrChan, errChan)
		log.Println("Building")
		RunCommand(buildCommand, stdoutChan, stderrChan, errChan)
		log.Println("Running unit tests")
		RunCommand(unitTestsCommand, stdoutChan, stderrChan, errChan)
		log.Println("Running integration tests")
		RunCommand(integrationTestsCommand, stdoutChan, stderrChan, errChan)

		*done = true
	}(&done)

	for {
		select {
		case stdout := <-stdoutChan:
			log.Println("STDOUT", stdout)
		case stderr := <-stderrChan:
			log.Println("STDERR", stderr)
		case err := <-errChan:
			log.Fatal("ERR", err)
		default:
			if done {
				return
			}
		}
	}
}
