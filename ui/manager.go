package gui

import (
	"os"
	"time"

	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/veandco/go-sdl2/sdl"
)

// Platform covers mouse/keyboard/gamepad inputs, cursor shape, timing, windowing.
type Platform interface {
	// ShouldStop is regularly called as the abort condition for the program loop.
	ShouldStop() bool
	// ProcessEvents is called once per render loop to dispatch any pending events.
	ProcessEvents()
	// DisplaySize returns the dimension of the display.
	DisplaySize() [2]float32
	// FramebufferSize returns the dimension of the framebuffer.
	FramebufferSize() [2]float32
	// NewFrame marks the begin of a render pass. It must update the imgui IO state according to user input (mouse, keyboard, ...)
	NewFrame()
	// PostRender marks the completion of one render pass. Typically this causes the display buffer to be swapped.
	PostRender()
	// ClipboardText returns the current text of the clipboard, if available.
	ClipboardText() (string, error)
	// SetClipboardText sets the text as the current text of the clipboard.
	SetClipboardText(text string)
}

type clipboard struct {
	platform Platform
}

var guiVars WindowVariables
var clearColor [4]float32

func (board clipboard) Text() (string, error) {
	return board.platform.ClipboardText()
}

func (board clipboard) SetText(text string) {
	board.platform.SetClipboardText(text)
}

// Renderer covers rendering imgui draw data.
type Renderer interface {
	// PreRender causes the display buffer to be prepared for new output.
	PreRender(clearColor [4]float32)
	// Render draws the provided imgui draw data.
	Render(displaySize [2]float32, framebufferSize [2]float32, drawData imgui.DrawData)
}

// WindowVariables holds boolean variables for all the windows
type WindowVariables struct {
	showDemoWindow   bool
	showAboutImGui   bool
	showAboutKuplung bool
	showMetrics      bool
}

// InitGUIManager will initialize GUI variables
func InitGUIManager(io imgui.IO, p *SDL) {
	imgui.CurrentIO().SetClipboard(clipboard{platform: p})

	clearColor = [4]float32{70.0 / 255.0, 70.0 / 255.0, 70.0 / 255.0, 1.0}

	guiVars.showDemoWindow = false
	guiVars.showAboutImGui = false
	guiVars.showAboutKuplung = false
	guiVars.showMetrics = false
}

// InitGUIManagerPlatform will initialize imgui and all other components
func InitGUIManagerPlatform(window *sdl.Window, io imgui.IO) (platform *SDL) {
	platform = NewSDL(io, window)
	defer platform.Dispose()
	return platform
}

// InitGUIManagerRenderer will initialize imgui and all other components
func InitGUIManagerRenderer(io imgui.IO) (renderer *OpenGL3) {
	renderer = NewOpenGL3(io)
	defer renderer.Dispose()
	return renderer
}

// UIRenderStart implements the main program loop of the demo. It returns when the platform signals to stop.
// This demo application shows some basic features of ImGui, as well as exposing the standard demo window.
func UIRenderStart(p Platform, r Renderer) {
	p.ProcessEvents()

	// Signal start of a new frame
	p.NewFrame()
	imgui.NewFrame()

	// Main Menu
	{
		imgui.BeginMainMenuBar()

		if imgui.BeginMenu("File") {
			imgui.Separator()
			if imgui.MenuItemV("Quit", "Cmd+Q", false, true) {
				os.Exit(3)
			}
			imgui.EndMenu()
		}

		if imgui.BeginMenu("Scene") {
			imgui.EndMenu()
		}

		if imgui.BeginMenu("View") {
			imgui.EndMenu()
		}

		if imgui.BeginMenu("Help") {
			if imgui.MenuItem("Metrics") {
				guiVars.showMetrics = true
			}
			if imgui.MenuItem("About ImGui") {
				guiVars.showAboutImGui = true
			}
			if imgui.MenuItem("About Kuplung") {
				guiVars.showAboutKuplung = true
			}
			imgui.Separator()
			if imgui.MenuItem("ImGui Demo Window") {
				guiVars.showDemoWindow = true
			}
			imgui.EndMenu()
		}

		imgui.EndMainMenuBar()
	}

	if guiVars.showAboutImGui {
		ShowAboutImGui(&guiVars.showAboutImGui)
	}

	if guiVars.showAboutKuplung {
		ShowAboutKuplung(&guiVars.showAboutKuplung)
	}

	if guiVars.showDemoWindow {
		imgui.ShowDemoWindow(&guiVars.showDemoWindow)
	}

	if guiVars.showMetrics {
		ShowMetrics(&guiVars.showMetrics)
	}

	// Rendering
	imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

	r.PreRender(clearColor)
	// A this point, the application could perform its own rendering...
	// app.RenderScene()

	r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
	p.PostRender()

	// sleep to avoid 100% CPU usage for this demo
	<-time.After(time.Millisecond * 25)
}

// UIRenderEnd implements the main program loop of the demo. It returns when the platform signals to stop.
// This demo application shows some basic features of ImGui, as well as exposing the standard demo window.
func UIRenderEnd(p Platform, r Renderer) {
	r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
	p.PostRender()

	// sleep to avoid 100% CPU usage for this demo
	<-time.After(time.Millisecond * 25)
}

// ShowAboutImGui show ImGui About screen
func ShowAboutImGui(open *bool) {
	if imgui.BeginV("About ImGui", open, imgui.WindowFlagsAlwaysAutoResize) {
		imgui.Text("ImGui " + imgui.Version())
		imgui.Separator()
		imgui.Text("By Omar Cornut and all github contributors.")
		imgui.Text("ImGui is licensed under the MIT License, see LICENSE for more information.")
		imgui.Separator()
		imgui.Text("Go binding by Inky Blackness")
		imgui.Text("https://github.com/inkyblackness/imgui-go/")
		imgui.End()
	}
}

// ShowAboutKuplung show Kuplung About screen
func ShowAboutKuplung(open *bool) {
	var sett = settings.GetSettings()
	if imgui.BeginV("About Kuplung", open, imgui.WindowFlagsAlwaysAutoResize) {
		imgui.Text("Kuplung " + sett.App.ApplicationVersion)
		imgui.Separator()
		imgui.Text("By supudo.net + github.com/supudo")
		imgui.Text("Whatever license...")
		imgui.Separator()
		imgui.Text("Hold mouse wheel to rotate around")
		imgui.Text("Left Alt + Mouse wheel to increase/decrease the FOV")
		imgui.Text("Left Shift + Mouse wheel to increase/decrease the FOV")
		imgui.Text("By supudo.net + github.com/supudo")
		imgui.End()
	}
}

// ShowMetrics shows application metrics
func ShowMetrics(open *bool) {
	if imgui.BeginV("Scene stats", open, imgui.WindowFlagsAlwaysAutoResize|imgui.WindowFlagsNoTitleBar|imgui.WindowFlagsNoResize|imgui.WindowFlagsNoSavedSettings) {
		// imgui.Text("OpenGL version: 4.1 (" + gl.GoStr(gl.GetString(gl.VERSION)) + ")")
		// imgui.Text("GLSL version: 4.10 (" + gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)) + ")")
		// imgui.Text("Vendor: " + gl.GoStr(gl.GetString(gl.VENDOR)))
		// imgui.Text("Renderer: " + gl.GoStr(gl.GetString(gl.RENDERER)))
		imgui.End()
	}
}

// Run implements the main program loop of the demo. It returns when the platform signals to stop.
// This demo application shows some basic features of ImGui, as well as exposing the standard demo window.
func Run(p Platform, r Renderer) {
	imgui.CurrentIO().SetClipboard(clipboard{platform: p})

	clearColor := [4]float32{70.0 / 255.0, 70.0 / 255.0, 70.0 / 255.0, 1.0}

	var guiVars WindowVariables
	guiVars.showDemoWindow = false
	guiVars.showAboutImGui = false
	guiVars.showAboutKuplung = false
	guiVars.showMetrics = false

	for !p.ShouldStop() {
		p.ProcessEvents()

		// Signal start of a new frame
		p.NewFrame()
		imgui.NewFrame()

		// Main Menu
		{
			imgui.BeginMainMenuBar()

			if imgui.BeginMenu("File") {
				imgui.Separator()
				if imgui.MenuItemV("Quit", "Cmd+Q", false, true) {
					os.Exit(3)
				}
				imgui.EndMenu()
			}

			if imgui.BeginMenu("Scene") {
				imgui.EndMenu()
			}

			if imgui.BeginMenu("View") {
				imgui.EndMenu()
			}

			if imgui.BeginMenu("Help") {
				if imgui.MenuItem("Metrics") {
					guiVars.showMetrics = true
				}
				if imgui.MenuItem("About ImGui") {
					guiVars.showAboutImGui = true
				}
				if imgui.MenuItem("About Kuplung") {
					guiVars.showAboutKuplung = true
				}
				imgui.Separator()
				if imgui.MenuItem("ImGui Demo Window") {
					guiVars.showDemoWindow = true
				}
				imgui.EndMenu()
			}

			imgui.EndMainMenuBar()
		}

		if guiVars.showAboutImGui {
			ShowAboutImGui(&guiVars.showAboutImGui)
		}

		if guiVars.showAboutKuplung {
			ShowAboutKuplung(&guiVars.showAboutKuplung)
		}

		if guiVars.showDemoWindow {
			imgui.ShowDemoWindow(&guiVars.showDemoWindow)
		}

		if guiVars.showMetrics {
			ShowMetrics(&guiVars.showMetrics)
		}

		// Rendering
		imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

		r.PreRender(clearColor)
		// A this point, the application could perform its own rendering...
		// app.RenderScene()

		r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
		p.PostRender()

		// sleep to avoid 100% CPU usage for this demo
		<-time.After(time.Millisecond * 25)
	}
}
