package main

import (
	"flag"
	"github.com/pojntfx/dibs/pkg/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// Config is a dibs configuration
type Config struct {
	Targets []struct {
		Name string `yaml:"name"`
		Helm struct {
			Src  string `yaml:"src"`
			Dist string `yaml:"dist"`
		}
		DockerManifest string `yaml:"dockerManifest"`
		Platforms      []struct {
			Identifier string `yaml:"identifier"`
			Paths      struct {
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
				Publish          string `yaml:"publish"`
				Start            string `yaml:"start"`
			}
			Docker struct {
				Build            dockerConfig `yaml:"build"`
				UnitTests        dockerConfig `yaml:"unitTests"`
				IntegrationTests dockerConfig `yaml:"integrationTests"`
				ChartTests       dockerConfig `yaml:"chartTests"`
				Publish          dockerConfig `yaml:"publish"`
			}
		}
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
		buildManifest    bool
		buildChart       bool
		unitTests        bool
		integrationTests bool
		imageTests       bool
		chartTests       bool
		publish          bool
		pushBinary       bool
		pushImage        bool
		pushManifest     bool
		pushChart        bool
		docker           bool
		target           string
		platform         string
		skipTests        bool
	)

	flag.StringVar(&configFilePath, "configFile", "dibs.yaml", "The config file to use")
	flag.StringVar(&context, "context", "", "The config file to use")
	flag.BoolVar(&docker, "docker", false, "Run in Docker")
	flag.BoolVar(&dev, "dev", false, "Start the development flow for the project")
	flag.BoolVar(&skipTests, "skipTests", false, "Skip the tests for the project")
	flag.BoolVar(&generateSources, "generateSources", false, "Generate the sources for the project")
	flag.BoolVar(&build, "build", false, "Build the project")
	flag.BoolVar(&buildImage, "buildImage", false, "Build the Docker image of the project")
	flag.BoolVar(&buildManifest, "buildManifest", false, `Build a Docker manifest.
It will add all images of the specified platforms; to add all, set -platform to "*".`)
	flag.BoolVar(&unitTests, "unitTests", false, "Run the unit tests of the project")
	flag.BoolVar(&integrationTests, "integrationTests", false, "Run the integration tests of the project")
	flag.BoolVar(&imageTests, "imageTests", false, "Run the image tests of the project")
	flag.BoolVar(&chartTests, "chartTests", false, "Run the chart tests of the project")
	flag.BoolVar(&publish, "publish", false, "Publish the project")
	flag.BoolVar(&pushImage, "pushImage", false, "Push the Docker image of the project")
	flag.BoolVar(&pushManifest, "pushManifest", false, "Push the Docker manifest of the project")
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
	flag.StringVar(&target, "target", "linux", `The name of the target to use.
This may also be set with the TARGET env variable; a value of "*" runs all targets.`)
	flag.StringVar(&platform, "platform", "linux/amd64", `The identifier of the platform to use.
This may also be set with the TARGETPLATFORM env variable; a value of "*" runs for all platforms.`)
	flag.Parse()

	// Normalize the environment and pass on env variables
	if targetFromEnv := os.Getenv("TARGET"); targetFromEnv != "" {
		target = targetFromEnv
	}
	if platformFromEnv := os.Getenv("TARGETPLATFORM"); platformFromEnv != "" {
		platform = platformFromEnv
	}
	envVariablesToSet := [][]string{
		{"DOCKER_CLI_EXPERIMENTAL", "enabled"},
		{"DOCKER_BUILDKIT", "1"},
	}
	for _, envVariableToSet := range envVariablesToSet {
		if err := os.Setenv(envVariableToSet[0], envVariableToSet[1]); err != nil {
			log.Fatal(err)
		}
	}

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

	configs := Config{}
	if err := yaml.Unmarshal(configFile, &configs); err != nil {
		log.Fatal(err)
	}

	stdoutChan, stderrChan := make(chan string), make(chan string)

	for _, targetConfig := range configs.Targets {
		if targetConfig.Name == target || target == "*" {
			if err := os.Setenv("TARGET", target); err != nil {
				log.Fatal(err)
			}

			if buildManifest {
				var images []string

				for _, platformConfig := range targetConfig.Platforms {
					if platformConfig.Identifier == platform || platform == "*" {
						images = append(images, platformConfig.Docker.Build.Tag)
					}
				}

				d := utils.NewDockerManager(context, stdoutChan, stderrChan)

				go handleStdoutAndStderr(stdoutChan, stderrChan)

				if err := d.BuildManifest(targetConfig.DockerManifest, images); err != nil {
					log.Fatal(err)
				}
			}

			if pushManifest {
				d := utils.NewDockerManager(context, stdoutChan, stderrChan)

				go handleStdoutAndStderr(stdoutChan, stderrChan)

				if err := d.PushManifest(targetConfig.DockerManifest); err != nil {
					log.Fatal(err)
				}
			}

			if buildChart {
				if err := os.MkdirAll(filepath.Join(context, targetConfig.Helm.Dist), 0777); err != nil {
					log.Fatal(err)
				}

				h := utils.NewHelmManager(context, stdoutChan, stderrChan)

				go handleStdoutAndStderr(stdoutChan, stderrChan)

				if err := h.Build(filepath.Join(context, targetConfig.Helm.Src), filepath.Join(targetConfig.Helm.Dist)); err != nil {
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
					filepath.Join(context, targetConfig.Helm.Dist),
					filepath.Join(os.TempDir(), "dibs-push-chart-repo"),
				); err != nil {
					log.Fatal(err)
				}
			}

			for _, platformConfig := range targetConfig.Platforms {
				if platformConfig.Identifier == platform || platform == "*" {
					if err := os.Setenv("TARGETPLATFORM", platformConfig.Identifier); err != nil {
						log.Fatal(err)
					}

					if dev {
						go handleStdoutAndStderr(stdoutChan, stderrChan)

						commandFlow := utils.NewCommandFlow([]string{
							platformConfig.Commands.GenerateSources,
							platformConfig.Commands.Build,
							platformConfig.Commands.UnitTests,
							platformConfig.Commands.IntegrationTests,
							platformConfig.Commands.Start,
						}, context, stdoutChan, stderrChan)
						if skipTests {
							commandFlow = utils.NewCommandFlow([]string{
								platformConfig.Commands.GenerateSources,
								platformConfig.Commands.Build,
								platformConfig.Commands.Start,
							}, context, stdoutChan, stderrChan)
						}

						interrupt := make(chan os.Signal, 2)
						signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
						go func() {
							<-interrupt

							// Allow manually killing the process
							go func() {
								<-interrupt

								os.Exit(1)
							}()

							log.Println("Gracefully stopping command flow (this might take a few seconds)")

							if err := commandFlow.Stop(); err != nil {
								log.Fatal(err)
							}

							os.Exit(0) // The path watcher is blocking
						}()

						if err := commandFlow.Start(); err != nil {
							log.Fatal(err)
						}

						eventChan := make(chan string)

						pathWatcher := utils.NewPathWatcher(filepath.Join(context, platformConfig.Paths.Watch), filepath.Join(context, platformConfig.Paths.Include), eventChan)

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
						runCommandWithLog(platformConfig.Commands.GenerateSources, context, stdoutChan, stderrChan)
					}

					if build {
						if docker {
							d := utils.NewDockerManager(context, stdoutChan, stderrChan)

							go handleStdoutAndStderr(stdoutChan, stderrChan)

							if err := d.Build(filepath.Join(context, platformConfig.Docker.Build.File), filepath.Join(context, platformConfig.Docker.Build.Context), platformConfig.Docker.Build.Tag); err != nil {
								log.Fatal(err)
							}

							if err := os.MkdirAll(filepath.Join(context, platformConfig.Paths.AssetOut, ".."), 0777); err != nil {
								log.Fatal(err)
							}

							if err := d.CopyFromImage(platformConfig.Docker.Build.Tag, platformConfig.Paths.AssetInImage, filepath.Join(context, platformConfig.Paths.AssetOut)); err != nil {
								log.Fatal(err)
							}
						} else {
							runCommandWithLog(platformConfig.Commands.Build, context, stdoutChan, stderrChan)
						}
					}

					if buildImage {
						d := utils.NewDockerManager(context, stdoutChan, stderrChan)

						go handleStdoutAndStderr(stdoutChan, stderrChan)

						if err := d.Build(filepath.Join(context, platformConfig.Docker.Build.File), filepath.Join(context, platformConfig.Docker.Build.Context), platformConfig.Docker.Build.Tag); err != nil {
							log.Fatal(err)
						}
					}

					if unitTests {
						if docker {
							buildAndRunDockerContainer("", context, platformConfig.Docker.UnitTests, false, stdoutChan, stderrChan)
						} else {
							runCommandWithLog(platformConfig.Commands.UnitTests, context, stdoutChan, stderrChan)
						}
					}

					if integrationTests {
						if docker {
							buildAndRunDockerContainer("", context, platformConfig.Docker.IntegrationTests, false, stdoutChan, stderrChan)
						} else {
							runCommandWithLog(platformConfig.Commands.IntegrationTests, context, stdoutChan, stderrChan)
						}
					}

					if imageTests {
						runCommandWithLog(platformConfig.Commands.ImageTests, context, stdoutChan, stderrChan)
					}

					if chartTests {
						if docker {
							buildAndRunDockerContainer("", context, platformConfig.Docker.ChartTests, true, stdoutChan, stderrChan)
						} else {
							runCommandWithLog(platformConfig.Commands.ChartTests, context, stdoutChan, stderrChan)
						}
					}

					if publish {
						if docker {
							buildAndRunDockerContainer("", context, platformConfig.Docker.Publish, false, stdoutChan, stderrChan)
						} else {
							runCommandWithLog(platformConfig.Commands.Publish, context, stdoutChan, stderrChan)
						}
					}

					if pushImage {
						d := utils.NewDockerManager(context, stdoutChan, stderrChan)

						go handleStdoutAndStderr(stdoutChan, stderrChan)

						if err := d.Push(platformConfig.Docker.Build.Tag); err != nil {
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
							filepath.Join(context, platformConfig.Paths.GitRepoRoot),
							filepath.Join(context, platformConfig.Paths.AssetOut),
						); err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		}
	}
}
