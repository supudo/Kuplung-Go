package renderers

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/types"
)

// RendererForward ...
type RendererForward struct {
	window interfaces.Window
}

// NewRendererForward ...
func NewRendererForward(window interfaces.Window) *RendererForward {
	rend := &RendererForward{}
	rend.window = window
	return rend
}

// Render ...
func (rend *RendererForward) Render(rp types.RenderProperties, meshModelFaces []*meshes.ModelFace, matrixGrid mgl32.Mat4, camPos mgl32.Vec3) {
}

// Dispose ...
func (rend *RendererForward) Dispose() {
}
