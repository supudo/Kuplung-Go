package gui

import (
	"fmt"
	"os"

	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
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
	imgui.Separator()
	imgui.Text(fmt.Sprintf("Application average %.3f ms/frame (%.1f FPS)", 1000.0/imgui.CurrentIO().Framerate(), imgui.CurrentIO().Framerate()))
	imgui.Text(fmt.Sprintf("%d vertices, %d indices (%d triangles)", imgui.CurrentIO().MetricsRenderVertices(), imgui.CurrentIO().MetricsRenderIndices(), imgui.CurrentIO().MetricsRenderIndices()/3))
}

func (context *Context) popupRecentFileImportedDoesntExists(open *bool) {
	if *open {
		imgui.OpenPopup("Warning")
	}
	sett := settings.GetSettings()
	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(sett.AppWindow.SDLWindowWidth)/2 - 200, Y: float32(sett.AppWindow.SDLWindowHeight)/2 - 100}, imgui.ConditionAlways, imgui.Vec2{X: 0.5, Y: 0.5})
	imgui.SetNextWindowFocus()
	if imgui.BeginPopupModalV("Warning", open, imgui.WindowFlagsAlwaysAutoResize|imgui.WindowFlagsNoResize|imgui.WindowFlagsNoTitleBar) {
		imgui.Text("This file no longer exists!")
		if imgui.ButtonV("OK", imgui.Vec2{X: 140, Y: 0}) {
			var recents []*types.FBEntity
			for i := 0; i < len(context.GuiVars.recentFilesImported); i++ {
				if _, err := os.Stat(context.GuiVars.recentFilesImported[i].Path); !os.IsNotExist(err) {
					recents = append(recents, context.GuiVars.recentFilesImported[i])
				}
			}
			context.GuiVars.recentFilesImported = recents
			settings.SaveRecentFilesImported(context.GuiVars.recentFilesImported)
			*open = false
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}
}
