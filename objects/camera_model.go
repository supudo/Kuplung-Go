package objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// CameraModel ...
type CameraModel struct {
	window interfaces.Window

	shaderProgram uint32
	glVAO         uint32

	glUniformMVPMatrix           int32
	glUniformInnerLightDirection int32
	glUniformColor               int32

	MatrixModel mgl32.Mat4
	model       types.MeshModel

	PositionX, PositionY, PositionZ                                  types.ObjectCoordinate
	RotateX, RotateY, RotateZ                                        types.ObjectCoordinate
	RotateCenterX, RotateCenterY, RotateCenterZ                      types.ObjectCoordinate
	InnerLightDirectionX, InnerLightDirectionY, InnerLightDirectionZ types.ObjectCoordinate
	ColorR, ColorG, ColorB                                           types.ObjectCoordinate

	ShowCameraObject bool
	ShowInWire       bool
}

// InitCameraModel ...
func InitCameraModel(window interfaces.Window, model types.MeshModel) *CameraModel {
	sett := settings.GetSettings()
	gl := window.OpenGL()

	cameraModel := &CameraModel{}
	cameraModel.window = window
	cameraModel.model = model

	vertexShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/camera.vert")
	fragmentShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/camera.frag")

	var err error
	cameraModel.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[CameraModel] Can't load the camera shaders: %v", err)
	}

	cameraModel.glUniformMVPMatrix = gl.GLGetUniformLocation(cameraModel.shaderProgram, gl.Str("u_MVPMatrix\x00"))
	cameraModel.glUniformColor = gl.GLGetUniformLocation(cameraModel.shaderProgram, gl.Str("fs_color\x00"))
	cameraModel.glUniformInnerLightDirection = gl.GLGetUniformLocation(cameraModel.shaderProgram, gl.Str("fs_innerLightDirection\x00"))

	return cameraModel
}

// InitProperties ...
func (cm *CameraModel) InitProperties() {
	cm.PositionX = types.ObjectCoordinate{Animate: false, Point: -6.0}
	cm.PositionY = types.ObjectCoordinate{Animate: false, Point: -2.0}
	cm.PositionZ = types.ObjectCoordinate{Animate: false, Point: 3.0}

	cm.RotateX = types.ObjectCoordinate{Animate: false, Point: 0.0}
	cm.RotateY = types.ObjectCoordinate{Animate: false, Point: 0.0}
	cm.RotateZ = types.ObjectCoordinate{Animate: false, Point: 300.0}

	cm.RotateCenterX = types.ObjectCoordinate{Animate: false, Point: 0.0}
	cm.RotateCenterY = types.ObjectCoordinate{Animate: false, Point: 35.0}
	cm.RotateCenterZ = types.ObjectCoordinate{Animate: false, Point: 0.0}

	cm.InnerLightDirectionX = types.ObjectCoordinate{Animate: false, Point: 1.0}
	cm.InnerLightDirectionY = types.ObjectCoordinate{Animate: false, Point: 0.055}
	cm.InnerLightDirectionZ = types.ObjectCoordinate{Animate: false, Point: 0.206}

	cm.ColorR = types.ObjectCoordinate{Animate: false, Point: 0.61}
	cm.ColorG = types.ObjectCoordinate{Animate: false, Point: 0.61}
	cm.ColorB = types.ObjectCoordinate{Animate: false, Point: 0.61}

	cm.MatrixModel = mgl32.Ident4()

	cm.ShowCameraObject = true
	cm.ShowInWire = false
}

// InitBuffers ...
func (cm *CameraModel) InitBuffers() {
	gl := cm.window.OpenGL()

	cm.glVAO = gl.GenVertexArrays(1)[0]

	gl.BindVertexArray(cm.glVAO)

	vboVertices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(cm.model.Vertices)*3*4, gl.Ptr(cm.model.Vertices), oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

	vboNormals := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboNormals)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(cm.model.Normals)*3*4, gl.Ptr(cm.model.Normals), oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

	vboIndices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, vboIndices)
	gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, int(cm.model.CountIndices)*4, gl.Ptr(cm.model.Indices), oglconsts.STATIC_DRAW)

	gl.BindVertexArray(0)

	gl.DeleteBuffers([]uint32{vboVertices, vboNormals, vboIndices})

	gl.CheckForOpenGLErrors("CameraModel")
}

// Render ...
func (cm *CameraModel) Render(gridMatrix mgl32.Mat4) {
	if cm.ShowCameraObject {
		gl := cm.window.OpenGL()
		rsett := settings.GetRenderingSettings()

		gl.UseProgram(cm.shaderProgram)

		cm.MatrixModel = mgl32.Ident4()
		if rsett.Grid.WorldGridFixedWithWorld {
			cm.MatrixModel = gridMatrix
		}
		cm.MatrixModel = cm.MatrixModel.Mul4(mgl32.Scale3D(1, 1, 1))
		cm.MatrixModel = cm.MatrixModel.Mul4(mgl32.Translate3D(0, 0, 0))
		cm.MatrixModel = cm.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(cm.RotateX.Point), mgl32.Vec3{1, 0, 0}))
		cm.MatrixModel = cm.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(cm.RotateY.Point), mgl32.Vec3{0, 1, 0}))
		cm.MatrixModel = cm.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(cm.RotateZ.Point), mgl32.Vec3{0, 0, 1}))
		cm.MatrixModel = cm.MatrixModel.Mul4(mgl32.Translate3D(0, 0, 0))
		cm.MatrixModel = cm.MatrixModel.Mul4(mgl32.Translate3D(cm.PositionX.Point, cm.PositionY.Point, cm.PositionZ.Point))

		cm.MatrixModel = cm.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(cm.RotateCenterX.Point), mgl32.Vec3{1, 0, 0}))
		cm.MatrixModel = cm.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(cm.RotateCenterY.Point), mgl32.Vec3{0, 1, 0}))
		cm.MatrixModel = cm.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(cm.RotateCenterZ.Point), mgl32.Vec3{0, 0, 1}))

		mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(cm.MatrixModel))
		gl.GLUniformMatrix4fv(cm.glUniformMVPMatrix, 1, false, &mvpMatrix[0])

		gl.Uniform3f(cm.glUniformColor, cm.ColorR.Point, cm.ColorG.Point, cm.ColorB.Point)
		gl.Uniform3f(cm.glUniformInnerLightDirection, cm.InnerLightDirectionX.Point, cm.InnerLightDirectionY.Point, cm.InnerLightDirectionZ.Point)

		// draw
		gl.BindVertexArray(cm.glVAO)
		if cm.ShowInWire {
			gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.LINE)
		}
		gl.DrawElements(oglconsts.TRIANGLES, cm.model.CountIndices, oglconsts.UNSIGNED_INT, 0)
		if cm.ShowInWire {
			gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.FILL)
		}

		gl.BindVertexArray(0)
		gl.UseProgram(0)
	}
}

// Dispose will cleanup everything
func (cm *CameraModel) Dispose() {
	gl := cm.window.OpenGL()

	gl.DeleteVertexArrays([]uint32{cm.glVAO})
	gl.DeleteProgram(cm.shaderProgram)
}
