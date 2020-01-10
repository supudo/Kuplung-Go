package gui

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/settings"
)

func (context *Context) dialogSceneStats(open *bool) {
	sett := settings.GetSettings()
	gl := context.window.OpenGL()
	imgui.SetNextWindowPosV(imgui.Vec2{X: 10, Y: float32(sett.AppWindow.SDLWindowHeight) - 60}, imgui.ConditionAlways, imgui.Vec2{X: 0.5, Y: 0.5})
	imgui.SetNextWindowBgAlpha(0.3)
	imgui.BeginV("Scene Stats", open, imgui.WindowFlagsAlwaysAutoResize|imgui.WindowFlagsNoTitleBar|imgui.WindowFlagsNoResize|imgui.WindowFlagsNoSavedSettings)
	imgui.Text("OpenGL version: 4.1 (" + gl.GetOpenGLVersion() + ")")
	imgui.Text("GLSL version: 4.10 (" + gl.GetShadingLanguageVersion() + ")")
	imgui.Text("Vendor: " + gl.GetVendorName())
	imgui.Text("Renderer: " + gl.GetRendererName())
	// imgui.Separator()
	// imgui.Text("Mouse Position: (%.1f, %.1f)", ImGui::GetIO().MousePos.x, ImGui::GetIO().MousePos.y);
	// imgui.Separator()
	// imgui.Text("Application average %.3f ms/frame (%.1f FPS)", 1000.0f / ImGui::GetIO().Framerate, ImGui::GetIO().Framerate);
	// imgui.Text("%d vertices, %d indices (%d triangles)", sett. ImGui::GetIO().MetricsRenderVertices, ImGui::GetIO().MetricsRenderIndices, ImGui::GetIO().MetricsRenderIndices / 3);
}
