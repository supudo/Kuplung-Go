package main

import (
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/settings"
	gui "github.com/supudo/Kuplung-Go/ui"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	window, err := InitSDL()
	if err != nil {
		settings.LogError("[main] InitSDL error: %v", err)
	}
	defer window.Destroy()
	defer sdl.Quit()

	SetSDLFlags()

	errOpenGL := SetupOpenGL(window)
	if errOpenGL != nil {
		settings.LogError("[main] SetupOpenGL : %v", errOpenGL)
	}

	platform, renderer := SetupEnvironment(window)

	// surface, err := window.GetSurface()
	// if err != nil {
	// 	settings.LogError("[main] window.GetSurface error: %v", err)
	// }
	// surface.FillRect(nil, 0)

	// rect := sdl.Rect{0, 0, 200, 200}
	// surface.FillRect(&rect, 0xffff0000)
	// window.UpdateSurface()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		gui.UIRenderStart(platform, renderer)
		gui.UIRenderEnd(platform, renderer)
	}
}

// InitSDL will initialize SDL window
func InitSDL() (*sdl.Window, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}
	sett := settings.GetSettings()
	window, err := sdl.CreateWindow("Kuplung "+sett.App.ApplicationVersion, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight), sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN)
	if err != nil {
		settings.LogError("[main] Failed to create window: %v", err)
	}
	return window, nil
}

// SetSDLFlags will setup all SDL flags
func SetSDLFlags() {
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
}

// SetupOpenGL will setup OpenGL context
func SetupOpenGL(window *sdl.Window) error {
	glContext, err := window.GLCreateContext()
	if err != nil {
		settings.LogWarn("[SetupEnvironment] Failed to create OpenGL context: %v", err)
		return err
	}
	defer sdl.GLDeleteContext(glContext)

	err = window.GLMakeCurrent(glContext)
	if err != nil {
		settings.LogWarn("[SetupEnvironment] Failed to set current OpenGL context: %v", err)
		return err
	}

	_ = sdl.GLSetSwapInterval(1)

	err = gl.Init()
	if err != nil {
		settings.LogWarn("[SetupEnvironment] Failed to initialize OpenGL: %v", err)
		return err
	}

	return nil
}

// SetupEnvironment will setup the platform and the renderer
func SetupEnvironment(window *sdl.Window) (platform *gui.SDL, renderer *gui.OpenGL3) {
	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	runtime.LockOSThread()

	platform = gui.InitGUIManagerPlatform(window, io)
	gui.InitGUIManager(io, platform)

	renderer = gui.InitGUIManagerRenderer(io)

	return platform, renderer
}
