package rendering

import (
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/objects"
)

// RenderManager is the main structure for rendering
type RenderManager struct {
	window interfaces.Window

	cube *objects.Cube
}

// NewRenderManager will return an instance of the rendering manager
func NewRenderManager(window interfaces.Window) *RenderManager {
	rm := &RenderManager{}
	rm.window = window

	rm.initCube()

	return rm
}

// Render handles rendering of all scene objects
func (rm *RenderManager) Render() {
	rm.cube.Render()
}

// Dispose will cleanup everything
func (rm *RenderManager) Dispose() {
	rm.cube.Dispose()
}

func (rm *RenderManager) initCube() {
	rm.cube = objects.CubeInit(rm.window)
}
