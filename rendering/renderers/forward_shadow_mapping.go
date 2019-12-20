package renderers

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/types"
)

// RendererForwardShadowMapping ...
type RendererForwardShadowMapping struct {
	window interfaces.Window
}

// NewRendererForwardShadowMapping ...
func NewRendererForwardShadowMapping(window interfaces.Window) *RendererForwardShadowMapping {
	rend := &RendererForwardShadowMapping{}
	rend.window = window
	return rend
}

// Render ...
func (rend *RendererForwardShadowMapping) Render(rp types.RenderProperties, meshModelFaces []*meshes.ModelFace, matrixGrid mgl32.Mat4, camPos mgl32.Vec3) {
}

// Dispose ...
func (rend *RendererForwardShadowMapping) Dispose() {
}
