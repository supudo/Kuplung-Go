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

	minx, maxx   float32
	miny, maxy   float32
	minz, maxz   float32
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
		minx:   0,
		maxx:   0,
		miny:   0,
		maxy:   0,
		minz:   0,
		maxz:   0,
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
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(bb.meshModel.Vertices)*3*4, gl.Ptr(bb.meshModel.Vertices), oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

	// indices
	bb.dataIndices = []uint32{0, 1, 2, 3, 4, 5, 6, 7, 0, 4, 1, 5, 2, 6, 3, 7}
	vboIndices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, vboIndices)
	gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, int(bb.meshModel.CountIndices)*4, gl.Ptr(bb.meshModel.Indices), oglconsts.STATIC_DRAW)

	bb.minx = 0.0
	bb.maxx = 0.0
	bb.miny = 0.0
	bb.maxy = 0.0
	bb.minz = 0.0
	bb.maxz = 0.0
	if len(bb.meshModel.Vertices) > 0 {
		bb.minx = bb.meshModel.Vertices[0].X()
		bb.maxx = bb.meshModel.Vertices[0].X()
		bb.miny = bb.meshModel.Vertices[0].Y()
		bb.maxy = bb.meshModel.Vertices[0].Y()
		bb.minz = bb.meshModel.Vertices[0].Z()
		bb.maxz = bb.meshModel.Vertices[0].Z()
	}
	for i := 0; i < len(bb.meshModel.Vertices); i++ {
		if bb.meshModel.Vertices[i].X() < bb.minx {
			bb.minx = bb.meshModel.Vertices[i].X()
		}
		if bb.meshModel.Vertices[i].X() > bb.maxx {
			bb.maxx = bb.meshModel.Vertices[i].X()
		}
		if bb.meshModel.Vertices[i].Y() < bb.miny {
			bb.miny = bb.meshModel.Vertices[i].Y()
		}
		if bb.meshModel.Vertices[i].Y() > bb.maxy {
			bb.maxy = bb.meshModel.Vertices[i].Y()
		}
		if bb.meshModel.Vertices[i].Z() < bb.minz {
			bb.minz = bb.meshModel.Vertices[i].Z()
		}
		if bb.meshModel.Vertices[i].Z() > bb.maxz {
			bb.maxz = bb.meshModel.Vertices[i].Z()
		}
	}

	padding := rsett.General.BoundingBoxPadding
	if bb.minx > 0 {
		bb.minx = bb.minx + padding
	} else {
		bb.minx = bb.minx - padding
	}
	if bb.maxx > 0 {
		bb.maxx = bb.maxx + padding
	} else {
		bb.maxx = bb.maxx - padding
	}
	if bb.miny > 0 {
		bb.miny = bb.miny + padding
	} else {
		bb.miny = bb.miny - padding
	}
	if bb.maxy > 0 {
		bb.maxy = bb.maxy + padding
	} else {
		bb.maxy = bb.maxy - padding
	}
	if bb.minz > 0 {
		bb.minz = bb.minz + padding
	} else {
		bb.minz = bb.minz - padding
	}
	if bb.maxz > 0 {
		bb.maxz = bb.maxz + padding
	} else {
		bb.maxz = bb.maxz - padding
	}

	bb.size = mgl32.Vec3{bb.maxx - bb.minx, bb.maxy - bb.miny, bb.maxz - bb.minz}
	bb.center = mgl32.Vec3{((bb.minx + bb.maxx) / 2) * 0.5, ((bb.miny + bb.maxy) / 2) * 0.5, ((bb.minz + bb.maxz) / 2) * 0.5}
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

		// gl.DrawElements(oglconsts.LINE_LOOP, 4, oglconsts.UNSIGNED_INT, 4)
		// gl.DrawElements(oglconsts.LINE_LOOP, 4, oglconsts.UNSIGNED_INT, (4 * 4))
		// gl.DrawElements(oglconsts.LINES, 8, oglconsts.UNSIGNED_INT, (8 * 4))
		gl.DrawElements(oglconsts.TRIANGLES, int32(len(bb.dataIndices)), oglconsts.UNSIGNED_INT, 0)
		gl.BindVertexArray(0)
		gl.UseProgram(0)
	}
}

// Dispose will cleanup everything
func (bb *BoundingBox) Dispose() {
	gl := bb.window.OpenGL()
	gl.DeleteVertexArrays([]uint32{bb.glVAO})
	gl.DeleteProgram(bb.shaderProgram)
}
