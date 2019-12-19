package gui

import (
	"os"

	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// DrawMainMenu ...
func (context *Context) DrawMainMenu() {
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

	imgui.EndMainMenuBar()
}
