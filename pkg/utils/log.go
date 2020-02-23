package utils

import (
	"fmt"

	"gitlab.com/bloom42/libs/rz-go"
	"gitlab.com/bloom42/libs/rz-go/log"
)

// LogError handles a non-platform-specific non-fatal error
func LogError(message string, err error, output ...string) {
	if output != nil {
		fmt.Println(output)
	}
	log.Error(message, rz.Err(err))
}

// LogErrorFatal handles a non-platform-specific fatal error
func LogErrorFatal(message string, err error, output ...string) {
	if output != nil {
		fmt.Println(output)
	}
	log.Fatal(message, rz.Err(err))
}

// LogErrorFatalPlatformSpecific handles a fatal error
func LogErrorFatalPlatformSpecific(message string, err error, platform string, output ...string) {
	if output != nil {
		fmt.Println(output)
	}
	log.Fatal(message, rz.String("platform", platform), rz.Err(err))
}

// LogErrorFatalWithProfile handles a profile-specific fatal error
func LogErrorFatalWithProfile(message string, err error, platform, profile string, output ...string) {
	if output != nil {
		fmt.Println(output)
	}
	log.Fatal(message, rz.String("platform", platform), rz.String("profile", profile), rz.Err(err))
}

// LogErrorFatalPlatformNotFound handles a fatal error if one or more platforms can't be found
func LogErrorFatalPlatformNotFound(platform interface{}, err error) {
	log.Fatal("Platform(s) not found in configuration file", rz.Any("platforms", platform), rz.Err(err))
}

// LogErrorFatalCouldNotStopModule handles a fatal error a module can't be stopped
func LogErrorFatalCouldStopModule(err error) {
	log.Fatal("Could not stop module", rz.Err(err))
}

// LogErrorInfo handles a non-fatal error
func LogErrorInfo(message string, err error, platform, output string) {
	fmt.Println(output)
	log.Info(message, rz.String("platform", platform), rz.Err(err))
}

// LogErrorCouldNotBindFlag handles flag binding errors
func LogErrorCouldNotBindFlag(err error) {
	log.Fatal("Could bind flag", rz.Err(err))
}

// LogForModule logs for a module
func LogForModule(message, module string) {
	log.Info(message, rz.String("module", module))
}

// LogErrorForModuleFatal handles fatal errors for a module
func LogErrorForModuleFatal(message string, err error, module string) {
	log.Fatal(message, rz.Err(err), rz.String("module", module))
}

// LogErrorForModule handles errors for a module
func LogErrorForModule(message string, err error, module string) {
	log.Error(message, rz.Err(err), rz.String("module", module))
}
