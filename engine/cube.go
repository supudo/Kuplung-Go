package engine

import (
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/veandco/go-sdl2/sdl"
)

const windowWidth = 800
const windowHeight = 600

// Cube ...
type Cube struct {
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
func CubeInit() *Cube {
	cube := &Cube{}

	cube.version = "#version 410"

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
	projectionUniform := gl.GetUniformLocation(cube.program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(cube.program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(cube.program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(cube.program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(cube.program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	sett := settings.GetSettings()
	cube.newTexture(sett.App.CurrentPath + "/../Resources/resources/textures/square.png")

	// Configure the vertex data
	gl.GenVertexArrays(1, &cube.vao)
	gl.BindVertexArray(cube.vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(cube.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(cube.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	// Configure global settings

	cube.angle = 0.0
	cube.previousTime = sdl.GetTicks()

	return cube
}

// CubeRender ...
func (cube *Cube) CubeRender() {
	// Update
	time := sdl.GetTicks()
	elapsed := time - cube.previousTime
	cube.previousTime = time

	cube.angle += float32(elapsed)
	cube.model = mgl32.HomogRotate3D(float32(cube.angle), mgl32.Vec3{0, 1, 0})

	// Render
	gl.UseProgram(cube.program)
	gl.UniformMatrix4fv(cube.modelUniform, 1, false, &cube.model[0])

	gl.BindVertexArray(cube.vao)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, cube.texture)

	gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
}

func (cube *Cube) newProgram(vertexShaderSource, fragmentShaderSource string) {
	cube.program = gl.CreateProgram()
	vertHandle := gl.CreateShader(gl.VERTEX_SHADER)
	fragHandle := gl.CreateShader(gl.FRAGMENT_SHADER)

	glShaderSource := func(handle uint32, source string) {
		csource, free := gl.Strs(source + "\x00")
		defer free()

		gl.ShaderSource(handle, 1, csource, nil)
	}

	glShaderSource(vertHandle, vertexShaderSource)
	glShaderSource(fragHandle, fragmentShaderSource)
	gl.CompileShader(vertHandle)
	gl.CompileShader(fragHandle)
	gl.AttachShader(cube.program, vertHandle)
	gl.AttachShader(cube.program, fragHandle)
	gl.LinkProgram(cube.program)
}

func (cube *Cube) newTexture(file string) {
	imgFile, err := os.Open(file)
	if err != nil {
		settings.LogError("[CubeInit] Texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		settings.LogError("[CubeInit] Texture %q can't be read: %v", file, err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		settings.LogError("[CubeInit] Texture %q has unsupported stride.", file)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
}
