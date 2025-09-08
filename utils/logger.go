package utils

import (
	"fmt"
)

func LogInfo(msg string) {
	fmt.Printf("\033[36m%s\033[0m\n", msg) // cyan
}

func LogWarn(msg string) {
	fmt.Printf("\033[33m%s\033[0m\n", msg) // yellow
}

func LogSuccess(msg string) {
	fmt.Printf("\033[32m%s\033[0m\n", msg) // green
}
