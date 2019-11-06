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
	err := gl.Init()
	if err != nil {
		settings.LogError("[setupOpenGL] Failed to initialize OpenGL: %v", err)
	}

	glContext, err := window.GLCreateContext()
	if err != nil {
		settings.LogError("[setupOpenGL] Failed to create OpenGL context: %v", err)
	}
	defer sdl.GLDeleteContext(glContext)

	err = window.GLMakeCurrent(glContext)
	if err != nil {
		settings.LogError("[setupOpenGL] Failed to set current OpenGL context: %v", err)
	}

	err = sdl.GLSetSwapInterval(1)
	if err != nil {
		settings.LogError("[setupOpenGL] Failed to set swap interval: %v", err)
	}
}
