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

	camera *objects.Camera

	cube       *objects.Cube
	wgrid      *objects.WorldGrid
	axisLabels *objects.AxisLabels

	gridSize int32

	doProgress func(float32)
	fileParser *parsers.ParserManager

	systemModels map[string]types.MeshModel
}

// NewRenderManager will return an instance of the rendering manager
func NewRenderManager(window interfaces.Window, doProgress func(float32)) *RenderManager {
	rsett := settings.GetRenderingSettings()
	ahPosition := float32(rsett.WorldGridSizeSquares)
	rm := &RenderManager{}
	rm.window = window
	rm.gridSize = rsett.WorldGridSizeSquares
	rm.doProgress = doProgress
	rm.initParserManager()
	rm.initSystemModels()
	rm.initCamera()
	rm.initCube()
	rm.initWorldGrid()
	rm.initAxisLabels(ahPosition)
	return rm
}

// Render handles rendering of all scene objects
func (rm *RenderManager) Render() {
	rsett := settings.GetRenderingSettings()

	w, h := rm.window.Size()
	rm.window.OpenGL().Viewport(0, 0, int32(w), int32(h))

	rsett.MatrixProjection = mgl32.Perspective(mgl32.DegToRad(rsett.Fov), rsett.RatioWidth/rsett.RatioHeight, rsett.PlaneClose, rsett.PlaneFar)
	rm.camera.Render()
	rsett.MatrixCamera = rm.camera.MatrixCamera

	ahPosition := float32(rsett.WorldGridSizeSquares / 2)

	if rsett.WorldGridSizeSquares != rm.gridSize {
		rm.gridSize = rsett.WorldGridSizeSquares
		rm.wgrid.GridSize = rsett.WorldGridSizeSquares
		rm.wgrid.InitBuffers(rsett.WorldGridSizeSquares, 1.0)
		rm.axisLabels.InitBuffers()
	}

	if rsett.ShowCube {
		rm.cube.Render()
	}

	if rsett.ShowGrid {
		rm.wgrid.ActAsMirror = rsett.ActAsMirror
		rm.wgrid.Render()
	}

	if rsett.ShowAxisHelpers {
		rm.axisLabels.Render(ahPosition)
	}
}

// Dispose will cleanup everything
func (rm *RenderManager) Dispose() {
	rm.cube.Dispose()
	rm.wgrid.Dispose()
	rm.camera.Dispose()
	rm.axisLabels.Dispose()
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
}

func (rm *RenderManager) initCamera() {
	rm.camera = objects.InitCamera(rm.window)
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
