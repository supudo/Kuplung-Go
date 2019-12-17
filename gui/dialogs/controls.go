package dialogs

import (
	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/gui/helpers"
	"github.com/supudo/Kuplung-Go/rendering"
	"github.com/supudo/Kuplung-Go/settings"
)

// ViewControls ...
type ViewControls struct {
	selectedObject      int
	selectedObjectLight int

	heightTopPanel float32

	fovAnimated bool

	tabCamera1, tabCamera2, tabCamera3 bool
}

// NewViewControls ...
func NewViewControls() *ViewControls {
	view := &ViewControls{
		selectedObject:      -1,
		selectedObjectLight: -1,

		heightTopPanel: 160,

		fovAnimated: false,

		tabCamera1: true,
		tabCamera2: false,
		tabCamera3: false,
	}
	trigger.On("selectedObject", view.setSelectedObject)
	trigger.On("selectedObjectLight", view.setSelectedObjectLight)
	return view
}

// Render ...
func (view *ViewControls) Render(open, isFrame *bool, rm *rendering.RenderManager) {
	sett := settings.GetSettings()
	rsett := settings.GetRenderingSettings()

	imgui.SetNextWindowSizeV(imgui.Vec2{X: 300, Y: 600}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: float32(sett.AppWindow.SDLWindowWidth - 310), Y: 28}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	_ = imgui.BeginV("GUI Controls", open, imgui.WindowFlagsResizeFromAnySide)

	imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
	imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: .4, Y: .2, Z: .2, W: 1})
	imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: .9, Y: .2, Z: .2, W: 1})
	if imgui.ButtonV("Reset values to default", imgui.Vec2{X: -1, Y: 0}) {
		settings.ResetRenderSettings()
		// TODO: reset settings in objects
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
			helpers.AddControlsSlider("Field of view", 1, 1.0, -180.0, 180.0, false, &view.fovAnimated, &rsett.General.Fov, true, isFrame)
			imgui.Separator()
			imgui.Text("Ratio")
			if imgui.IsItemHovered() {
				imgui.SetTooltip("W & H")
			}
			imgui.SliderFloat("W##105", &rsett.General.RatioWidth, 0.0, 5.0)
			imgui.SliderFloat("H##106", &rsett.General.RatioHeight, 0.0, 5.0)
			imgui.Separator()
			imgui.Text("Planes")
			if imgui.IsItemHovered() {
				imgui.SetTooltip("Far & Close")
			}
			imgui.SliderFloat("Close##108", &rsett.General.PlaneClose, 0.0, 1000.0)
			imgui.SliderFloat("Far##107", &rsett.General.PlaneFar, 0.0, 1000.0)
			imgui.Separator()
			imgui.Text("Gamma")
			if imgui.IsItemHovered() {
				imgui.SetTooltip("Gamma correction")
			}
			imgui.SliderFloat("##109", &rsett.General.GammaCoeficient, 1.0, 4.0)
			//mgui.PopStyleVar()
			imgui.TreePop()
		}
		if imgui.TreeNodeV("Editor Artefacts", imgui.TreeNodeFlagsCollapsingHeader) {
			imgui.Checkbox("Axis Helpers", &rsett.Axis.ShowAxisHelpers)
			imgui.Checkbox("Z Axis", &rsett.Axis.ShowZAxis)
			imgui.TreePop()
		}
		if imgui.TreeNodeV("Rays", imgui.TreeNodeFlagsCollapsingHeader) {
			if imgui.Checkbox("Show Rays", &rsett.General.ShowPickRays) {
				settings.SaveRenderingSettings()
			}
			if imgui.Checkbox("Single Ray", &rsett.General.ShowPickRaysSingle) {
				settings.SaveRenderingSettings()
			}
			if imgui.TreeNodeV("Add Ray", imgui.TreeNodeFlagsCollapsingHeader|imgui.TreeNodeFlagsDefaultOpen) {
				imgui.Text("Origin")
				if imgui.ButtonV("Set to camera position", imgui.Vec2{X: imgui.WindowWidth() * 0.75, Y: 0}) {
					rsett.General.RayOriginX = rm.Camera.PositionX.Point
					rsett.General.RayOriginY = rm.Camera.PositionY.Point
					rsett.General.RayOriginZ = rm.Camera.PositionZ.Point
				}
				imgui.InputTextV("X##9920", &rsett.General.RayOriginXS, imgui.InputTextFlagsCharsDecimal, nil)
				imgui.InputTextV("Y##9921", &rsett.General.RayOriginYS, imgui.InputTextFlagsCharsDecimal, nil)
				imgui.InputTextV("Z##9922", &rsett.General.RayOriginZS, imgui.InputTextFlagsCharsDecimal, nil)
				imgui.Checkbox("Animate", &rsett.General.RayAnimate)
				if rsett.General.RayAnimate {
					imgui.Text("Direction")
					imgui.SliderFloat("X##9930", &rsett.General.RayDirectionX, -1.0, 1.0)
					imgui.SliderFloat("Y##9931", &rsett.General.RayDirectionY, -1.0, 1.0)
					imgui.SliderFloat("Z##9932", &rsett.General.RayDirectionZ, -1.0, 1.0)
				} else {
					imgui.Text("Direction")
					imgui.InputTextV("X##9930", &rsett.General.RayDirectionXS, imgui.InputTextFlagsCharsDecimal, nil)
					imgui.InputTextV("Y##9931", &rsett.General.RayDirectionYS, imgui.InputTextFlagsCharsDecimal, nil)
					imgui.InputTextV("Z##9932", &rsett.General.RayDirectionZS, imgui.InputTextFlagsCharsDecimal, nil)
					if imgui.ButtonV("Draw", imgui.Vec2{X: imgui.WindowWidth() * 0.75, Y: 0}) {
						rsett.General.RayDraw = true
					}
				}
			}
			imgui.TreePop()
		}
	case 1:
		if imgui.BeginTabBarV("cameraTabs", imgui.TabBarFlagsNoCloseWithMiddleMouseButton|imgui.TabBarFlagsNoTooltip) {
			if imgui.BeginTabItem("Look At") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
				imgui.Text("Look-At Matrix")
				imgui.PopStyleColorV(1)
				imgui.Separator()
				imgui.Text("Eye")
				helpers.AddControlsSliderSameLine("X", 1, 1.0, -rsett.General.PlaneFar, rsett.General.PlaneFar, false, nil, &rm.Camera.EyeSettings.ViewEye[0], true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 2, 1.0, -rsett.General.PlaneFar, rsett.General.PlaneFar, false, nil, &rm.Camera.EyeSettings.ViewEye[1], true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 3, 1.0, -rsett.General.PlaneFar, rsett.General.PlaneFar, false, nil, &rm.Camera.EyeSettings.ViewEye[2], true, isFrame)
				imgui.Separator()
				imgui.Text("Center")
				helpers.AddControlsSliderSameLine("X", 4, 1.0, -10.0, 10.0, false, nil, &rm.Camera.EyeSettings.ViewCenter[0], true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 5, 1.0, -10.0, 10.0, false, nil, &rm.Camera.EyeSettings.ViewCenter[1], true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 6, 1.0, 0.0, 45.0, false, nil, &rm.Camera.EyeSettings.ViewCenter[2], true, isFrame)
				imgui.Separator()
				imgui.Text("Up")
				helpers.AddControlsSliderSameLine("X", 7, 1.0, -1.0, 1.0, false, nil, &rm.Camera.EyeSettings.ViewUp[0], true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 8, 1.0, -1.0, 1.0, false, nil, &rm.Camera.EyeSettings.ViewUp[1], true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 9, 1.0, -1.0, 1.0, false, nil, &rm.Camera.EyeSettings.ViewUp[2], true, isFrame)
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Rotate") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
				imgui.Text("Rotate object around axis")
				imgui.PopStyleColorV(1)
				helpers.AddControlsSliderSameLine("X", 13, 1.0, 0.0, 360.0, true, &rm.Camera.RotateX.Animate, &rm.Camera.RotateX.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 14, 1.0, 0.0, 360.0, true, &rm.Camera.RotateY.Animate, &rm.Camera.RotateY.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 15, 1.0, 0.0, 360.0, true, &rm.Camera.RotateZ.Animate, &rm.Camera.RotateZ.Point, true, isFrame)
				imgui.Separator()
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
				imgui.Text("Rotate object around center")
				imgui.PopStyleColorV(1)
				helpers.AddControlsSliderSameLine("X", 16, 1.0, -180.0, 180.0, true, &rm.Camera.RotateCenterX.Animate, &rm.Camera.RotateCenterX.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 17, 1.0, -180.0, 180.0, true, &rm.Camera.RotateCenterY.Animate, &rm.Camera.RotateCenterY.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 18, 1.0, -180.0, 180.0, true, &rm.Camera.RotateCenterZ.Animate, &rm.Camera.RotateCenterZ.Point, true, isFrame)
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Translate") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
				imgui.Text("Move object by axis")
				imgui.PopStyleColorV(1)
				helpers.AddControlsSliderSameLine("X", 19, 0.05, float32(-2*rsett.Grid.WorldGridSizeSquares), float32(2*rsett.Grid.WorldGridSizeSquares), true, &rm.Camera.PositionX.Animate, &rm.Camera.PositionX.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 20, 0.05, float32(-2*rsett.Grid.WorldGridSizeSquares), float32(2*rsett.Grid.WorldGridSizeSquares), true, &rm.Camera.PositionY.Animate, &rm.Camera.PositionY.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 21, 0.05, float32(-2*rsett.Grid.WorldGridSizeSquares), float32(2*rsett.Grid.WorldGridSizeSquares), true, &rm.Camera.PositionZ.Animate, &rm.Camera.PositionZ.Point, true, isFrame)
				imgui.EndTabItem()
			}
			imgui.EndTabBar()
		}
	case 2:
		if imgui.BeginTabBarV("cameraModelTabs", imgui.TabBarFlagsNoCloseWithMiddleMouseButton|imgui.TabBarFlagsNoTooltip) {
			if imgui.BeginTabItem("General") {
				imgui.Checkbox("Camera", &rm.CameraModel.ShowCameraObject)
				imgui.Checkbox("Wire", &rm.CameraModel.ShowInWire)
				imgui.Separator()
				imgui.Text("Inner Light Direction")
				helpers.AddControlsSliderSameLine("X", 1, 0.001, -1.0, 1.0, true, &rm.CameraModel.InnerLightDirectionX.Animate, &rm.CameraModel.InnerLightDirectionX.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 2, 0.001, -1.0, 1.0, true, &rm.CameraModel.InnerLightDirectionY.Animate, &rm.CameraModel.InnerLightDirectionY.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 3, 0.001, -1.0, 1.0, true, &rm.CameraModel.InnerLightDirectionZ.Animate, &rm.CameraModel.InnerLightDirectionZ.Point, true, isFrame)
				imgui.Separator()
				imgui.Text("ModelFace Color")
				helpers.AddControlsSliderSameLine("X", 13, 0.01, 0.0, 1.0, true, &rm.CameraModel.ColorR.Animate, &rm.CameraModel.ColorR.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 14, 0.01, 0.0, 1.0, true, &rm.CameraModel.ColorG.Animate, &rm.CameraModel.ColorG.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 15, 0.01, 0.0, 1.0, true, &rm.CameraModel.ColorB.Animate, &rm.CameraModel.ColorB.Point, true, isFrame)
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Position") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
				imgui.Text("Move object by axis")
				imgui.PopStyleColorV(1)
				helpers.AddControlsSliderSameLine("X", 4, 0.05, float32(-2*rsett.Grid.WorldGridSizeSquares), float32(2*rsett.Grid.WorldGridSizeSquares), true, &rm.CameraModel.PositionX.Animate, &rm.CameraModel.PositionX.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 5, 0.05, float32(-2*rsett.Grid.WorldGridSizeSquares), float32(2*rsett.Grid.WorldGridSizeSquares), true, &rm.CameraModel.PositionY.Animate, &rm.CameraModel.PositionY.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 6, 0.05, float32(-2*rsett.Grid.WorldGridSizeSquares), float32(2*rsett.Grid.WorldGridSizeSquares), true, &rm.CameraModel.PositionZ.Animate, &rm.CameraModel.PositionZ.Point, true, isFrame)
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Rotate") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
				imgui.Text("Rotate object around axis")
				imgui.PopStyleColorV(1)
				helpers.AddControlsSliderSameLine("X", 7, 1.0, 0.0, 360.0, true, &rm.CameraModel.RotateX.Animate, &rm.CameraModel.RotateX.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 8, 1.0, 0.0, 360.0, true, &rm.CameraModel.RotateY.Animate, &rm.CameraModel.RotateY.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 9, 1.0, 0.0, 360.0, true, &rm.CameraModel.RotateZ.Animate, &rm.CameraModel.RotateZ.Point, true, isFrame)
				imgui.Separator()
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: .6, Y: .2, Z: .2, W: 1})
				imgui.Text("Rotate object around center")
				imgui.PopStyleColorV(1)
				helpers.AddControlsSliderSameLine("X", 10, 1.0, -180.0, 180.0, true, &rm.CameraModel.RotateCenterX.Animate, &rm.CameraModel.RotateCenterX.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 11, 1.0, -180.0, 180.0, true, &rm.CameraModel.RotateCenterY.Animate, &rm.CameraModel.RotateCenterY.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 12, 1.0, -180.0, 180.0, true, &rm.CameraModel.RotateCenterZ.Animate, &rm.CameraModel.RotateCenterZ.Point, true, isFrame)
				imgui.EndTabItem()
			}
			imgui.EndTabBar()
		}
	case 3:
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
		imgui.Text("General World Grid settings")
		imgui.PopStyleColor()
		imgui.Text("Grid Size")
		if imgui.IsItemHovered() {
			imgui.SetTooltip("In squares")
		}
		imgui.SliderInt("##109", &rsett.Grid.WorldGridSizeSquares, 0, 100)
		imgui.Separator()
		imgui.Checkbox("Grid fixed with World", &rsett.Grid.WorldGridFixedWithWorld)
		imgui.Checkbox("Show Grid", &rsett.Grid.ShowGrid)
		imgui.Checkbox("Act as mirror", &rsett.Grid.ActAsMirror)
	case 5:
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
		imgui.Text("Skybox")
		imgui.PopStyleColor()
		if imgui.BeginCombo("##skybox", rm.SkyBox.SkyboxItems[rsett.SkyBox.SkyboxSelectedItem].Title) {
			var i int32
			for i = 0; i < int32(len(rm.SkyBox.SkyboxItems)); i++ {
				sksel := (i == rsett.SkyBox.SkyboxSelectedItem)
				if imgui.SelectableV(rm.SkyBox.SkyboxItems[i].Title, sksel, 0, imgui.Vec2{0, 0}) {
					rsett.SkyBox.SkyboxSelectedItem = i
				}
				if sksel {
					imgui.SetItemDefaultFocus()
				}
			}
			imgui.EndCombo()
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
