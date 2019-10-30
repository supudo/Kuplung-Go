package settings

import (
	"fmt"
	"log"
	"os"
)

// LogError logs an error and exits the application
func LogError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Fatalf(msg)
	os.Exit(3)
}
