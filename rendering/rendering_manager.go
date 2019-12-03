package rendering

import (
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/objects"
)

// RenderManager ...
type RenderManager struct {
	glWrapper interfaces.OpenGL

	cube *objects.Cube
}

// NewRenderManager ...
func NewRenderManager(gl interfaces.OpenGL) *RenderManager {
	rm := &RenderManager{}
	rm.glWrapper = gl
	rm.cube = objects.CubeInit(gl)
	return rm
}

// Render ...
func (rm *RenderManager) Render() {
	rm.cube.Render()
}
