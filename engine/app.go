package engine

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/veandco/go-sdl2/sdl"
)

// KuplungRun is the main loop
func KuplungRun(kuplung *Kuplung) {
	runtime.LockOSThread()

	window := NewKuplungWindow()
	defer window.Window.Destroy()

	setupOpenGL(window.Window)

	kuplung.NewKuplung(window)
}

func setupOpenGL(window *sdl.Window) {
	glContext, err := window.GLCreateContext()
	if err != nil {
		settings.LogError("[SetupEnvironment] Failed to create OpenGL context: %v", err)
	}
	defer sdl.GLDeleteContext(glContext)

	err = window.GLMakeCurrent(glContext)
	if err != nil {
		settings.LogError("[SetupEnvironment] Failed to set current OpenGL context: %v", err)
	}

	_ = sdl.GLSetSwapInterval(1)

	err = gl.Init()
	if err != nil {
		settings.LogError("[SetupEnvironment] Failed to initialize OpenGL: %v", err)
	}
}
