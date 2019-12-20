package renderers

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/types"
)

// RendererShadowMapping ...
type RendererShadowMapping struct {
	window interfaces.Window
}

// NewRendererShadowMapping ...
func NewRendererShadowMapping(window interfaces.Window) *RendererShadowMapping {
	rend := &RendererShadowMapping{}
	rend.window = window
	return rend
}

// Render ...
func (rend *RendererShadowMapping) Render(rp types.RenderProperties, meshModelFaces []*meshes.ModelFace, matrixGrid mgl32.Mat4, camPos mgl32.Vec3) {
}

// Dispose ...
func (rend *RendererShadowMapping) Dispose() {
}
