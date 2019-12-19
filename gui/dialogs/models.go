package dialogs

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/types"
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
	imgui.SetNextWindowSizeV(imgui.Vec2{X: 300, Y: 660}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: 10, Y: 28}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	if imgui.BeginV("Models", open, imgui.WindowFlagsResizeFromAnySide) {
		if imgui.BeginTabBarV("cameraTabs", imgui.TabBarFlagsNoCloseWithMiddleMouseButton|imgui.TabBarFlagsNoTooltip) {
			if imgui.BeginTabItem("Create") {
				view.drawShapes()
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Models") {
				imgui.EndTabItem()
			}
			imgui.EndTabBar()
		}
		imgui.End()
	}
}

func (view *ViewModels) drawShapes() {
	if imgui.ButtonV("Triangle", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeTriangle)
	}
	if imgui.ButtonV("Cone", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeCone)
	}
	if imgui.ButtonV("Cube", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeCube)
	}
	if imgui.ButtonV("Cylinder", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeCylinder)
	}
	if imgui.ButtonV("Grid", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeGrid)
	}
	if imgui.ButtonV("Ico Sphere", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeIcoSphere)
	}
	if imgui.ButtonV("Plane", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypePlane)
	}
	if imgui.ButtonV("Torus", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeTorus)
	}
	if imgui.ButtonV("Tube", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeTube)
	}
	if imgui.ButtonV("UV Sphere", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeUVSphere)
	}
	if imgui.ButtonV("Monkey Head", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeMonkeyHead)
	}

	imgui.Separator()
	imgui.Separator()

	if imgui.ButtonV("Epcot", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeEpcot)
	}
	if imgui.ButtonV("Brick Wall", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeBrickWall)
	}
	if imgui.ButtonV("Plane Objects", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypePlaneObjects)
	}
	if imgui.ButtonV("Plane Objects - Large Plane", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypePlaneObjectsLargePlane)
	}
	if imgui.ButtonV("Material Ball", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeMaterialBall)
	}
	if imgui.ButtonV("Material Ball - Blender", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddShape, types.ShapeTypeMaterialBallBlender)
	}

	imgui.Separator()
	imgui.Separator()

	if imgui.ButtonV("Directional (Sun)", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddLight, types.LightSourceTypeDirectional)
	}
	if imgui.ButtonV("Point (Light bulb)", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddLight, types.LightSourceTypePoint)
	}
	if imgui.ButtonV("Spot (Flashlight)", imgui.Vec2{X: imgui.WindowWidth(), Y: 0}) {
		_, _ = trigger.Fire(types.ActionGuiAddLight, types.LightSourceTypeSpot)
	}
}
