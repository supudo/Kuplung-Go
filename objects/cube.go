package objects

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine/constants"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/veandco/go-sdl2/sdl"
)

// Cube ...
type Cube struct {
	glWrapper interfaces.OpenGL

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
func CubeInit(gl interfaces.OpenGL) *Cube {
	sett := settings.GetSettings()

	cube := &Cube{}

	cube.version = "#version 330"
	cube.glWrapper = gl

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

	// Configure the vertex and fragment shaders
	program, err := cube.newProgram(vertexShader, fragmentShader)
	if err != nil {
		settings.LogError("[Cube] New program error: %v", err)
	}

	gl.UseProgram(program)

	w, h := sett.AppWindow.SDLWindowWidth, sett.AppWindow.SDLWindowHeight

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
	cube.texture, err = cube.newTexture(sett.App.CurrentPath + "/../Resources/resources/textures/square.png")
	if err != nil {
		settings.LogError("[Cube] can't load texture: %v", err)
	}

	// Configure the vertex data
	vao := gl.GenVertexArrays(1)[0]
	gl.BindVertexArray(vao)

	vbo := gl.GenBuffers(1)[0]
	gl.BindBuffer(constants.ARRAY_BUFFER, vbo)
	gl.BufferData(constants.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), constants.STATIC_DRAW)

	vertAttrib := uint32(gl.GLGetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, constants.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GLGetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, constants.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	cube.angle = 0.0
	cube.previousTime = sdl.GetTicks()

	return cube
}

// Render ...
func (cube *Cube) Render() {
	gl := cube.glWrapper

	gl.Enable(constants.DEPTH_TEST)
	gl.DepthFunc(constants.LESS)
	gl.Clear(constants.COLOR_BUFFER_BIT | constants.DEPTH_BUFFER_BIT)

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

	gl.ActiveTexture(constants.TEXTURE0)
	gl.BindTexture(constants.TEXTURE_2D, cube.texture)

	gl.DrawArrays(constants.TRIANGLES, 0, 6*2*3)
}

func (cube *Cube) newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	gl := cube.glWrapper

	vertexShader, err := cube.compileShader(vertexShaderSource, constants.VERTEX_SHADER)
	if err != nil {
		settings.LogError("[Cube] Compile shader error: %v", err)
	}

	fragmentShader, err := cube.compileShader(fragmentShaderSource, constants.FRAGMENT_SHADER)
	if err != nil {
		settings.LogError("[Cube] Compile shader error: %v", err)
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, constants.LINK_STATUS, &status)
	if status == constants.FALSE {
		var logLength int32
		gl.GetProgramiv(program, constants.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GLGetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func (cube *Cube) compileShader(source string, shaderType uint32) (uint32, error) {
	gl := cube.glWrapper

	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.GLShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GLGetShaderiv(shader, constants.COMPILE_STATUS, &status)
	if status == constants.FALSE {
		var logLength int32
		gl.GLGetShaderiv(shader, constants.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GLGetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func (cube *Cube) newTexture(file string) (uint32, error) {
	gl := cube.glWrapper

	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	cube.texture = gl.GenTextures(1)[0]
	gl.ActiveTexture(constants.TEXTURE0)
	gl.BindTexture(constants.TEXTURE_2D, texture)
	gl.TexParameteri(constants.TEXTURE_2D, constants.TEXTURE_MIN_FILTER, constants.LINEAR)
	gl.TexParameteri(constants.TEXTURE_2D, constants.TEXTURE_MAG_FILTER, constants.LINEAR)
	gl.TexParameteri(constants.TEXTURE_2D, constants.TEXTURE_WRAP_S, constants.CLAMP_TO_EDGE)
	gl.TexParameteri(constants.TEXTURE_2D, constants.TEXTURE_WRAP_T, constants.CLAMP_TO_EDGE)
	gl.TexImage2D(
		constants.TEXTURE_2D,
		0,
		constants.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		constants.RGBA,
		constants.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}
