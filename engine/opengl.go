package engine

import (
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/supudo/Kuplung-Go/settings"
)

// OpenGL wraps the native GL API into a common interface.
type OpenGL struct {
}

// NewOpenGL initializes the GL bindings and returns an OpenGL instance.
func NewOpenGL() *OpenGL {
	opengl := &OpenGL{}

	if err := gl.Init(); err != nil {
		settings.LogError("[NewOpenGL] Error in intializer : %v", err)
	}

	return opengl
}

// ActiveTexture implements the interfaces.OpenGL interface.
func (native *OpenGL) ActiveTexture(texture uint32) {
	gl.ActiveTexture(texture)
}

// AttachShader implements the interfaces.OpenGL interface.
func (native *OpenGL) AttachShader(program uint32, shader uint32) {
	gl.AttachShader(program, shader)
}

// BindAttribLocation implements the interfaces.OpenGL interface.
func (native *OpenGL) BindAttribLocation(program uint32, index uint32, name string) {
	gl.BindAttribLocation(program, index, gl.Str(name+"\x00"))
}

// BindBuffer implements the interfaces.OpenGL interface.
func (native *OpenGL) BindBuffer(target uint32, buffer uint32) {
	gl.BindBuffer(target, buffer)
}

// BindSampler implements the interfaces.OpenGL interface.
func (native *OpenGL) BindSampler(unit uint32, sampler uint32) {
	gl.BindSampler(unit, sampler)
}

// BindTexture implements the interfaces.OpenGL interface.
func (native *OpenGL) BindTexture(target uint32, texture uint32) {
	gl.BindTexture(target, texture)
}

// BindVertexArray implements the interfaces.OpenGL interface.
func (native *OpenGL) BindVertexArray(array uint32) {
	gl.BindVertexArray(array)
}

// BlendEquation implements the interfaces.OpenGL interface.
func (native *OpenGL) BlendEquation(mode uint32) {
	gl.BlendEquation(mode)
}

// BlendEquationSeparate implements the interfaces.OpenGL interface.
func (native *OpenGL) BlendEquationSeparate(modeRGB uint32, modeAlpha uint32) {
	gl.BlendEquationSeparate(modeRGB, modeAlpha)
}

// BlendFunc implements the interfaces.OpenGL interface.
func (native *OpenGL) BlendFunc(sfactor uint32, dfactor uint32) {
	gl.BlendFunc(sfactor, dfactor)
}

// BlendFuncSeparate implements the interfaces.OpenGL interface.
func (native *OpenGL) BlendFuncSeparate(srcRGB uint32, dstRGB uint32, srcAlpha uint32, dstAlpha uint32) {
	gl.BlendFuncSeparate(srcRGB, dstRGB, srcAlpha, dstAlpha)
}

// BufferData implements the interfaces.OpenGL interface.
func (native *OpenGL) BufferData(target uint32, size int, data interface{}, usage uint32) {
	dataPtr, isPtr := data.(unsafe.Pointer)
	if isPtr {
		gl.BufferData(target, size, dataPtr, usage)
	} else {
		gl.BufferData(target, size, gl.Ptr(data), usage)
	}
}

// Clear implements the interfaces.OpenGL interface.
func (native *OpenGL) Clear(mask uint32) {
	gl.Clear(mask)
}

// ClearColor implements the interfaces.OpenGL interface.
func (native *OpenGL) ClearColor(red float32, green float32, blue float32, alpha float32) {
	gl.ClearColor(red, green, blue, alpha)
}

// CompileShader implements the interfaces.OpenGL interface.
func (native *OpenGL) CompileShader(shader uint32) {
	gl.CompileShader(shader)
}

// CreateProgram implements the interfaces.OpenGL interface.
func (native *OpenGL) CreateProgram() uint32 {
	return gl.CreateProgram()
}

// CreateShader implements the interfaces.OpenGL interface.
func (native *OpenGL) CreateShader(shaderType uint32) uint32 {
	return gl.CreateShader(shaderType)
}

// DeleteBuffers implements the interfaces.OpenGL interface.
func (native *OpenGL) DeleteBuffers(buffers []uint32) {
	gl.DeleteBuffers(int32(len(buffers)), &buffers[0])
}

// DeleteProgram implements the interfaces.OpenGL interface.
func (native *OpenGL) DeleteProgram(program uint32) {
	gl.DeleteProgram(program)
}

// DeleteShader implements the interfaces.OpenGL interface.
func (native *OpenGL) DeleteShader(shader uint32) {
	gl.DeleteShader(shader)
}

// DeleteTextures implements the interfaces.OpenGL interface.
func (native *OpenGL) DeleteTextures(textures []uint32) {
	gl.DeleteTextures(int32(len(textures)), &textures[0])
}

// DeleteVertexArrays implements the interfaces.OpenGL interface.
func (native *OpenGL) DeleteVertexArrays(arrays []uint32) {
	gl.DeleteVertexArrays(int32(len(arrays)), &arrays[0])
}

// Disable implements the interfaces.OpenGL interface.
func (native *OpenGL) Disable(capability uint32) {
	gl.Disable(capability)
}

// DrawArrays implements the interfaces.OpenGL interface.
func (native *OpenGL) DrawArrays(mode uint32, first int32, count int32) {
	gl.DrawArrays(mode, first, count)
}

// DrawElements implements the interfaces.OpenGL interface.
func (native *OpenGL) DrawElements(mode uint32, count int32, elementType uint32, indices uintptr) {
	gl.DrawElements(mode, count, elementType, unsafe.Pointer(indices)) // nolint: govet,gas
}

// Enable implements the interfaces.OpenGL interface.
func (native *OpenGL) Enable(capability uint32) {
	gl.Enable(capability)
}

// EnableVertexAttribArray implements the interfaces.OpenGL interface.
func (native *OpenGL) EnableVertexAttribArray(index uint32) {
	gl.EnableVertexAttribArray(index)
}

// GenerateMipmap implements the interfaces.OpenGL interface.
func (native *OpenGL) GenerateMipmap(target uint32) {
	gl.GenerateMipmap(target)
}

// GenBuffers implements the interfaces.OpenGL interface.
func (native *OpenGL) GenBuffers(n int32) []uint32 {
	buffers := make([]uint32, n)
	gl.GenBuffers(n, &buffers[0])
	return buffers
}

// GenTextures implements the interfaces.OpenGL interface.
func (native *OpenGL) GenTextures(n int32) []uint32 {
	ids := make([]uint32, n)
	gl.GenTextures(n, &ids[0])
	return ids
}

// GenVertexArrays implements the interfaces.OpenGL interface.
func (native *OpenGL) GenVertexArrays(n int32) []uint32 {
	ids := make([]uint32, n)
	gl.GenVertexArrays(n, &ids[0])
	return ids
}

// GetAttribLocation implements the interfaces.OpenGL interface.
func (native *OpenGL) GetAttribLocation(program uint32, name string) int32 {
	return gl.GetAttribLocation(program, gl.Str(name+"\x00"))
}

// GetError implements the interfaces.OpenGL interface.
func (native *OpenGL) GetError() uint32 {
	return gl.GetError()
}

// GetIntegerv implements the interfaces.OpenGL interface.
func (native *OpenGL) GetIntegerv(name uint32, data *int32) {
	gl.GetIntegerv(name, data)
}

// GetProgramInfoLog implements the interfaces.OpenGL interface.
func (native *OpenGL) GetProgramInfoLog(program uint32) string {
	logLength := native.GetProgramParameter(program, gl.INFO_LOG_LENGTH)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
	return log
}

// GetProgramParameter implements the interfaces.OpenGL interface.
func (native *OpenGL) GetProgramParameter(program uint32, param uint32) int32 {
	result := int32(0)
	gl.GetProgramiv(program, param, &result)
	return result
}

// GetShaderInfoLog implements the interfaces.OpenGL interface.
func (native *OpenGL) GetShaderInfoLog(shader uint32) string {
	logLength := native.GetShaderParameter(shader, gl.INFO_LOG_LENGTH)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
	return log
}

// GetShaderParameter implements the interfaces.OpenGL interface.
func (native *OpenGL) GetShaderParameter(shader uint32, param uint32) int32 {
	result := int32(0)
	gl.GetShaderiv(shader, param, &result)
	return result
}

// GetUniformLocation implements the interfaces.OpenGL interface.
func (native *OpenGL) GetUniformLocation(program uint32, name string) int32 {
	return gl.GetUniformLocation(program, gl.Str(name+"\x00"))
}

// IsEnabled implements the OpenGL interface.
func (native *OpenGL) IsEnabled(capability uint32) bool {
	return gl.IsEnabled(capability)
}

// LinkProgram implements the interfaces.OpenGL interface.
func (native *OpenGL) LinkProgram(program uint32) {
	gl.LinkProgram(program)
}

// PixelStorei implements the OpenGL interface.
func (native *OpenGL) PixelStorei(name uint32, param int32) {
	gl.PixelStorei(name, param)
}

// PolygonMode implements the interfaces.OpenGL interface.
func (native *OpenGL) PolygonMode(face uint32, mode uint32) {
	gl.PolygonMode(face, mode)
}

// ReadPixels implements the interfaces.OpenGL interface.
func (native *OpenGL) ReadPixels(x int32, y int32, width int32, height int32, format uint32, pixelType uint32, pixels interface{}) {
	gl.ReadPixels(x, y, width, height, format, pixelType, gl.Ptr(pixels))
}

// Scissor implements the interfaces.OpenGL interface.
func (native *OpenGL) Scissor(x, y int32, width, height int32) {
	gl.Scissor(x, y, width, height)
}

// ShaderSource implements the interfaces.OpenGL interface.
func (native *OpenGL) ShaderSource(shader uint32, source string) {
	csources, free := gl.Strs(source + "\x00")
	defer free()

	gl.ShaderSource(shader, 1, csources, nil)
}

// TexImage2D implements the interfaces.OpenGL interface.
func (native *OpenGL) TexImage2D(target uint32, level int32, internalFormat uint32, width int32, height int32,
	border int32, format uint32, xtype uint32, pixels interface{}) {
	ptr, isPointer := pixels.(unsafe.Pointer)
	if isPointer {
		gl.TexImage2D(target, level, int32(internalFormat), width, height, border, format, xtype, ptr)
	} else {
		gl.TexImage2D(target, level, int32(internalFormat), width, height, border, format, xtype, gl.Ptr(pixels))
	}
}

// TexParameteri implements the interfaces.OpenGL interface.
func (native *OpenGL) TexParameteri(target uint32, pname uint32, param int32) {
	gl.TexParameteri(target, pname, param)
}

// Uniform1i implements the interfaces.OpenGL interface.
func (native *OpenGL) Uniform1i(location int32, value int32) {
	gl.Uniform1i(location, value)
}

// Uniform4fv implements the interfaces.OpenGL interface.
func (native *OpenGL) Uniform4fv(location int32, value *[4]float32) {
	gl.Uniform4fv(location, 1, &value[0])
}

// UniformMatrix4fvImgui implements the interfaces.OpenGL interface.
func (native *OpenGL) UniformMatrix4fvImgui(location int32, transpose bool, value *[16]float32) {
	count := int32(1)
	gl.UniformMatrix4fv(location, count, transpose, &value[0])
}

// UniformMatrix4fv implements the interfaces.OpenGL interface.
func (native *OpenGL) UniformMatrix4fv(location int32, transpose bool, value *float32) {
	count := int32(1)
	gl.UniformMatrix4fv(location, count, transpose, value)
}

// UseProgram implements the interfaces.OpenGL interface.
func (native *OpenGL) UseProgram(program uint32) {
	gl.UseProgram(program)
}

// VertexAttribOffset implements the interfaces.OpenGL interface.
func (native *OpenGL) VertexAttribOffset(index uint32, size int32, attribType uint32, normalized bool, stride int32, offset int) {
	gl.VertexAttribPointer(index, size, attribType, normalized, stride, gl.PtrOffset(offset))
}

// Viewport implements the interfaces.OpenGL interface.
func (native *OpenGL) Viewport(x int32, y int32, width int32, height int32) {
	gl.Viewport(x, y, width, height)
}

// GetOpenGLVersion will return the OpenGL version
func (native *OpenGL) GetOpenGLVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

// GetShadingLanguageVersion will return the shading language version
func (native *OpenGL) GetShadingLanguageVersion() string {
	return gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
}

// GetVendorName will return the shading language version
func (native *OpenGL) GetVendorName() string {
	return gl.GoStr(gl.GetString(gl.VENDOR))
}

// GetRendererName will return the shading language version
func (native *OpenGL) GetRendererName() string {
	return gl.GoStr(gl.GetString(gl.RENDERER))
}

// BindFragDataLocation binds a user-defined varying out variable to a fragment shader color number
func (native *OpenGL) BindFragDataLocation(program uint32, color uint32, name string) {
	gl.BindFragDataLocation(program, color, gl.Str(name+"\x00"))
}
