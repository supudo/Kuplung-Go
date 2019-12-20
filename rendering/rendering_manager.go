package rendering

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/engine/parsers"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/objects"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// RenderManager is the main structure for rendering
type RenderManager struct {
	window interfaces.Window

	Camera *objects.Camera

	cube        *objects.Cube
	wgrid       *objects.WorldGrid
	axisLabels  *objects.AxisLabels
	CameraModel *objects.CameraModel
	miniAxis    *objects.MiniAxis
	SkyBox      *objects.SkyBox

	gridSize int32

	doProgress func(float32)
	fileParser *parsers.ParserManager

	systemModels map[string]types.MeshModel

	MeshModelFaces []*meshes.ModelFace
	LightSources   []*objects.Light

	UIAmbientLightX, UIAmbientLightY, UIAmbientLightZ float32

	SolidLightDirectionX, SolidLightDirectionY, SolidLightDirectionZ float32

	SolidLightMaterialColor, SolidLightAmbient mgl32.Vec3
	SolidLightDiffuse, SolidLightSpecular      mgl32.Vec3

	SolidLightAmbientStrength, SolidLightDiffuseStrength, SolidLightSpecularStrength float32

	SolidLightMaterialColorColorPicker, SolidLightAmbientColorPicker bool
	SolidLightDiffuseColorPicker, SolidLightSpecularColorPicker      bool
}

// NewRenderManager will return an instance of the rendering manager
func NewRenderManager(window interfaces.Window, doProgress func(float32)) *RenderManager {
	rsett := settings.GetRenderingSettings()
	ahPosition := float32(rsett.Grid.WorldGridSizeSquares)

	rm := &RenderManager{}
	rm.window = window
	rm.doProgress = doProgress

	rm.initSettings()
	rm.initParserManager()
	rm.initSystemModels()
	rm.initCamera()
	rm.initCube()
	rm.initWorldGrid()
	rm.initAxisLabels(ahPosition)
	rm.initCameraModel()
	rm.initMiniAxis()
	rm.initSkyBox()

	trigger.On(types.ActionGuiAddShape, rm.addShape)
	trigger.On(types.ActionGuiAddLight, rm.addLight)
	trigger.On(types.ActionGuiActionFileNew, rm.clearScene)

	return rm
}

// ResetSettings ..
func (rm *RenderManager) ResetSettings() {
	settings.ResetRenderSettings()
	rm.initSettings()
	rm.Camera.InitProperties()
	rm.CameraModel.InitProperties()
	rm.wgrid.InitProperties()
	for i := 0; i < len(rm.LightSources); i++ {
		rm.LightSources[i].InitProperties(rm.LightSources[i].LightType)
	}
}

// Render handles rendering of all scene objects
func (rm *RenderManager) Render() {
	rsett := settings.GetRenderingSettings()

	w, h := rm.window.Size()
	rm.window.OpenGL().Viewport(0, 0, int32(w), int32(h))

	rsett.MatrixProjection = mgl32.Perspective(mgl32.DegToRad(rsett.General.Fov), rsett.General.RatioWidth/rsett.General.RatioHeight, rsett.General.PlaneClose, rsett.General.PlaneFar)
	rm.Camera.Render()
	rsett.MatrixCamera = rm.Camera.MatrixCamera

	ahPosition := float32(rsett.Grid.WorldGridSizeSquares / 2)

	if rsett.Grid.WorldGridSizeSquares != rm.gridSize {
		rm.gridSize = rsett.Grid.WorldGridSizeSquares
		rm.wgrid.GridSize = rsett.Grid.WorldGridSizeSquares
		rm.wgrid.InitBuffers(rsett.Grid.WorldGridSizeSquares, 1.0)
		rm.axisLabels.InitBuffers()
	}

	if rsett.General.ShowCube {
		rm.cube.Render()
	}

	if rsett.Grid.ShowGrid {
		rm.wgrid.ActAsMirror = rsett.Grid.ActAsMirror
		rm.wgrid.Render()
	}

	if rsett.Axis.ShowAxisHelpers {
		rm.axisLabels.Render(ahPosition)
	}

	rm.CameraModel.Render(rm.wgrid.MatrixModel)
	rm.miniAxis.Render()
	rm.SkyBox.Render()

	for i := 0; i < len(rm.MeshModelFaces); i++ {
		rm.MeshModelFaces[i].Render()
	}

	for i := 0; i < len(rm.LightSources); i++ {
		rm.LightSources[i].Render()
	}
}

// Dispose will cleanup everything
func (rm *RenderManager) Dispose() {
	rm.cube.Dispose()
	rm.wgrid.Dispose()
	rm.Camera.Dispose()
	rm.axisLabels.Dispose()
	rm.CameraModel.Dispose()
	rm.miniAxis.Dispose()
	rm.SkyBox.Dispose()
	for i := 0; i < len(rm.MeshModelFaces); i++ {
		rm.MeshModelFaces[i].Dispose()
	}
	for i := 0; i < len(rm.LightSources); i++ {
		rm.LightSources[i].Dispose()
	}
}

func (rm *RenderManager) initSettings() {
	rsett := settings.GetRenderingSettings()

	rm.gridSize = rsett.Grid.WorldGridSizeSquares

	rm.UIAmbientLightX = 0.2
	rm.UIAmbientLightY = 0.2
	rm.UIAmbientLightZ = 0.2

	rm.SolidLightDirectionX = 0.0
	rm.SolidLightDirectionY = 1.0
	rm.SolidLightDirectionZ = 0.0

	rm.SolidLightMaterialColor = mgl32.Vec3{0.0, 0.7, 0.0}
	rm.SolidLightAmbient = mgl32.Vec3{1.0}
	rm.SolidLightDiffuse = mgl32.Vec3{1.0}
	rm.SolidLightSpecular = mgl32.Vec3{1.0}

	rm.SolidLightAmbientStrength = 0.3
	rm.SolidLightDiffuseStrength = 1.0
	rm.SolidLightSpecularStrength = 0.0

	rm.SolidLightMaterialColorColorPicker = false
	rm.SolidLightAmbientColorPicker = false
	rm.SolidLightDiffuseColorPicker = false
	rm.SolidLightSpecularColorPicker = false
}

func (rm *RenderManager) initParserManager() {
	rm.fileParser = parsers.NewParserManager(rm.doProgress)
}

func (rm *RenderManager) initSystemModels() {
	sett := settings.GetSettings()
	rm.systemModels = make(map[string]types.MeshModel)
	rm.systemModels["axis_x_plus"] = rm.fileParser.Parse(sett.App.CurrentPath+"axis_helpers/x_plus.obj", nil)[0]
	rm.systemModels["axis_x_minus"] = rm.fileParser.Parse(sett.App.CurrentPath+"axis_helpers/x_minus.obj", nil)[0]
	rm.systemModels["axis_y_plus"] = rm.fileParser.Parse(sett.App.CurrentPath+"axis_helpers/y_plus.obj", nil)[0]
	rm.systemModels["axis_y_minus"] = rm.fileParser.Parse(sett.App.CurrentPath+"axis_helpers/y_minus.obj", nil)[0]
	rm.systemModels["axis_z_plus"] = rm.fileParser.Parse(sett.App.CurrentPath+"axis_helpers/z_plus.obj", nil)[0]
	rm.systemModels["axis_z_minus"] = rm.fileParser.Parse(sett.App.CurrentPath+"axis_helpers/z_minus.obj", nil)[0]
	rm.systemModels["camera"] = rm.fileParser.Parse(sett.App.CurrentPath+"gui/camera.obj", nil)[0]
	rm.systemModels["light_directional"] = rm.fileParser.Parse(sett.App.CurrentPath+"gui/light_directional.obj", nil)[0]
	rm.systemModels["light_point"] = rm.fileParser.Parse(sett.App.CurrentPath+"gui/light_point.obj", nil)[0]
	rm.systemModels["light_spot"] = rm.fileParser.Parse(sett.App.CurrentPath+"gui/light_spot.obj", nil)[0]
}

func (rm *RenderManager) initCamera() {
	rm.Camera = objects.InitCamera(rm.window)
}

func (rm *RenderManager) initCube() {
	rm.cube = objects.CubeInit(rm.window)
}

func (rm *RenderManager) initWorldGrid() {
	rm.wgrid = objects.InitWorldGrid(rm.window)
}

func (rm *RenderManager) initAxisLabels(ahPosition float32) {
	rm.axisLabels = objects.InitAxisLabels(rm.window)
	models := []types.MeshModel{
		rm.systemModels["axis_x_plus"],
		rm.systemModels["axis_x_minus"],
		rm.systemModels["axis_y_plus"],
		rm.systemModels["axis_y_minus"],
		rm.systemModels["axis_z_plus"],
		rm.systemModels["axis_z_minus"]}
	rm.axisLabels.SetModels(models, ahPosition)
}

func (rm *RenderManager) initCameraModel() {
	rm.CameraModel = objects.InitCameraModel(rm.window, rm.systemModels["camera"])
	rm.CameraModel.InitProperties()
	rm.CameraModel.InitBuffers()
}

func (rm *RenderManager) initMiniAxis() {
	rm.miniAxis = objects.InitMiniAxis(rm.window)
	rm.miniAxis.InitProperties()
	rm.miniAxis.InitBuffers()
}

func (rm *RenderManager) initSkyBox() {
	rm.SkyBox = objects.InitSkyBox(rm.window)
	rm.SkyBox.InitBuffers()
}

func (rm *RenderManager) addShape(shape types.ShapeType) {
	go rm.addShapeAsync(shape)
}

func (rm *RenderManager) addShapeAsync(shape types.ShapeType) {
	_, _ = trigger.Fire(types.ActionParsingShow)
	shapeName := ""
	switch shape {
	case types.ShapeTypeCone:
		shapeName = "cone"
	case types.ShapeTypeCube:
		shapeName = "cube"
	case types.ShapeTypeCylinder:
		shapeName = "cylinder"
	case types.ShapeTypeGrid:
		shapeName = "grid"
	case types.ShapeTypeIcoSphere:
		shapeName = "ico_sphere"
	case types.ShapeTypeMonkeyHead:
		shapeName = "monkey_head"
	case types.ShapeTypePlane:
		shapeName = "plane"
	case types.ShapeTypeTriangle:
		shapeName = "triangle"
	case types.ShapeTypeTorus:
		shapeName = "torus"
	case types.ShapeTypeTube:
		shapeName = "tube"
	case types.ShapeTypeUVSphere:
		shapeName = "uv_sphere"
	case types.ShapeTypeBrickWall:
		shapeName = "brick_wall"
	case types.ShapeTypePlaneObjects:
		shapeName = "plane_objects"
	case types.ShapeTypePlaneObjectsLargePlane:
		shapeName = "plane_objects_large"
	case types.ShapeTypeMaterialBall:
		shapeName = "MaterialBall"
	case types.ShapeTypeMaterialBallBlender:
		shapeName = "MaterialBallBlender"
	case types.ShapeTypeEpcot:
		shapeName = "epcot"
	}
	sett := settings.GetSettings()
	mmodel := rm.fileParser.Parse(sett.App.CurrentPath+"shapes/"+shapeName+".obj", nil)[0]
	mesh := meshes.NewModelFace(rm.window, mmodel)
	rm.MeshModelFaces = append(rm.MeshModelFaces, mesh)
	_, _ = trigger.Fire(types.ActionParsingHide)
}

func (rm *RenderManager) addLight(shape types.LightSourceType) {
	lightObject := objects.InitLight(rm.window)
	lightObject.InitProperties(shape)
	switch shape {
	case types.LightSourceTypeDirectional:
		lightObject.Title = fmt.Sprintf("Directional %v", len(rm.LightSources)+1)
		lightObject.Description = "Directional area light source"
		lightObject.SetModel(rm.systemModels["light_directional"])
	case types.LightSourceTypePoint:
		lightObject.Title = fmt.Sprintf("Point %v", len(rm.LightSources)+1)
		lightObject.Description = "Omnidirectional point light source"
		lightObject.SetModel(rm.systemModels["light_point"])
	case types.LightSourceTypeSpot:
		lightObject.Title = fmt.Sprintf("Spot %v", len(rm.LightSources)+1)
		lightObject.Description = "Directional cone light source"
		lightObject.SetModel(rm.systemModels["light_spot"])
	}
	lightObject.InitBuffers()
	rm.LightSources = append(rm.LightSources, lightObject)
}

func (rm *RenderManager) clearScene() {
	for i := 0; i < len(rm.MeshModelFaces); i++ {
		rm.MeshModelFaces[i].Dispose()
	}
	for i := 0; i < len(rm.LightSources); i++ {
		rm.LightSources[i].Dispose()
	}
	rm.MeshModelFaces = nil
	rm.LightSources = nil
	rm.ResetSettings()
	_, _ = trigger.Fire(types.ActionClearGuiControls)
}
