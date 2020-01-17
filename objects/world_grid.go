package objects

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// WorldGrid ...
type WorldGrid struct {
	window interfaces.Window

	shaderProgram uint32
	glVAO         uint32
	vboVertices   uint32
	vboColors     uint32
	vboIndices    uint32

	glUniformMVPMatrix     int32
	glAttributeActAsMirror int32
	glAttributeAlpha       int32

	MatrixModel mgl32.Mat4

	indicesCount uint32

	ActAsMirror            bool
	actAsMirrorNeedsChange bool
	zIndex                 int
	Transparency           float32

	GridSize int32

	PositionX, PositionY, PositionZ types.ObjectCoordinate
	ScaleX, ScaleY, ScaleZ          types.ObjectCoordinate
	RotateX, RotateY, RotateZ       types.ObjectCoordinate
}

// InitWorldGrid ...
func InitWorldGrid(window interfaces.Window) *WorldGrid {
	sett := settings.GetSettings()
	gl := window.OpenGL()

	grid := &WorldGrid{}
	grid.window = window

	vertexShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/grid2d.vert")
	fragmentShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/grid2d.frag")

	var err error
	grid.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[WorldGrid] Can't load the grid shaders: %v", err)
	}

	grid.glAttributeActAsMirror = gl.GLGetUniformLocation(grid.shaderProgram, gl.Str("a_actAsMirror\x00"))
	grid.glAttributeAlpha = gl.GLGetUniformLocation(grid.shaderProgram, gl.Str("a_alpha\x00"))
	grid.glUniformMVPMatrix = gl.GLGetUniformLocation(grid.shaderProgram, gl.Str("u_MVPMatrix\x00"))

	grid.InitProperties()
	grid.InitBuffers(grid.GridSize, 1)

	gl.CheckForOpenGLErrors("WorldGrid")

	return grid
}

// InitProperties ...
func (grid *WorldGrid) InitProperties() {
	rsett := settings.GetRenderingSettings()
	grid.ActAsMirror = rsett.Grid.ActAsMirror
	grid.actAsMirrorNeedsChange = true
	grid.zIndex = 0
	grid.MatrixModel = mgl32.Ident4()
	grid.GridSize = rsett.Grid.WorldGridSizeSquares
	grid.Transparency = 1.0
}

// InitBuffers ...
func (grid *WorldGrid) InitBuffers(gridSize int32, unitSize float32) {
	gl := grid.window.OpenGL()
	grid.glVAO = gl.GenVertexArrays(1)[0]
	gl.BindVertexArray(grid.glVAO)
	grid.GridSize = gridSize

	if !grid.ActAsMirror {
		grid.actAsMirrorNeedsChange = true
		if grid.GridSize%2 == 0 {
			grid.GridSize++
		}

		dataVertices := []float32{}
		dataColors := []float32{}
		dataIndices := []uint32{}

		indiceCounter := uint32(0)
		perLines := float32(math.Ceil(float64(grid.GridSize) / 2))

		// +
		for z := perLines; z > 0; z-- {
			dataVertices = append(dataVertices, -perLines*unitSize, 0, -z*unitSize)
			dataVertices = append(dataVertices, perLines*unitSize, 0, -z*unitSize)
			dataColors = append(dataColors, 0.7, 0.7, 0.7)
			dataColors = append(dataColors, 0.7, 0.7, 0.7)
			dataIndices = append(dataIndices, indiceCounter)
			indiceCounter++
			dataIndices = append(dataIndices, indiceCounter)
			indiceCounter++
		}

		// X
		dataVertices = append(dataVertices, -perLines*unitSize, 0, 0)
		dataVertices = append(dataVertices, perLines*unitSize, 0, 0)
		dataColors = append(dataColors, 1, 0, 0)
		dataColors = append(dataColors, 1, 0, 0)
		dataIndices = append(dataIndices, indiceCounter)
		indiceCounter++
		dataIndices = append(dataIndices, indiceCounter)
		indiceCounter++
		// X

		for z := float32(1.0); z <= perLines; z++ {
			dataVertices = append(dataVertices, -perLines*unitSize, 0, z*unitSize)
			dataVertices = append(dataVertices, perLines*unitSize, 0, z*unitSize)
			dataColors = append(dataColors, 0.7, 0.7, 0.7)
			dataColors = append(dataColors, 0.7, 0.7, 0.7)
			dataIndices = append(dataIndices, indiceCounter)
			indiceCounter++
			dataIndices = append(dataIndices, indiceCounter)
			indiceCounter++
		}

		// -
		for x := perLines; x > 0; x-- {
			dataVertices = append(dataVertices, -x*unitSize, 0, -perLines*unitSize)
			dataVertices = append(dataVertices, -x*unitSize, 0, perLines*unitSize)
			dataColors = append(dataColors, 0.7, 0.7, 0.7)
			dataColors = append(dataColors, 0.7, 0.7, 0.7)
			dataIndices = append(dataIndices, indiceCounter)
			indiceCounter++
			dataIndices = append(dataIndices, indiceCounter)
			indiceCounter++
		}

		// Z
		dataVertices = append(dataVertices, 0, 0, perLines*unitSize)
		dataVertices = append(dataVertices, 0, 0, -perLines*unitSize)
		dataColors = append(dataColors, 0, 0, 1)
		dataColors = append(dataColors, 0, 0, 1)
		dataIndices = append(dataIndices, indiceCounter)
		indiceCounter++
		dataIndices = append(dataIndices, indiceCounter)
		indiceCounter++
		// Z

		for x := float32(1.0); x <= perLines; x++ {
			dataVertices = append(dataVertices, x*unitSize, 0, -perLines*unitSize)
			dataVertices = append(dataVertices, x*unitSize, 0, perLines*unitSize)
			dataColors = append(dataColors, 0.7, 0.7, 0.7)
			dataColors = append(dataColors, 0.7, 0.7, 0.7)
			dataIndices = append(dataIndices, indiceCounter)
			indiceCounter++
			dataIndices = append(dataIndices, indiceCounter)
			indiceCounter++
		}

		grid.zIndex = len(dataVertices)

		// Y
		dataVertices = append(dataVertices, 0, perLines*unitSize, 0)
		dataVertices = append(dataVertices, 0, -perLines*unitSize, 0)
		dataColors = append(dataColors, 0, 1, 0)
		dataColors = append(dataColors, 0, 1, 0)
		dataIndices = append(dataIndices, indiceCounter)
		indiceCounter++
		dataIndices = append(dataIndices, indiceCounter)
		indiceCounter++

		dataVertices = append(dataVertices, 0, perLines*unitSize, 0)
		dataVertices = append(dataVertices, 0, -perLines*unitSize, 0)
		dataColors = append(dataColors, 0, 1, 0)
		dataColors = append(dataColors, 0, 1, 0)
		dataIndices = append(dataIndices, indiceCounter)
		indiceCounter++
		dataIndices = append(dataIndices, indiceCounter)
		grid.indicesCount = indiceCounter

		grid.vboVertices = gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, grid.vboVertices)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(dataVertices)*4, gl.Ptr(dataVertices), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		grid.vboColors = gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, grid.vboColors)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(dataColors)*4, gl.Ptr(dataColors), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		grid.vboIndices = gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, grid.vboIndices)
		gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, len(dataIndices)*4, gl.Ptr(dataIndices), oglconsts.STATIC_DRAW)
	} else {
		grid.actAsMirrorNeedsChange = false
		planePoint := float32(grid.GridSize / 2)
		dataVertices := []float32{
			planePoint, 0.0, planePoint,
			planePoint, 0.0, -1 * planePoint,
			-1 * planePoint, 0.0, -1 * planePoint,
			-1 * planePoint, 0.0, planePoint,
			planePoint, 0.0, planePoint,
			-1 * planePoint, 0.0, -1 * planePoint}
		dataColors := []float32{
			0.7, 0.7, 0.7,
			0.7, 0.7, 0.7,
			0.7, 0.7, 0.7,
			0.7, 0.7, 0.7}
		dataIndices := []uint32{0, 1, 2, 3, 4, 5}
		grid.indicesCount = 6

		grid.vboVertices = gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, grid.vboVertices)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(dataVertices)*4, gl.Ptr(dataVertices), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(grid.vboVertices)
		gl.VertexAttribPointer(grid.vboVertices, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		grid.vboColors = gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, grid.vboColors)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(dataColors)*4, gl.Ptr(dataColors), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(grid.vboColors)
		gl.VertexAttribPointer(grid.vboColors, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		grid.vboIndices = gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, grid.vboIndices)
		gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, len(dataIndices)*4, gl.Ptr(dataIndices), oglconsts.STATIC_DRAW)
	}

	gl.BindVertexArray(0)

	gl.CheckForOpenGLErrors("WorldGrid-InitBuffers")
}

// Render ...
func (grid *WorldGrid) Render() {
	gl := grid.window.OpenGL()
	rsett := settings.GetRenderingSettings()

	gl.Enable(oglconsts.DEPTH_TEST)
	gl.DepthFunc(oglconsts.LESS)
	gl.Disable(oglconsts.BLEND)
	gl.BlendFunc(oglconsts.SRC_ALPHA, oglconsts.ONE_MINUS_SRC_ALPHA)

	grid.MatrixModel = mgl32.Ident4()
	grid.MatrixModel = grid.MatrixModel.Mul4(mgl32.Scale3D(1, 1, 1))
	grid.MatrixModel = grid.MatrixModel.Mul4(mgl32.Translate3D(0, 0, 0))
	grid.MatrixModel = grid.MatrixModel.Mul4(mgl32.HomogRotate3D(0, mgl32.Vec3{1, 0, 0}))
	grid.MatrixModel = grid.MatrixModel.Mul4(mgl32.HomogRotate3D(0, mgl32.Vec3{0, 1, 0}))
	grid.MatrixModel = grid.MatrixModel.Mul4(mgl32.HomogRotate3D(0, mgl32.Vec3{0, 0, 1}))
	grid.MatrixModel = grid.MatrixModel.Mul4(mgl32.Translate3D(0, 0, 0))
	grid.MatrixModel = grid.MatrixModel.Mul4(mgl32.Translate3D(0, 0, 0))

	mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(grid.MatrixModel))

	// settings.LogInfo("Projection Matrix : %v", rsett.MatrixProjection.String())
	// settings.LogInfo("Camera Matrix : %v", rsett.MatrixCamera.String())
	// settings.LogInfo("Model Matrix : %v", grid.matrixModel.String())
	// settings.LogInfo("MVP Matrix : %v", mvpMatrix.String())
	// settings.LogError("--------")

	if grid.ActAsMirror && grid.actAsMirrorNeedsChange {
		grid.InitBuffers(grid.GridSize, 1)
	} else if !grid.ActAsMirror && !grid.actAsMirrorNeedsChange {
		grid.InitBuffers(grid.GridSize, 1)
	}

	gl.UseProgram(grid.shaderProgram)
	gl.BindVertexArray(grid.glVAO)
	gl.GLUniformMatrix4fv(grid.glUniformMVPMatrix, 1, false, &mvpMatrix[0])

	if !grid.ActAsMirror {
		gl.LineWidth(1.0)
		gl.Uniform1i(grid.glAttributeAlpha, 1.0)
		gl.Uniform1i(grid.glAttributeActAsMirror, 0)
		gl.DrawElements(oglconsts.LINES, int32(grid.indicesCount), oglconsts.UNSIGNED_INT, uintptr(0))
	} else {
		gl.Enable(oglconsts.DEPTH_TEST)
		gl.BlendFunc(oglconsts.SRC_ALPHA, oglconsts.ONE_MINUS_SRC_ALPHA)
		gl.Enable(oglconsts.BLEND)
		gl.Uniform1i(grid.glAttributeAlpha, 1.0)
		gl.Uniform1i(grid.glAttributeActAsMirror, 1)
		gl.DepthMask(false)
		gl.DrawElements(oglconsts.TRIANGLES, int32(grid.indicesCount), oglconsts.UNSIGNED_INT, 0)
		gl.DepthMask(true)
		gl.Disable(oglconsts.BLEND)
		gl.Enable(oglconsts.DEPTH_TEST)
	}

	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

// Dispose will cleanup everything
func (grid *WorldGrid) Dispose() {
	gl := grid.window.OpenGL()

	gl.DeleteBuffers([]uint32{grid.vboVertices, grid.vboColors, grid.vboIndices})
	gl.DeleteVertexArrays([]uint32{grid.glVAO})
	gl.DeleteProgram(grid.shaderProgram)
}
