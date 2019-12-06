package dialogs

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/gui/helpers"
	"github.com/supudo/Kuplung-Go/settings"
)

// ViewControls ...
type ViewControls struct {
	selectedObject      int
	selectedObjectLight int

	heightTopPanel float32

	fov         float32
	fovAnimated bool
}

// NewViewControls ...
func NewViewControls() *ViewControls {
	view := &ViewControls{
		selectedObject:      -1,
		selectedObjectLight: -1,

		heightTopPanel: 160,

		fov:         45.0,
		fovAnimated: false,
	}
	trigger.On("selectedObject", view.setSelectedObject)
	trigger.On("selectedObjectLight", view.setSelectedObjectLight)
	return view
}

// Render ...
func (view *ViewControls) Render(open, isFrame *bool) {
	sett := settings.GetSettings()

	imgui.SetNextWindowSizeV(imgui.Vec2{X: 300, Y: 600}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(sett.AppWindow.SDLWindowWidth - 310), Y: 28}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	_ = imgui.Begin("GUI Controls")

	imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
	imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: .4, Y: .2, Z: .2, W: 1})
	imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: .9, Y: .2, Z: .2, W: 1})
	if imgui.ButtonV("Reset values to default", imgui.Vec2{X: -1, Y: 0}) {
		// TODO: reset all settings
	}
	imgui.PopStyleColorV(3)

	imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{X: 0, Y: 6})
	imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{X: 20, Y: 0})
	imgui.PushStyleColor(imgui.StyleColorFrameBg, imgui.Vec4{X: 1.0, Y: 0.0, Z: 0.0, W: 1.0})
	imgui.PushItemWidth(imgui.WindowWidth() * .95)
	imgui.BeginChildV("Global Items", imgui.Vec2{X: 0, Y: view.heightTopPanel}, true, 0)
	for i := 0; i < 7; i++ {
		switch i {
		case 0:
			if imgui.SelectableV("General", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
				view.selectedObject = i
				view.selectedObjectLight = -1
				imgui.Text("...")
			}
		case 1:
			if imgui.SelectableV("Camera", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
				view.selectedObject = i
				view.selectedObjectLight = -1
			}
		case 2:
			if imgui.SelectableV("Camera Model", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
				view.selectedObject = i
				view.selectedObjectLight = -1
			}
		case 3:
			if imgui.SelectableV("Grid", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
				view.selectedObject = i
				view.selectedObjectLight = -1
			}
		case 4:
			if imgui.SelectableV("Scene Lights", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
				view.selectedObject = i
				view.selectedObjectLight = -1
			}
		case 5:
			if imgui.SelectableV("Skybox", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
				view.selectedObject = i
				view.selectedObjectLight = -1
			}
		case 6:
			if imgui.SelectableV("Lights", view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
				view.selectedObject = i
				view.selectedObjectLight = -1
				// TODO: enumerate light sources and add sub-items for each one
			}
		}
	}
	imgui.EndChild()
	imgui.PopItemWidth()
	imgui.PopStyleColor()
	imgui.PopStyleVarV(2)

	sc := float32(1.0 / 255.0)
	imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 89.0 * sc, Y: 91.0 * sc, Z: 94 * sc, W: 1})
	imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: 119.0 * sc, Y: 122.0 * sc, Z: 124.0 * sc, W: 1})
	imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: .0, Y: .0, Z: .0, W: 1})
	imgui.ButtonV("###splitterGUI", imgui.Vec2{X: -1, Y: 4})
	imgui.PopStyleColorV(3)
	// TODO: get mouse delta up/down
	// if imgui.IsMouseDown(0) {
	// 	view.heightTopPanel += 4
	// }
	if imgui.IsItemHovered() {
		imgui.SetMouseCursor(imgui.MouseCursorResizeNS)
	} else {
		imgui.SetMouseCursor(imgui.MouseCursorNone)
	}

	imgui.BeginChildV("Properties Page", imgui.Vec2{X: 0, Y: 0}, false, 0)
	imgui.PushItemWidth(imgui.WindowWidth() * .75)
	switch view.selectedObject {
	case 0:
		if imgui.TreeNodeV("View Options", imgui.TreeNodeFlagsCollapsingHeader) {
			imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{X: 20, Y: 0})
			helpers.AddSliderF32("Field of view", 1, 1.0, -180, 180, true, true, &view.fovAnimated, isFrame, &view.fov)
			imgui.PopStyleVar()
		}
	}
	imgui.PopItemWidth()
	imgui.EndChild()

	imgui.End()
}

func (view *ViewControls) setSelectedObject(s int) {
	view.selectedObject = s
}

func (view *ViewControls) setSelectedObjectLight(s int) {
	view.selectedObjectLight = s
}
