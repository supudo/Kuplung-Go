package gui

import (
	"fmt"
	"os"

	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
	"github.com/supudo/Kuplung-Go/utilities"
)

// DrawMainMenu ...
func (context *Context) DrawMainMenu() {
	sett := settings.GetSettings()
	rsett := settings.GetRenderingSettings()

	// Main Menu
	imgui.BeginMainMenuBar()

	if imgui.BeginMenu("File") {
		// TODO: add FA icons
		// lbl := ""
		// lbl += fmt.Sprintf("%#U", unicode.ToLower(fonts.FA_ICON_FILE_O)) + " "
		// lbl += "New"
		if imgui.MenuItem("New") {
			_, _ = trigger.Fire(types.ActionGuiActionFileNew)
		}
		imgui.Separator()
		if imgui.MenuItemV("Quit", "Cmd+Q", false, true) {
			os.Exit(3)
		}
		imgui.EndMenu()
	}

	if imgui.BeginMenu("Scene") {
		if imgui.BeginMenu("Add Light") {
			if imgui.MenuItem("Directional (Sun)") {
				_, _ = trigger.Fire(types.ActionGuiAddLight, types.LightSourceTypeDirectional)
			}
			if imgui.MenuItem("Point (Light bulb)") {
				_, _ = trigger.Fire(types.ActionGuiAddLight, types.LightSourceTypePoint)
			}
			if imgui.MenuItem("Spot (Flashlight)") {
				_, _ = trigger.Fire(types.ActionGuiAddLight, types.LightSourceTypeSpot)
			}
			imgui.EndMenu()
		}
		imgui.Separator()
		if imgui.BeginMenu("Scene Rendering") {
			if imgui.MenuItemV("Solid", "", rsett.General.SelectedViewModelSkin == types.ViewModelSkinSolid, true) {
				rsett.General.SelectedViewModelSkin = types.ViewModelSkinSolid
			}
			if imgui.MenuItemV("Material", "", rsett.General.SelectedViewModelSkin == types.ViewModelSkinMaterial, true) {
				rsett.General.SelectedViewModelSkin = types.ViewModelSkinMaterial
			}
			if imgui.MenuItemV("Texture", "", rsett.General.SelectedViewModelSkin == types.ViewModelSkinTexture, true) {
				rsett.General.SelectedViewModelSkin = types.ViewModelSkinTexture
			}
			if imgui.MenuItemV("Wireframe", "", rsett.General.SelectedViewModelSkin == types.ViewModelSkinWireframe, true) {
				rsett.General.SelectedViewModelSkin = types.ViewModelSkinWireframe
			}
			if imgui.MenuItemV("Rendered", "", rsett.General.SelectedViewModelSkin == types.ViewModelSkinRendered, true) {
				rsett.General.SelectedViewModelSkin = types.ViewModelSkinRendered
			}
			imgui.Separator()
			imgui.MenuItemV("Render - Depth", "", rsett.General.RenderingDepth, true)
			imgui.EndMenu()
		}
		imgui.Separator()
		imgui.MenuItemV("Render Image", "", context.GuiVars.showImageSave, true)
		imgui.MenuItemV("Renderer UI", "", context.GuiVars.showRendererUI, true)

		imgui.EndMenu()
	}

	if imgui.BeginMenu("View") {
		if imgui.MenuItem("Models") {
			context.GuiVars.showModels = !context.GuiVars.showModels
		}
		if imgui.MenuItem("Controls") {
			context.GuiVars.showControls = !context.GuiVars.showControls
		}
		imgui.Separator()
		if imgui.MenuItem("Log") {
			context.GuiVars.showLog = !context.GuiVars.showLog
		}
		imgui.Separator()
		if imgui.MenuItem("Cube") {
			rsett.General.ShowCube = !rsett.General.ShowCube
		}
		imgui.EndMenu()
	}

	if imgui.BeginMenu("Help") {
		if imgui.MenuItem("Metrics") {
			context.GuiVars.showMetrics = !context.GuiVars.showMetrics
		}
		if imgui.MenuItem("About ImGui") {
			context.GuiVars.showAboutImGui = !context.GuiVars.showAboutImGui
		}
		if imgui.MenuItem("About Kuplung") {
			context.GuiVars.showAboutKuplung = !context.GuiVars.showAboutKuplung
		}
		imgui.Separator()
		if imgui.MenuItem("ImGui Demo Window") {
			context.GuiVars.showDemoWindow = !context.GuiVars.showDemoWindow
		}
		imgui.EndMenu()
	}

	imgui.Text(fmt.Sprintf("  | [%.4f ms/frame] %d objs, %d verts, %d indices (%d tris, %d faces) | %v", sett.MemSettings.NbResult, sett.MemSettings.TotalObjects, sett.MemSettings.TotalVertices, sett.MemSettings.TotalIndices, sett.MemSettings.TotalTriangles, sett.MemSettings.TotalFaces, utilities.GetUsage()))

	imgui.EndMainMenuBar()
}
