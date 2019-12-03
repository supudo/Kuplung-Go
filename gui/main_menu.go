package gui

import (
	"os"

	"github.com/inkyblackness/imgui-go"
)

// DrawMainMenu ...
func (context *Context) DrawMainMenu() {
	// Main Menu
	imgui.BeginMainMenuBar()

	if imgui.BeginMenu("File") {
		if imgui.MenuItem("New") {
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
			context.guiVars.showModels = !context.guiVars.showModels
		}
		if imgui.MenuItem("Controls") {
			context.guiVars.showControls = !context.guiVars.showControls
		}
		imgui.EndMenu()
	}

	if imgui.BeginMenu("Help") {
		if imgui.MenuItem("Metrics") {
			context.guiVars.showMetrics = !context.guiVars.showMetrics
		}
		if imgui.MenuItem("About ImGui") {
			context.guiVars.showAboutImGui = !context.guiVars.showAboutImGui
		}
		if imgui.MenuItem("About Kuplung") {
			context.guiVars.showAboutKuplung = !context.guiVars.showAboutKuplung
		}
		imgui.Separator()
		if imgui.MenuItem("ImGui Demo Window") {
			context.guiVars.showDemoWindow = !context.guiVars.showDemoWindow
		}
		imgui.EndMenu()
	}

	imgui.EndMainMenuBar()
}
