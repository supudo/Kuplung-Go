package interfaces

import "unsafe"

// OpenGL ...
type OpenGL interface {
	ActiveTexture(texture uint32)
	AttachShader(program uint32, shader uint32)

	BindAttribLocation(program uint32, index uint32, name string)
	BindBuffer(target uint32, buffer uint32)
	BindSampler(unit uint32, sampler uint32)
	BindTexture(target uint32, texture uint32)
	BindVertexArray(array uint32)
	BlendEquation(mode uint32)
	BlendEquationSeparate(modeRGB uint32, modeAlpha uint32)
	BlendFunc(sfactor uint32, dfactor uint32)
	BlendFuncSeparate(srcRGB uint32, dstRGB uint32, srcAlpha uint32, dstAlpha uint32)
	BufferData(target uint32, size int, data interface{}, usage uint32)

	Clear(mask uint32)
	ClearColor(red float32, green float32, blue float32, alpha float32)

	CompileShader(shader uint32)

	CreateProgram() uint32
	CreateShader(shaderType uint32) uint32

	DeleteBuffers(buffers []uint32)
	DeleteProgram(program uint32)
	DeleteShader(shader uint32)
	DeleteTextures(textures []uint32)
	DeleteVertexArrays(arrays []uint32)
	Disable(cap uint32)

	DrawArrays(mode uint32, first int32, count int32)
	DrawElements(mode uint32, count int32, elementType uint32, indices uintptr)

	Enable(cap uint32)
	EnableVertexAttribArray(index uint32)

	GenerateMipmap(target uint32)
	GenBuffers(n int32) []uint32
	GenTextures(n int32) []uint32
	GenVertexArrays(n int32) []uint32
	GenQueries(n int32) []uint32

	GetAttribLocation(program uint32, name string) int32
	GetError() uint32
	GetIntegerv(name uint32, data *int32)
	GetShaderInfoLog(shader uint32) string
	GetShaderParameter(shader uint32, param uint32) int32
	GetProgramInfoLog(program uint32) string
	GetProgramParameter(program uint32, param uint32) int32
	GetUniformLocation(program uint32, name string) int32

	IsEnabled(cap uint32) bool

	LinkProgram(program uint32)

	PixelStorei(name uint32, param int32)
	PolygonMode(face uint32, mode uint32)

	ReadPixels(x int32, y int32, width int32, height int32, format uint32, pixelType uint32, pixels interface{})

	Scissor(x, y int32, width, height int32)
	ShaderSource(shader uint32, source string)

	TexImage2D(target uint32, level int32, internalFormat uint32, width int32, height int32, border int32, format uint32, xtype uint32, pixels interface{})
	TexParameteri(target uint32, pname uint32, param int32)

	Uniform1i(location int32, value int32)
	Uniform1f(location int32, value1 float32)
	Uniform3f(location int32, value1 float32, value2 float32, value3 float32)
	Uniform4fv(location int32, value *[4]float32)
	UniformMatrix4fv(location int32, transpose bool, value *[16]float32)
	GLUniformMatrix4fv(location int32, count int32, transpose bool, value *float32)
	UseProgram(program uint32)

	VertexAttribOffset(index uint32, size int32, attribType uint32, normalized bool, stride int32, offset int)
	Viewport(x int32, y int32, width int32, height int32)

	GetOpenGLVersion() string
	GetShadingLanguageVersion() string
	GetVendorName() string
	GetRendererName() string

	BindFragDataLocation(program uint32, color uint32, name string)
	DepthFunc(xfunc uint32)

	Str(str string) *uint8
	Strs(strs ...string) (cstrs **uint8, free func())
	Ptr(data interface{}) unsafe.Pointer
	PtrOffset(offset int) unsafe.Pointer
	VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer)
	GLGetUniformLocation(program uint32, name *uint8) int32
	GLGetAttribLocation(program uint32, name *uint8) int32
	GLBindFragDataLocation(program uint32, color uint32, name *uint8)
	GLGetShaderiv(shader uint32, pname uint32, params *int32)
	GLShaderSource(shader uint32, count int32, xstring **uint8, length *int32)
	GLGetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8)
	GetProgramiv(program uint32, pname uint32, params *int32)
	GLGetProgramInfoLog(program uint32, bufSize int32, length *int32, infoLog *uint8)

	LogOpenGLReturn() string
	LogOpenGLError()
	LogOpenGLWarn()

	LineWidth(float32)
	DepthMask(bool)

	BeginConditionalRender(id uint32, mode uint32)
	EndConditionalRender()
}
