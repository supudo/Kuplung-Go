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

	window.setupOpenGL()

	kuplung.NewKuplung(window)
}

func (window *KuplungWindow) setupOpenGL() {
	err := gl.Init()
	if err != nil {
		settings.LogError("[setupOpenGL] Failed to initialize OpenGL: %v", err)
	}

	glContext, err := window.Window.GLCreateContext()
	if err != nil {
		settings.LogError("[setupOpenGL] Failed to create OpenGL context: %v", err)
	}
	defer sdl.GLDeleteContext(glContext)

	err = window.Window.GLMakeCurrent(glContext)
	if err != nil {
		settings.LogError("[setupOpenGL] Failed to set current OpenGL context: %v", err)
	}

	err = sdl.GLSetSwapInterval(1)
	if err != nil {
		settings.LogError("[setupOpenGL] Failed to set swap interval: %v", err)
	}
}
