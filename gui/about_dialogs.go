package gui

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/settings"
)

// ShowAboutImGui ...
func (context *Context) ShowAboutImGui(open *bool) {
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

// ShowAboutKuplung ...
func (context *Context) ShowAboutKuplung(open *bool) {
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

// ShowMetrics ...
func (context *Context) ShowMetrics(open *bool) {
	if imgui.BeginV("Scene stats", open, imgui.WindowFlagsAlwaysAutoResize|imgui.WindowFlagsNoTitleBar|imgui.WindowFlagsNoResize|imgui.WindowFlagsNoSavedSettings) {
		gl := context.window.OpenGL()
		imgui.Text("OpenGL version: 4.1 (" + gl.GetOpenGLVersion() + ")")
		imgui.Text("GLSL version: 4.10 (" + gl.GetShadingLanguageVersion() + ")")
		imgui.Text("Vendor: " + gl.GetVendorName())
		imgui.Text("Renderer: " + gl.GetRendererName())
		imgui.End()
	}
}
