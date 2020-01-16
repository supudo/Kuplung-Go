package renderers

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/objects"
	"github.com/supudo/Kuplung-Go/types"
)

// RendererDefered ...
type RendererDefered struct {
	window interfaces.Window

	shaderProgramGeometryPass     uint32
	shaderProgramLightingPass     uint32
	shaderProgramLightBox         uint32
	glGeometryPassTextureDiffuse  int32
	glGeometryPassTextureSpecular int32

	matrixProjection, matrixCamera mgl32.Mat4

	GLSLLightSourceNumberDirectional uint32
	GLSLLightSourceNumberPoint       uint32
	GLSLLightSourceNumberSpot        uint32

	mfLightsDirectional []*types.ModelFaceLightSourceDirectional
	mfLightsPoint       []*types.ModelFaceLightSourcePoint
	mfLightsSpot        []*types.ModelFaceLightSourceSpot

	gBuffer, gPosition, gNormal, gAlbedoSpec uint32

	NRLIGHTS       uint16
	lightPositions []mgl32.Vec3
	lightColors    []mgl32.Vec3

	quadVAO uint32
	quadVBO uint32
	cubeVAO uint32
	cubeVBO uint32
}

// NewRendererDefered ...
func NewRendererDefered(window interfaces.Window) *RendererDefered {
	rend := &RendererDefered{}
	rend.window = window
	rend.NRLIGHTS = 32
	rend.GLSLLightSourceNumberDirectional = 0
	rend.GLSLLightSourceNumberPoint = 0
	rend.GLSLLightSourceNumberSpot = 0
	return rend
}

// Render ...
func (rend *RendererDefered) Render(rp types.RenderProperties, meshModelFaces []*meshes.ModelFace, matrixGrid mgl32.Mat4, camPos mgl32.Vec3, selectedModel int32, lightSources []*objects.Light) {
}

// Dispose ...
func (rend *RendererDefered) Dispose() {
	gl := rend.window.OpenGL()

	gl.DeleteProgram(rend.shaderProgramGeometryPass)
	gl.DeleteProgram(rend.shaderProgramLightingPass)
	gl.DeleteProgram(rend.shaderProgramLightBox)
}
