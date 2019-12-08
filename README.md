# godibs

System for distributed multi-module, multi-architecture development with Go.

[![pipeline status](https://gitlab.com/pojntfx/godibs/badges/master/pipeline.svg)](https://gitlab.com/pojntfx/godibs/commits/master)

## Installation

```bash
% go get github.com/pojntfx/godibs
```

## Usage

### Server

```bash
% godibs server --help
Start the server

Usage:
  godibs server [flags]

Flags:
      --dir-repos string   Directory in which the Git repos should be stored (default "/tmp/godibs/gitrepos/8f1e43ab-f35f-4d72-a4cb-3ccc861417fb")
  -h, --help               help for server
      --path string        HTTP path prefix for the served Git repos (default "/repos")
      --port string        Port on which the Git repos should be served (default "25000")

Global Flags:
      --redis-prefix string   Redis channel prefix (default "godibs")
      --redis-url string      URL of the Redis instance to use (default "localhost:6379")
```

### Client

```bash
% godibs client --help
Start the client

Usage:
  godibs client [flags]

Flags:
      --cmd-build string      Command to run to build the module (default "go build ./...")
      --cmd-start string      Command to run to start the module (default "go run main.go")
      --cmd-test string       Command to run to test the module (default "go test ./...")
      --dir-pull string       Directory to pull the modules to (default "/tmp/godibs/pull/fefd9ca1-8a0b-4e16-8a62-78a13bcec255")
      --dir-push string       Temporary directory to put the module into before pushing (default "/tmp/godibs/push/fefd9ca1-8a0b-4e16-8a62-78a13bcec255")
      --dir-src string        Directory in which the source code of the module to push resides (default ".")
      --dir-watch string      Directory to watch for changes (default ".")
      --git-base-url string   Base URL of the sync remote (default "http://localhost:25000/repos")
  -h, --help                  help for client
      --modules-file string   Go module file of the module to push (default "go.mod")
      --modules-pull string   Comma-seperated list of the names of the modules to pull
      --regex-ignore string   Regular expression for files to ignore (default "*.pb.go")

Global Flags:
      --redis-prefix string   Redis channel prefix (default "godibs")
      --redis-url string      URL of the Redis instance to use (default "localhost:6379")
```

## License

godibs (c) 2019 Felicitas Pojtinger

SPDX-License-Identifier: AGPL-3.0
