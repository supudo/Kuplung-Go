package settings

import (
	"fmt"
	"log"
	"os"
)

// LogError logs an error and exits the application
func LogError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Fatalf("[Kuplung] " + msg)
	os.Exit(3)
}

// LogWarn logs a warning
func LogWarn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[Kuplung] " + msg)
}
