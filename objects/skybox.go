package objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
)

// SkyBox ...
type SkyBox struct {
	window interfaces.Window

	shaderProgram uint32
	glVAO         uint32

	glVSMatrixView       int32
	glVSMatrixProjection int32

	matrixModel mgl32.Mat4
	vboTexture  uint32

	gridSize           int32
	SkyboxSelectedItem int32
	SkyboxItems        []SkyboxItem
}

// SkyboxItem ...
type SkyboxItem struct {
	Title  string
	Images []string
}

// InitSkyBox ...
func InitSkyBox(window interfaces.Window) *SkyBox {
	rsett := settings.GetRenderingSettings()

	skyBox := &SkyBox{}
	skyBox.window = window

	skyBox.gridSize = rsett.Grid.WorldGridSizeSquares
	skyBox.SkyboxSelectedItem = 0

	skyBox.SkyboxItems = append(skyBox.SkyboxItems, SkyboxItem{Title: "-- No box --", Images: []string{"", "", "", "", "", ""}})
	skyBox.SkyboxItems = append(skyBox.SkyboxItems, SkyboxItem{Title: "Lake Mountain", Images: []string{"lake_mountain_right.jpg", "lake_mountain_left.jpg", "lake_mountain_top.jpg", "lake_mountain_bottom.jpg", "lake_mountain_back.jpg", "lake_mountain_front.jpg"}})
	skyBox.SkyboxItems = append(skyBox.SkyboxItems, SkyboxItem{Title: "Fire Planet", Images: []string{"fire_planet_right.jpg", "fire_planet_left.jpg", "fire_planet_top.jpg", "fire_planet_bottom.jpg", "fire_planet_back.jpg", "fire_planet_front.jpg"}})
	skyBox.SkyboxItems = append(skyBox.SkyboxItems, SkyboxItem{Title: "Stormy Days", Images: []string{"stormydays_right.jpg", "stormydays_left.jpg", "stormydays_top.jpg", "stormydays_bottom.jpg", "stormydays_back.jpg", "stormydays_front.jpg"}})

	return skyBox
}

// InitBuffers ...
func (sb *SkyBox) InitBuffers() {
	sett := settings.GetSettings()
	gl := sb.window.OpenGL()

	vertexShader := engine.GetShaderSource(sett.App.CurrentPath + "/../Resources/resources/shaders/skybox.vert")
	fragmentShader := engine.GetShaderSource(sett.App.CurrentPath + "/../Resources/resources/shaders/skybox.frag")

	var err error
	sb.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[SkyBox] Can't load the camera shaders: %v", err)
	}

	sb.glVSMatrixView = gl.GLGetUniformLocation(sb.shaderProgram, gl.Str("vs_MatrixView\x00"))
	sb.glVSMatrixProjection = gl.GLGetUniformLocation(sb.shaderProgram, gl.Str("vs_MatrixProjection\x00"))

	if sb.SkyboxSelectedItem > 0 {
		sb.glVAO = gl.GenVertexArrays(1)[0]
		gl.BindVertexArray(sb.glVAO)

		// Positions
		skyboxVertices := []float32{
			-1.0, 1.0, -1.0, -1.0, -1.0, -1.0, 1.0, -1.0, -1.0, 1.0, -1.0, -1.0, 1.0, 1.0, -1.0, -1.0, 1.0, -1.0,
			-1.0, -1.0, 1.0, -1.0, -1.0, -1.0, -1.0, 1.0, -1.0, -1.0, 1.0, -1.0, -1.0, 1.0, 1.0, -1.0, -1.0, 1.0,
			1.0, -1.0, -1.0, 1.0, -1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, -1.0, 1.0, -1.0, -1.0,
			-1.0, -1.0, 1.0, -1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, -1.0, 1.0, -1.0, -1.0, 1.0,
			-1.0, 1.0, -1.0, 1.0, 1.0, -1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, -1.0, 1.0, 1.0, -1.0, 1.0, -1.0,
			-1.0, -1.0, -1.0, -1.0, -1.0, 1.0, 1.0, -1.0, -1.0, 1.0, -1.0, -1.0, -1.0, -1.0, 1.0, 1.0, -1.0, 1.0}

		for i := 0; i < len(skyboxVertices); i++ {
			skyboxVertices[i] *= float32(sb.gridSize) * 10.0
		}

		vboVertices := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(skyboxVertices)*4, gl.Ptr(skyboxVertices), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		sb.vboTexture = engine.LoadCubemapTexture(gl, sb.SkyboxItems[sb.SkyboxSelectedItem].Images)

		gl.BindVertexArray(0)
	}
}

// Render ...
func (sb *SkyBox) Render() {
	rsett := settings.GetRenderingSettings()
	if rsett.SkyBox.SkyboxSelectedItem != sb.SkyboxSelectedItem {
		sb.SkyboxSelectedItem = rsett.SkyBox.SkyboxSelectedItem
		sb.InitBuffers()
	}
	if sb.SkyboxSelectedItem > 0 {
		gl := sb.window.OpenGL()
		sett := settings.GetSettings()

		gl.BindVertexArray(sb.glVAO)
		gl.UseProgram(sb.shaderProgram)

		var currentDepthFuncMode int32
		gl.GetIntegerv(oglconsts.DEPTH_FUNC, &currentDepthFuncMode)

		gl.DepthFunc(oglconsts.LEQUAL)

		gl.GLUniformMatrix4fv(sb.glVSMatrixView, 1, false, &rsett.MatrixCamera[0])

		matrixProjection := mgl32.Perspective(rsett.General.Fov, sett.AppWindow.SDLWindowWidth/sett.AppWindow.SDLWindowHeight, rsett.General.PlaneClose, rsett.General.PlaneFar)
		gl.GLUniformMatrix4fv(sb.glVSMatrixProjection, 1, false, &matrixProjection[0])

		gl.BindTexture(oglconsts.TEXTURE_CUBE_MAP, sb.vboTexture)

		gl.DrawArrays(oglconsts.TRIANGLES, 0, 36)

		gl.UseProgram(0)
		gl.BindVertexArray(0)

		gl.DepthFunc(uint32(currentDepthFuncMode))
	}
}

// Dispose will cleanup everything
func (sb *SkyBox) Dispose() {
	gl := sb.window.OpenGL()

	gl.DeleteVertexArrays([]uint32{sb.glVAO})
	gl.DeleteProgram(sb.shaderProgram)
}
