package rendering

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/engine/parsers"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/objects"
	"github.com/supudo/Kuplung-Go/rendering/renderers"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// RenderManager is the main structure for rendering
type RenderManager struct {
	Window interfaces.Window

	Camera      *objects.Camera
	cube        *objects.Cube
	wgrid       *objects.WorldGrid
	axisLabels  *objects.AxisLabels
	CameraModel *objects.CameraModel
	miniAxis    *objects.MiniAxis
	SkyBox      *objects.SkyBox

	rendererDefered              *renderers.RendererDefered
	rendererForward              *renderers.RendererForward
	rendererForwardShadowMapping *renderers.RendererForwardShadowMapping
	rendererShadowMapping        *renderers.RendererShadowMapping
	rendererSimple               *renderers.RendererSimple

	gridSize int32

	doProgress func(float32)
	fileParser *parsers.ParserManager

	systemModels map[string]types.MeshModel

	MeshModelFaces []*meshes.ModelFace
	LightSources   []*objects.Light

	RenderProps types.RenderProperties

	SceneSelectedModelObject int32
}

// NewRenderManager will return an instance of the rendering manager
func NewRenderManager(window interfaces.Window, doProgress func(float32)) *RenderManager {
	rsett := settings.GetRenderingSettings()
	ahPosition := float32(rsett.Grid.WorldGridSizeSquares)

	rm := &RenderManager{}
	rm.Window = window
	rm.doProgress = doProgress
	rm.SceneSelectedModelObject = -1

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
	rm.initRenderers()

	trigger.On(types.ActionGuiAddShape, rm.addShape)
	trigger.On(types.ActionGuiAddLight, rm.addLight)
	trigger.On(types.ActionGuiActionFileNew, rm.clearScene)
	trigger.On(types.ActionImportFile, rm.fileImport)

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
	sett := settings.GetSettings()
	rsett := settings.GetRenderingSettings()

	w, h := rm.Window.Size()
	rm.Window.OpenGL().Viewport(0, 0, int32(w), int32(h))

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

	if rsett.General.ShowAllVisualArtefacts && rsett.Grid.ShowGrid {
		rm.wgrid.ActAsMirror = rsett.Grid.ActAsMirror
		rm.wgrid.Render()
	}

	if rsett.General.ShowAllVisualArtefacts && rsett.Axis.ShowAxisHelpers {
		rm.axisLabels.Render(ahPosition)
	}

	if rsett.General.ShowAllVisualArtefacts {
		rm.CameraModel.Render(rm.wgrid.MatrixModel)
		rm.miniAxis.Render()
		rm.SkyBox.Render()

		for i := 0; i < len(rm.LightSources); i++ {
			rm.LightSources[i].Render()
		}
	}

	if sett.Components.ShouldRecompileShaders && sett.App.RendererType == types.InAppRendererTypeForward {
		rm.rendererForward.CompileShaders()
	}

	switch sett.App.RendererType {
	case types.InAppRendererTypeDeferred:
		rm.rendererDefered.Render(rm.RenderProps, rm.MeshModelFaces, rm.wgrid.MatrixModel, rm.Camera.CameraPosition)
	case types.InAppRendererTypeForward:
		rm.rendererForward.Render(rm.RenderProps, rm.MeshModelFaces, rm.wgrid.MatrixModel, rm.Camera.CameraPosition, rm.SceneSelectedModelObject, rm.LightSources)
	case types.InAppRendererTypeForwardShadowMapping:
		rm.rendererForwardShadowMapping.Render(rm.RenderProps, rm.MeshModelFaces, rm.wgrid.MatrixModel, rm.Camera.CameraPosition)
	case types.InAppRendererTypeShadowMapping:
		rm.rendererShadowMapping.Render(rm.RenderProps, rm.MeshModelFaces, rm.wgrid.MatrixModel, rm.Camera.CameraPosition)
	case types.InAppRendererTypeSimple:
		rm.rendererSimple.Render(rm.RenderProps, rm.MeshModelFaces, rm.wgrid.MatrixModel, rm.Camera.CameraPosition)
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
	rm.rendererSimple.Dispose()
}

func (rm *RenderManager) initSettings() {
	rsett := settings.GetRenderingSettings()

	rm.gridSize = rsett.Grid.WorldGridSizeSquares

	rm.RenderProps.UIAmbientLightX = 0.2
	rm.RenderProps.UIAmbientLightY = 0.2
	rm.RenderProps.UIAmbientLightZ = 0.2

	rm.RenderProps.SolidLightDirectionX = 0.0
	rm.RenderProps.SolidLightDirectionY = 1.0
	rm.RenderProps.SolidLightDirectionZ = 0.0

	rm.RenderProps.SolidLightMaterialColor = mgl32.Vec3{0.0, 0.7, 0.0}
	rm.RenderProps.SolidLightAmbient = mgl32.Vec3{1.0}
	rm.RenderProps.SolidLightDiffuse = mgl32.Vec3{1.0}
	rm.RenderProps.SolidLightSpecular = mgl32.Vec3{1.0}

	rm.RenderProps.SolidLightAmbientStrength = 0.3
	rm.RenderProps.SolidLightDiffuseStrength = 1.0
	rm.RenderProps.SolidLightSpecularStrength = 0.0

	rm.RenderProps.SolidLightMaterialColorColorPicker = false
	rm.RenderProps.SolidLightAmbientColorPicker = false
	rm.RenderProps.SolidLightDiffuseColorPicker = false
	rm.RenderProps.SolidLightSpecularColorPicker = false

	rm.SceneSelectedModelObject = -1
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
	rm.Camera = objects.InitCamera(rm.Window)
}

func (rm *RenderManager) initCube() {
	rm.cube = objects.CubeInit(rm.Window)
}

func (rm *RenderManager) initWorldGrid() {
	rm.wgrid = objects.InitWorldGrid(rm.Window)
}

func (rm *RenderManager) initAxisLabels(ahPosition float32) {
	rm.axisLabels = objects.InitAxisLabels(rm.Window)
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
	rm.CameraModel = objects.InitCameraModel(rm.Window, rm.systemModels["camera"])
	rm.CameraModel.InitProperties()
	rm.CameraModel.InitBuffers()
}

func (rm *RenderManager) initMiniAxis() {
	rm.miniAxis = objects.InitMiniAxis(rm.Window)
	rm.miniAxis.InitProperties()
	rm.miniAxis.InitBuffers()
}

func (rm *RenderManager) initSkyBox() {
	rm.SkyBox = objects.InitSkyBox(rm.Window)
	rm.SkyBox.InitBuffers()
}

func (rm *RenderManager) addShape(shape types.ShapeType) {
	parsingChan := make(chan []types.MeshModel)
	go rm.addShapeAsync(parsingChan, shape)
	mmodels := <-parsingChan

	for i := 0; i < len(mmodels); i++ {
		mesh := meshes.NewModelFace(rm.Window, mmodels[i])
		mesh.InitProperties()
		mesh.InitBuffers()
		rm.MeshModelFaces = append(rm.MeshModelFaces, mesh)

		sett := settings.GetSettings()
		sett.MemSettings.TotalVertices += mesh.MeshModel.CountVertices
		sett.MemSettings.TotalIndices += mesh.MeshModel.CountIndices
		sett.MemSettings.TotalTriangles += mesh.MeshModel.CountVertices / 3
		sett.MemSettings.TotalFaces += mesh.MeshModel.CountVertices / 6
		sett.MemSettings.TotalObjects++
	}
}

func (rm *RenderManager) addShapeAsync(parsingChannel chan []types.MeshModel, shape types.ShapeType) {
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
	mmodels := rm.fileParser.Parse(sett.App.CurrentPath+"shapes/"+shapeName+".obj", nil)
	_, _ = trigger.Fire(types.ActionParsingHide)
	parsingChannel <- mmodels
}

func (rm *RenderManager) addLight(shape types.LightSourceType) {
	lightObject := objects.InitLight(rm.Window)
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
	sett := settings.GetSettings()
	sett.MemSettings.TotalVertices = 0
	sett.MemSettings.TotalIndices = 0
	sett.MemSettings.TotalTriangles = 0
	sett.MemSettings.TotalFaces = 0
	sett.MemSettings.TotalObjects = 0
	rm.ResetSettings()
	_, _ = trigger.Fire(types.ActionClearGuiControls)
}

func (rm *RenderManager) initRenderers() {
	rm.rendererDefered = renderers.NewRendererDefered(rm.Window)
	rm.rendererForward = renderers.NewRendererForward(rm.Window)
	rm.rendererForwardShadowMapping = renderers.NewRendererForwardShadowMapping(rm.Window)
	rm.rendererShadowMapping = renderers.NewRendererShadowMapping(rm.Window)
	rm.rendererSimple = renderers.NewRendererSimple(rm.Window)
}

func (rm *RenderManager) fileImport(entity *types.FBEntity, setts []string, itype types.ImportExportFormat) {
	parsingChan := make(chan []types.MeshModel)
	go rm.fileImportAsync(parsingChan, entity, setts, itype)
	mmodels := <-parsingChan

	for i := 0; i < len(mmodels); i++ {
		mesh := meshes.NewModelFace(rm.Window, mmodels[i])
		mesh.InitProperties()
		mesh.InitBuffers()
		rm.MeshModelFaces = append(rm.MeshModelFaces, mesh)

		sett := settings.GetSettings()
		sett.MemSettings.TotalVertices += mesh.MeshModel.CountVertices
		sett.MemSettings.TotalIndices += mesh.MeshModel.CountIndices
		sett.MemSettings.TotalTriangles += mesh.MeshModel.CountVertices / 3
		sett.MemSettings.TotalFaces += mesh.MeshModel.CountVertices / 6
		sett.MemSettings.TotalObjects++
	}
}

func (rm *RenderManager) fileImportAsync(parsingChannel chan []types.MeshModel, entity *types.FBEntity, setts []string, itype types.ImportExportFormat) {
	_, _ = trigger.Fire(types.ActionParsingShow)
	mmodels := rm.fileParser.Parse(entity.Path, nil)
	_, _ = trigger.Fire(types.ActionParsingHide)
	parsingChannel <- mmodels
}
