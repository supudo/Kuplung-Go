package settings

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// LogError logs an error and exits the application
func LogError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	_, fn, line, _ := runtime.Caller(1)
	log.Fatalf("[Kuplung] %s:%d - %v", fn, line, msg)
	os.Exit(3)
}

// LogWarn logs a warning
func LogWarn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[Kuplung] " + msg)
}