package objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
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

	matrixModel mgl32.Mat4

	actAsMirror    bool
	gridSizeVertex int32
	zIndex         int
}

// InitWorldGrid ...
func InitWorldGrid(window interfaces.Window) *WorldGrid {
	sett := settings.GetSettings()
	rsett := settings.GetRenderingSettings()
	gl := window.OpenGL()

	grid := &WorldGrid{}

	grid.window = window
	grid.actAsMirror = rsett.ActAsMirror
	grid.gridSizeVertex = rsett.WorldGridSizeSquares
	grid.zIndex = 0
	grid.matrixModel = mgl32.Ident4()

	vertexShader := engine.GetShaderSource(sett.App.CurrentPath + "/../Resources/resources/shaders/grid2d.vert")
	fragmentShader := engine.GetShaderSource(sett.App.CurrentPath + "/../Resources/resources/shaders/grid2d.frag")

	var err error
	grid.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[WorldGrid] Can't load the grid shaders: %v", err)
	}

	grid.glAttributeActAsMirror = gl.GLGetUniformLocation(grid.shaderProgram, gl.Str("a_actAsMirror\x00"))
	grid.glAttributeAlpha = gl.GLGetUniformLocation(grid.shaderProgram, gl.Str("a_alpha\x00"))
	grid.glUniformMVPMatrix = gl.GLGetUniformLocation(grid.shaderProgram, gl.Str("u_MVPMatrix\x00"))

	grid.glVAO = gl.GenVertexArrays(1)[0]
	gl.BindVertexArray(grid.glVAO)

	if grid.gridSizeVertex%2 == 0 {
		grid.gridSizeVertex++
	}

	gridVerticesData := []float32{}
	gridColorsData := []float32{}
	//gridIndicesData := []float32{}

	unitSize := float32(1.0)
	gridMinus := float32(grid.gridSizeVertex / 2)
	var i, j float32
	var x, y, z float32
	for i = 0; i < float32(grid.gridSizeVertex*2); i++ {
		for j = 0; j < float32(grid.gridSizeVertex); j++ {
			if i < float32(grid.gridSizeVertex) {
				x = (j - gridMinus) * unitSize
				y = 0
				z = (i - gridMinus) * unitSize
				gridVerticesData = append(gridVerticesData, x, y, z)
				if z < 0 || z > 0 {
					gridColorsData = append(gridColorsData, 0.7, 0.7, 0.7)
				} else {
					gridColorsData = append(gridColorsData, 1.0, 0.0, 0.0)
				}
			} else {
				x = (i - float32(grid.gridSizeVertex) - gridMinus) * unitSize
				y = 0
				z = (j - gridMinus) * unitSize
				if x < 0 || x > 0 {
					gridColorsData = append(gridColorsData, 0.7, 0.7, 0.7)
				} else {
					gridColorsData = append(gridColorsData, .0, 0.0, 1.0)
				}
			}
		}
	}

	grid.zIndex = len(gridVerticesData)

	x = float32(0.0)
	y = float32(-1.0 * gridMinus)
	z = float32(0.0)
	gridVerticesData = append(gridVerticesData, x, y, z)
	gridColorsData = append(gridColorsData, 0.0, 1.0, 0.0)

	x = float32(0.0)
	y = float32(gridMinus)
	z = float32(0.0)
	gridVerticesData = append(gridVerticesData, x, y, z)
	gridColorsData = append(gridColorsData, 0.0, 1.0, 0.0)

	grid.vboVertices = gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, grid.vboVertices)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(gridVerticesData)*3, gl.Ptr(gridVerticesData), oglconsts.STATIC_DRAW)

	grid.vboColors = gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, grid.vboVertices)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(gridColorsData)*3, gl.Ptr(gridColorsData), oglconsts.STATIC_DRAW)

	// grid.vboIndices = gl.GenBuffers(1)[0]
	// gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, grid.vboVertices)
	// gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, len(gridIndicesData), gl.Ptr(gridIndicesData), oglconsts.STATIC_DRAW)

	return grid
}

// Render ...
func (grid *WorldGrid) Render() {
	gl := grid.window.OpenGL()
	rsett := settings.GetRenderingSettings()

	gl.Enable(oglconsts.DEPTH_TEST)
	gl.DepthFunc(oglconsts.LESS)
	gl.Disable(oglconsts.BLEND)
	gl.BlendFunc(oglconsts.SRC_ALPHA, oglconsts.ONE_MINUS_SRC_ALPHA)

	w, h := grid.window.Size()
	gl.Viewport(0, 0, int32(w), int32(h))

	grid.matrixModel = mgl32.Ident4()
	grid.matrixModel = grid.matrixModel.Mul4(mgl32.Scale3D(1, 1, 1))
	grid.matrixModel = grid.matrixModel.Mul4(mgl32.Translate3D(0, 0, 0))
	grid.matrixModel = grid.matrixModel.Mul4(mgl32.HomogRotate3D(0, mgl32.Vec3{1, 0, 0}))
	grid.matrixModel = grid.matrixModel.Mul4(mgl32.HomogRotate3D(0, mgl32.Vec3{0, 1, 0}))
	grid.matrixModel = grid.matrixModel.Mul4(mgl32.HomogRotate3D(0, mgl32.Vec3{0, 0, 1}))
	grid.matrixModel = grid.matrixModel.Mul4(mgl32.Translate3D(0, 0, 0))
	grid.matrixModel = grid.matrixModel.Mul4(mgl32.Translate3D(0, 0, 0))

	mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(grid.matrixModel))

	gl.UseProgram(grid.shaderProgram)
	gl.GLUniformMatrix4fv(grid.glUniformMVPMatrix, 1, false, &mvpMatrix[0])
	gl.BindVertexArray(grid.glVAO)
	gl.LineWidth(1.0)
	gl.Uniform1i(grid.glAttributeAlpha, 1.0)
	gl.Uniform1i(grid.glAttributeActAsMirror, 0)
	var i int32
	for i = 0; i < grid.gridSizeVertex*2; i++ {
		gl.DrawArrays(oglconsts.LINE_STRIP, grid.gridSizeVertex*i, grid.gridSizeVertex)
	}
	for i = 0; i < grid.gridSizeVertex; i++ {
		gl.DrawArrays(oglconsts.LINE_STRIP, 0, grid.gridSizeVertex)
	}
	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

// Dispose will cleanup everything
func (grid *WorldGrid) Dispose() {
	gl := grid.window.OpenGL()

	gl.DeleteVertexArrays([]uint32{grid.glVAO})
	gl.DeleteProgram(grid.shaderProgram)
}
