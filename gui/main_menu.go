package gui

import (
	"fmt"
	"os"
	"unicode"

	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/gui/fonts"
)

// DrawMainMenu ...
func (context *Context) DrawMainMenu() {
	// Main Menu
	imgui.BeginMainMenuBar()

	if imgui.BeginMenu("File") {
		// TODO: add FA icons
		lbl := ""
		lbl += fmt.Sprintf("%#U", unicode.ToLower(fonts.FA_ICON_FILE_O)) + " "
		lbl += "New"
		if imgui.MenuItem(lbl) {
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
			context.GuiVars.GlobalVars.ShowCube = !context.GuiVars.GlobalVars.ShowCube
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
