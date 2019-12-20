package gui

import (
	"fmt"

	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/settings"
)

// ShowAboutImGui ...
func (context *Context) showParsing(open *bool) {
	if *open {
		imgui.OpenPopup("Kuplung Parsing")
	}
	sett := settings.GetSettings()
	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(sett.AppWindow.SDLWindowWidth) / 2, Y: float32(sett.AppWindow.SDLWindowHeight) / 2}, imgui.ConditionAlways, imgui.Vec2{X: 0.5, Y: 0.5})
	imgui.SetNextWindowFocus()
	if imgui.BeginPopupModalV("Kuplung Parsing", open, imgui.WindowFlagsAlwaysAutoResize|imgui.WindowFlagsNoResize|imgui.WindowFlagsNoTitleBar) {
		imgui.PushStyleColor(imgui.StyleColorPlotHistogram, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
		imgui.Text(fmt.Sprintf("Processing ... %.2f%%", context.GuiVars.ParsingPercentage))
		imgui.ProgressBarV(context.GuiVars.ParsingPercentage/100.0, imgui.Vec2{X: 0.0, Y: 0.0}, "")
		imgui.PopStyleColor()
		imgui.EndPopup()
	}
}
