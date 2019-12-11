package rendering

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/objects"
	"github.com/supudo/Kuplung-Go/settings"
)

// RenderManager is the main structure for rendering
type RenderManager struct {
	window interfaces.Window

	camera *objects.Camera

	cube  *objects.Cube
	wgrid *objects.WorldGrid

	gridSize int32
}

// NewRenderManager will return an instance of the rendering manager
func NewRenderManager(window interfaces.Window) *RenderManager {
	rm := &RenderManager{}
	rsett := settings.GetRenderingSettings()
	rm.window = window
	rm.gridSize = rsett.WorldGridSizeSquares
	rm.initCamera()
	rm.initCube()
	rm.initWorldGrid()
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

	if rsett.WorldGridSizeSquares != rm.gridSize {
		rm.gridSize = rsett.WorldGridSizeSquares
		rm.wgrid.GridSize = rsett.WorldGridSizeSquares
		rm.wgrid.InitBuffers(rsett.WorldGridSizeSquares, 1.0)
	}

	if rsett.ShowCube {
		rm.cube.Render()
	}

	if rsett.ShowGrid {
		rm.wgrid.ActAsMirror = rsett.ActAsMirror
		rm.wgrid.Render()
	}
}

// Dispose will cleanup everything
func (rm *RenderManager) Dispose() {
	rm.cube.Dispose()
	rm.wgrid.Dispose()
	rm.camera.Dispose()
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
