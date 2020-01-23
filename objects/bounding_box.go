package objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// BoundingBox ...
type BoundingBox struct {
	window interfaces.Window

	meshModel types.MeshModel

	matrixTransform mgl32.Mat4

	MinX, MaxX   float32
	MinY, MaxY   float32
	MinZ, MaxZ   float32
	size, center mgl32.Vec3

	dataVertices []float32
	dataIndices  []uint32

	shaderProgram                uint32
	shaderVertex, shaderFragment uint32
	glVAO                        uint32

	glUniformMVPMatrix int32
	glUniformColor     int32
}

// InitBoundingBox ...
func InitBoundingBox(window interfaces.Window) *BoundingBox {
	bb := &BoundingBox{
		window: window,
		MinX:   0,
		MaxX:   0,
		MinY:   0,
		MaxY:   0,
		MinZ:   0,
		MaxZ:   0,
	}
	return bb
}

// InitShaderProgram ...
func (bb *BoundingBox) InitShaderProgram() {
	sett := settings.GetSettings()
	gl := bb.window.OpenGL()

	vertexShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/bounding_box.vert")
	fragmentShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/bounding_box.frag")

	var err error
	bb.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[BoundingBox] Can't load the bounding box shaders: %v", err)
	}

	bb.glUniformMVPMatrix = gl.GLGetUniformLocation(bb.shaderProgram, gl.Str("u_MVPMatrix\x00"))
	bb.glUniformColor = gl.GLGetUniformLocation(bb.shaderProgram, gl.Str("fs_color\x00"))

	gl.CheckForOpenGLErrors("BoundingBox")
}

// InitBuffers ...
func (bb *BoundingBox) InitBuffers(meshModel types.MeshModel) {
	rsett := settings.GetRenderingSettings()
	gl := bb.window.OpenGL()
	bb.meshModel = meshModel

	bb.glVAO = gl.GenVertexArrays(1)[0]

	gl.BindVertexArray(bb.glVAO)

	// vertices
	bb.dataVertices = []float32{-0.5, -0.5, -0.5, 1.0, 0.5, -0.5, -0.5, 1.0, 0.5, 0.5, -0.5, 1.0, -0.5, 0.5, -0.5, 1.0, -0.5, -0.5, 0.5, 1.0, 0.5, -0.5, 0.5, 1.0, 0.5, 0.5, 0.5, 1.0, -0.5, 0.5, 0.5, 1.0}
	vboVertices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(bb.dataVertices)*3*4, gl.Ptr(bb.dataVertices), oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

	// indices
	bb.dataIndices = []uint32{0, 1, 2, 3, 4, 5, 6, 7, 0, 4, 1, 5, 2, 6, 3, 7}
	vboIndices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, vboIndices)
	gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, len(bb.dataIndices)*4, gl.Ptr(bb.dataIndices), oglconsts.STATIC_DRAW)

	bb.MinX = 0.0
	bb.MaxX = 0.0
	bb.MinY = 0.0
	bb.MaxY = 0.0
	bb.MinZ = 0.0
	bb.MaxZ = 0.0
	if len(bb.meshModel.Vertices) > 0 {
		bb.MinX = bb.meshModel.Vertices[0].X()
		bb.MaxX = bb.meshModel.Vertices[0].X()
		bb.MinY = bb.meshModel.Vertices[0].Y()
		bb.MaxY = bb.meshModel.Vertices[0].Y()
		bb.MinZ = bb.meshModel.Vertices[0].Z()
		bb.MaxZ = bb.meshModel.Vertices[0].Z()
	}
	for i := 0; i < len(bb.meshModel.Vertices); i++ {
		if bb.meshModel.Vertices[i].X() < bb.MinX {
			bb.MinX = bb.meshModel.Vertices[i].X()
		}
		if bb.meshModel.Vertices[i].X() > bb.MaxX {
			bb.MaxX = bb.meshModel.Vertices[i].X()
		}
		if bb.meshModel.Vertices[i].Y() < bb.MinY {
			bb.MinY = bb.meshModel.Vertices[i].Y()
		}
		if bb.meshModel.Vertices[i].Y() > bb.MaxY {
			bb.MaxY = bb.meshModel.Vertices[i].Y()
		}
		if bb.meshModel.Vertices[i].Z() < bb.MinZ {
			bb.MinZ = bb.meshModel.Vertices[i].Z()
		}
		if bb.meshModel.Vertices[i].Z() > bb.MaxZ {
			bb.MaxZ = bb.meshModel.Vertices[i].Z()
		}
	}

	padding := rsett.General.BoundingBoxPadding
	if bb.MinX > 0 {
		bb.MinX = bb.MinX + padding
	} else {
		bb.MinX = bb.MinX - padding
	}
	if bb.MaxX > 0 {
		bb.MaxX = bb.MaxX + padding
	} else {
		bb.MaxX = bb.MaxX - padding
	}
	if bb.MinY > 0 {
		bb.MinY = bb.MinY + padding
	} else {
		bb.MinY = bb.MinY - padding
	}
	if bb.MaxY > 0 {
		bb.MaxY = bb.MaxY + padding
	} else {
		bb.MaxY = bb.MaxY - padding
	}
	if bb.MinZ > 0 {
		bb.MinZ = bb.MinZ + padding
	} else {
		bb.MinZ = bb.MinZ - padding
	}
	if bb.MaxZ > 0 {
		bb.MaxZ = bb.MaxZ + padding
	} else {
		bb.MaxZ = bb.MaxZ - padding
	}

	bb.size = mgl32.Vec3{bb.MaxX - bb.MinX, bb.MaxY - bb.MinY, bb.MaxZ - bb.MinZ}
	bb.center = mgl32.Vec3{((bb.MinX + bb.MaxX) / 2) * 0.5, ((bb.MinY + bb.MaxY) / 2) * 0.5, ((bb.MinZ + bb.MaxZ) / 2) * 0.5}
	mtxs := mgl32.Ident4().Mul4(mgl32.Scale3D(bb.size.X(), bb.size.Y(), bb.size.Z()))
	mtxt := mgl32.Ident4().Mul4(mgl32.Translate3D(bb.center.X(), bb.center.Y(), bb.center.Z()))
	bb.matrixTransform = mtxs.Mul4(mtxt)

	rsett.General.BoundingBoxRefresh = false
	gl.BindVertexArray(0)
	gl.DeleteBuffers([]uint32{vboVertices, vboIndices})

	gl.CheckForOpenGLErrors("BoundingBox-InitBuffers")
}

// Render ...
func (bb *BoundingBox) Render(mtxModel mgl32.Mat4, outlineColor mgl32.Vec4) {
	rsett := settings.GetRenderingSettings()
	gl := bb.window.OpenGL()

	if rsett.General.BoundingBoxRefresh {
		bb.InitBuffers(bb.meshModel)
	}
	if bb.glVAO > 0 {
		gl.UseProgram(bb.shaderProgram)
		gl.BindVertexArray(bb.glVAO)

		mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(mtxModel)).Mul4(bb.matrixTransform)
		gl.GLUniformMatrix4fv(bb.glUniformMVPMatrix, 1, false, &mvpMatrix[0])
		gl.Uniform3f(bb.glUniformColor, outlineColor.X(), outlineColor.Y(), outlineColor.Z())

		gl.DrawElementsOffset(oglconsts.LINE_LOOP, 4, oglconsts.UNSIGNED_INT, 0)
		gl.DrawElementsOffset(oglconsts.LINE_LOOP, 4, oglconsts.UNSIGNED_INT, 4*4)
		gl.DrawElementsOffset(oglconsts.LINES, 8, oglconsts.UNSIGNED_INT, 8*4)
		gl.BindVertexArray(0)
		gl.UseProgram(0)

		gl.CheckForOpenGLErrors("BoundingBox-Render")
	}
}

// Dispose will cleanup everything
func (bb *BoundingBox) Dispose() {
	gl := bb.window.OpenGL()
	gl.DeleteVertexArrays([]uint32{bb.glVAO})
	gl.DeleteProgram(bb.shaderProgram)
}
