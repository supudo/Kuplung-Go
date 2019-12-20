package meshes

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/types"
)

// ModelFace ...
type ModelFace struct {
	window interfaces.Window

	shaderProgram uint32

	VertexSphereVisible, VertexSphereIsSphere, VertexSphereShowWireframes bool

	VertexSphereRadius   float32
	VertexSphereSegments int32
	VertexSphereColor    mgl32.Vec4

	glVAO uint32

	vboTextureAmbient, vboTextureDiffuse, vboTextureSpecular                          uint32
	vboTextureSpecularExp, vboTextureDissolve, vboTextureBump, vboTextureDisplacement uint32

	occQuery uint32

	Model types.MeshModel
}

// NewModelFace ...
func NewModelFace(window interfaces.Window, model types.MeshModel) *ModelFace {
	mesh := &ModelFace{}
	mesh.window = window
	mesh.Model = model
	return mesh
}

// Render ...
func (mesh *ModelFace) Render() {
}

// Dispose ...
func (mesh *ModelFace) Dispose() {
	gl := mesh.window.OpenGL()
	gl.DeleteProgram(mesh.glVAO)
}
