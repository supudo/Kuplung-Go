package dialogs

import (
	"github.com/inkyblackness/imgui-go"
)

// ViewModels ...
type ViewModels struct {
}

// NewViewModels ...
func NewViewModels() *ViewModels {
	return &ViewModels{}
}

// Render ...
func (view *ViewModels) Render(open, isFrame *bool) {
	imgui.SetNextWindowSizeV(imgui.Vec2{X: 300, Y: 600}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: 10, Y: 28}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	if imgui.BeginV("Models", open, imgui.WindowFlagsResizeFromAnySide) {
		imgui.Text("Available Scene Models")
		imgui.End()
	}
}
