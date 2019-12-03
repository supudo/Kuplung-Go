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
			context.guiVars.showMetrics = true
		}
		if imgui.MenuItem("About ImGui") {
			context.guiVars.showAboutImGui = true
		}
		if imgui.MenuItem("About Kuplung") {
			context.guiVars.showAboutKuplung = true
		}
		imgui.Separator()
		if imgui.MenuItem("ImGui Demo Window") {
			context.guiVars.showDemoWindow = true
		}
		imgui.EndMenu()
	}

	imgui.EndMainMenuBar()

	if context.guiVars.showAboutImGui {
		context.ShowAboutImGui()
	}

	if context.guiVars.showAboutKuplung {
		context.ShowAboutKuplung()
	}

	if context.guiVars.showDemoWindow {
		imgui.ShowDemoWindow(&context.guiVars.showDemoWindow)
	}

	if context.guiVars.showMetrics {
		context.ShowMetrics()
	}
}
