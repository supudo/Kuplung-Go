package app

import (
	"runtime"
	"time"

	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/interfaces"
)

// Run creates a native OpenGL window, initializes it with the given function and
// then runs the event loop until the window shall be closed.
// The provided deferrer is a channel of tasks that can be injected into the event loop.
// When the channel is closed, the loop is stopped and the window is closed.
func Run(initializer func(interfaces.Window), title string, framesPerSecond float64, deferrer <-chan func()) (err error) {
	runtime.LockOSThread()

	var window *engine.KuplungWindow
	window = engine.NewKuplungWindow(title)
	defer window.Close()

	initializer(window)

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
		window.Update()
	}

	return
}
