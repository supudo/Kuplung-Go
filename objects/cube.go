package objects

import (
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/veandco/go-sdl2/sdl"
)

// Cube ...
type Cube struct {
	window interfaces.Window

	angle        float32
	previousTime uint32
	program      uint32
	texture      uint32

	modelUniform int32
	model        mgl32.Mat4

	vao uint32

	version string
}

// CubeInit ...
func CubeInit(window interfaces.Window) *Cube {
	sett := settings.GetSettings()

	cube := &Cube{}

	cube.version = "#version 330"
	cube.window = window

	vertexShader := cube.version + `
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
void main()
{
	fragTexCoord = vertTexCoord;
	gl_Position = projection * camera * model * vec4(vert, 1);
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

	gl := window.OpenGL()

	// Configure the vertex and fragment shaders
	program := cube.newProgram(vertexShader, fragmentShader)

	gl.UseProgram(program)

	gl.LogOpenGLError()

	w, h := window.Size()

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(w/h), 0.1, 10.0)
	projectionUniform := gl.GLGetUniformLocation(program, gl.Str("projection\x00"))
	gl.GLUniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GLGetUniformLocation(program, gl.Str("camera\x00"))
	gl.GLUniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GLGetUniformLocation(program, gl.Str("model\x00"))
	gl.GLUniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GLGetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.GLBindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	cube.texture = cube.newTexture(sett.App.CurrentPath + "/../Resources/resources/textures/square.png")

	// Configure the vertex data
	vao := gl.GenVertexArrays(1)[0]

	gl.BindVertexArray(vao)

	vbo := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vbo)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), oglconsts.STATIC_DRAW)

	vertAttrib := uint32(gl.GLGetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, oglconsts.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GLGetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, oglconsts.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	cube.angle = 0.0
	cube.previousTime = sdl.GetTicks()

	return cube
}

// Render ...
func (cube *Cube) Render() {
	gl := cube.window.OpenGL()

	gl.Enable(oglconsts.DEPTH_TEST)
	gl.DepthFunc(oglconsts.LESS)
	gl.Clear(oglconsts.COLOR_BUFFER_BIT | oglconsts.DEPTH_BUFFER_BIT)

	// Update
	time := sdl.GetTicks()
	elapsed := time - cube.previousTime
	cube.previousTime = time

	cube.angle += float32(elapsed)
	cube.model = mgl32.HomogRotate3D(float32(cube.angle), mgl32.Vec3{0, 1, 0})

	// Render
	gl.UseProgram(cube.program)
	gl.GLUniformMatrix4fv(cube.modelUniform, 1, false, &cube.model[0])

	gl.BindVertexArray(cube.vao)

	gl.ActiveTexture(oglconsts.TEXTURE0)
	gl.BindTexture(oglconsts.TEXTURE_2D, cube.texture)

	gl.DrawArrays(oglconsts.TRIANGLES, 0, 6*2*3)
}

// Dispose will cleanup everything
func (cube *Cube) Dispose() {
	gl := cube.window.OpenGL()

	gl.DeleteTextures([]uint32{cube.texture})
	gl.DeleteVertexArrays([]uint32{cube.vao})
	gl.DeleteProgram(cube.program)
}

func (cube *Cube) newProgram(vertexShaderSource, fragmentShaderSource string) uint32 {
	gl := cube.window.OpenGL()

	vertexShader := cube.compileShader(vertexShaderSource, oglconsts.VERTEX_SHADER)
	fragmentShader := cube.compileShader(fragmentShaderSource, oglconsts.FRAGMENT_SHADER)

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, oglconsts.LINK_STATUS, &status)
	if status == oglconsts.FALSE {
		var logLength int32
		gl.GetProgramiv(program, oglconsts.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GLGetProgramInfoLog(program, logLength, nil, gl.Str(log))

		settings.LogError("[Cube] Failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program
}

func (cube *Cube) compileShader(source string, shaderType uint32) uint32 {
	gl := cube.window.OpenGL()

	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.GLShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GLGetShaderiv(shader, oglconsts.COMPILE_STATUS, &status)
	if status == oglconsts.FALSE {
		var logLength int32
		gl.GLGetShaderiv(shader, oglconsts.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GLGetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		settings.LogError("[Cube] Failed to compile shader %v : %v", source, log)
	}

	return shader
}

func (cube *Cube) newTexture(file string) uint32 {
	gl := cube.window.OpenGL()

	imgFile, err := os.Open(file)
	if err != nil {
		settings.LogError("[Cube] Texture file not found: %v", err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		settings.LogError("[Cube] Can't decode texture: %v", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		settings.LogError("[Cube] Texture unsupported stride!")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	cube.texture = gl.GenTextures(1)[0]
	gl.ActiveTexture(oglconsts.TEXTURE0)
	gl.BindTexture(oglconsts.TEXTURE_2D, texture)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MIN_FILTER, oglconsts.LINEAR)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MAG_FILTER, oglconsts.LINEAR)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_WRAP_S, oglconsts.CLAMP_TO_EDGE)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_WRAP_T, oglconsts.CLAMP_TO_EDGE)
	gl.TexImage2D(
		oglconsts.TEXTURE_2D,
		0,
		oglconsts.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		oglconsts.RGBA,
		oglconsts.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}
