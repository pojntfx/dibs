# dibs

System for distributed polyglot, multi-module, multi-architecture development and CI/CD.

[![pipeline status](https://gitlab.com/pojntfx/dibs/badges/master/pipeline.svg)](https://gitlab.com/pojntfx/dibs/commits/master)

## Installation

### Prebuilt Binaries

Prebuilt binaries are available on the [releases page](https://github.com/pojntfx/dibs/releases/latest).

### Go Package

A Go package [is available](https://godoc.org/github.com/pojntfx/dibs).

### Docker Image

A Docker image is available on [Docker Hub](https://hub.docker.com/r/pojntfx/dibs).

### Helm Chart

A Helm chart is available in [@pojntfx's Helm chart repository](https://pojntfx.github.io/charts/).

## Usage

```bash
% dibs
System for distributed polyglot, multi-module, multi-architecture development and CI/CD

For full functionality, it requires the following binaries to be in PATH:

- "docker":     https://www.docker.com/
- "kubectl":    https://kubernetes.io/docs/reference/kubectl/
- "helm"        https://helm.sh/
- "skaffold"    https://skaffold.dev/
- "ghr"         https://github.com/tcnksm/ghr
- "cr"          https://github.com/helm/chart-releaser

If you want to support Dockerized multi-architecture builds, you'll also have to setup "qemu-user-static": https://github.com/multiarch/qemu-user-static

Usage:
  dibs [command]

Available Commands:
  dev         Develop the project
  help        Help about any command
  install     Install and start the project
  pipeline    Individual pipeline building blocks
  uninstall   Stop and uninstall the project

Flags:
  -f, --config-file string   Configuration file to use (default ".dibs.yml")
  -e, --executor string      Executor to run on ("docker"|"native") (default "native")
  -h, --help                 help for dibs
  -p, --platform string      Platform to run on ("all" runs on all platforms specified in configuration file) (default "all")

Use "dibs [command] --help" for more information about a command.
```

## License

dibs (c) 2019 Felix Pojtinger

SPDX-License-Identifier: AGPL-3.0
