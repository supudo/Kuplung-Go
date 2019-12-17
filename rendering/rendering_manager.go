package rendering

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine/parsers"
	"github.com/supudo/Kuplung-Go/interfaces"
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
}

// NewRenderManager will return an instance of the rendering manager
func NewRenderManager(window interfaces.Window, doProgress func(float32)) *RenderManager {
	rsett := settings.GetRenderingSettings()
	ahPosition := float32(rsett.Grid.WorldGridSizeSquares)
	rm := &RenderManager{}
	rm.window = window
	rm.gridSize = rsett.Grid.WorldGridSizeSquares
	rm.doProgress = doProgress
	rm.initParserManager()
	rm.initSystemModels()
	rm.initCamera()
	rm.initCube()
	rm.initWorldGrid()
	rm.initAxisLabels(ahPosition)
	rm.initCameraModel()
	rm.initMiniAxis()
	rm.initSkyBox()
	return rm
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
}

func (rm *RenderManager) initParserManager() {
	rm.fileParser = parsers.NewParserManager(rm.doProgress)
}

func (rm *RenderManager) initSystemModels() {
	sett := settings.GetSettings()
	rm.systemModels = make(map[string]types.MeshModel)
	rm.systemModels["axis_x_plus"] = rm.fileParser.Parse(sett.App.CurrentPath+"/../Resources/resources/axis_helpers/x_plus.obj", nil)[0]
	rm.systemModels["axis_x_minus"] = rm.fileParser.Parse(sett.App.CurrentPath+"/../Resources/resources/axis_helpers/x_minus.obj", nil)[0]
	rm.systemModels["axis_y_plus"] = rm.fileParser.Parse(sett.App.CurrentPath+"/../Resources/resources/axis_helpers/y_plus.obj", nil)[0]
	rm.systemModels["axis_y_minus"] = rm.fileParser.Parse(sett.App.CurrentPath+"/../Resources/resources/axis_helpers/y_minus.obj", nil)[0]
	rm.systemModels["axis_z_plus"] = rm.fileParser.Parse(sett.App.CurrentPath+"/../Resources/resources/axis_helpers/z_plus.obj", nil)[0]
	rm.systemModels["axis_z_minus"] = rm.fileParser.Parse(sett.App.CurrentPath+"/../Resources/resources/axis_helpers/z_minus.obj", nil)[0]
	rm.systemModels["camera"] = rm.fileParser.Parse(sett.App.CurrentPath+"/../Resources/resources/gui/camera.obj", nil)[0]
	rm.systemModels["light_directional"] = rm.fileParser.Parse(sett.App.CurrentPath+"/../Resources/resources/gui/light_directional.obj", nil)[0]
	rm.systemModels["light_point"] = rm.fileParser.Parse(sett.App.CurrentPath+"/../Resources/resources/gui/light_point.obj", nil)[0]
	rm.systemModels["light_spot"] = rm.fileParser.Parse(sett.App.CurrentPath+"/../Resources/resources/gui/light_spot.obj", nil)[0]
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
