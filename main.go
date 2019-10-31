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
	platform, renderer := InitKuplung()

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
	}
}

// InitKuplung will initialize SDL and all components
func InitKuplung() (p *gui.SDL, r *gui.OpenGL3) {
	sett := settings.GetSettings()

	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	runtime.LockOSThread()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		settings.LogError("[main] Failed to initialize SDL2: %v", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Kuplung "+sett.App.ApplicationVersion, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight), sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN)
	if err != nil {
		settings.LogError("[main] Failed to create window: %v", err)
	}
	defer window.Destroy()

	var platform = gui.InitGUIManagerPlatform(window, io)
	gui.InitGUIManager(io, platform)

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

	_ = sdl.SetHint(sdl.HINT_MAC_CTRL_CLICK_EMULATE_RIGHT_CLICK, "1")
	_ = sdl.SetHint(sdl.HINT_VIDEO_HIGHDPI_DISABLED, "0")

	glContext, err := window.GLCreateContext()
	if err != nil {
		settings.LogError("[main] Failed to create OpenGL context: %v", err)
	}

	err = window.GLMakeCurrent(glContext)
	if err != nil {
		settings.LogError("[main] Failed to set current OpenGL context: %v", err)
	}

	_ = sdl.GLSetSwapInterval(1)

	err = gl.Init()
	if err != nil {
		settings.LogError("[main] Failed to initialize OpenGL: %v", err)
	}

	var renderer = gui.InitGUIManagerRenderer(io)

	return platform, renderer
}
