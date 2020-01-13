package objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
)

// RayLine ...
type RayLine struct {
	window interfaces.Window

	shaderProgram      uint32
	glVAO              uint32
	glUniformMVPMatrix int32

	MatrixModel mgl32.Mat4

	AxisSize int32
	x, y, z  float32
}

// NewLightRay ...
func NewLightRay(window interfaces.Window) *RayLine {
	sett := settings.GetSettings()
	gl := window.OpenGL()

	lrModel := &RayLine{}
	lrModel.window = window

	vertexShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/light_ray.vert")
	fragmentShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/light_ray.frag")

	var err error
	lrModel.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[RayLine] Can't load the light ray shaders: %v", err)
	}

	lrModel.glUniformMVPMatrix = gl.GLGetUniformLocation(lrModel.shaderProgram, gl.Str("u_MVPMatrix\x00"))

	lrModel.x = 999
	lrModel.y = 999
	lrModel.z = 999

	return lrModel
}

// InitBuffers ...
func (rl *RayLine) InitBuffers(vecFrom, vecTo mgl32.Vec3) {
	gl := rl.window.OpenGL()

	rl.glVAO = gl.GenVertexArrays(1)[0]

	gl.BindVertexArray(rl.glVAO)

	dataVertices := []float32{vecFrom.X(), vecFrom.Y(), vecFrom.Z(), vecTo.X(), vecTo.Y(), vecTo.Z()}
	dataColors := []float32{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	dataIndices := []uint32{0, 1, 2, 3, 4, 5}

	vboVertices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(dataVertices)*4, gl.Ptr(dataVertices), oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 4, gl.PtrOffset(0))

	vboColors := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboColors)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(dataColors)*4, gl.Ptr(dataColors), oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, oglconsts.FLOAT, false, 4, gl.PtrOffset(0))

	vboIndices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, vboIndices)
	gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, len(dataIndices)*4, gl.Ptr(dataIndices), oglconsts.STATIC_DRAW)

	gl.BindVertexArray(0)

	gl.DeleteBuffers([]uint32{vboVertices, vboColors, vboIndices})
}

// Render ...
func (rl *RayLine) Render(matrixModel mgl32.Mat4) {
	gl := rl.window.OpenGL()
	rsett := settings.GetRenderingSettings()

	gl.UseProgram(rl.shaderProgram)
	gl.BindVertexArray(rl.glVAO)

	mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(matrixModel))
	gl.GLUniformMatrix4fv(rl.glUniformMVPMatrix, 1, false, &mvpMatrix[0])

	gl.LineWidth(2.5)

	gl.Enable(oglconsts.DEPTH_TEST)
	gl.DepthFunc(oglconsts.LESS)
	gl.Disable(oglconsts.BLEND)

	gl.DrawArrays(oglconsts.LINE_STRIP, 0, 2)

	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

// Dispose will cleanup everything
func (rl *RayLine) Dispose() {
	gl := rl.window.OpenGL()

	gl.DeleteVertexArrays([]uint32{rl.glVAO})
	gl.DeleteProgram(rl.shaderProgram)
}
