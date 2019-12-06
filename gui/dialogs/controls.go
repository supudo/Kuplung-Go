package dialogs

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/settings"
)

// ViewControls ...
type ViewControls struct {
	selectedObject      int
	selectedObjectLight int
}

// NewViewControls ...
func NewViewControls() *ViewControls {
	view := &ViewControls{
		selectedObject:      -1,
		selectedObjectLight: -1,
	}
	trigger.On("selectedObject", view.setSelectedObject)
	trigger.On("selectedObjectLight", view.setSelectedObjectLight)
	return view
}

// Render ...
func (view *ViewControls) Render(open *bool) {
	sett := settings.GetSettings()

	imgui.SetNextWindowSizeV(imgui.Vec2{X: 300, Y: 600}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(sett.AppWindow.SDLWindowWidth - 310), Y: 28}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	if imgui.BeginV("Controls", open, imgui.WindowFlagsResizeFromAnySide) {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
		imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: .4, Y: .2, Z: .2, W: 1})
		imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: .9, Y: .2, Z: .2, W: 1})
		if imgui.ButtonV("Reset values to default", imgui.Vec2{X: -1, Y: 0}) {
			// TODO: reset all settings
		}
		imgui.PopStyleColorV(3)

		imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{X: 0, Y: 6})
		imgui.PushStyleColor(imgui.StyleColorFrameBg, imgui.Vec4{X: 1.0, Y: 0.0, Z: 0.0, W: 1.0})
		imgui.PushItemWidth(imgui.WindowWidth() * .95)
		for i := 0; i < 7; i++ {
			switch i {
			case 0:
				if imgui.SelectableV("G1", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
					view.selectedObject = i
					view.selectedObjectLight = -1
					imgui.Text("...")
				}
			case 1:
				if imgui.SelectableV("C", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
					view.selectedObject = i
					view.selectedObjectLight = -1
				}
			case 2:
				if imgui.SelectableV("G2", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
					view.selectedObject = i
					view.selectedObjectLight = -1
				}
			case 3:
				if imgui.SelectableV("SL", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
					view.selectedObject = i
					view.selectedObjectLight = -1
				}
			case 4:
				if imgui.SelectableV("C", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
					view.selectedObject = i
					view.selectedObjectLight = -1
				}
			case 5:
				if imgui.SelectableV("S", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
					view.selectedObject = i
					view.selectedObjectLight = -1
				}
			case 6:
				if imgui.SelectableV("L", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
					view.selectedObject = i
					view.selectedObjectLight = -1
					// TODO: enumerate light sources and add sub-items for each one
				}
			}
		}
		imgui.PopItemWidth()
		imgui.PopStyleColor()
		imgui.PopStyleVar()

		imgui.End()
	}
}

func (view *ViewControls) setSelectedObject(s int) {
	view.selectedObject = s
}

func (view *ViewControls) setSelectedObjectLight(s int) {
	view.selectedObjectLight = s
}
