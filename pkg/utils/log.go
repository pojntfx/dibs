package utils

import (
	"fmt"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

func PipeLogErrorFatalNonPlatformSpecific(message string, err error, output ...string) {
	if output != nil {
		fmt.Println(output)
	}
	log.Fatal(message, rz.Err(err))
}

func PipeLogErrorFatal(message string, err error, platform string, output ...string) {
	if output != nil {
		fmt.Println(output)
	}
	log.Fatal(message, rz.String("platform", platform), rz.Err(err))
}

func PipeLogErrorFatalWithProfile(message string, err error, platform, profile string, output ...string) {
	if output != nil {
		fmt.Println(output)
	}
	log.Fatal(message, rz.String("platform", platform), rz.String("profile", profile), rz.Err(err))
}

func PipeLogErrorFatalPlatformNotFound(platform interface{}, err error) {
	log.Fatal("Platform(s) not found in configuration file", rz.Any("platforms", platform), rz.Err(err))
}

func PipeLogErrorFatalCouldNotParseIP(ip string) {
	log.Fatal("Could not parse IP", rz.String("ip", ip))
}

func PipeLogErrorInfo(message string, err error, platform, output string) {
	fmt.Println(output)
	log.Info(message, rz.String("platform", platform), rz.Err(err))
}
