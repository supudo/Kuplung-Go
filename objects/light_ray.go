package objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
)

// LightRay ...
type LightRay struct {
	window interfaces.Window

	shaderProgram      uint32
	glVAO              uint32
	glUniformMVPMatrix int32

	matrixModel mgl32.Mat4

	AxisSize int32
	x, y, z  float32
}

// InitLightRay ...
func InitLightRay(window interfaces.Window) *LightRay {
	sett := settings.GetSettings()
	gl := window.OpenGL()

	lrModel := &LightRay{}
	lrModel.window = window

	vertexShader := engine.GetShaderSource(sett.App.CurrentPath + "/../Resources/resources/shaders/light_ray.vert")
	fragmentShader := engine.GetShaderSource(sett.App.CurrentPath + "/../Resources/resources/shaders/light_ray.frag")

	var err error
	lrModel.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[LightRay] Can't load the light ray shaders: %v", err)
	}

	lrModel.glUniformMVPMatrix = gl.GLGetUniformLocation(lrModel.shaderProgram, gl.Str("u_MVPMatrix\x00"))

	lrModel.x = 999
	lrModel.y = 999
	lrModel.z = 999

	return lrModel
}

// InitBuffers ...
func (lr *LightRay) InitBuffers(position, direction mgl32.Vec3, simple bool) {
	if (position.X() > lr.x || position.X() < lr.x) && (position.Y() > lr.y || position.Y() < lr.y) && (position.Z() > lr.z || position.Z() < lr.z) {
		gl := lr.window.OpenGL()

		lr.x = position.X()
		lr.y = position.Y()
		lr.z = position.Z()

		lr.glVAO = gl.GenVertexArrays(1)[0]

		gl.BindVertexArray(lr.glVAO)

		lr.AxisSize = 12
		if simple {
			lr.AxisSize = 2
		}
		var vertices []float32
		if simple {
			vertices = append(vertices, position.X())
			vertices = append(vertices, position.Y())
			vertices = append(vertices, position.Z())

			vertices = append(vertices, direction.X())
			vertices = append(vertices, direction.Y())
			vertices = append(vertices, direction.Z())
		} else {
			gVertexBufferData := []float32{-0.866024971, 0, 0.5, 0.866024971, 0, -0.5, 0, 0, -1, 0, 0, -1, -0.866024971, 0, -0.5, -0.866024971, 0, 0.5, -0.866024971, 0, 0.5, 0, 0, 1, 0.866024971, 0, 0.5, 0.866024971, 0, 0.5, 0.866024971, 0, -0.5, -0.866024971, 0, 0.5}
			for i := 0; i < len(gVertexBufferData); i++ {
				vertices = append(vertices, gVertexBufferData[i])
			}
		}

		// vertices
		vboVertices := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		vboIndices := gl.GenBuffers(1)[0]
		if !simple {
			var indices []uint32
			var i uint32
			for i = 0; i < 12; i++ {
				indices = append(indices, i)
			}
			// indices
			gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, vboIndices)
			gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), oglconsts.STATIC_DRAW)
		}

		gl.BindVertexArray(0)
		gl.DeleteBuffers([]uint32{vboVertices, vboIndices})
	}
}

// Render ...
func (lr *LightRay) Render(matrixModel mgl32.Mat4) {
	gl := lr.window.OpenGL()
	rsett := settings.GetRenderingSettings()

	gl.UseProgram(lr.shaderProgram)

	gl.LineWidth(5.5)

	mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(matrixModel))
	gl.GLUniformMatrix4fv(lr.glUniformMVPMatrix, 1, false, &mvpMatrix[0])

	gl.BindVertexArray(lr.glVAO)
	if lr.AxisSize > 2 {
		gl.DrawElements(oglconsts.TRIANGLES, lr.AxisSize, oglconsts.UNSIGNED_INT, 0)
	} else {
		gl.DrawArrays(oglconsts.LINES, 0, lr.AxisSize)
	}

	gl.BindVertexArray(0)

	gl.UseProgram(0)
}

// Dispose will cleanup everything
func (lr *LightRay) Dispose() {
	gl := lr.window.OpenGL()

	gl.DeleteVertexArrays([]uint32{lr.glVAO})
	gl.DeleteProgram(lr.shaderProgram)
}
