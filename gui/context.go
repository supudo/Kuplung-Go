package gui

import (
	"fmt"
	"os"
	"time"

	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
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

	mouseButtonWasDown [3]bool
	mouseButtonIsDown  [3]bool

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

	guiVars WindowVariables
}

// WindowVariables holds boolean variables for all the windows
type WindowVariables struct {
	showDemoWindow   bool
	showAboutImGui   bool
	showAboutKuplung bool
	showMetrics      bool
}

// NewContext initializes a new UI context based on the provided OpenGL window.
func NewContext(window interfaces.Window, param ContextParameters) *Context {
	imgui.SetAssertHandler(nil)
	context := &Context{
		imguiContext: imgui.CreateContext(nil),
		window:       window,
	}

	context.guiVars.showDemoWindow = false
	context.guiVars.showAboutImGui = false
	context.guiVars.showAboutKuplung = false
	context.guiVars.showMetrics = false

	err := context.createDeviceObjects(param)
	if err != nil {
		context.Destroy()
		context = nil
		settings.LogError("[gui context] Error initialized ImGui Context: %v", err)
	}

	return context
}

// Destroy cleans up the resources of the graphical user interface.
func (context *Context) Destroy() {
	context.destroyDeviceObjects(context.window.OpenGL())
	context.imguiContext.Destroy()
}

// MouseButtonChanged is called when the state of a button has changed.
func (context *Context) MouseButtonChanged(buttonIndex int, down bool) {
	if (buttonIndex >= 0) && (buttonIndex < len(context.mouseButtonIsDown)) {
		context.mouseButtonIsDown[buttonIndex] = down
		if down {
			context.mouseButtonWasDown[buttonIndex] = down
		}
	}
}

// IsUsingMouse returns true if the UI is using the mouse.
// The application should not process mouse events in this case.
func (context Context) IsUsingMouse() bool {
	return imgui.CurrentIO().WantCaptureMouse()
}

// IsUsingKeyboard returns true if the UI is currently capturing keyboard input.
// The application should not process keyboard input events in this case.
func (context Context) IsUsingKeyboard() bool {
	return imgui.CurrentIO().WantTextInput()
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

	for i := 0; i < len(context.mouseButtonWasDown); i++ {
		down := context.mouseButtonWasDown[i] || context.mouseButtonIsDown[i]
		io.SetMouseButtonDown(i, down)
		context.mouseButtonWasDown[i] = false
	}

	imgui.NewFrame()
}

// Render must be called at the end of rendering.
func (context *Context) Render() {
	context.DrawMainMenu()

	imgui.Render()
	context.renderDrawData(imgui.RenderedDrawData())
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
	image := fontAtlas.TextureDataAlpha8()

	context.fontTexture = gl.GenTextures(1)[0]
	gl.BindTexture(engine.TEXTURE_2D, context.fontTexture)
	gl.TexParameteri(engine.TEXTURE_2D, engine.TEXTURE_MIN_FILTER, engine.LINEAR)
	gl.TexParameteri(engine.TEXTURE_2D, engine.TEXTURE_MAG_FILTER, engine.LINEAR)
	gl.PixelStorei(engine.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(engine.TEXTURE_2D, 0, engine.RED, int32(image.Width), int32(image.Height),
		0, engine.RED, engine.UNSIGNED_BYTE, image.Pixels)

	io.Fonts().SetTextureID(TextureIDForSimpleTexture(context.fontTexture))

	gl.BindTexture(engine.TEXTURE_2D, 0)
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

	displayWidth, displayHeight := context.window.Size()

	// Backup GL state
	var lastActiveTexture int32
	gl.GetIntegerv(engine.ACTIVE_TEXTURE, &lastActiveTexture)
	gl.ActiveTexture(engine.TEXTURE0)
	var lastProgram int32
	gl.GetIntegerv(engine.CURRENT_PROGRAM, &lastProgram)
	var lastTexture int32
	gl.GetIntegerv(engine.TEXTURE_BINDING_2D, &lastTexture)
	var lastSampler int32
	gl.GetIntegerv(engine.SAMPLER_BINDING, &lastSampler)
	var lastArrayBuffer int32
	gl.GetIntegerv(engine.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	var lastElementArrayBuffer int32
	gl.GetIntegerv(engine.ELEMENT_ARRAY_BUFFER_BINDING, &lastElementArrayBuffer)
	var lastVertexArray int32
	gl.GetIntegerv(engine.VERTEX_ARRAY_BINDING, &lastVertexArray)
	var lastPolygonMode [2]int32
	gl.GetIntegerv(engine.POLYGON_MODE, &lastPolygonMode[0])
	var lastViewport [4]int32
	gl.GetIntegerv(engine.VIEWPORT, &lastViewport[0])
	var lastScissorBox [4]int32
	gl.GetIntegerv(engine.SCISSOR_BOX, &lastScissorBox[0])
	var lastBlendSrcRgb int32
	gl.GetIntegerv(engine.BLEND_SRC_RGB, &lastBlendSrcRgb)
	var lastBlendDstRgb int32
	gl.GetIntegerv(engine.BLEND_DST_RGB, &lastBlendDstRgb)
	var lastBlendSrcAlpha int32
	gl.GetIntegerv(engine.BLEND_SRC_ALPHA, &lastBlendSrcAlpha)
	var lastBlendDstAlpha int32
	gl.GetIntegerv(engine.BLEND_DST_ALPHA, &lastBlendDstAlpha)
	var lastBlendEquationRgb int32
	gl.GetIntegerv(engine.BLEND_EQUATION_RGB, &lastBlendEquationRgb)
	var lastBlendEquationAlpha int32
	gl.GetIntegerv(engine.BLEND_EQUATION_ALPHA, &lastBlendEquationAlpha)
	lastEnableBlend := gl.IsEnabled(engine.BLEND)
	lastEnableCullFace := gl.IsEnabled(engine.CULL_FACE)
	lastEnableDepthTest := gl.IsEnabled(engine.DEPTH_TEST)
	lastEnableScissorTest := gl.IsEnabled(engine.SCISSOR_TEST)

	// Setup render state: alpha-blending enabled, no face culling, no depth testing, scissor enabled, polygon fill
	gl.Enable(engine.BLEND)
	gl.BlendEquation(engine.FUNC_ADD)
	gl.BlendFunc(engine.SRC_ALPHA, engine.ONE_MINUS_SRC_ALPHA)
	gl.Disable(engine.CULL_FACE)
	gl.Disable(engine.DEPTH_TEST)
	gl.Enable(engine.SCISSOR_TEST)
	gl.PolygonMode(engine.FRONT_AND_BACK, engine.FILL)

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
	gl.BindBuffer(engine.ARRAY_BUFFER, context.vboHandle)
	gl.EnableVertexAttribArray(uint32(context.attribLocationPosition))
	gl.EnableVertexAttribArray(uint32(context.attribLocationUV))
	gl.EnableVertexAttribArray(uint32(context.attribLocationColor))
	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	gl.VertexAttribOffset(uint32(context.attribLocationPosition), 2, engine.FLOAT, false, int32(vertexSize), vertexOffsetPos)
	gl.VertexAttribOffset(uint32(context.attribLocationUV), 2, engine.FLOAT, false, int32(vertexSize), vertexOffsetUv)
	gl.VertexAttribOffset(uint32(context.attribLocationColor), 4, engine.UNSIGNED_BYTE, true, int32(vertexSize), vertexOffsetCol)
	indexSize := imgui.IndexBufferLayout()
	drawType := engine.UNSIGNED_SHORT
	if indexSize == 4 {
		drawType = engine.UNSIGNED_INT
	}

	// Draw
	for _, list := range drawData.CommandLists() {
		var indexBufferOffset uintptr

		vertexBuffer, vertexBufferSize := list.VertexBuffer()
		gl.BindBuffer(engine.ARRAY_BUFFER, context.vboHandle)
		gl.BufferData(engine.ARRAY_BUFFER, vertexBufferSize, vertexBuffer, engine.STREAM_DRAW)

		indexBuffer, indexBufferSize := list.IndexBuffer()
		gl.BindBuffer(engine.ELEMENT_ARRAY_BUFFER, context.elementsHandle)
		gl.BufferData(engine.ELEMENT_ARRAY_BUFFER, indexBufferSize, indexBuffer, engine.STREAM_DRAW)

		for _, cmd := range list.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(list)
			} else {
				textureID := cmd.TextureID()
				imageType := ImageTypeFromID(textureID)
				gl.Uniform1i(context.attribLocationType, int32(imageType))
				switch imageType {
				case ImageTypeSimpleTexture:
					gl.ActiveTexture(engine.TEXTURE0 + uint32(0))
					gl.BindTexture(engine.TEXTURE_2D, uint32(textureID))
				// case ImageTypeBitmapTexture:
				// 	palette, bitmap := bitmapTextureQuery(textureID)
				// 	gl.ActiveTexture(engine.TEXTURE0 + uint32(0))
				// 	gl.BindTexture(engine.TEXTURE_2D, bitmap)
				// 	gl.ActiveTexture(engine.TEXTURE0 + uint32(1))
				// 	gl.BindTexture(engine.TEXTURE_2D, palette)
				default:
					gl.ActiveTexture(engine.TEXTURE0 + uint32(0))
					gl.BindTexture(engine.TEXTURE_2D, 0)
					gl.ActiveTexture(engine.TEXTURE0 + uint32(1))
					gl.BindTexture(engine.TEXTURE_2D, 0)
				}
				clipRect := cmd.ClipRect()
				gl.Scissor(int32(clipRect.X), int32(displayHeight)-int32(clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))
				gl.DrawElements(engine.TRIANGLES, int32(cmd.ElementCount()), uint32(drawType), indexBufferOffset)
			}
			indexBufferOffset += uintptr(cmd.ElementCount() * indexSize)
		}
	}
	gl.DeleteVertexArrays([]uint32{vaoHandle})

	// Restore modified GL state
	gl.UseProgram(uint32(lastProgram))
	gl.BindTexture(engine.TEXTURE_2D, uint32(lastTexture))
	gl.BindSampler(0, uint32(lastSampler))
	gl.ActiveTexture(uint32(lastActiveTexture))
	gl.BindVertexArray(uint32(lastVertexArray))
	gl.BindBuffer(engine.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindBuffer(engine.ELEMENT_ARRAY_BUFFER, uint32(lastElementArrayBuffer))
	gl.BlendEquationSeparate(uint32(lastBlendEquationRgb), uint32(lastBlendEquationAlpha))
	gl.BlendFuncSeparate(uint32(lastBlendSrcRgb), uint32(lastBlendDstRgb), uint32(lastBlendSrcAlpha), uint32(lastBlendDstAlpha))
	if lastEnableBlend {
		gl.Enable(engine.BLEND)
	} else {
		gl.Disable(engine.BLEND)
	}
	if lastEnableCullFace {
		gl.Enable(engine.CULL_FACE)
	} else {
		gl.Disable(engine.CULL_FACE)
	}
	if lastEnableDepthTest {
		gl.Enable(engine.DEPTH_TEST)
	} else {
		gl.Disable(engine.DEPTH_TEST)
	}
	if lastEnableScissorTest {
		gl.Enable(engine.SCISSOR_TEST)
	} else {
		gl.Disable(engine.SCISSOR_TEST)
	}
	gl.PolygonMode(engine.FRONT_AND_BACK, uint32(lastPolygonMode[0]))
	gl.Viewport(lastViewport[0], lastViewport[1], lastViewport[2], lastViewport[3])
	gl.Scissor(lastScissorBox[0], lastScissorBox[1], lastScissorBox[2], lastScissorBox[3])
}

// SetMousePosition must be called to report the current mouse position.
func (context *Context) SetMousePosition(x, y float32) {
	imgui.CurrentIO().SetMousePosition(imgui.Vec2{X: x, Y: y})
}

// MouseScroll must be
func (context *Context) MouseScroll(dx, dy float32) {
	imgui.CurrentIO().AddMouseWheelDelta(dx, dy)
}

// GUI elements

// DrawMainMenu ...
func (context *Context) DrawMainMenu() {
	// Main Menu
	imgui.BeginMainMenuBar()

	if imgui.BeginMenu("File") {
		imgui.Separator()
		if imgui.MenuItemV("Quit", "Cmd+Q", false, true) {
			os.Exit(3)
		}
		imgui.EndMenu()
	}

	if imgui.BeginMenu("Scene") {
		imgui.EndMenu()
	}

	if imgui.BeginMenu("View") {
		imgui.EndMenu()
	}

	if imgui.BeginMenu("Help") {
		if imgui.MenuItem("Metrics") {
			context.guiVars.showMetrics = true
		}
		if imgui.MenuItem("About ImGui") {
			context.guiVars.showAboutImGui = true
		}
		if imgui.MenuItem("About Kuplung") {
			context.guiVars.showAboutKuplung = true
		}
		imgui.Separator()
		if imgui.MenuItem("ImGui Demo Window") {
			context.guiVars.showDemoWindow = true
		}
		imgui.EndMenu()
	}

	imgui.EndMainMenuBar()

	if context.guiVars.showAboutImGui {
		showAboutImGui(&context.guiVars.showAboutImGui)
	}

	if context.guiVars.showAboutKuplung {
		showAboutKuplung(&context.guiVars.showAboutKuplung)
	}

	if context.guiVars.showDemoWindow {
		imgui.ShowDemoWindow(&context.guiVars.showDemoWindow)
	}

	if context.guiVars.showMetrics {
		showMetrics(&context.guiVars.showMetrics)
	}
}

func showAboutImGui(open *bool) {
	if imgui.BeginV("About ImGui", open, imgui.WindowFlagsAlwaysAutoResize) {
		imgui.Text("ImGui " + imgui.Version())
		imgui.Separator()
		imgui.Text("By Omar Cornut and all github contributors.")
		imgui.Text("ImGui is licensed under the MIT License, see LICENSE for more information.")
		imgui.Separator()
		imgui.Text("Go binding by Inky Blackness")
		imgui.Text("https://github.com/inkyblackness/imgui-go/")
		imgui.End()
	}
}

func showAboutKuplung(open *bool) {
	var sett = settings.GetSettings()
	if imgui.BeginV("About Kuplung", open, imgui.WindowFlagsAlwaysAutoResize) {
		imgui.Text("Kuplung " + sett.App.ApplicationVersion)
		imgui.Separator()
		imgui.Text("By supudo.net + github.com/supudo")
		imgui.Text("Whatever license...")
		imgui.Separator()
		imgui.Text("Hold mouse wheel to rotate around")
		imgui.Text("Left Alt + Mouse wheel to increase/decrease the FOV")
		imgui.Text("Left Shift + Mouse wheel to increase/decrease the FOV")
		imgui.Text("By supudo.net + github.com/supudo")
		imgui.End()
	}
}

func showMetrics(open *bool) {
	if imgui.BeginV("Scene stats", open, imgui.WindowFlagsAlwaysAutoResize|imgui.WindowFlagsNoTitleBar|imgui.WindowFlagsNoResize|imgui.WindowFlagsNoSavedSettings) {
		// imgui.Text("OpenGL version: 4.1 (" + gl.GoStr(gl.GetString(gl.VERSION)) + ")")
		// imgui.Text("GLSL version: 4.10 (" + gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)) + ")")
		// imgui.Text("Vendor: " + gl.GoStr(gl.GetString(gl.VENDOR)))
		// imgui.Text("Renderer: " + gl.GoStr(gl.GetString(gl.RENDERER)))
		imgui.End()
		// version := gl.GoStr(gl.GetString(gl.VERSION))
		// log.Fatalf("OpenGL version %v", version)
	}
}
