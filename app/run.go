package app

import (
	"runtime"
	"time"

	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
)

// Run ...
func Run(initializer func(interfaces.Window), title string, deferrer <-chan func()) (err error) {
	runtime.LockOSThread()

	trigger.On("log", addToLog)

	var window *engine.KuplungWindow
	window = engine.NewKuplungWindow(title)
	defer window.Close()

	initializer(window)
	sett := settings.GetSettings()

	stopLoop := false
	for !window.ShouldClose() && !stopLoop {
		select {
		case task, ok := <-deferrer:
			if ok {
				task()
			} else {
				stopLoop = true
			}
		case <-time.After(time.Millisecond):
		}
		if !sett.MemSettings.QuitApplication {
			window.Update()
		}
	}

	return
}

func addToLog(msg string) {
	sett := settings.GetSettings()
	if len(sett.MemSettings.LogBuffer) > sett.MemSettings.LogBufferLimit {
		sett.MemSettings.LogBuffer = ""
	}
	sett.MemSettings.LogBuffer += "\n" + msg
}
