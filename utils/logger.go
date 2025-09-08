package utils

import (
	"fmt"
)

var loggingEnabled = true

func SetLoggingEnabled(enabled bool) {
	loggingEnabled = enabled
}

func LogInfo(msg string) {
	if !loggingEnabled {
		return
	}
	fmt.Printf("\033[36m%s\033[0m\n", msg) // cyan
}

func LogWarn(msg string) {
	if !loggingEnabled {
		return
	}
	fmt.Printf("\033[33m%s\033[0m\n", msg) // yellow
}

func LogSuccess(msg string) {
	if !loggingEnabled {
		return
	}
	fmt.Printf("\033[32m%s\033[0m\n", msg) // green
}

func LogCritical(msg string) {
	if !loggingEnabled {
		return
	}
	fmt.Printf("\033[1;32m%s\033[0m\n", msg) // bold bright green
}
