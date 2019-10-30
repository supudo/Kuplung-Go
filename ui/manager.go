package gui

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/veandco/go-sdl2/sdl"
)

// InitGUIManagerPlatform will initialize imgui and all other components
func InitGUIManagerPlatform(window *sdl.Window, io imgui.IO) *SDL {
	//platform := NewSDL(io, window)
	platform, _ := NewSDL2(io)
	defer platform.Dispose()
	return platform
}

// InitGUIManagerRenderer will initialize imgui and all other components
func InitGUIManagerRenderer(io imgui.IO) *OpenGL3 {
	renderer := NewOpenGL3(io)
	defer renderer.Dispose()
	return renderer
}

// RenderGUIStart handles GUI rendering
func RenderGUIStart(platform *SDL, renderer *OpenGL3) {
	UIRenderStart(platform, renderer)
}

// RenderGUIEnd handles GUI rendering
func RenderGUIEnd(platform *SDL, renderer *OpenGL3) {
	UIRenderEnd(platform, renderer)
}
