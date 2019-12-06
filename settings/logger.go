package settings

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/sadlil/go-trigger"
)

// LogError logs an error and exits the application
func LogError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	_, fn, line, _ := runtime.Caller(1)
	_, _ = trigger.Fire("log", msg)
	log.Fatalf("[Kuplung] %s:%d - %v", fn, line, msg)
	os.Exit(3)
}

// LogWarn logs a warning
func LogWarn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[Kuplung] " + msg)
	_, _ = trigger.Fire("log", msg)
}

// LogInfo logs a info message
func LogInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println("[Kuplung] [" + time.Now().String() + "] " + msg)
	_, _ = trigger.Fire("log", msg)
}
