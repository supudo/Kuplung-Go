package dialogs

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"

	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/gui/helpers"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/rendering"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// ViewModels ...
type ViewModels struct {
	selectedObject   int32
	selectedTabScene int32

	heightTopPanel float32

	showTextureWindowAmbient      bool
	showTextureWindowDiffuse      bool
	showTextureWindowDissolve     bool
	showTextureWindowBump         bool
	showTextureWindowDisplacement bool
	showTextureWindowSpecular     bool
	showTextureWindowSpecularExp  bool

	showTextureAmbient      bool
	showTextureDiffuse      bool
	showTextureDissolve     bool
	showTextureBump         bool
	showTextureDisplacement bool
	showTextureSpecular     bool
	showTextureSpecularExp  bool

	showUVEditor bool

	textureAmbientWidth       int32
	textureAmbientHeight      int32
	textureDiffuseWidth       int32
	textureDiffuseHeight      int32
	textureDissolveWidth      int32
	textureDissolveHeight     int32
	textureBumpWidth          int32
	textureBumpHeight         int32
	textureDisplacementWidth  int32
	textureDisplacementHeight int32
	textureSpecularWidth      int32
	textureSpecularHeight     int32
	textureSpecularExpWidth   int32
	textureSpecularExpHeight  int32

	vboTextureAmbient      uint32
	vboTextureDiffuse      uint32
	vboTextureDissolve     uint32
	vboTextureBump         uint32
	vboTextureDisplacement uint32
	vboTextureSpecular     uint32
	vboTextureSpecularExp  uint32
}

// NewViewModels ...
func NewViewModels() *ViewModels {
	return &ViewModels{
		selectedObject:   -1,
		selectedTabScene: -1,
		heightTopPanel:   170,
	}
}

// Render ...
func (view *ViewModels) Render(open, isFrame *bool, rm *rendering.RenderManager) {
	sett := settings.GetSettings()

	imgui.SetNextWindowSizeV(imgui.Vec2{X: 300, Y: float32(sett.AppWindow.SDLWindowHeight - 40)}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: 10, Y: 28}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	imgui.PushStyleColor(imgui.StyleColorTab, imgui.Vec4{X: 153 / 255.0, Y: 61 / 255.0, Z: 61 / 255.0, W: 1.0})
	imgui.PushStyleColor(imgui.StyleColorTabHovered, imgui.Vec4{X: 179 / 255.0, Y: 54 / 255.0, Z: 54 / 255.0, W: 1.0})
	imgui.PushStyleColor(imgui.StyleColorTabActive, imgui.Vec4{X: 204 / 255.0, Y: 41 / 255.0, Z: 41 / 255.0, W: 1.0})

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

	imgui.PopStyleColorV(3)
}

func (view *ViewModels) drawModels(isFrame *bool, rm *rendering.RenderManager) {
	rsett := settings.GetRenderingSettings()
	halfGridSize := float32(rsett.Grid.WorldGridSizeSquares / 2)
	gridSize := float32(rsett.Grid.WorldGridSizeSquares)

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
		if imgui.SelectableV(rm.MeshModelFaces[i].MeshModel.ModelTitle, view.selectedObject == i, 0, imgui.Vec2{X: 0, Y: 0}) {
			view.selectedObject = i
			rm.MeshModelFaces[i].IsModelSelected = true
			rm.SceneSelectedModelObject = i
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

	imgui.BeginChildV("Properties Pane", imgui.Vec2{X: 0, Y: 0}, false, 0)
	imgui.PushItemWidth(imgui.WindowWidth() * .75)
	if view.selectedObject > -1 && int32(len(rm.MeshModelFaces)) > view.selectedObject {
		mmf := rm.MeshModelFaces[view.selectedObject]
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
		imgui.Text("OBJ File:")
		imgui.SameLine()
		imgui.Text(mmf.MeshModel.File)
		imgui.Text("ModelFace:")
		imgui.SameLine()
		imgui.Text(mmf.MeshModel.ModelTitle)
		imgui.Text("Material:")
		imgui.SameLine()
		imgui.Text(mmf.MeshModel.MaterialTitle)
		imgui.Text("Vertices:")
		imgui.SameLine()
		imgui.Text(fmt.Sprintf("%d", mmf.MeshModel.CountVertices))
		imgui.Text("Normals:")
		imgui.SameLine()
		imgui.Text(fmt.Sprintf("%d", mmf.MeshModel.CountNormals))
		imgui.Text("Indices:")
		imgui.SameLine()
		imgui.Text(fmt.Sprintf("%d", mmf.MeshModel.CountIndices))

		if len(mmf.MeshModel.ModelMaterial.TextureAmbient.Image) > 0 ||
			len(mmf.MeshModel.ModelMaterial.TextureDiffuse.Image) > 0 ||
			len(mmf.MeshModel.ModelMaterial.TextureBump.Image) > 0 ||
			len(mmf.MeshModel.ModelMaterial.TextureDisplacement.Image) > 0 ||
			len(mmf.MeshModel.ModelMaterial.TextureDissolve.Image) > 0 ||
			len(mmf.MeshModel.ModelMaterial.TextureSpecular.Image) > 0 ||
			len(mmf.MeshModel.ModelMaterial.TextureSpecularExp.Image) > 0 {
			imgui.Separator()
			imgui.Text("Textures")
		}
		imgui.PopStyleColor()

		view.showTextureLine("##001", types.MaterialTextureTypeAmbient, &view.showTextureWindowAmbient, &view.showTextureAmbient, rm)
		view.showTextureLine("##002", types.MaterialTextureTypeDiffuse, &view.showTextureWindowDiffuse, &view.showTextureDiffuse, rm)
		view.showTextureLine("##003", types.MaterialTextureTypeDissolve, &view.showTextureWindowDissolve, &view.showTextureDissolve, rm)
		view.showTextureLine("##004", types.MaterialTextureTypeBump, &view.showTextureWindowBump, &view.showTextureBump, rm)
		view.showTextureLine("##005", types.MaterialTextureTypeDisplacement, &view.showTextureWindowDisplacement, &view.showTextureDisplacement, rm)
		view.showTextureLine("##006", types.MaterialTextureTypeSpecular, &view.showTextureWindowSpecular, &view.showTextureSpecular, rm)
		view.showTextureLine("##007", types.MaterialTextureTypeSpecularExp, &view.showTextureWindowSpecularExp, &view.showTextureSpecularExp, rm)

		if view.showTextureWindowAmbient {
			view.showTextureImage(mmf, types.MaterialTextureTypeAmbient, "Ambient", &view.showTextureWindowAmbient, &view.showTextureAmbient, &view.vboTextureAmbient, &view.textureAmbientWidth, &view.textureAmbientHeight, rm)
		}

		if view.showTextureWindowDiffuse {
			view.showTextureImage(mmf, types.MaterialTextureTypeDiffuse, "Diffuse", &view.showTextureWindowDiffuse, &view.showTextureDiffuse, &view.vboTextureDiffuse, &view.textureDiffuseWidth, &view.textureDiffuseHeight, rm)
		}

		if view.showTextureWindowDissolve {
			view.showTextureImage(mmf, types.MaterialTextureTypeDissolve, "Dissolve", &view.showTextureWindowDissolve, &view.showTextureDissolve, &view.vboTextureDissolve, &view.textureDissolveWidth, &view.textureDissolveHeight, rm)
		}

		if view.showTextureWindowBump {
			view.showTextureImage(mmf, types.MaterialTextureTypeBump, "Bump", &view.showTextureWindowBump, &view.showTextureBump, &view.vboTextureBump, &view.textureBumpWidth, &view.textureBumpHeight, rm)
		}

		if view.showTextureWindowDisplacement {
			view.showTextureImage(mmf, types.MaterialTextureTypeDisplacement, "Height", &view.showTextureWindowDisplacement, &view.showTextureDisplacement, &view.vboTextureDisplacement, &view.textureDisplacementWidth, &view.textureDisplacementHeight, rm)
		}

		if view.showTextureWindowSpecular {
			view.showTextureImage(mmf, types.MaterialTextureTypeSpecular, "Specular", &view.showTextureWindowSpecular, &view.showTextureSpecular, &view.vboTextureSpecular, &view.textureSpecularWidth, &view.textureSpecularHeight, rm)
		}

		if view.showTextureWindowSpecularExp {
			view.showTextureImage(mmf, types.MaterialTextureTypeSpecularExp, "SpecularExp", &view.showTextureWindowSpecularExp, &view.showTextureSpecularExp, &view.vboTextureSpecularExp, &view.textureSpecularExpWidth, &view.textureSpecularExpHeight, rm)
		}

		KDCMReadOnly := false

		imgui.Separator()
		if imgui.BeginTabBarV("Model Properties", imgui.TabBarFlagsNoCloseWithMiddleMouseButton|imgui.TabBarFlagsNoTooltip) {
			if imgui.BeginTabItem("General") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
				imgui.Text("Properties")
				imgui.PopStyleColor()

				// cel shading
				imgui.Checkbox("Cel Shading", &rm.MeshModelFaces[view.selectedObject].CelShading)
				imgui.Checkbox("Wireframe", &rm.MeshModelFaces[view.selectedObject].Wireframe)
				imgui.Checkbox("Edit Mode", &rm.MeshModelFaces[view.selectedObject].EditMode)
				imgui.Checkbox("Shadows", &rm.MeshModelFaces[view.selectedObject].ShowShadows)
				// alpha
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: rm.MeshModelFaces[view.selectedObject].Alpha})
				imgui.Text("Alpha Blending")
				imgui.PopStyleColor()
				helpers.AddControlsFloatSlider("", 1, 0.0, 1.0, &rm.MeshModelFaces[view.selectedObject].Alpha)

				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Scale") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
				imgui.Text("Scale ModelFace")
				imgui.PopStyleColor()
				// TODO: Gizmo
				imgui.Checkbox("Scale all", &rm.MeshModelFaces[view.selectedObject].Scale0)
				if rm.MeshModelFaces[view.selectedObject].Scale0 {
					imgui.Checkbox("", &KDCMReadOnly)
					imgui.SameLine()
					imgui.SliderFloat("##001", &rm.MeshModelFaces[view.selectedObject].ScaleX.Point, 0.05, halfGridSize)
					imgui.SameLine()
					imgui.Text("X")
					imgui.Checkbox("", &KDCMReadOnly)
					imgui.SameLine()
					imgui.SliderFloat("##001", &rm.MeshModelFaces[view.selectedObject].ScaleY.Point, 0.05, halfGridSize)
					imgui.SameLine()
					imgui.Text("Y")
					imgui.Checkbox("", &KDCMReadOnly)
					imgui.SameLine()
					imgui.SliderFloat("##001", &rm.MeshModelFaces[view.selectedObject].ScaleZ.Point, 0.05, halfGridSize)
					imgui.SameLine()
					imgui.Text("Z")
				} else {
					helpers.AddControlsSliderSameLine("X", 1, 0.05, 0.0, halfGridSize, true, &rm.MeshModelFaces[view.selectedObject].ScaleX.Animate, &rm.MeshModelFaces[view.selectedObject].ScaleX.Point, false, isFrame)
					helpers.AddControlsSliderSameLine("Y", 2, 0.05, 0.0, halfGridSize, true, &rm.MeshModelFaces[view.selectedObject].ScaleY.Animate, &rm.MeshModelFaces[view.selectedObject].ScaleY.Point, false, isFrame)
					helpers.AddControlsSliderSameLine("Z", 3, 0.05, 0.0, halfGridSize, true, &rm.MeshModelFaces[view.selectedObject].ScaleZ.Animate, &rm.MeshModelFaces[view.selectedObject].ScaleZ.Point, false, isFrame)
				}
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Rotate") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
				imgui.Text("Rotate model around axis")
				imgui.PopStyleColor()
				// TODO: Gizmo
				helpers.AddControlsSliderSameLine("X", 4, 1.0, -180.0, 180.0, true, &rm.MeshModelFaces[view.selectedObject].RotateX.Animate, &rm.MeshModelFaces[view.selectedObject].RotateX.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 5, 1.0, -180.0, 180.0, true, &rm.MeshModelFaces[view.selectedObject].RotateY.Animate, &rm.MeshModelFaces[view.selectedObject].RotateY.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 6, 1.0, -180.0, 180.0, true, &rm.MeshModelFaces[view.selectedObject].RotateZ.Animate, &rm.MeshModelFaces[view.selectedObject].RotateZ.Point, true, isFrame)
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Translate") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
				imgui.Text("Move model by axis")
				imgui.PopStyleColor()
				// TODO: Gizmo
				helpers.AddControlsSliderSameLine("X", 7, 0.5, (-1 * gridSize), gridSize, true, &rm.MeshModelFaces[view.selectedObject].PositionX.Animate, &rm.MeshModelFaces[view.selectedObject].PositionX.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 8, 0.5, (-1 * gridSize), gridSize, true, &rm.MeshModelFaces[view.selectedObject].PositionY.Animate, &rm.MeshModelFaces[view.selectedObject].PositionY.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 9, 0.5, (-1 * gridSize), gridSize, true, &rm.MeshModelFaces[view.selectedObject].PositionZ.Animate, &rm.MeshModelFaces[view.selectedObject].PositionZ.Point, true, isFrame)
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Displace") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
				imgui.Text("Displace model")
				imgui.PopStyleColor()
				helpers.AddControlsSliderSameLine("X", 10, 0.5, (-1 * gridSize), gridSize, true, &rm.MeshModelFaces[view.selectedObject].DisplaceX.Animate, &rm.MeshModelFaces[view.selectedObject].DisplaceX.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Y", 11, 0.5, (-1 * gridSize), gridSize, true, &rm.MeshModelFaces[view.selectedObject].DisplaceY.Animate, &rm.MeshModelFaces[view.selectedObject].DisplaceY.Point, true, isFrame)
				helpers.AddControlsSliderSameLine("Z", 12, 0.5, (-1 * gridSize), gridSize, true, &rm.MeshModelFaces[view.selectedObject].DisplaceZ.Animate, &rm.MeshModelFaces[view.selectedObject].DisplaceZ.Point, true, isFrame)
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Material") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
				imgui.Text("Material of the model")
				imgui.PopStyleColor()
				if imgui.Button("Material Editor") {
					rm.MeshModelFaces[view.selectedObject].ShowMaterialEditor = true
				}
				imgui.Checkbox("Parallax Mapping", &rm.MeshModelFaces[view.selectedObject].ParallaxMapping)
				imgui.Separator()
				imgui.Checkbox("Use Tessellation", &rm.MeshModelFaces[view.selectedObject].UseTessellation)
				if rm.MeshModelFaces[view.selectedObject].UseTessellation {
					imgui.Checkbox("Culling", &rm.MeshModelFaces[view.selectedObject].UseCullFace)
					helpers.AddControlsIntegerSlider("Subdivision", 24, 0, 100, &rm.MeshModelFaces[view.selectedObject].TessellationSubdivision)
					imgui.Separator()
					if mmf.MeshModel.ModelMaterial.TextureDisplacement.UseTexture {
						helpers.AddControlsSlider("Displacement", 15, 0.05, -2.0, 2.0, true, &rm.MeshModelFaces[view.selectedObject].DisplacementHeightScale.Animate, &rm.MeshModelFaces[view.selectedObject].DisplacementHeightScale.Point, false, isFrame)
						imgui.Separator()
					}
				} else {
					imgui.Separator()
				}
				helpers.AddControlsSlider("Refraction", 13, 0.05, -10.0, 10.0, true, &rm.MeshModelFaces[view.selectedObject].MaterialRefraction.Animate, &rm.MeshModelFaces[view.selectedObject].MaterialRefraction.Point, true, isFrame)
				helpers.AddControlsSlider("Specular Exponent", 14, 10.0, 0.0, 1000.0, true, &rm.MeshModelFaces[view.selectedObject].MaterialSpecularExp.Animate, &rm.MeshModelFaces[view.selectedObject].MaterialSpecularExp.Point, true, isFrame)
				imgui.Separator()
				helpers.AddControlColor3("Ambient Color", &rm.MeshModelFaces[view.selectedObject].MaterialAmbient.Color, &rm.MeshModelFaces[view.selectedObject].MaterialAmbient.ColorPickerOpen)
				helpers.AddControlColor3("Diffuse Color", &rm.MeshModelFaces[view.selectedObject].MaterialDiffuse.Color, &rm.MeshModelFaces[view.selectedObject].MaterialDiffuse.ColorPickerOpen)
				helpers.AddControlColor3("Specular Color", &rm.MeshModelFaces[view.selectedObject].MaterialSpecular.Color, &rm.MeshModelFaces[view.selectedObject].MaterialSpecular.ColorPickerOpen)
				helpers.AddControlColor3("Emission Color", &rm.MeshModelFaces[view.selectedObject].MaterialEmission.Color, &rm.MeshModelFaces[view.selectedObject].MaterialEmission.ColorPickerOpen)
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Effects") {
				if imgui.TreeNodeV("PBR", imgui.TreeNodeFlagsCollapsingHeader) {
					imgui.BeginGroup()
					imgui.Checkbox("Use PBR", &rm.MeshModelFaces[view.selectedObject].RenderingPBR)
					helpers.AddControlsSlider("Metallic", 13, 0.0000001, 0.0, 1.0, false, nil, &rm.MeshModelFaces[view.selectedObject].RenderingPBRMetallic, true, isFrame)
					helpers.AddControlsSlider("Rougness", 14, 0.0000001, 0.0, 1.0, false, nil, &rm.MeshModelFaces[view.selectedObject].RenderingPBRRoughness, true, isFrame)
					helpers.AddControlsSlider("AO", 15, 0.0000001, 0.0, 1.0, false, nil, &rm.MeshModelFaces[view.selectedObject].RenderingPBRAO, true, isFrame)
					imgui.EndGroup()
					imgui.TreePop()
				}
				if imgui.TreeNodeV("Gaussian Blur", imgui.TreeNodeFlagsCollapsingHeader) {
					imgui.BeginGroup()
					gbtitle := "No Blur"
					switch rm.MeshModelFaces[view.selectedObject].EffectGBlurMode {
					case 0:
						gbtitle = "No Blur"
					case 1:
						gbtitle = "Horizontal"
					case 2:
						gbtitle = "Vertical"
					}
					if imgui.BeginCombo("Mode##228", gbtitle) {
						if imgui.SelectableV("No Blur", rm.MeshModelFaces[view.selectedObject].EffectGBlurMode == 0, 0, imgui.Vec2{X: 0, Y: 0}) {
							rm.MeshModelFaces[view.selectedObject].EffectGBlurMode = 0
						}
						if imgui.SelectableV("Horizontal", rm.MeshModelFaces[view.selectedObject].EffectGBlurMode == 1, 0, imgui.Vec2{X: 0, Y: 0}) {
							rm.MeshModelFaces[view.selectedObject].EffectGBlurMode = 1
						}
						if imgui.SelectableV("Vertical", rm.MeshModelFaces[view.selectedObject].EffectGBlurMode == 2, 0, imgui.Vec2{X: 0, Y: 0}) {
							rm.MeshModelFaces[view.selectedObject].EffectGBlurMode = 2
						}
						imgui.EndCombo()
					}
					helpers.AddControlsSlider("Radius", 16, 0.0, 0.0, 1000.0, true, &rm.MeshModelFaces[view.selectedObject].EffectGBlurRadius.Animate, &rm.MeshModelFaces[view.selectedObject].EffectGBlurRadius.Point, true, isFrame)
					helpers.AddControlsSlider("Width", 17, 0.0, 0.0, 1000.0, true, &rm.MeshModelFaces[view.selectedObject].EffectGBlurWidth.Animate, &rm.MeshModelFaces[view.selectedObject].EffectGBlurWidth.Point, true, isFrame)
					imgui.EndGroup()
					imgui.TreePop()
				}
				if imgui.TreeNodeV("Tone Mapping", imgui.TreeNodeFlagsCollapsingHeader) {
					imgui.BeginGroup()
					imgui.Checkbox("ACES Film Rec2020", &rm.MeshModelFaces[view.selectedObject].EffectToneMappingACESFilmRec2020)
					imgui.Checkbox("HDR", &rm.MeshModelFaces[view.selectedObject].EffectHDRTonemapping)
					imgui.EndGroup()
					imgui.TreePop()
				}
				imgui.EndTabItem()
			}
			if imgui.BeginTabItem("Illumination") {
				imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: 0, Z: 0, W: 1})
				imgui.Text("Illumination type")
				imgui.PopStyleColor()

				illumModels := []string{"[0] Color on and Ambient off",
					"[1] Color on and Ambient on",
					"[2] Highlight on",
					"[3] Reflection on and Raytrace on",
					"[4] Transparency: Glass on\n    Reflection: Raytrace on",
					"[5] Reflection: Fresnel on\n    Raytrace on",
					"[6] Transparency: Refraction on\n    Reflection: Fresnel off\n    Raytrace on",
					"[7] Transparency: Refraction on\n    Reflection: Fresnel on\n    Raytrace on",
					"[8] Reflection on\n    Raytrace off",
					"[9] Transparency: Glass on\n    Reflection: Raytrace off",
					"[10] Casts shadows onto invisible surfaces"}

				ititle := ""
				im := uint32(0)
				for im = 0; im < uint32(len(illumModels)); im++ {
					if im == rm.MeshModelFaces[view.selectedObject].MaterialIlluminationModel {
						ititle = illumModels[im]
					}
				}

				if imgui.BeginCombo("Mode##228", ititle) {
					for im = 0; im < uint32(len(illumModels)); im++ {
						if imgui.SelectableV(illumModels[im], rm.MeshModelFaces[view.selectedObject].MaterialIlluminationModel == im, 0, imgui.Vec2{X: 0, Y: 0}) {
							rm.MeshModelFaces[view.selectedObject].MaterialIlluminationModel = im
						}
					}
					imgui.EndCombo()
				}
				imgui.EndTabItem()
			}
			imgui.EndTabBar()
		}
	}
	imgui.PopItemWidth()
	imgui.EndChild()
}

func (view *ViewModels) showTextureLine(chkLabel string, texType uint16, showWindow *bool, loadTexture *bool, rm *rendering.RenderManager) {
	mmf := rm.MeshModelFaces[view.selectedObject]
	var image, title string
	var useTexture *bool
	switch texType {
	case types.MaterialTextureTypeAmbient:
		title = "Ambient"
		useTexture = &mmf.MeshModel.ModelMaterial.TextureAmbient.UseTexture
		image = mmf.MeshModel.ModelMaterial.TextureAmbient.Image
	case types.MaterialTextureTypeDiffuse:
		title = "Diffuse"
		useTexture = &mmf.MeshModel.ModelMaterial.TextureDiffuse.UseTexture
		image = mmf.MeshModel.ModelMaterial.TextureDiffuse.Image
	case types.MaterialTextureTypeDissolve:
		title = "Dissolve"
		useTexture = &mmf.MeshModel.ModelMaterial.TextureDissolve.UseTexture
		image = mmf.MeshModel.ModelMaterial.TextureDissolve.Image
	case types.MaterialTextureTypeBump:
		title = "Normal"
		useTexture = &mmf.MeshModel.ModelMaterial.TextureBump.UseTexture
		image = mmf.MeshModel.ModelMaterial.TextureBump.Image
	case types.MaterialTextureTypeSpecular:
		title = "Specular"
		useTexture = &mmf.MeshModel.ModelMaterial.TextureSpecular.UseTexture
		image = mmf.MeshModel.ModelMaterial.TextureSpecular.Image
	case types.MaterialTextureTypeSpecularExp:
		title = "Specular Exp"
		useTexture = &mmf.MeshModel.ModelMaterial.TextureSpecularExp.UseTexture
		image = mmf.MeshModel.ModelMaterial.TextureSpecularExp.Image
	case types.MaterialTextureTypeDisplacement:
		title = "Displacement"
		useTexture = &mmf.MeshModel.ModelMaterial.TextureDisplacement.UseTexture
		image = mmf.MeshModel.ModelMaterial.TextureDisplacement.Image
	case types.MaterialTextureTypeUndefined:
		title = "Undefined"
		image = ""
	}

	if len(image) > 0 {
		imgui.Checkbox(chkLabel, useTexture)
		if imgui.IsItemHovered() {
			imgui.SetTooltip(fmt.Sprintf("Show/Hide %v texture", title))
		}
		imgui.SameLine()
		if imgui.Button(fmt.Sprintf("X%v", chkLabel)) {
			*loadTexture = false
			switch texType {
			case types.MaterialTextureTypeAmbient:
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureAmbient.UseTexture = false
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureAmbient.Image = ""
			case types.MaterialTextureTypeDiffuse:
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureDiffuse.UseTexture = false
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureDiffuse.Image = ""
			case types.MaterialTextureTypeDissolve:
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureDissolve.UseTexture = false
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureDissolve.Image = ""
			case types.MaterialTextureTypeBump:
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureBump.UseTexture = false
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureBump.Image = ""
			case types.MaterialTextureTypeSpecular:
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureSpecular.UseTexture = false
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureSpecular.Image = ""
			case types.MaterialTextureTypeSpecularExp:
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureSpecularExp.UseTexture = false
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureSpecularExp.Image = ""
			case types.MaterialTextureTypeDisplacement:
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureDisplacement.UseTexture = false
				rm.MeshModelFaces[view.selectedObject].MeshModel.ModelMaterial.TextureDisplacement.Image = ""
			case types.MaterialTextureTypeUndefined:
				break
			}
		}
		imgui.SameLine()
		if imgui.Button(fmt.Sprintf("V%v", chkLabel)) {
			*showWindow = !*showWindow
			*loadTexture = true
		}
		imgui.SameLine()
		if imgui.Button(fmt.Sprintf("UV%v", chkLabel)) {
			// TODO: UV Editor
			// this.componentUVEditor.setModel(rm.MeshModelFaces[view.selectedObject], texType, "", std::bind(&DialogControlsModels::processTexture, this, std::placeholders::_1));
			// this.showUVEditor = true;
		}
		imgui.SameLine()
		imageFile := filepath.Base(image)
		imgui.Text(fmt.Sprintf("%v: %v", title, imageFile))
	} else {
		btnLabel := "Add Texture " + types.GetMaterialTextureName(texType)
		if imgui.Button(btnLabel) {
			// TODO: UV Editor
			// this.showUVEditor = true;
			// this.componentUVEditor.setModel(rm.MeshModelFaces[view.selectedObject], texType, "", std::bind(&DialogControlsModels::processTexture, this, std::placeholders::_1));
		}
	}
}

func (view *ViewModels) showTextureImage(mmf *meshes.ModelFace, textype uint16, title string, showWindow, genTexture *bool, vboBuffer *uint32, width, height *int32, rm *rendering.RenderManager) {
	var wWidth, wHeight, tWidth, tHeight int32

	sett := settings.GetSettings()
	wWidth = int32(sett.AppWindow.SDLWindowWidth)
	wHeight = int32(sett.AppWindow.SDLWindowHeight)

	img := ""
	if textype == types.MaterialTextureTypeAmbient {
		img = mmf.MeshModel.ModelMaterial.TextureAmbient.Image
	} else if textype == types.MaterialTextureTypeDiffuse {
		img = mmf.MeshModel.ModelMaterial.TextureDiffuse.Image
	} else if textype == types.MaterialTextureTypeDissolve {
		img = mmf.MeshModel.ModelMaterial.TextureDissolve.Image
	} else if textype == types.MaterialTextureTypeBump {
		img = mmf.MeshModel.ModelMaterial.TextureBump.Image
	} else if textype == types.MaterialTextureTypeDisplacement {
		img = mmf.MeshModel.ModelMaterial.TextureDisplacement.Image
	} else if textype == types.MaterialTextureTypeSpecular {
		img = mmf.MeshModel.ModelMaterial.TextureSpecular.Image
	} else {
		img = mmf.MeshModel.ModelMaterial.TextureSpecularExp.Image
	}

	if *genTexture {
		*vboBuffer = view.createTextureBuffer(img, width, height, rm)
	}

	tWidth = *width + 20
	if tWidth > wWidth {
		tWidth = wWidth - 20
	}

	tHeight = *height + 20
	if tHeight > wHeight {
		tHeight = wHeight - 40
	}

	imgui.SetNextWindowSizeV(imgui.Vec2{X: float32(tWidth), Y: float32(tHeight)}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: 60, Y: 40}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	title = title + " Texture"
	if imgui.BeginV(title, showWindow, imgui.WindowFlagsResizeFromAnySide) {
		imgui.Text(fmt.Sprintf("Image: %v", img))
		imgui.Text(fmt.Sprintf("Image dimensions: %d x %d", *width, *height))
		imgui.Separator()
		imgui.Image(imgui.TextureID(*vboBuffer), imgui.Vec2{X: float32(*width), Y: float32(*height)})
		imgui.End()
	}
	*genTexture = false
}

func (view *ViewModels) createTextureBuffer(imageFile string, width, height *int32, rm *rendering.RenderManager) uint32 {
	gl := rm.Window.OpenGL()

	imgFile, err := os.Open(imageFile)
	if err != nil {
		settings.LogError("[DialogModels] Texture file (%v) not found: %v", imageFile, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		settings.LogError("[DialogModels] Can't decode texture (%v): %v", imageFile, err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		settings.LogError("[DialogModels] Texture unsupported stride! (%v)", imageFile)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	*width = int32(rgba.Rect.Size().X)
	*height = int32(rgba.Rect.Size().Y)

	vboBuffer := gl.GenTextures(1)[0]
	gl.ActiveTexture(oglconsts.TEXTURE0)
	gl.BindTexture(oglconsts.TEXTURE_2D, vboBuffer)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MIN_FILTER, oglconsts.LINEAR)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MAG_FILTER, oglconsts.LINEAR)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_WRAP_S, oglconsts.REPEAT)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_WRAP_T, oglconsts.REPEAT)
	gl.TexImage2D(
		oglconsts.TEXTURE_2D,
		0,
		oglconsts.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		oglconsts.RGBA,
		oglconsts.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return vboBuffer
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
