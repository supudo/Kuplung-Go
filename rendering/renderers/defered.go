package renderers

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/types"
)

// RendererDefered ...
type RendererDefered struct {
	window interfaces.Window
}

// NewRendererDefered ...
func NewRendererDefered(window interfaces.Window) *RendererDefered {
	rend := &RendererDefered{}
	rend.window = window
	return rend
}

// Render ...
func (rend *RendererDefered) Render(rp types.RenderProperties, meshModelFaces []*meshes.ModelFace, matrixGrid mgl32.Mat4, camPos mgl32.Vec3) {
}

// Dispose ...
func (rend *RendererDefered) Dispose() {
}
