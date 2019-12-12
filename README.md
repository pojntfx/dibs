# dibs

System for distributed polyglot, multi-module and multi-architecture development.

[![pipeline status](https://gitlab.com/pojntfx/dibs/badges/master/pipeline.svg)](https://gitlab.com/pojntfx/dibs/commits/master)

## Installation

### From Prebuilt Binaries

Prebuilt binaries are available on the [releases page](https://github.com/pojntfx/dibs/releases/latest).

### From Go

```bash
% go get github.com/pojntfx/dibs
```

## Usage

```bash
% dibs
System for distributed polyglot, multi-module and multi-architecture development

Usage:
  dibs [command]

Available Commands:
  help        Help about any command
  pipeline    Pipeline building blocks

Flags:
  -e, --executor string       Executor to run on ("docker"|"native") (default "native")
  -h, --help                  help for dibs
  -p, --platform string       Platform to run on ("all" runs on all platforms specified in configuration file) (default "all")
  -c, --redis-prefix string   Redis channel prefix to use (default "dibs")
  -u, --redis-url string      URL of the Redis instance to use (default "localhost:6379")

Use "dibs [command] --help" for more information about a command.
```

## License

dibs (c) 2019 Felix Pojtinger

SPDX-License-Identifier: AGPL-3.0
