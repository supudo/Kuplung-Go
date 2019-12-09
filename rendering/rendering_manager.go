package rendering

import (
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/objects"
	"github.com/supudo/Kuplung-Go/settings"
)

// RenderManager is the main structure for rendering
type RenderManager struct {
	window interfaces.Window

	cube  *objects.Cube
	wgrid *objects.WorldGrid
}

// NewRenderManager will return an instance of the rendering manager
func NewRenderManager(window interfaces.Window) *RenderManager {
	rm := &RenderManager{}
	rm.window = window
	rm.initCube()
	rm.initWorldGrid()
	return rm
}

// Render handles rendering of all scene objects
func (rm *RenderManager) Render() {
	rsett := settings.GetRenderingSettings()
	if rsett.ShowCube {
		rm.cube.Render()
	}
	if rsett.ShowGrid {
		rm.wgrid.Render()
	}
}

// Dispose will cleanup everything
func (rm *RenderManager) Dispose() {
	rm.cube.Dispose()
	rm.wgrid.Dispose()
}

func (rm *RenderManager) initCube() {
	rm.cube = objects.CubeInit(rm.window)
}

func (rm *RenderManager) initWorldGrid() {
	rm.wgrid = objects.InitWorldGrid(rm.window)
}
