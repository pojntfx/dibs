# dibs

System for distributed polyglot, multi-module, multi-architecture development and CI/CD.

[![pipeline status](https://gitlab.com/pojntfx/dibs/badges/master/pipeline.svg)](https://gitlab.com/pojntfx/dibs/commits/master)

[![introduction video](https://img.youtube.com/vi/fUfW-z6fWZs/maxresdefault.jpg)](https://youtu.be/fUfW-z6fWZs)

## Overview

dibs, short for `di`stributed `b`uild `s`ystem, enables polyglot, multi-module, multi-architecture development and CI/CD without the configuration overhead one would normally need. During development, it can be called from either [Skaffold](https://skaffold.dev/) for cloud-native development or be used natively with `dibs -dev`.

## Installation

### Prebuilt Binaries

Prebuilt binaries are available on the [releases page](https://github.com/pojntfx/dibs/releases/latest).

### Go Package

A Go package [is available](https://pkg.go.dev/github.com/pojntfx/dibs?tab=doc).

## Usage

dibs is configured by using a [config file](./test-app/dibs.yaml).

To use dibs with GitLab CI/CD, see the [example GitLab CI/CD configuration file](./.gitlab-ci.yml).

```bash
% dibs -help
Usage of dibs:
  -build
    	Build the project
  -buildChart
    	Build the Helm chart of the project
  -buildImage
    	Build the Docker image of the project
  -buildManifest
    	Build a Docker manifest.
    	It will add all images of the specified platforms; to add all, set -platform to "*".
  -chartTests
    	Run the chart tests of the project
  -configFile string
    	The config file to use (default "dibs.yaml")
  -context string
    	The config file to use
  -dev
    	Start the development flow for the project
  -docker
    	Run in Docker
  -generateSources
    	Generate the sources for the project
  -imageTests
    	Run the image tests of the project
  -integrationTests
    	Run the integration tests of the project
  -platform string
    	The identifier of the platform to use.
    	This may also be set with the TARGETPLATFORM env variable; a value of "*" runs for all platforms. (default "linux/amd64")
  -publish
    	Publish the project
  -pushBinary
    	Push the binary of the project.
    	This command requires the following env variables to be set:
    	- DIBS_GITHUB_USER_NAME
    	- DIBS_GITHUB_TOKEN
    	- DIBS_GITHUB_REPOSITORY
  -pushChart
    	Push the Helm chart of the project.
    	This command requires the following env variables to be set:
    	- DIBS_GIT_USER_NAME
    	- DIBS_GIT_USER_EMAIL
    	- DIBS_GIT_COMMIT_MESSAGE
    	- DIBS_GITHUB_USER_NAME
    	- DIBS_GITHUB_TOKEN
    	- DIBS_GITHUB_REPOSITORY_NAME
    	- DIBS_GITHUB_REPOSITORY_URL
    	- DIBS_GITHUB_PAGES_URL
  -pushImage
    	Push the Docker image of the project
  -pushManifest
    	Push the Docker manifest of the project
  -skipTests
    	Skip the tests for the project
  -target string
    	The name of the target to use.
    	This may also be set with the TARGET env variable; a value of "*" runs all targets. (default "linux")
  -unitTests
    	Run the unit tests of the project
```

## License

dibs (c) 2020 Felix Pojtinger

SPDX-License-Identifier: AGPL-3.0
