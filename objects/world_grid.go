package objects

import (
	"image"
	"image/draw"
	"os"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/veandco/go-sdl2/sdl"
)

// WorldGrid ...
type WorldGrid struct {
	window interfaces.Window

	angle        float32
	previousTime float32
	program      uint32
	texture      uint32

	modelUniform      int32
	model             mgl32.Mat4
	projectionUniform int32
	fov               float32

	vao uint32

	version string
}

// InitWorldGrid ...
func InitWorldGrid(window interfaces.Window) *WorldGrid {
	sett := settings.GetSettings()
	rsett := settings.GetRenderingSettings()

	wgrid := &WorldGrid{}

	wgrid.version = "#version 410"
	wgrid.window = window

	vertexShader := wgrid.version + `
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
	fragmentShader := wgrid.version + `
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

	wgrid.fov = rsett.Fov

	gl := window.OpenGL()

	// Configure the vertex and fragment shaders
	wgrid.program = wgrid.newProgram(vertexShader, fragmentShader)

	gl.UseProgram(wgrid.program)

	projection := mgl32.Perspective(mgl32.DegToRad(wgrid.fov), rsett.RatioWidth/rsett.RatioHeight, rsett.PlaneClose, rsett.PlaneFar)
	wgrid.projectionUniform = gl.GLGetUniformLocation(wgrid.program, gl.Str("projection\x00"))
	gl.GLUniformMatrix4fv(wgrid.projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GLGetUniformLocation(wgrid.program, gl.Str("camera\x00"))
	gl.GLUniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	wgrid.modelUniform = gl.GLGetUniformLocation(wgrid.program, gl.Str("model\x00"))
	gl.GLUniformMatrix4fv(wgrid.modelUniform, 1, false, &model[0])

	textureUniform := gl.GLGetUniformLocation(wgrid.program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.GLBindFragDataLocation(wgrid.program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	wgrid.texture = wgrid.newTexture(sett.App.CurrentPath + "/../Resources/resources/textures/square.png")

	// Configure the vertex data
	wgrid.vao = gl.GenVertexArrays(1)[0]

	gl.BindVertexArray(wgrid.vao)

	vbo := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vbo)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), oglconsts.STATIC_DRAW)

	vertAttrib := uint32(gl.GLGetAttribLocation(wgrid.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, oglconsts.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GLGetAttribLocation(wgrid.program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, oglconsts.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	wgrid.angle = 0.0
	wgrid.previousTime = float32(sdl.GetTicks())

	return wgrid
}

// Render ...
func (grid *WorldGrid) Render() {
	gl := grid.window.OpenGL()
	rsett := settings.GetRenderingSettings()

	gl.Enable(oglconsts.DEPTH_TEST)
	gl.DepthFunc(oglconsts.LESS)
	gl.Disable(oglconsts.BLEND)
	gl.BlendFunc(oglconsts.SRC_ALPHA, oglconsts.ONE_MINUS_SRC_ALPHA)
	gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.FILL)

	// Update
	sdlTime := float32(sdl.GetTicks())
	elapsed := (sdlTime - grid.previousTime) / 1000
	grid.previousTime = sdlTime

	grid.angle += float32(elapsed)
	grid.model = mgl32.HomogRotate3D(float32(grid.angle), mgl32.Vec3{0, 1, 0})

	if grid.fov != rsett.Fov {
		projection := mgl32.Perspective(mgl32.DegToRad(grid.fov), rsett.RatioWidth/rsett.RatioHeight, rsett.PlaneClose, rsett.PlaneFar)
		gl.GLUniformMatrix4fv(grid.projectionUniform, 1, false, &projection[0])
		grid.fov = rsett.Fov
	}

	w, h := grid.window.Size()
	gl.Viewport(0, 0, int32(w), int32(h))

	// Render
	gl.UseProgram(grid.program)
	gl.GLUniformMatrix4fv(grid.modelUniform, 1, false, &grid.model[0])

	gl.BindVertexArray(grid.vao)

	gl.ActiveTexture(oglconsts.TEXTURE0)
	gl.BindTexture(oglconsts.TEXTURE_2D, grid.texture)

	gl.DrawArrays(oglconsts.TRIANGLES, 0, 6*2*3)
}

// Dispose will cleanup everything
func (grid *WorldGrid) Dispose() {
	gl := grid.window.OpenGL()

	gl.DeleteTextures([]uint32{grid.texture})
	gl.DeleteVertexArrays([]uint32{grid.vao})
	gl.DeleteProgram(grid.program)
}

func (grid *WorldGrid) newProgram(vertexShaderSource, fragmentShaderSource string) uint32 {
	gl := grid.window.OpenGL()

	vertexShader := grid.compileShader(vertexShaderSource, oglconsts.VERTEX_SHADER)
	fragmentShader := grid.compileShader(fragmentShaderSource, oglconsts.FRAGMENT_SHADER)

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

func (grid *WorldGrid) compileShader(source string, shaderType uint32) uint32 {
	gl := grid.window.OpenGL()

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

func (grid *WorldGrid) newTexture(file string) uint32 {
	gl := grid.window.OpenGL()

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
	grid.texture = gl.GenTextures(1)[0]
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
