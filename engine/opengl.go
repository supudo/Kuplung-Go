package engine

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/settings"
)

// OpenGL wraps the native GL API into a common interface.
type OpenGL struct {
	reportedErrors []string
}

// NewOpenGL initializes the GL bindings and returns an OpenGL instance.
func NewOpenGL() *OpenGL {
	opengl := &OpenGL{}

	if err := gl.Init(); err != nil {
		settings.LogError("[NewOpenGL] Error in intializer : %v", err)
	}

	return opengl
}

// CheckForOpenGLErrors implements the interfaces.OpenGL interface.
func (native *OpenGL) CheckForOpenGLErrors(message string) {
	sett := settings.GetSettings()
	if sett.Rendering.ShowGLErrors {
		err := gl.GetError()
		if err != oglconsts.NO_ERROR {
			errMessage := fmt.Sprintf("[GLError] [%v] glError = %X!", message, err)
			for _, n := range native.reportedErrors {
				if errMessage == n {
					return
				}
			}
			settings.LogWarn(errMessage)
			native.reportedErrors = append(native.reportedErrors, errMessage)
		}
	}
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

// BindFramebuffer implements the interfaces.OpenGL interface.
func (native *OpenGL) BindFramebuffer(target, buffer uint32) {
	gl.BindFramebuffer(target, buffer)
}

// BindRenderbuffer implements the interfaces.OpenGL interface.
func (native *OpenGL) BindRenderbuffer(target uint32, renderbuffer uint32) {
	gl.BindRenderbuffer(target, renderbuffer)
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

// BlitFramebuffer implements the interfaces.OpenGL interface.
func (native *OpenGL) BlitFramebuffer(srcX0 int32, srcY0 int32, srcX1 int32, srcY1 int32, dstX0 int32, dstY0 int32, dstX1 int32, dstY1 int32, mask uint32, filter uint32) {
	gl.BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1, mask, filter)
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

// DrawElementsOffset implements the interfaces.OpenGL interface.
func (native *OpenGL) DrawElementsOffset(mode uint32, count int32, elementType uint32, offset int) {
	gl.DrawElements(mode, count, elementType, gl.PtrOffset(offset))
}

// DrawBuffers implements the interfaces.OpenGL interface.
func (native *OpenGL) DrawBuffers(n int32, bufs *uint32) {
	gl.DrawBuffers(n, bufs)
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

// GenFramebuffers implements the interfaces.OpenGL interface.
func (native *OpenGL) GenFramebuffers(n int32) []uint32 {
	ids := make([]uint32, n)
	gl.GenFramebuffers(n, &ids[0])
	return ids
}

// GenQueries implements the interfaces.OpenGL interface.
func (native *OpenGL) GenQueries(n int32) []uint32 {
	ids := make([]uint32, n)
	gl.GenQueries(n, &ids[0])
	return ids
}

// GenRenderbuffers implements the interfaces.OpenGL interface.
func (native *OpenGL) GenRenderbuffers(n int32) []uint32 {
	ids := make([]uint32, n)
	gl.GenRenderbuffers(n, &ids[0])
	return ids
}

// BeginQuery delimit the boundaries of a query object.
func (native *OpenGL) BeginQuery(target uint32, id uint32) {
	gl.BeginQuery(target, id)
}

// EndQuery implements the interfaces.OpenGL interface.
func (native *OpenGL) EndQuery(target uint32) {
	gl.EndQuery(target)
}

// IsQuery implements the interfaces.OpenGL interface.
func (native *OpenGL) IsQuery(id uint32) bool {
	return gl.IsQuery(id)
}

// ColorMask implements the interfaces.OpenGL interface.
func (native *OpenGL) ColorMask(red bool, green bool, blue bool, alpha bool) {
	gl.ColorMask(red, green, blue, alpha)
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

// TexParameterf implements the interfaces.OpenGL interface.
func (native *OpenGL) TexParameterf(target uint32, pname uint32, param float32) {
	gl.TexParameterf(target, pname, param)
}

// FramebufferTexture2D implements the interfaces.OpenGL interface.
func (native *OpenGL) FramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	gl.FramebufferTexture2D(target, attachment, textarget, texture, level)
}

// Uniform1i implements the interfaces.OpenGL interface.
func (native *OpenGL) Uniform1i(location int32, value int32) {
	gl.Uniform1i(location, value)
}

// Uniform2i implements the interfaces.OpenGL interface.
func (native *OpenGL) Uniform2i(location int32, v0 int32, v1 int32) {
	gl.Uniform2i(location, v0, v1)
}

// Uniform3i implements the interfaces.OpenGL interface.
func (native *OpenGL) Uniform3i(location int32, v0 int32, v1 int32, v2 int32) {
	gl.Uniform3i(location, v0, v1, v2)
}

// Uniform1f implements the interfaces.OpenGL interface.
func (native *OpenGL) Uniform1f(location int32, value1 float32) {
	gl.Uniform1f(location, value1)
}

// Uniform3f implements the interfaces.OpenGL interface.
func (native *OpenGL) Uniform3f(location int32, value1 float32, value2 float32, value3 float32) {
	gl.Uniform3f(location, value1, value2, value3)
}

// Uniform3fv implements the interfaces.OpenGL interface.
func (native *OpenGL) Uniform3fv(location int32, count int32, value *float32) {
	gl.Uniform3fv(location, count, value)
}

// Uniform4f implements the interfaces.OpenGL interface.
func (native *OpenGL) Uniform4f(location int32, v0 float32, v1 float32, v2 float32, v3 float32) {
	gl.Uniform4f(location, v0, v1, v2, v3)
}

// Uniform4fv implements the interfaces.OpenGL interface.
func (native *OpenGL) Uniform4fv(location int32, value *[4]float32) {
	gl.Uniform4fv(location, 1, &value[0])
}

// UniformMatrix4fv implements the interfaces.OpenGL interface.
func (native *OpenGL) UniformMatrix4fv(location int32, transpose bool, value *[16]float32) {
	count := int32(1)
	gl.UniformMatrix4fv(location, count, transpose, &value[0])
}

// UniformMatrix3fv implements the interfaces.OpenGL interface.
func (native *OpenGL) UniformMatrix3fv(location int32, count int32, transpose bool, value *float32) {
	gl.UniformMatrix3fv(location, count, transpose, value)
}

// GLUniformMatrix4fv implements the interfaces.OpenGL interface.
func (native *OpenGL) GLUniformMatrix4fv(location int32, count int32, transpose bool, value *float32) {
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

// DepthFunc specify the value used for depth buffer comparisons
func (native *OpenGL) DepthFunc(xfunc uint32) {
	gl.DepthFunc(xfunc)
}

// Str takes a null-terminated Go string and returns its GL-compatible address. This function reaches
// into Go string storage in an unsafe way so the caller must ensure the string is not garbage collected.
func (native *OpenGL) Str(str string) *uint8 {
	return gl.Str(str)
}

// Strs takes a list of Go strings (with or without null-termination) and returns their C counterpart.
func (native *OpenGL) Strs(strs ...string) (cstrs **uint8, free func()) {
	return gl.Strs(strs[0])
}

// PtrOffset takes a pointer offset and returns a GL-compatible pointer. Useful for functions such as
// glVertexAttribPointer that take pointer parameters indicating an offset rather than an absolute memory address.
func (native *OpenGL) PtrOffset(offset int) unsafe.Pointer {
	return gl.PtrOffset(offset)
}

// VertexAttribPointer define an array of generic vertex attribute data
func (native *OpenGL) VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
	gl.VertexAttribPointer(index, size, xtype, normalized, stride, pointer)
}

// Ptr takes a slice or pointer (to a singular scalar value or the first element of an array or slice)
// and returns its GL-compatible address.
func (native *OpenGL) Ptr(data interface{}) unsafe.Pointer {
	return gl.Ptr(data)
}

// GLGetUniformLocation Returns the location of a uniform variable
func (native *OpenGL) GLGetUniformLocation(program uint32, name *uint8) int32 {
	return gl.GetUniformLocation(program, name)
}

// GLGetAttribLocation Returns the location of an attribute variable
func (native *OpenGL) GLGetAttribLocation(program uint32, name *uint8) int32 {
	return gl.GetAttribLocation(program, name)
}

// GLBindFragDataLocation bind a user-defined varying out variable to a fragment shader color number
func (native *OpenGL) GLBindFragDataLocation(program uint32, color uint32, name *uint8) {
	gl.BindFragDataLocation(program, color, name)
}

// GLGetShaderiv Returns a parameter from a shader object
func (native *OpenGL) GLGetShaderiv(shader uint32, pname uint32, params *int32) {
	gl.GetShaderiv(shader, pname, params)
}

// GLShaderSource returns the source code string from a shader object
func (native *OpenGL) GLShaderSource(shader uint32, count int32, xstring **uint8, length *int32) {
	gl.ShaderSource(shader, count, xstring, length)
}

// GLGetShaderInfoLog Returns the information log for a shader object
func (native *OpenGL) GLGetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8) {
	gl.GetShaderInfoLog(shader, bufSize, length, infoLog)
}

// GetProgramiv Returns a parameter from a program object
func (native *OpenGL) GetProgramiv(program uint32, pname uint32, params *int32) {
	gl.GetProgramiv(program, pname, params)
}

// GLGetProgramInfoLog Returns the information log for a program object
func (native *OpenGL) GLGetProgramInfoLog(program uint32, bufSize int32, length *int32, infoLog *uint8) {
	gl.GetProgramInfoLog(program, bufSize, length, infoLog)
}

// LogOpenGLReturn will return the OpenGL error
func (native *OpenGL) LogOpenGLReturn() string {
	if err := gl.GetError(); err != oglconsts.NO_ERROR {
		return fmt.Sprintf("[OpenGL Erorr] %v", err)
	}
	return ""
}

// LogOpenGLError will output the OpenGL error and exit the application
func (native *OpenGL) LogOpenGLError() {
	if err := gl.GetError(); err != oglconsts.NO_ERROR {
		settings.LogError("[OpenGL Error] Error occured: %v", err)
	}
}

// LogOpenGLWarn will output the OpenGL error
func (native *OpenGL) LogOpenGLWarn() {
	if err := gl.GetError(); err != oglconsts.NO_ERROR {
		settings.LogWarn("[OpenGL Error] Error occured: %v", err)
	}
}

// LineWidth will set the line width
func (native *OpenGL) LineWidth(lw float32) {
	gl.LineWidth(lw)
}

// DepthMask will set the depth mask
func (native *OpenGL) DepthMask(mask bool) {
	gl.DepthMask(mask)
}

// BeginConditionalRender will set the depth mask
func (native *OpenGL) BeginConditionalRender(id uint32, mode uint32) {
	gl.BeginConditionalRender(id, mode)
}

// EndConditionalRender will set the depth mask
func (native *OpenGL) EndConditionalRender() {
	gl.EndConditionalRender()
}

// PatchParameteri specifies the parameters for patch primitives
func (native *OpenGL) PatchParameteri(pname uint32, value int32) {
	gl.PatchParameteri(pname, value)
}

// GetQueryObjectui64v specifies the parameters for patch primitives
func (native *OpenGL) GetQueryObjectui64v(id uint32, pname uint32, params *uint64) {
	gl.GetQueryObjectui64v(id, pname, params)
}

// RenderbufferStorage implements the interfaces.OpenGL interface.
func (native *OpenGL) RenderbufferStorage(target uint32, internalformat uint32, width int32, height int32) {
	gl.RenderbufferStorage(target, internalformat, width, height)
}

// FramebufferRenderbuffer implements the interfaces.OpenGL interface.
func (native *OpenGL) FramebufferRenderbuffer(target uint32, attachment uint32, renderbuffertarget uint32, renderbuffer uint32) {
	gl.FramebufferRenderbuffer(target, attachment, renderbuffertarget, renderbuffer)
}

// CheckFramebufferStatus implements the interfaces.OpenGL interface.
func (native *OpenGL) CheckFramebufferStatus(target uint32) uint32 {
	return gl.CheckFramebufferStatus(target)
}
