package dialogs

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/settings"
)

// ViewControls ...
type ViewControls struct {
}

// NewViewControls ...
func NewViewControls() *ViewControls {
	return &ViewControls{}
}

// Render ...
func (view *ViewControls) Render(open *bool) {
	sett := settings.GetSettings()

	imgui.SetNextWindowSizeV(imgui.Vec2{X: 300, Y: 600}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(sett.AppWindow.SDLWindowWidth - 310), Y: 28}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	if imgui.BeginV("Controls", open, imgui.WindowFlagsResizeFromAnySide) {
		imgui.Text("Settings")
		imgui.End()
	}
}
