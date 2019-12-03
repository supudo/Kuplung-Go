package engine

import (
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/veandco/go-sdl2/sdl"
)

const windowWidth = 800
const windowHeight = 600

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
	cube := &Cube{}

	cube.version = "#version 410"
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
`
	fragmentShader := cube.version + `
uniform sampler2D tex;
in vec2 fragTexCoord;
out vec4 outputColor;
void main() {
	outputColor = texture(tex, fragTexCoord);
}
`
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

	cube.newProgram(vertexShader, fragmentShader)

	gl.UseProgram(cube.program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(cube.program, "projection")
	gl.UniformMatrix4fv(projectionUniform, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(cube.program, "camera")
	gl.UniformMatrix4fv(cameraUniform, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(cube.program, "model")
	gl.UniformMatrix4fv(modelUniform, false, &model[0])

	textureUniform := gl.GetUniformLocation(cube.program, "tex")
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(cube.program, 0, "outputColor")

	// Load the texture
	sett := settings.GetSettings()
	cube.newTexture(sett.App.CurrentPath + "/../Resources/resources/textures/square.png")

	// Configure the vertex data
	cube.vao = gl.GenVertexArrays(1)[0]
	gl.BindVertexArray(cube.vao)

	var vbo uint32
	vbo = gl.GenBuffers(1)[0]
	gl.BindBuffer(ARRAY_BUFFER, vbo)
	gl.BufferData(ARRAY_BUFFER, len(cubeVertices)*4, cubeVertices, STREAM_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(cube.program, "vert"))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribOffset(vertAttrib, 3, FLOAT, false, 5*4, 0)

	texCoordAttrib := uint32(gl.GetAttribLocation(cube.program, "vertTexCoord"))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribOffset(texCoordAttrib, 2, FLOAT, false, 5*4, 3*4)

	// Configure global settings

	cube.angle = 0.0
	cube.previousTime = sdl.GetTicks()

	return cube
}

// Render ...
func (cube *Cube) Render() {
	gl := cube.glWrapper

	// Update
	time := sdl.GetTicks()
	elapsed := time - cube.previousTime
	cube.previousTime = time

	cube.angle += float32(elapsed)
	cube.model = mgl32.HomogRotate3D(float32(cube.angle), mgl32.Vec3{0, 1, 0})

	// Render
	gl.UseProgram(cube.program)
	gl.UniformMatrix4fv(cube.modelUniform, false, &cube.model[0])

	gl.BindVertexArray(cube.vao)

	gl.ActiveTexture(TEXTURE0)
	gl.BindTexture(TEXTURE_2D, cube.texture)

	gl.DrawArrays(TRIANGLES, 0, 6*2*3)
}

func (cube *Cube) newProgram(vertexShaderSource, fragmentShaderSource string) {
	gl := cube.glWrapper

	p, err := LinkNewStandardProgram(gl, vertexShaderSource, fragmentShaderSource)
	if err != nil {
		settings.LogError("[Cube] - Can't initiate a new program: %v", err)
	}
	cube.program = p
}

func (cube *Cube) newTexture(file string) {
	gl := cube.glWrapper

	imgFile, err := os.Open(file)
	if err != nil {
		settings.LogError("[newTexture] Texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		settings.LogError("[newTexture] Texture %q can't be read: %v", file, err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		settings.LogError("[newTexture] Texture %q has unsupported stride.", file)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	texture = gl.GenTextures(1)[0]
	gl.ActiveTexture(TEXTURE0)
	gl.BindTexture(TEXTURE_2D, texture)
	gl.TexParameteri(TEXTURE_2D, TEXTURE_MIN_FILTER, LINEAR)
	gl.TexParameteri(TEXTURE_2D, TEXTURE_MAG_FILTER, LINEAR)
	gl.TexParameteri(TEXTURE_2D, TEXTURE_WRAP_S, CLAMP_TO_EDGE)
	gl.TexParameteri(TEXTURE_2D, TEXTURE_WRAP_T, CLAMP_TO_EDGE)
	gl.TexImage2D(TEXTURE_2D, 0, RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, RGBA, UNSIGNED_BYTE, rgba.Pix)
}
