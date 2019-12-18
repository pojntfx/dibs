package utils

import (
	"fmt"
	"gitlab.com/z0mbie42/rz-go/v2"
	"gitlab.com/z0mbie42/rz-go/v2/log"
)

// Handle a non-platform-specific fatal error in a pipe
func PipeLogErrorFatalNonPlatformSpecific(message string, err error, output ...string) {
	if output != nil {
		fmt.Println(output)
	}
	log.Fatal(message, rz.Err(err))
}

// Handle a fatal error in a pipe
func PipeLogErrorFatal(message string, err error, platform string, output ...string) {
	if output != nil {
		fmt.Println(output)
	}
	log.Fatal(message, rz.String("platform", platform), rz.Err(err))
}

// Handle a profile-specific fatal error in a pipe
func PipeLogErrorFatalWithProfile(message string, err error, platform, profile string, output ...string) {
	if output != nil {
		fmt.Println(output)
	}
	log.Fatal(message, rz.String("platform", platform), rz.String("profile", profile), rz.Err(err))
}

// Handle a fatal error in a pipe if one or more platforms can't be found
func PipeLogErrorFatalPlatformNotFound(platform interface{}, err error) {
	log.Fatal("Platform(s) not found in configuration file", rz.Any("platforms", platform), rz.Err(err))
}

// Handle a fatal error in a pipe if an IP can't be parsed
func PipeLogErrorFatalCouldNotParseIP(ip string) {
	log.Fatal("Could not parse IP", rz.String("ip", ip))
}

// Handle a non-fatal error in a pipe
func PipeLogErrorInfo(message string, err error, platform, output string) {
	fmt.Println(output)
	log.Info(message, rz.String("platform", platform), rz.Err(err))
}
