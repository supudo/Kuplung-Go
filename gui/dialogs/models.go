package dialogs

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/rendering"
	"github.com/supudo/Kuplung-Go/types"
)

// ViewModels ...
type ViewModels struct {
	selectedObject int32

	heightTopPanel float32
}

// NewViewModels ...
func NewViewModels() *ViewModels {
	return &ViewModels{
		selectedObject: -1,
		heightTopPanel: 160,
	}
}

// Render ...
func (view *ViewModels) Render(open, isFrame *bool, rm *rendering.RenderManager) {
	imgui.SetNextWindowSizeV(imgui.Vec2{X: 300, Y: 660}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: 10, Y: 28}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	if imgui.BeginV("Models", open, imgui.WindowFlagsResizeFromAnySide) {
		if imgui.BeginTabBarV("cameraTabs", imgui.TabBarFlagsNoCloseWithMiddleMouseButton|imgui.TabBarFlagsNoTooltip) {
			if imgui.BeginTabItem("Create") {
				view.drawShapes()
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Models") {
				if len(rm.MeshModelFaces) == 0 {
					imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
					imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: .4, Y: .2, Z: .2, W: 1})
					imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: .9, Y: .2, Z: .2, W: 1})
					imgui.Text("No models in the current scene.")
					imgui.PopStyleColorV(3)
				} else {
					view.drawModels(isFrame, rm)
				}
				imgui.EndTabItem()
			}
			imgui.EndTabBar()
		}
		imgui.End()
	}
}

func (view *ViewModels) drawModels(isFrame *bool, rm *rendering.RenderManager) {
	imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
	imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: .4, Y: .2, Z: .2, W: 1})
	imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: .9, Y: .2, Z: .2, W: 1})
	if imgui.ButtonV("Reset values to default", imgui.Vec2{X: -1, Y: 0}) {
		for i := 0; i < len(rm.MeshModelFaces); i++ {
			rm.MeshModelFaces[i].InitProperties()
		}
	}
	imgui.PopStyleColorV(3)

	imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{X: 0, Y: 6})
	imgui.PushStyleVarVec2(imgui.StyleVarWindowPadding, imgui.Vec2{X: 20, Y: 0})
	imgui.PushStyleColor(imgui.StyleColorFrameBg, imgui.Vec4{X: 1.0, Y: 0.0, Z: 0.0, W: 1.0})
	imgui.PushItemWidth(imgui.WindowWidth() * .95)
	imgui.BeginChildV("Scene Items", imgui.Vec2{X: 0, Y: view.heightTopPanel}, true, 0)
	var i int32
	for i = 0; i < int32(len(rm.MeshModelFaces)); i++ {
		mmf := rm.MeshModelFaces[i]
		if imgui.SelectableV(mmf.MeshModel.ModelTitle, view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
			view.selectedObject = i
		}
	}
	imgui.EndChild()
	imgui.PopItemWidth()
	imgui.PopStyleColor()
	imgui.PopStyleVarV(2)
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
