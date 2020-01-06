package gui

import (
	"fmt"
	"time"

	"github.com/inkyblackness/imgui-go"
	"github.com/sadlil/go-trigger"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/gui/components"
	"github.com/supudo/Kuplung-Go/gui/dialogs"
	"github.com/supudo/Kuplung-Go/gui/fonts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/rendering"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
	"github.com/veandco/go-sdl2/sdl"
)

// ContextParameters describes how to create the context.
type ContextParameters struct {
	// FontFile is the filename of a .TTF file to load instead of using default.
	FontFile string
	// FontSize is the requested size of the font. Defaults to 12.
	FontSize float32
}

// BitmapTextureQuery resolves the texture and the palette to be used for a bitmap.
type BitmapTextureQuery func(id imgui.TextureID) (palette uint32, texture uint32)

// Context describes a scope for a graphical user interface.
// It is based on ImGui.
type Context struct {
	imguiContext *imgui.Context
	window       interfaces.Window

	lastRenderTime time.Time
	isFrame        bool

	fontTexture            uint32
	shaderHandle           uint32
	attribLocationType     int32
	attribLocationTex      int32
	attribLocationPal      int32
	attribLocationProjMtx  int32
	attribLocationPosition int32
	attribLocationUV       int32
	attribLocationColor    int32
	vboHandle              uint32
	elementsHandle         uint32

	GuiVars WindowVariables

	viewControls *dialogs.ViewControls
	viewModels   *dialogs.ViewModels

	componentLog *components.ComponentLog

	fontFA imgui.Font
	fontMD imgui.Font
}

// WindowVariables holds boolean variables for all the windows
type WindowVariables struct {
	showModels   bool
	showControls bool

	showLog bool

	showDemoWindow   bool
	showAboutImGui   bool
	showAboutKuplung bool
	showMetrics      bool

	showParsing bool

	showImageSave  bool
	showRendererUI bool

	ParsingPercentage float32
}

// NewContext initializes a new UI context based on the provided OpenGL window.
func NewContext(window interfaces.Window, param ContextParameters) *Context {
	imgui.SetAssertHandler(nil)
	context := &Context{
		imguiContext: imgui.CreateContext(nil),
		window:       window,

		isFrame: false,

		viewControls: dialogs.NewViewControls(),
		viewModels:   dialogs.NewViewModels(),

		componentLog: components.NewComponentLog(),
	}

	context.GuiVars.showModels = true
	context.GuiVars.showControls = true

	context.GuiVars.showLog = false

	context.GuiVars.showDemoWindow = false
	context.GuiVars.showAboutImGui = false
	context.GuiVars.showAboutKuplung = false
	context.GuiVars.showMetrics = false

	context.GuiVars.showParsing = false

	context.GuiVars.ParsingPercentage = 0.0

	err := context.createDeviceObjects(param)
	if err != nil {
		context.Destroy()
		context = nil
		settings.LogError("[gui context] Error initialized ImGui Context: %v", err)
	}

	context.setKeyMapping()

	trigger.On(types.ActionParsingShow, func() {
		context.GuiVars.showParsing = true
	})
	trigger.On(types.ActionParsingHide, func() {
		context.GuiVars.showParsing = false
	})

	return context
}

func (context *Context) setKeyMapping() {
	keys := map[int]int{
		imgui.KeyTab:        sdl.SCANCODE_TAB,
		imgui.KeyLeftArrow:  sdl.SCANCODE_LEFT,
		imgui.KeyRightArrow: sdl.SCANCODE_RIGHT,
		imgui.KeyUpArrow:    sdl.SCANCODE_UP,
		imgui.KeyDownArrow:  sdl.SCANCODE_DOWN,
		imgui.KeyPageUp:     sdl.SCANCODE_PAGEUP,
		imgui.KeyPageDown:   sdl.SCANCODE_PAGEDOWN,
		imgui.KeyHome:       sdl.SCANCODE_HOME,
		imgui.KeyEnd:        sdl.SCANCODE_END,
		imgui.KeyInsert:     sdl.SCANCODE_INSERT,
		imgui.KeyDelete:     sdl.SCANCODE_DELETE,
		imgui.KeyBackspace:  sdl.SCANCODE_BACKSPACE,
		imgui.KeySpace:      sdl.SCANCODE_BACKSPACE,
		imgui.KeyEnter:      sdl.SCANCODE_RETURN,
		imgui.KeyEscape:     sdl.SCANCODE_ESCAPE,
		imgui.KeyA:          sdl.SCANCODE_A,
		imgui.KeyC:          sdl.SCANCODE_C,
		imgui.KeyV:          sdl.SCANCODE_V,
		imgui.KeyX:          sdl.SCANCODE_X,
		imgui.KeyY:          sdl.SCANCODE_Y,
		imgui.KeyZ:          sdl.SCANCODE_Z,
	}

	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	io := imgui.CurrentIO()
	for imguiKey, nativeKey := range keys {
		io.KeyMap(imguiKey, nativeKey)
	}
}

// Destroy cleans up the resources of the graphical user interface.
func (context *Context) Destroy() {
	context.destroyDeviceObjects(context.window.OpenGL())
	context.imguiContext.Destroy()
}

// NewFrame must be called at the start of rendering.
func (context *Context) NewFrame() {
	io := imgui.CurrentIO()

	windowWidth, windowHeight := context.window.Size()
	io.SetDisplaySize(imgui.Vec2{X: float32(windowWidth), Y: float32(windowHeight)})

	now := time.Now()
	if !context.lastRenderTime.IsZero() {
		elapsed := now.Sub(context.lastRenderTime)
		io.SetDeltaTime(float32(elapsed.Seconds()))
	}
	context.lastRenderTime = now

	imgui.NewFrame()
}

// Render must be called at the end of rendering.
func (context *Context) Render() {
	imgui.Render()
	context.renderDrawData(imgui.RenderedDrawData())
}

// DrawGUI ...
func (context *Context) DrawGUI(isFrame bool, rm *rendering.RenderManager) {
	context.DrawMainMenu()

	if context.GuiVars.showControls {
		context.viewControls.Render(&context.GuiVars.showControls, &isFrame, rm)
	}

	if context.GuiVars.showModels {
		context.viewModels.Render(&context.GuiVars.showModels, &isFrame)
	}

	if context.GuiVars.showLog {
		context.componentLog.Render(&context.GuiVars.showLog)
	}

	if context.GuiVars.showAboutImGui {
		context.showAboutImGui(&context.GuiVars.showAboutImGui)
	}

	if context.GuiVars.showAboutKuplung {
		context.showAboutKuplung(&context.GuiVars.showAboutKuplung)
	}

	if context.GuiVars.showDemoWindow {
		imgui.ShowDemoWindow(&context.GuiVars.showDemoWindow)
	}

	if context.GuiVars.showMetrics {
		context.showMetrics(&context.GuiVars.showMetrics)
	}

	if context.GuiVars.showParsing {
		context.showParsing(&context.GuiVars.showParsing)
	}
}

// IsUsingKeyboard returns true if the UI is currently capturing keyboard input.
// The application should not process keyboard input events in this case.
func (context Context) IsUsingKeyboard() bool {
	return imgui.CurrentIO().WantTextInput()
}

// IsUsingMouse returns true if the UI is using the mouse.
// The application should not process mouse events in this case.
func (context Context) IsUsingMouse() bool {
	return imgui.CurrentIO().WantCaptureMouse()
}

// MouseScroll must be
func (context *Context) MouseScroll(dx, dy float32) {
	imgui.CurrentIO().AddMouseWheelDelta(dx, dy)
}

func (context *Context) createDeviceObjects(param ContextParameters) (err error) {
	gl := context.window.OpenGL()
	glslVersion := "#version 150"

	vertexShaderSource := glslVersion + `
uniform mat4 ProjMtx;
in vec2 Position;
in vec2 UV;
in vec4 Color;
out vec2 Frag_UV;
out vec4 Frag_Color;
void main()
{
	Frag_UV = UV;
	Frag_Color = Color;
	gl_Position = ProjMtx * vec4(Position.xy,0,1);
}
`
	fragmentShaderSource := glslVersion + `
uniform int ImageType;
uniform sampler2D Texture;
uniform sampler2D Palette;
in vec2 Frag_UV;
in vec4 Frag_Color;
out vec4 Out_Color;
void main()
{
	if (ImageType == 1)
	{
		vec4 pixel = texture(Texture, Frag_UV.st);
		Out_Color = Frag_Color * texture(Palette, vec2(pixel.r, 0.5));
	}
	else
	{
		Out_Color = vec4(Frag_Color.rgb, Frag_Color.a * texture( Texture, Frag_UV.st).r);
	}
}
`

	context.shaderHandle, err = engine.LinkNewStandardProgram(gl, vertexShaderSource, fragmentShaderSource)
	if err != nil {
		return
	}

	context.attribLocationType = gl.GetUniformLocation(context.shaderHandle, "ImageType")
	context.attribLocationTex = gl.GetUniformLocation(context.shaderHandle, "Texture")
	context.attribLocationPal = gl.GetUniformLocation(context.shaderHandle, "Palette")
	context.attribLocationProjMtx = gl.GetUniformLocation(context.shaderHandle, "ProjMtx")
	context.attribLocationPosition = gl.GetAttribLocation(context.shaderHandle, "Position")
	context.attribLocationUV = gl.GetAttribLocation(context.shaderHandle, "UV")
	context.attribLocationColor = gl.GetAttribLocation(context.shaderHandle, "Color")

	buffers := gl.GenBuffers(2)
	context.vboHandle = buffers[0]
	context.elementsHandle = buffers[1]

	return context.createFontsTexture(gl, param)
}

func (context *Context) createFontsTexture(gl interfaces.OpenGL, param ContextParameters) error {
	io := imgui.CurrentIO()
	fontAtlas := io.Fonts()
	_ = fontAtlas.AddFontDefault()
	if len(param.FontFile) > 0 {
		fontSize := float32(16.0)
		if param.FontSize > 0.0 {
			fontSize = param.FontSize
		}
		font := fontAtlas.AddFontFromFileTTF(param.FontFile, fontSize)
		if font == imgui.DefaultFont {
			return fmt.Errorf("could not load font <%s>", param.FontFile)
		}
	}

	sett := settings.GetSettings()
	fileFA := sett.App.CurrentPath + "fonts/fontawesome-webfont.ttf"
	fontConfigFA := imgui.NewFontConfig()
	fontConfigFA.SetGlyphMaxAdvanceX(float32(fonts.FA_ICON_MIN))
	fontConfigFA.SetGlyphMinAdvanceX(float32(fonts.FA_ICON_MAX))
	fontConfigFA.SetMergeMode(true)
	fontConfigFA.SetPixelSnapH(true)
	var builderFA imgui.GlyphRangesBuilder
	builderFA.Add(rune(fonts.FA_ICON_MIN), rune(fonts.FA_ICON_MAX))
	rangesFA := builderFA.Build()
	context.fontFA = fontAtlas.AddFontFromFileTTFV(fileFA, 14, fontConfigFA, rangesFA.GlyphRanges)
	fontConfigFA.Delete()

	fileMD := sett.App.CurrentPath + "fonts/material-icons-regular.ttf"
	fontConfigMD := imgui.NewFontConfig()
	fontConfigMD.SetGlyphMaxAdvanceX(float32(fonts.MD_ICON_MIN))
	fontConfigMD.SetGlyphMinAdvanceX(float32(fonts.MD_ICON_MAX))
	fontConfigMD.SetMergeMode(true)
	fontConfigMD.SetPixelSnapH(true)
	var builderMD imgui.GlyphRangesBuilder
	builderMD.Add(rune(fonts.MD_ICON_MIN), rune(fonts.MD_ICON_MAX))
	rangesMD := builderMD.Build()
	context.fontMD = fontAtlas.AddFontFromFileTTFV(fileMD, 14, fontConfigMD, rangesMD.GlyphRanges)
	fontConfigMD.Delete()

	image := fontAtlas.TextureDataAlpha8()

	context.fontTexture = gl.GenTextures(1)[0]
	gl.BindTexture(oglconsts.TEXTURE_2D, context.fontTexture)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MIN_FILTER, oglconsts.LINEAR)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MAG_FILTER, oglconsts.LINEAR)
	gl.PixelStorei(oglconsts.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(oglconsts.TEXTURE_2D, 0, oglconsts.RED, int32(image.Width), int32(image.Height),
		0, oglconsts.RED, oglconsts.UNSIGNED_BYTE, image.Pixels)

	io.Fonts().SetTextureID(TextureIDForSimpleTexture(context.fontTexture))

	gl.BindTexture(oglconsts.TEXTURE_2D, 0)

	return nil
}

func (context *Context) destroyDeviceObjects(gl interfaces.OpenGL) {
	if context.vboHandle != 0 {
		gl.DeleteBuffers([]uint32{context.vboHandle})
	}
	context.vboHandle = 0
	if context.elementsHandle != 0 {
		gl.DeleteBuffers([]uint32{context.elementsHandle})
	}
	context.elementsHandle = 0

	if context.shaderHandle != 0 {
		gl.DeleteProgram(context.shaderHandle)
	}
	context.shaderHandle = 0

	if context.fontTexture != 0 {
		gl.DeleteTextures([]uint32{context.fontTexture})
		imgui.CurrentIO().Fonts().SetTextureID(0)
		context.fontTexture = 0
	}
}

func (context *Context) renderDrawData(drawData imgui.DrawData) {
	gl := context.window.OpenGL()
	sett := settings.GetSettings()
	displayWidth, displayHeight := int(sett.AppWindow.SDLWindowWidth), int(sett.AppWindow.SDLWindowHeight)

	// Avoid rendering when minimized, scale coordinates for retina displays (screen coordinates != framebuffer coordinates)
	fbWidth, fbHeight := context.window.Size()
	if (fbWidth <= 0) || (fbHeight <= 0) {
		return
	}
	drawData.ScaleClipRects(imgui.Vec2{
		X: float32(fbWidth / displayWidth),
		Y: float32(fbHeight / displayHeight),
	})

	// Backup GL state
	var lastActiveTexture int32
	gl.GetIntegerv(oglconsts.ACTIVE_TEXTURE, &lastActiveTexture)
	gl.ActiveTexture(oglconsts.TEXTURE0)
	var lastProgram int32
	gl.GetIntegerv(oglconsts.CURRENT_PROGRAM, &lastProgram)
	var lastTexture int32
	gl.GetIntegerv(oglconsts.TEXTURE_BINDING_2D, &lastTexture)
	var lastSampler int32
	gl.GetIntegerv(oglconsts.SAMPLER_BINDING, &lastSampler)
	var lastArrayBuffer int32
	gl.GetIntegerv(oglconsts.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	var lastElementArrayBuffer int32
	gl.GetIntegerv(oglconsts.ELEMENT_ARRAY_BUFFER_BINDING, &lastElementArrayBuffer)
	var lastVertexArray int32
	gl.GetIntegerv(oglconsts.VERTEX_ARRAY_BINDING, &lastVertexArray)
	var lastPolygonMode [2]int32
	gl.GetIntegerv(oglconsts.POLYGON_MODE, &lastPolygonMode[0])
	var lastViewport [4]int32
	gl.GetIntegerv(oglconsts.VIEWPORT, &lastViewport[0])
	var lastScissorBox [4]int32
	gl.GetIntegerv(oglconsts.SCISSOR_BOX, &lastScissorBox[0])
	var lastBlendSrcRgb int32
	gl.GetIntegerv(oglconsts.BLEND_SRC_RGB, &lastBlendSrcRgb)
	var lastBlendDstRgb int32
	gl.GetIntegerv(oglconsts.BLEND_DST_RGB, &lastBlendDstRgb)
	var lastBlendSrcAlpha int32
	gl.GetIntegerv(oglconsts.BLEND_SRC_ALPHA, &lastBlendSrcAlpha)
	var lastBlendDstAlpha int32
	gl.GetIntegerv(oglconsts.BLEND_DST_ALPHA, &lastBlendDstAlpha)
	var lastBlendEquationRgb int32
	gl.GetIntegerv(oglconsts.BLEND_EQUATION_RGB, &lastBlendEquationRgb)
	var lastBlendEquationAlpha int32
	gl.GetIntegerv(oglconsts.BLEND_EQUATION_ALPHA, &lastBlendEquationAlpha)
	lastEnableBlend := gl.IsEnabled(oglconsts.BLEND)
	lastEnableCullFace := gl.IsEnabled(oglconsts.CULL_FACE)
	lastEnableDepthTest := gl.IsEnabled(oglconsts.DEPTH_TEST)
	lastEnableScissorTest := gl.IsEnabled(oglconsts.SCISSOR_TEST)

	// Setup render state: alpha-blending enabled, no face culling, no depth testing, scissor enabled, polygon fill
	gl.Enable(oglconsts.BLEND)
	gl.BlendEquation(oglconsts.FUNC_ADD)
	gl.BlendFunc(oglconsts.SRC_ALPHA, oglconsts.ONE_MINUS_SRC_ALPHA)
	gl.Disable(oglconsts.CULL_FACE)
	gl.Disable(oglconsts.DEPTH_TEST)
	gl.Enable(oglconsts.SCISSOR_TEST)
	gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.FILL)

	// Setup viewport, orthographic projection matrix
	// Our visible imgui space lies from draw_data->DisplayPos (top left) to draw_data->DisplayPos+data_data->DisplaySize (bottom right).
	// DisplayMin is typically (0,0) for single viewport apps.
	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
	orthoProjection := [16]float32{
		2.0 / float32(displayWidth), 0.0, 0.0, 0.0,
		0.0, 2.0 / float32(-displayHeight), 0.0, 0.0,
		0.0, 0.0, -1.0, 0.0,
		-1.0, 1.0, 0.0, 1.0,
	}
	gl.UseProgram(context.shaderHandle)
	gl.Uniform1i(context.attribLocationTex, 0)
	gl.UniformMatrix4fv(context.attribLocationProjMtx, false, &orthoProjection)
	gl.BindSampler(0, 0) // Rely on combined texture/sampler state.

	vaoHandle := gl.GenVertexArrays(1)[0]
	gl.BindVertexArray(vaoHandle)
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, context.vboHandle)
	gl.EnableVertexAttribArray(uint32(context.attribLocationPosition))
	gl.EnableVertexAttribArray(uint32(context.attribLocationUV))
	gl.EnableVertexAttribArray(uint32(context.attribLocationColor))
	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	gl.VertexAttribOffset(uint32(context.attribLocationPosition), 2, oglconsts.FLOAT, false, int32(vertexSize), vertexOffsetPos)
	gl.VertexAttribOffset(uint32(context.attribLocationUV), 2, oglconsts.FLOAT, false, int32(vertexSize), vertexOffsetUv)
	gl.VertexAttribOffset(uint32(context.attribLocationColor), 4, oglconsts.UNSIGNED_BYTE, true, int32(vertexSize), vertexOffsetCol)
	indexSize := imgui.IndexBufferLayout()
	drawType := oglconsts.UNSIGNED_SHORT
	if indexSize == 4 {
		drawType = oglconsts.UNSIGNED_INT
	}

	// Draw
	for _, list := range drawData.CommandLists() {
		var indexBufferOffset uintptr

		vertexBuffer, vertexBufferSize := list.VertexBuffer()
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, context.vboHandle)
		gl.BufferData(oglconsts.ARRAY_BUFFER, vertexBufferSize, vertexBuffer, oglconsts.STREAM_DRAW)

		indexBuffer, indexBufferSize := list.IndexBuffer()
		gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, context.elementsHandle)
		gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, indexBufferSize, indexBuffer, oglconsts.STREAM_DRAW)

		for _, cmd := range list.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(list)
			} else {
				gl.BindTexture(oglconsts.TEXTURE_2D, uint32(cmd.TextureID()))
				clipRect := cmd.ClipRect()
				gl.Scissor(int32(clipRect.X), int32(fbHeight)-int32(clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))
				gl.DrawElements(oglconsts.TRIANGLES, int32(cmd.ElementCount()), uint32(drawType), indexBufferOffset)
			}
			indexBufferOffset += uintptr(cmd.ElementCount() * indexSize)
		}
	}
	gl.DeleteVertexArrays([]uint32{vaoHandle})

	// Restore modified GL state
	gl.UseProgram(uint32(lastProgram))
	gl.BindTexture(oglconsts.TEXTURE_2D, uint32(lastTexture))
	gl.BindSampler(0, uint32(lastSampler))
	gl.ActiveTexture(uint32(lastActiveTexture))
	gl.BindVertexArray(uint32(lastVertexArray))
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, uint32(lastElementArrayBuffer))
	gl.BlendEquationSeparate(uint32(lastBlendEquationRgb), uint32(lastBlendEquationAlpha))
	gl.BlendFuncSeparate(uint32(lastBlendSrcRgb), uint32(lastBlendDstRgb), uint32(lastBlendSrcAlpha), uint32(lastBlendDstAlpha))
	if lastEnableBlend {
		gl.Enable(oglconsts.BLEND)
	} else {
		gl.Disable(oglconsts.BLEND)
	}
	if lastEnableCullFace {
		gl.Enable(oglconsts.CULL_FACE)
	} else {
		gl.Disable(oglconsts.CULL_FACE)
	}
	if lastEnableDepthTest {
		gl.Enable(oglconsts.DEPTH_TEST)
	} else {
		gl.Disable(oglconsts.DEPTH_TEST)
	}
	if lastEnableScissorTest {
		gl.Enable(oglconsts.SCISSOR_TEST)
	} else {
		gl.Disable(oglconsts.SCISSOR_TEST)
	}
	gl.PolygonMode(oglconsts.FRONT_AND_BACK, uint32(lastPolygonMode[0]))
	gl.Viewport(lastViewport[0], lastViewport[1], lastViewport[2], lastViewport[3])
	gl.Scissor(lastScissorBox[0], lastScissorBox[1], lastScissorBox[2], lastScissorBox[3])
}

func (context *Context) renderDrawData2(drawData imgui.DrawData) {
	gl := context.window.OpenGL()

	displayWidth, displayHeight := context.window.Size()

	// Backup GL state
	var lastActiveTexture int32
	gl.GetIntegerv(oglconsts.ACTIVE_TEXTURE, &lastActiveTexture)
	gl.ActiveTexture(oglconsts.TEXTURE0)
	var lastProgram int32
	gl.GetIntegerv(oglconsts.CURRENT_PROGRAM, &lastProgram)
	var lastTexture int32
	gl.GetIntegerv(oglconsts.TEXTURE_BINDING_2D, &lastTexture)
	var lastSampler int32
	gl.GetIntegerv(oglconsts.SAMPLER_BINDING, &lastSampler)
	var lastArrayBuffer int32
	gl.GetIntegerv(oglconsts.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	var lastElementArrayBuffer int32
	gl.GetIntegerv(oglconsts.ELEMENT_ARRAY_BUFFER_BINDING, &lastElementArrayBuffer)
	var lastVertexArray int32
	gl.GetIntegerv(oglconsts.VERTEX_ARRAY_BINDING, &lastVertexArray)
	var lastPolygonMode [2]int32
	gl.GetIntegerv(oglconsts.POLYGON_MODE, &lastPolygonMode[0])
	var lastViewport [4]int32
	gl.GetIntegerv(oglconsts.VIEWPORT, &lastViewport[0])
	var lastScissorBox [4]int32
	gl.GetIntegerv(oglconsts.SCISSOR_BOX, &lastScissorBox[0])
	var lastBlendSrcRgb int32
	gl.GetIntegerv(oglconsts.BLEND_SRC_RGB, &lastBlendSrcRgb)
	var lastBlendDstRgb int32
	gl.GetIntegerv(oglconsts.BLEND_DST_RGB, &lastBlendDstRgb)
	var lastBlendSrcAlpha int32
	gl.GetIntegerv(oglconsts.BLEND_SRC_ALPHA, &lastBlendSrcAlpha)
	var lastBlendDstAlpha int32
	gl.GetIntegerv(oglconsts.BLEND_DST_ALPHA, &lastBlendDstAlpha)
	var lastBlendEquationRgb int32
	gl.GetIntegerv(oglconsts.BLEND_EQUATION_RGB, &lastBlendEquationRgb)
	var lastBlendEquationAlpha int32
	gl.GetIntegerv(oglconsts.BLEND_EQUATION_ALPHA, &lastBlendEquationAlpha)
	lastEnableBlend := gl.IsEnabled(oglconsts.BLEND)
	lastEnableCullFace := gl.IsEnabled(oglconsts.CULL_FACE)
	lastEnableDepthTest := gl.IsEnabled(oglconsts.DEPTH_TEST)
	lastEnableScissorTest := gl.IsEnabled(oglconsts.SCISSOR_TEST)

	// Setup render state: alpha-blending enabled, no face culling, no depth testing, scissor enabled, polygon fill
	gl.Enable(oglconsts.BLEND)
	gl.BlendEquation(oglconsts.FUNC_ADD)
	gl.BlendFunc(oglconsts.SRC_ALPHA, oglconsts.ONE_MINUS_SRC_ALPHA)
	gl.Disable(oglconsts.CULL_FACE)
	gl.Disable(oglconsts.DEPTH_TEST)
	gl.Enable(oglconsts.SCISSOR_TEST)
	gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.FILL)

	// Setup viewport, orthographic projection matrix
	gl.Viewport(0, 0, int32(displayWidth), int32(displayHeight))
	orthoProjection := [16]float32{
		2.0 / float32(displayWidth), 0.0, 0.0, 0.0,
		0.0, 2.0 / float32(-displayHeight), 0.0, 0.0,
		0.0, 0.0, -1.0, 0.0,
		-1.0, 1.0, 0.0, 1.0,
	}
	gl.UseProgram(context.shaderHandle)
	gl.Uniform1i(context.attribLocationTex, 0)
	gl.Uniform1i(context.attribLocationPal, 1)
	gl.UniformMatrix4fv(context.attribLocationProjMtx, false, &orthoProjection)
	gl.BindSampler(0, 0) // Rely on combined texture/sampler state.

	vaoHandle := gl.GenVertexArrays(1)[0]
	gl.BindVertexArray(vaoHandle)
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, context.vboHandle)
	gl.EnableVertexAttribArray(uint32(context.attribLocationPosition))
	gl.EnableVertexAttribArray(uint32(context.attribLocationUV))
	gl.EnableVertexAttribArray(uint32(context.attribLocationColor))
	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	gl.VertexAttribOffset(uint32(context.attribLocationPosition), 2, oglconsts.FLOAT, false, int32(vertexSize), vertexOffsetPos)
	gl.VertexAttribOffset(uint32(context.attribLocationUV), 2, oglconsts.FLOAT, false, int32(vertexSize), vertexOffsetUv)
	gl.VertexAttribOffset(uint32(context.attribLocationColor), 4, oglconsts.UNSIGNED_BYTE, true, int32(vertexSize), vertexOffsetCol)
	indexSize := imgui.IndexBufferLayout()
	drawType := oglconsts.UNSIGNED_SHORT
	if indexSize == 4 {
		drawType = oglconsts.UNSIGNED_INT
	}

	// Draw
	for _, list := range drawData.CommandLists() {
		var indexBufferOffset uintptr

		vertexBuffer, vertexBufferSize := list.VertexBuffer()
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, context.vboHandle)
		gl.BufferData(oglconsts.ARRAY_BUFFER, vertexBufferSize, vertexBuffer, oglconsts.STREAM_DRAW)

		indexBuffer, indexBufferSize := list.IndexBuffer()
		gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, context.elementsHandle)
		gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, indexBufferSize, indexBuffer, oglconsts.STREAM_DRAW)

		for _, cmd := range list.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(list)
			} else {
				textureID := cmd.TextureID()
				imageType := ImageTypeFromID(textureID)
				gl.Uniform1i(context.attribLocationType, int32(imageType))
				switch imageType {
				case ImageTypeSimpleTexture:
					gl.ActiveTexture(oglconsts.TEXTURE0 + uint32(0))
					gl.BindTexture(oglconsts.TEXTURE_2D, uint32(textureID))
				// case ImageTypeBitmapTexture:
				// 	palette, bitmap := bitmapTextureQuery(textureID)
				// 	gl.ActiveTexture(oglconsts.TEXTURE0 + uint32(0))
				// 	gl.BindTexture(oglconsts.TEXTURE_2D, bitmap)
				// 	gl.ActiveTexture(oglconsts.TEXTURE0 + uint32(1))
				// 	gl.BindTexture(oglconsts.TEXTURE_2D, palette)
				default:
					gl.ActiveTexture(oglconsts.TEXTURE0 + uint32(0))
					gl.BindTexture(oglconsts.TEXTURE_2D, 0)
					gl.ActiveTexture(oglconsts.TEXTURE0 + uint32(1))
					gl.BindTexture(oglconsts.TEXTURE_2D, 0)
				}
				clipRect := cmd.ClipRect()
				gl.Scissor(int32(clipRect.X), int32(displayHeight)-int32(clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))
				gl.DrawElements(oglconsts.TRIANGLES, int32(cmd.ElementCount()), uint32(drawType), indexBufferOffset)
			}
			indexBufferOffset += uintptr(cmd.ElementCount() * indexSize)
		}
	}
	gl.DeleteVertexArrays([]uint32{vaoHandle})

	// Restore modified GL state
	gl.UseProgram(uint32(lastProgram))
	gl.BindTexture(oglconsts.TEXTURE_2D, uint32(lastTexture))
	gl.BindSampler(0, uint32(lastSampler))
	gl.ActiveTexture(uint32(lastActiveTexture))
	gl.BindVertexArray(uint32(lastVertexArray))
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, uint32(lastElementArrayBuffer))
	gl.BlendEquationSeparate(uint32(lastBlendEquationRgb), uint32(lastBlendEquationAlpha))
	gl.BlendFuncSeparate(uint32(lastBlendSrcRgb), uint32(lastBlendDstRgb), uint32(lastBlendSrcAlpha), uint32(lastBlendDstAlpha))
	if lastEnableBlend {
		gl.Enable(oglconsts.BLEND)
	} else {
		gl.Disable(oglconsts.BLEND)
	}
	if lastEnableCullFace {
		gl.Enable(oglconsts.CULL_FACE)
	} else {
		gl.Disable(oglconsts.CULL_FACE)
	}
	if lastEnableDepthTest {
		gl.Enable(oglconsts.DEPTH_TEST)
	} else {
		gl.Disable(oglconsts.DEPTH_TEST)
	}
	if lastEnableScissorTest {
		gl.Enable(oglconsts.SCISSOR_TEST)
	} else {
		gl.Disable(oglconsts.SCISSOR_TEST)
	}
	gl.PolygonMode(oglconsts.FRONT_AND_BACK, uint32(lastPolygonMode[0]))
	gl.Viewport(lastViewport[0], lastViewport[1], lastViewport[2], lastViewport[3])
	gl.Scissor(lastScissorBox[0], lastScissorBox[1], lastScissorBox[2], lastScissorBox[3])
}
