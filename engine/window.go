package engine

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/veandco/go-sdl2/sdl"
)

// KuplungWindow ...
type KuplungWindow struct {
	*sdl.Window
}

// NewKuplungWindow ...
func NewKuplungWindow() *KuplungWindow {
	window := &KuplungWindow{}
	window.Window = initSDL()
	return window
}

func initSDL() *sdl.Window {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		settings.LogError("[initSDL] Failed to create window: %v", err)
	}

	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	_ = sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	_ = sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)
	_ = sdl.GLSetAttribute(sdl.GL_STENCIL_SIZE, 8)
	_ = sdl.GLSetAttribute(sdl.GL_ACCELERATED_VISUAL, 1)
	_ = sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
	_ = sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)

	sett := settings.GetSettings()
	window, err := sdl.CreateWindow("Kuplung "+sett.App.ApplicationVersion, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight), sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN)
	if err != nil {
		settings.LogError("[initSDL] Failed to create window: %v", err)
	}

	setupOpenGL(window)

	return window
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
