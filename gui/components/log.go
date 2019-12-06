package components

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/settings"
)

// ComponentLog ...
type ComponentLog struct {
}

// NewComponentLog ...
func NewComponentLog() *ComponentLog {
	return &ComponentLog{}
}

// Render ...
func (view *ComponentLog) Render(open *bool) {
	sett := settings.GetSettings()

	imgui.SetNextWindowSizeV(imgui.Vec2{X: sett.AppWindow.LogWidth, Y: sett.AppWindow.LogHeight}, imgui.ConditionFirstUseEver)
	x := float32(sett.AppWindow.SDLWindowWidth/2) - float32(sett.AppWindow.LogWidth/2)
	y := float32(sett.AppWindow.SDLWindowHeight - sett.AppWindow.LogHeight - 10)
	imgui.SetNextWindowPosV(imgui.Vec2{X: x, Y: y}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	if imgui.BeginV("Log", open, imgui.WindowFlagsResizeFromAnySide) {
		if imgui.Button("Clear") {
			view.clear()
		}
		imgui.SameLine()
		doCopy := imgui.Button("Copy")
		imgui.Separator()
		imgui.BeginChild("scrolling")
		imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{X: 0, Y: 1})
		if doCopy {
			// TODO: copy to clipboard
		}
		imgui.Text(sett.MemSettings.LogBuffer)
		imgui.SetScrollHereY(1.0)
		imgui.PopStyleVar()
		imgui.EndChild()
		imgui.End()
	}
}

func (view *ComponentLog) clear() {
	sett := settings.GetSettings()
	sett.MemSettings.LogBuffer = ""
}
