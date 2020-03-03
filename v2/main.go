package main

import (
	"flag"
	"github.com/pojntfx/dibs/v2/pkg/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Config is a dibs configuration
type Config struct {
	Paths struct {
		Watch        string `yaml:"watch"`
		Include      string `yaml:"include"`
		AssetInImage string `yaml:"assetInImage"`
		AssetOut     string `yaml:"assetOut"`
		GitRepoRoot  string `yaml:"gitRepoRoot"`
	}
	Commands struct {
		GenerateSources  string `yaml:"generateSources"`
		Build            string `yaml:"build"`
		UnitTests        string `yaml:"unitTests"`
		IntegrationTests string `yaml:"integrationTests"`
		ImageTests       string `yaml:"imageTests"`
		ChartTests       string `yaml:"chartTests"`
		Start            string `yaml:"start"`
	}
	Docker struct {
		Build            dockerConfig `yaml:"build"`
		UnitTests        dockerConfig `yaml:"unitTests"`
		IntegrationTests dockerConfig `yaml:"integrationTests"`
		ChartTests       dockerConfig `yaml:"chartTests"`
	}
	Helm struct {
		Src  string `yaml:"src"`
		Dist string `yaml:"dist"`
	}
}

type dockerConfig struct {
	File    string `yaml:"file"`
	Context string `yaml:"context"`
	Tag     string `yaml:"tag"`
}

func runCommandWithLog(execLine, dir string, stdoutChan, stderrChan chan string) {
	command := utils.NewManageableCommand(execLine, dir, stdoutChan, stderrChan)

	if err := command.Start(); err != nil {
		log.Fatal(err)
	}

	go handleStdoutAndStderr(stdoutChan, stderrChan)

	if err := command.Wait(); err != nil {
		if err.Error() == "exit status 2" { // -help
			return
		}
		log.Fatal(err)
	}
}

func buildAndRunDockerContainer(command, context string, config dockerConfig, privileged bool, stdoutChan, stderrChan chan string) {
	d := utils.NewDockerManager(context, stdoutChan, stderrChan)

	go handleStdoutAndStderr(stdoutChan, stderrChan)

	if err := d.Build(filepath.Join(context, config.File), filepath.Join(context, config.Context), config.Tag); err != nil {
		log.Fatal(err)
	}

	if err := d.Run(config.Tag, command, privileged); err != nil {
		log.Fatal(err)
	}
}

func handleStdoutAndStderr(stdoutChan, stderrChan chan string) {
	for {
		select {
		case stdout := <-stdoutChan:
			log.Println("STDOUT", stdout)
		case stderr := <-stderrChan:
			log.Println("STDERR", stderr)
		}
	}
}

func main() {
	var (
		configFilePath   string
		context          string
		dev              bool
		generateSources  bool
		build            bool
		buildImage       bool
		unitTests        bool
		integrationTests bool
		imageTests       bool
		chartTests       bool
		pushImage        bool
		docker           bool
		buildChart       bool
		pushChart        bool
		pushBinary       bool
	)

	flag.StringVar(&configFilePath, "configFile", "dibs.yaml", "The config file to use")
	flag.StringVar(&context, "context", "", "The config file to use")
	flag.BoolVar(&docker, "docker", false, "Run in Docker")
	flag.BoolVar(&dev, "dev", false, "Start the development flow for the project")
	flag.BoolVar(&generateSources, "generateSources", false, "Generate the sources for the project")
	flag.BoolVar(&build, "build", false, "Build the project")
	flag.BoolVar(&buildImage, "buildImage", false, "Build the Docker image of the project")
	flag.BoolVar(&unitTests, "unitTests", false, "Run the unit tests of the project")
	flag.BoolVar(&integrationTests, "integrationTests", false, "Run the integration tests of the project")
	flag.BoolVar(&imageTests, "imageTests", false, "Run the image tests of the project")
	flag.BoolVar(&chartTests, "chartTests", false, "Run the chart tests of the project")
	flag.BoolVar(&pushImage, "pushImage", false, "Push the Docker image of the project")
	flag.BoolVar(&buildChart, "buildChart", false, "Build the Helm chart of the project")
	flag.BoolVar(&pushChart, "pushChart", false, `Push the Helm chart of the project.
This command requires the following env variables to be set:
- DIBS_GIT_USER_NAME
- DIBS_GIT_USER_EMAIL
- DIBS_GIT_COMMIT_MESSAGE
- DIBS_GITHUB_USER_NAME
- DIBS_GITHUB_TOKEN
- DIBS_GITHUB_REPOSITORY_NAME
- DIBS_GITHUB_REPOSITORY_URL
- DIBS_GITHUB_PAGES_URL`)
	flag.BoolVar(&pushBinary, "pushBinary", false, `Push the binary of the project.
This command requires the following env variables to be set:
- DIBS_GITHUB_USER_NAME
- DIBS_GITHUB_TOKEN
- DIBS_GITHUB_REPOSITORY`)
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if context == "" {
		context = filepath.Join(pwd, configFilePath, "..")
	}

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
		}, context, stdoutChan, stderrChan)

		if err := commandFlow.Start(); err != nil {
			log.Fatal(err)
		}

		go handleStdoutAndStderr(stdoutChan, stderrChan)

		eventChan := make(chan string)

		pathWatcher := utils.NewPathWatcher(filepath.Join(context, config.Paths.Watch), filepath.Join(context, config.Paths.Include), eventChan)

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
		runCommandWithLog(config.Commands.GenerateSources, context, stdoutChan, stderrChan)
	}

	if build {
		if docker {
			d := utils.NewDockerManager(context, stdoutChan, stderrChan)

			go handleStdoutAndStderr(stdoutChan, stderrChan)

			if err := d.Build(filepath.Join(context, config.Docker.Build.File), filepath.Join(context, config.Docker.Build.Context), config.Docker.Build.Tag); err != nil {
				log.Fatal(err)
			}

			if err := os.MkdirAll(filepath.Join(context, config.Paths.AssetOut, ".."), 0777); err != nil {
				log.Fatal(err)
			}

			if err := d.CopyFromImage(config.Docker.Build.Tag, config.Paths.AssetInImage, filepath.Join(context, config.Paths.AssetOut)); err != nil {
				log.Fatal(err)
			}
		} else {
			runCommandWithLog(config.Commands.Build, context, stdoutChan, stderrChan)
		}
	}

	if buildImage {
		d := utils.NewDockerManager(context, stdoutChan, stderrChan)

		go handleStdoutAndStderr(stdoutChan, stderrChan)

		if err := d.Build(filepath.Join(context, config.Docker.Build.File), filepath.Join(context, config.Docker.Build.Context), config.Docker.Build.Tag); err != nil {
			log.Fatal(err)
		}
	}

	if unitTests {
		if docker {
			buildAndRunDockerContainer("", context, config.Docker.UnitTests, false, stdoutChan, stderrChan)
		} else {
			runCommandWithLog(config.Commands.UnitTests, context, stdoutChan, stderrChan)
		}
	}

	if integrationTests {
		if docker {
			buildAndRunDockerContainer("", context, config.Docker.IntegrationTests, false, stdoutChan, stderrChan)
		} else {
			runCommandWithLog(config.Commands.IntegrationTests, context, stdoutChan, stderrChan)
		}
	}

	if imageTests {
		runCommandWithLog(config.Commands.ImageTests, context, stdoutChan, stderrChan)
	}

	if chartTests {
		if docker {
			buildAndRunDockerContainer("", context, config.Docker.ChartTests, true, stdoutChan, stderrChan)
		} else {
			runCommandWithLog(config.Commands.ChartTests, context, stdoutChan, stderrChan)
		}
	}

	if pushImage {
		d := utils.NewDockerManager(context, stdoutChan, stderrChan)

		go handleStdoutAndStderr(stdoutChan, stderrChan)

		if err := d.Push(config.Docker.Build.Tag); err != nil {
			log.Fatal(err)
		}
	}

	if buildChart {
		if err := os.MkdirAll(filepath.Join(context, config.Helm.Dist), 0777); err != nil {
			log.Fatal(err)
		}

		h := utils.NewHelmManager(context, stdoutChan, stderrChan)

		go handleStdoutAndStderr(stdoutChan, stderrChan)

		if err := h.Build(filepath.Join(context, config.Helm.Src), filepath.Join(config.Helm.Dist)); err != nil {
			log.Fatal(err)
		}
	}

	if pushChart {
		h := utils.NewHelmManager(context, stdoutChan, stderrChan)

		go handleStdoutAndStderr(stdoutChan, stderrChan)

		if err := h.Push(
			os.Getenv("DIBS_GIT_USER_NAME"),
			os.Getenv("DIBS_GIT_USER_EMAIL"),
			os.Getenv("DIBS_GIT_COMMIT_MESSAGE"),
			os.Getenv("DIBS_GITHUB_USER_NAME"),
			os.Getenv("DIBS_GITHUB_TOKEN"),
			os.Getenv("DIBS_GITHUB_REPOSITORY_NAME"),
			os.Getenv("DIBS_GITHUB_REPOSITORY_URL"),
			os.Getenv("DIBS_GITHUB_PAGES_URL"),
			filepath.Join(context, config.Helm.Dist),
			filepath.Join(os.TempDir(), "dibs-push-chart-repo"),
		); err != nil {
			log.Fatal(err)
		}
	}

	if pushBinary {
		h := utils.NewBinaryManager(context, stdoutChan, stderrChan)

		go handleStdoutAndStderr(stdoutChan, stderrChan)

		if err := h.Push(
			os.Getenv("DIBS_GITHUB_USER_NAME"),
			os.Getenv("DIBS_GITHUB_TOKEN"),
			os.Getenv("DIBS_GITHUB_REPOSITORY"),
			filepath.Join(context, config.Paths.GitRepoRoot),
			filepath.Join(context, config.Paths.AssetOut),
		); err != nil {
			log.Fatal(err)
		}
	}
}
