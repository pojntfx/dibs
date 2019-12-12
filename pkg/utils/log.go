package utils

import (
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

func PipeLogErrorNonPlatformSpecific(message string, err error, output string) {
	log.Fatal(message, rz.String("output", output), rz.Err(err))
}

func PipeLogError(message string, err error, platform, output string) {
	log.Fatal(message, rz.String("platform", platform), rz.String("output", output), rz.Err(err))
}

func PipeLogErrorPlatformNotFound(platform interface{}, err error) {
	log.Fatal("Platform(s) not found in configuration file", rz.Any("platform", platform), rz.Err(err))
}
