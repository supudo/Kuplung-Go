package objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// AxisLabels ...
type AxisLabels struct {
	window interfaces.Window

	shaderProgram      uint32
	glUniformMVPMatrix int32

	ahPosition float32
	labels     []ModelLabel
}

// ModelLabel ...
type ModelLabel struct {
	glVAO        uint32
	model        types.MeshModel
	translatePos mgl32.Vec3
}

// InitAxisLabels ...
func InitAxisLabels(window interfaces.Window) *AxisLabels {
	sett := settings.GetSettings()
	gl := window.OpenGL()

	axisLabels := &AxisLabels{}
	axisLabels.window = window

	vertexShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/axis_labels.vert")
	fragmentShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/axis_labels.frag")

	var err error
	axisLabels.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[AxisLabels] Can't load the axis labels shaders: %v", err)
	}

	axisLabels.glUniformMVPMatrix = gl.GLGetUniformLocation(axisLabels.shaderProgram, gl.Str("u_MVPMatrix\x00"))

	gl.CheckForOpenGLErrors("AxisLabels")

	return axisLabels
}

// SetModels ...
func (al *AxisLabels) SetModels(models []types.MeshModel, position float32) {
	for i := 0; i < len(models); i++ {
		var color mgl32.Vec3
		lm := ModelLabel{}
		lm.model = models[i]
		switch lm.model.ModelTitle {
		case "XPlus":
			lm.translatePos = mgl32.Vec3{position, 0, 0}
			color = mgl32.Vec3{1, 0, 0}
		case "XMinus":
			lm.translatePos = mgl32.Vec3{-position, -1, 0}
			color = mgl32.Vec3{1, 0, 0}
		case "YPlus":
			lm.translatePos = mgl32.Vec3{0, position, 0}
			color = mgl32.Vec3{0, 1, 0}
		case "YMinus":
			lm.translatePos = mgl32.Vec3{0, -position, 0}
			color = mgl32.Vec3{0, 1, 0}
		case "ZPlus":
			lm.translatePos = mgl32.Vec3{0, 0, position}
			color = mgl32.Vec3{0, 0, 1}
		case "ZMinus":
			lm.translatePos = mgl32.Vec3{0, 0, -position}
			color = mgl32.Vec3{0, 0, 1}
		}

		for j := 0; j < len(models[i].Vertices); j++ {
			lm.model.Colors = append(lm.model.Colors, color)
		}
		lm.model.CountColors = int32(len(lm.model.Colors))

		al.labels = append(al.labels, lm)
	}
}

// InitBuffers ...
func (al *AxisLabels) InitBuffers() {
	gl := al.window.OpenGL()

	for i := 0; i < len(al.labels); i++ {
		switch al.labels[i].model.ModelTitle {
		case "XPlus":
			al.labels[i].translatePos = mgl32.Vec3{al.ahPosition + 1, 0, 0}
		case "XMinus":
			al.labels[i].translatePos = mgl32.Vec3{-al.ahPosition - 1, -1, 0}
		case "YPlus":
			al.labels[i].translatePos = mgl32.Vec3{0, al.ahPosition + 1, 0}
		case "YMinus":
			al.labels[i].translatePos = mgl32.Vec3{0, -al.ahPosition - 1, 0}
		case "ZPlus":
			al.labels[i].translatePos = mgl32.Vec3{0, 0, al.ahPosition + 1}
		case "ZMinus":
			al.labels[i].translatePos = mgl32.Vec3{0, 0, -al.ahPosition - 1}
		}

		al.labels[i].glVAO = gl.GenVertexArrays(1)[0]

		gl.BindVertexArray(al.labels[i].glVAO)

		vboVertices := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(al.labels[i].model.Vertices)*3*4, gl.Ptr(al.labels[i].model.Vertices), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		vboColors := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboColors)
		gl.BufferData(oglconsts.ARRAY_BUFFER, int(al.labels[i].model.CountColors)*3*4, gl.Ptr(al.labels[i].model.Colors), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		vboIndices := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, vboIndices)
		gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, int(al.labels[i].model.CountIndices)*4, gl.Ptr(al.labels[i].model.Indices), oglconsts.STATIC_DRAW)

		gl.BindVertexArray(0)

		gl.DeleteBuffers([]uint32{vboVertices, vboColors, vboIndices})
	}
}

// Render ...
func (al *AxisLabels) Render(position float32) {
	gl := al.window.OpenGL()
	rsett := settings.GetRenderingSettings()

	if position != al.ahPosition {
		al.ahPosition = position
		al.InitBuffers()
	}

	mvpMatrix := mgl32.Ident4()

	gl.UseProgram(al.shaderProgram)

	for i := 0; i < len(al.labels); i++ {
		mvpMatrix = rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(mgl32.Ident4().Mul4(mgl32.Translate3D(al.labels[i].translatePos.X(), al.labels[i].translatePos.Y(), al.labels[i].translatePos.Z()))))
		gl.BindVertexArray(al.labels[i].glVAO)
		gl.GLUniformMatrix4fv(al.glUniformMVPMatrix, 1, false, &mvpMatrix[0])
		gl.DrawElements(oglconsts.TRIANGLES, al.labels[i].model.CountIndices, oglconsts.UNSIGNED_INT, 0)
	}

	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

// Dispose will cleanup everything
func (al *AxisLabels) Dispose() {
	gl := al.window.OpenGL()

	for i := 0; i < len(al.labels); i++ {
		gl.DeleteVertexArrays([]uint32{al.labels[i].glVAO})
	}
	gl.DeleteProgram(al.shaderProgram)
}
