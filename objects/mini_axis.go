package objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// MiniAxis ...
type MiniAxis struct {
	window interfaces.Window

	shaderProgram      uint32
	glVAO              uint32
	glUniformMVPMatrix int32

	ShowAxis                  bool
	AxisSize                  int32
	RotateX, RotateY, RotateZ types.ObjectCoordinate
	matrixModel               mgl32.Mat4
}

// InitMiniAxis ...
func InitMiniAxis(window interfaces.Window) *MiniAxis {
	sett := settings.GetSettings()
	gl := window.OpenGL()

	miniAxis := &MiniAxis{}
	miniAxis.window = window

	vertexShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/axis.vert")
	fragmentShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/axis.frag")

	var err error
	miniAxis.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[MiniAxis] Can't load the mini axis shaders: %v", err)
	}

	miniAxis.glUniformMVPMatrix = gl.GLGetUniformLocation(miniAxis.shaderProgram, gl.Str("u_MVPMatrix\x00"))

	miniAxis.InitProperties()

	return miniAxis
}

// InitProperties ...
func (ma *MiniAxis) InitProperties() {
	ma.ShowAxis = true
	ma.AxisSize = 100

	ma.RotateX = types.ObjectCoordinate{Animate: false, Point: 0.0}
	ma.RotateY = types.ObjectCoordinate{Animate: false, Point: 0.0}
	ma.RotateZ = types.ObjectCoordinate{Animate: false, Point: 0.0}

	ma.matrixModel = mgl32.Ident4()
}

// InitBuffers ...
func (ma *MiniAxis) InitBuffers() {
	gl := ma.window.OpenGL()

	gVertexBufferData := []float32{
		// X
		-100, 0, 0,
		100, 0, 0,
		// Y
		0, -100, 0,
		0, 100, 0,
		// Z
		0, 0, -100,
		0, 0, 100}

	gColorBufferData := []float32{
		// X - red
		1.0, 0.0, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0,
		// Y - green
		0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0,
		// Z - blue
		0.0, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, 1.0}

	ma.glVAO = gl.GenVertexArrays(1)[0]

	gl.BindVertexArray(ma.glVAO)

	vboVertices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
	gl.BufferData(oglconsts.ARRAY_BUFFER, 18*4, gVertexBufferData, oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

	// colors
	vboColors := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboColors)
	gl.BufferData(oglconsts.ARRAY_BUFFER, 24*4, gColorBufferData, oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, oglconsts.FLOAT, false, 4*4, gl.PtrOffset(0))

	gl.BindVertexArray(0)

	gl.DeleteBuffers([]uint32{vboVertices, vboColors})
}

// Render ...
func (ma *MiniAxis) Render() {
	if ma.ShowAxis {
		gl := ma.window.OpenGL()
		rsett := settings.GetRenderingSettings()
		wi, hi := ma.window.Size()
		w := int32(wi)
		h := int32(hi)

		gl.UseProgram(ma.shaderProgram)

		axisW := int32(120)
		axisH := int32((h * axisW) / w)
		axisX := int32(10)
		axisY := int32(10)

		gl.Viewport(axisX, axisY, axisW, axisH)

		ma.matrixModel = mgl32.Ident4()
		ma.matrixModel = ma.matrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(ma.RotateX.Point), mgl32.Vec3{1, 0, 0}))
		ma.matrixModel = ma.matrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(ma.RotateY.Point), mgl32.Vec3{0, 1, 0}))
		ma.matrixModel = ma.matrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(ma.RotateZ.Point), mgl32.Vec3{0, 0, 1}))

		gl.LineWidth(2.0)

		mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(ma.matrixModel))

		gl.GLUniformMatrix4fv(ma.glUniformMVPMatrix, 1, false, &mvpMatrix[0])

		gl.BindVertexArray(ma.glVAO)
		gl.DrawArrays(oglconsts.LINES, 0, ma.AxisSize)
		gl.BindVertexArray(0)
		gl.UseProgram(0)

		gl.Viewport(0, 0, w, h)
	}
}

// Dispose will cleanup everything
func (ma *MiniAxis) Dispose() {
	gl := ma.window.OpenGL()

	gl.DeleteVertexArrays([]uint32{ma.glVAO})
	gl.DeleteProgram(ma.shaderProgram)
}
