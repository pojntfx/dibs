package main

import (
	"bufio"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

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
		Start            string `yaml:"start"`
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
func RunCommand(cmd *exec.Cmd, stdoutChan, stderrChan chan string, errChan chan error, runInBackground bool) {
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

	if runInBackground {
		if err := cmd.Start(); err != nil {
			errChan <- err
		}

		return
	}

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

	generateSourcesCommands, buildCommands, unitTestsCommands, integrationTestsCommands, startCommands :=
		strings.Fields(config.Commands.GenerateSources),
		strings.Fields(config.Commands.Build),
		strings.Fields(config.Commands.UnitTests),
		strings.Fields(config.Commands.IntegrationTests),
		strings.Fields(config.Commands.Start)
	generateSourcesCommand, buildCommand, unitTestsCommand, integrationTestsCommand, startCommand :=
		GetCommandWrappedInSh(generateSourcesCommands),
		GetCommandWrappedInSh(buildCommands),
		GetCommandWrappedInSh(unitTestsCommands),
		GetCommandWrappedInSh(integrationTestsCommands),
		GetCommandWrappedInSh(startCommands)
	stdoutChan, stderrChan, errChan := make(chan string), make(chan string), make(chan error)

	done := false
	go func(done *bool) {
		log.Println("Generating sources")
		RunCommand(generateSourcesCommand, stdoutChan, stderrChan, errChan, false)
		log.Println("Building")
		RunCommand(buildCommand, stdoutChan, stderrChan, errChan, false)
		log.Println("Running unit tests")
		RunCommand(unitTestsCommand, stdoutChan, stderrChan, errChan, false)
		log.Println("Running integration tests")
		RunCommand(integrationTestsCommand, stdoutChan, stderrChan, errChan, false)

		log.Println("Starting app")
		startCommand.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		RunCommand(startCommand, stdoutChan, stderrChan, errChan, true)

		_ = startCommand.Wait()

		*done = true

		close(stdoutChan)
		close(stderrChan)
		close(errChan)
	}(&done)

	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt

		go func() {
			<-interrupt

			os.Exit(1)
		}()

		processGroupID, err := syscall.Getpgid(startCommand.Process.Pid)
		if err != nil {
			log.Fatal(err)
		}

		if err := syscall.Kill(processGroupID, syscall.SIGKILL); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {
		case stdout := <-stdoutChan:
			if done {
				return
			}

			log.Println("STDOUT", stdout)
		case stderr := <-stderrChan:
			if done {
				return
			}

			log.Println("STDERR", stderr)
		case err := <-errChan:
			if done {
				return
			}

			log.Println("ERR", err)
		}
	}
}
