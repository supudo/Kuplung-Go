package objects

import (
	_ "image/png"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/veandco/go-sdl2/sdl"
)

// Cube ...
type Cube struct {
	window interfaces.Window

	angle             float32
	previousTime      float32
	program           uint32
	texture           uint32
	model             mgl32.Mat4
	modelUniform      int32
	projectionUniform int32
	cameraUniform     int32
	mvpMatrixUniform  int32
	vao               uint32
	version           string

	fov float32
}

// CubeInit ...
func CubeInit(window interfaces.Window) *Cube {
	sett := settings.GetSettings()
	rsett := settings.GetRenderingSettings()

	cube := &Cube{}

	cube.version = "#version 410 core"
	cube.window = window

	vertexShader := cube.version + `
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform mat4 mvp;
in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
void main()
{
	fragTexCoord = vertTexCoord;
	//gl_Position = projection * camera * model * vec4(vert, 1);
	gl_Position = mvp * vec4(vert, 1);
}
` + "\x00"
	fragmentShader := cube.version + `
uniform sampler2D tex;
in vec2 fragTexCoord;
out vec4 outputColor;
void main() {
	outputColor = texture(tex, fragTexCoord);
}
` + "\x00"
	cubeVertices := []float32{
		//  X, Y, Z, U, V
		// Bottom
		-1.0, -1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,

		// Top
		-1.0, 1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0, 1.0,

		// Front
		-1.0, -1.0, 1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,

		// Back
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 1.0,

		// Left
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,

		// Right
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
	}

	cube.fov = rsett.General.Fov

	gl := window.OpenGL()

	cube.model = mgl32.Ident4()

	// Configure the vertex and fragment shaders
	cube.program, _ = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)

	gl.UseProgram(cube.program)

	cube.projectionUniform = gl.GLGetUniformLocation(cube.program, gl.Str("projection\x00"))
	cube.cameraUniform = gl.GLGetUniformLocation(cube.program, gl.Str("camera\x00"))
	cube.modelUniform = gl.GLGetUniformLocation(cube.program, gl.Str("model\x00"))
	cube.mvpMatrixUniform = gl.GLGetUniformLocation(cube.program, gl.Str("mvp\x00"))

	textureUniform := gl.GLGetUniformLocation(cube.program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.GLBindFragDataLocation(cube.program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	cube.texture = engine.LoadTexture(gl, sett.App.CurrentPath+"textures/square.png")

	// Configure the vertex data
	cube.vao = gl.GenVertexArrays(1)[0]

	gl.BindVertexArray(cube.vao)

	vbo := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vbo)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), oglconsts.STATIC_DRAW)

	vertAttrib := uint32(gl.GLGetAttribLocation(cube.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, oglconsts.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GLGetAttribLocation(cube.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, oglconsts.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	cube.angle = 0.0
	cube.previousTime = float32(sdl.GetTicks())

	return cube
}

// Render ...
func (cube *Cube) Render() {
	gl := cube.window.OpenGL()
	rsett := settings.GetRenderingSettings()

	gl.Enable(oglconsts.DEPTH_TEST)
	gl.DepthFunc(oglconsts.LESS)
	gl.Disable(oglconsts.BLEND)
	gl.BlendFunc(oglconsts.SRC_ALPHA, oglconsts.ONE_MINUS_SRC_ALPHA)
	gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.FILL)

	// Update
	sdlTime := float32(sdl.GetTicks())
	elapsed := (sdlTime - cube.previousTime) / 1000
	cube.previousTime = sdlTime

	cube.angle += float32(elapsed)

	cube.model = mgl32.Ident4()
	cube.model = cube.model.Mul4(mgl32.Scale3D(1, 1, 1))
	cube.model = cube.model.Mul4(mgl32.Translate3D(0, 0, 0))
	cube.model = cube.model.Mul4(mgl32.HomogRotate3D(0, mgl32.Vec3{1, 0, 0}))
	cube.model = cube.model.Mul4(mgl32.HomogRotate3D(float32(cube.angle), mgl32.Vec3{0, 1, 0}))
	cube.model = cube.model.Mul4(mgl32.HomogRotate3D(0, mgl32.Vec3{0, 0, 1}))
	cube.model = cube.model.Mul4(mgl32.Translate3D(0, 0, 0))
	cube.model = cube.model.Mul4(mgl32.Translate3D(0, 0, 0))

	mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(cube.model))

	// settings.LogInfo("Projection Matrix : %v", rsett.MatrixProjection.String())
	// settings.LogInfo("Camera Matrix : %v", rsett.MatrixCamera.String())
	// settings.LogInfo("Model Matrix : %v", cube.model.String())
	// settings.LogInfo("MVP Matrix : %v", mvpMatrix.String())

	// Render
	gl.UseProgram(cube.program)
	gl.GLUniformMatrix4fv(cube.modelUniform, 1, false, &cube.model[0])
	gl.GLUniformMatrix4fv(cube.projectionUniform, 1, false, &rsett.MatrixProjection[0])
	gl.GLUniformMatrix4fv(cube.cameraUniform, 1, false, &rsett.MatrixCamera[0])
	gl.GLUniformMatrix4fv(cube.mvpMatrixUniform, 1, false, &mvpMatrix[0])

	gl.BindVertexArray(cube.vao)

	gl.ActiveTexture(oglconsts.TEXTURE0)
	gl.BindTexture(oglconsts.TEXTURE_2D, cube.texture)

	gl.DrawArrays(oglconsts.TRIANGLES, 0, 6*2*3)

	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

// Dispose will cleanup everything
func (cube *Cube) Dispose() {
	gl := cube.window.OpenGL()

	gl.DeleteTextures([]uint32{cube.texture})
	gl.DeleteVertexArrays([]uint32{cube.vao})
	gl.DeleteProgram(cube.program)
}
