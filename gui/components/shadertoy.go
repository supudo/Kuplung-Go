package components

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/objects"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/veandco/go-sdl2/sdl"
)

// ComponentShadertoy ...
type ComponentShadertoy struct {
	window interfaces.Window

	vboTexture   uint32
	windowWidth  float32
	windowHeight float32

	viewPaddingHorizontal float32
	viewPaddingVertical   float32

	textureWidth  int32
	textureHeight int32

	scrolling mgl32.Vec2

	engineShadertoy *objects.Shadertoy

	heightTopPanel      float32
	widthTexturesPanel  float32
	buttonCompileHeight float32

	shadertoyEditorText string

	channel0Cube, channel1Cube, channel2Cube, channel3Cube     bool
	texImage0, texImage1, texImage2, texImage3                 int32
	cubemapImage0, cubemapImage1, cubemapImage2, cubemapImage3 int32
}

// NewComponentShadertoy ...
func NewComponentShadertoy(window interfaces.Window) *ComponentShadertoy {
	comp := &ComponentShadertoy{}
	comp.window = window
	comp.scrolling = mgl32.Vec2{0.0, 0.0}
	comp.viewPaddingHorizontal = 20.0
	comp.viewPaddingVertical = 40.0
	comp.heightTopPanel = 200.0
	comp.widthTexturesPanel = 140.0
	comp.buttonCompileHeight = 44.0
	comp.textureWidth = int32(comp.windowWidth - comp.viewPaddingHorizontal)
	comp.textureHeight = int32(comp.windowHeight - comp.viewPaddingVertical)

	comp.channel0Cube = false
	comp.channel1Cube = false
	comp.channel2Cube = false
	comp.channel3Cube = false
	comp.cubemapImage0 = -1
	comp.cubemapImage1 = -1
	comp.cubemapImage2 = -1
	comp.cubemapImage3 = -1
	comp.texImage0 = -1
	comp.texImage1 = -1
	comp.texImage2 = -1
	comp.texImage3 = -1

	comp.shadertoyEditorText = `
void mainImage(out vec4 fragColor, in vec2 fragCoord) {
	vec2 uv = fragCoord.xy / iResolution.xy;
	fragColor = vec4(uv, 0.5 + 0.5 * sin(iGlobalTime), 1.0);
}
`
	w, h := window.Size()
	comp.engineShadertoy = objects.InitShadertoy(window)
	comp.engineShadertoy.InitShaderProgram(comp.shadertoyEditorText)
	comp.engineShadertoy.InitBuffers()
	comp.engineShadertoy.InitFBO(int32(w), int32(h), &comp.vboTexture)

	return comp
}

// Render ...
func (comp *ComponentShadertoy) Render(open *bool, deltaTime float32) {
	sett := settings.GetSettings()
	rsett := settings.GetRenderingSettings()
	mousex := rsett.Controls.MouseX
	mousey := rsett.Controls.MouseY

	imgui.SetNextWindowSizeV(imgui.Vec2{X: sett.AppWindow.LogWidth, Y: sett.AppWindow.LogHeight}, imgui.ConditionFirstUseEver)
	imgui.SetNextWindowPosV(imgui.Vec2{X: 40, Y: 40}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})

	if imgui.BeginV("Shadertoy.com", open, imgui.WindowFlagsResizeFromAnySide) {
		comp.windowWidth = imgui.WindowWidth()
		comp.windowHeight = imgui.WindowHeight()
		comp.textureWidth = int32(comp.windowWidth - comp.viewPaddingHorizontal)
		comp.textureHeight = int32(comp.windowHeight - comp.viewPaddingVertical)

		if int32(comp.heightTopPanel) < comp.engineShadertoy.TextureHeight || int32(comp.heightTopPanel) > comp.engineShadertoy.TextureHeight {
			comp.engineShadertoy.InitFBO(int32(comp.windowWidth), int32(comp.heightTopPanel), &comp.vboTexture)
		}
		comp.engineShadertoy.RenderToTexture(deltaTime, mousex, mousey, float32(sdl.GetTicks()/1000.0), &comp.vboTexture)

		imgui.BeginChildV("Preview", imgui.Vec2{X: 0, Y: comp.heightTopPanel}, true, 0)

		// drawlist := imgui.DrawList()
		// sp := imgui.CursorScreenPos()
		// offset := imgui.Vec2{X: sp.X - comp.scrolling.X(), Y: sp.Y - comp.scrolling.Y()}
		// // drawlist.ChannelsSetCurrent(0)
		// bbmin := offset
		// bbmax := imgui.Vec2{X: float32(comp.textureWidth) + offset.X, Y: float32(comp.textureHeight) + offset.Y}
		// // drawlist.AddImage(ImTextureID(intptr_t(comp.vboTexture)), bbmin, bbmax)
		imgui.Image(imgui.TextureID(comp.vboTexture), imgui.Vec2{X: float32(comp.windowWidth), Y: float32(comp.heightTopPanel)})
		// drawlist.ChannelsMerge()

		/*
			ImDrawList* draw_list = ImGui::GetWindowDrawList();
			ImVec2 offset = ImGui::GetCursorScreenPos() - this->scrolling;
			draw_list->ChannelsSetCurrent(0);
			ImVec2 bb_min = offset;
			ImVec2 bb_max = ImVec2(this->textureWidth, this->textureHeight) + offset;
			draw_list->AddImage(ImTextureID(intptr_t(this->vboTexture)), bb_min, bb_max);
			draw_list->ChannelsMerge();
		*/

		imgui.EndChild()

		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 89.0 * 255.0, Y: 91.0 * 255.0, Z: 94 * 255.0, W: 1})
		imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: 119.0 * 255.0, Y: 122.0 * 255.0, Z: 124.0 * 255.0, W: 1})
		imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: .0, Y: .0, Z: .0, W: 1})
		imgui.ButtonV("###splitterGUI", imgui.Vec2{X: -1.0, Y: 6.0})
		imgui.PopStyleColorV(3)
		if imgui.IsItemActive() {
			comp.heightTopPanel += imgui.CurrentIO().MouseDelta().Y
		}
		if imgui.IsItemHovered() {
			imgui.SetMouseCursor(imgui.MouseCursorResizeNS)
		} else {
			imgui.SetMouseCursor(imgui.MouseCursorNone)
		}

		// BEGIN editor
		imgui.BeginChildV("Editor", imgui.Vec2{X: 0, Y: 0}, true, 0)

		// buttons
		if imgui.ButtonV("COMPILE", imgui.Vec2{X: imgui.WindowWidth() * 0.85, Y: comp.buttonCompileHeight}) {
			comp.compileShader()
		}
		imgui.SameLine()
		if imgui.ButtonV("Paste", imgui.Vec2{X: imgui.WindowWidth() * 0.14, Y: comp.buttonCompileHeight}) {
			comp.getFromClipboard()
		}

		// BEGIN textures
		imgui.BeginChildV("Options", imgui.Vec2{X: comp.widthTexturesPanel, Y: 0}, false, 0)

		imgui.Text("Examples")

		if imgui.ButtonV("Artificial", imgui.Vec2{X: -1.0, Y: 0.0}) {
			comp.openExample(sett.App.AppFolder + "/shaders/shadertoy/4ljGW1.stoy")
		}
		if imgui.IsItemHovered() {
			imgui.SetTooltip("Artificial")
		}

		if imgui.ButtonV("Combustible\nVoronoi Layers", imgui.Vec2{X: -1.0, Y: 0.0}) {
			comp.openExample(sett.App.AppFolder + "/shaders/shadertoy/4tlSzl.stoy")
		}
		if imgui.IsItemHovered() {
			imgui.SetTooltip("Combustible Voronoi Layers")
		}

		if imgui.ButtonV("Seascape", imgui.Vec2{X: -1.0, Y: 0.0}) {
			comp.openExample(sett.App.AppFolder + "/shaders/shadertoy/Ms2SD1.stoy")
		}
		if imgui.IsItemHovered() {
			imgui.SetTooltip("Seascape")
		}

		if imgui.ButtonV("Star Nest", imgui.Vec2{X: -1.0, Y: 0.0}) {
			comp.openExample(sett.App.AppFolder + "/shaders/shadertoy/XlfGRj.stoy")
		}
		if imgui.IsItemHovered() {
			imgui.SetTooltip("Star Nest")
		}

		if imgui.ButtonV("Sun Surface", imgui.Vec2{X: -1.0, Y: 0.0}) {
			comp.openExample(sett.App.AppFolder + "/shaders/shadertoy/XlSSzK.stoy")
		}
		if imgui.IsItemHovered() {
			imgui.SetTooltip("Sun Surface")
		}

		imgui.Separator()

		textureImages := []string{
			" -- NONE -- ",
			"tex00.jpg", "tex01.jpg", "tex02.jpg", "tex03.jpg", "tex04.jpg", "tex05.jpg", "tex06.jpg", "tex07.jpg", "tex08.jpg", "tex09.jpg", "tex10.jpg",
			"tex11.jpg", "tex12.jpg", "tex13.jpg", "tex14.jpg", "tex15.jpg", "tex16.jpg", "tex17.jpg", "tex18.jpg", "tex19.jpg", "tex20.jpg"}

		cubemapImages := []string{
			" -- NONE -- ",
			"cube00_0.jpg", "cube00_1.jpg", "cube00_2.jpg", "cube00_3.jpg",
			"cube00_4.jpg", "cube00_5.jpg", "cube01_0.png", "cube01_1.png", "cube01_2.png", "cube01_3.png", "cube01_4.png", "cube01_5.png",
			"cube02_0.jpg", "cube02_1.jpg", "cube02_2.jpg", "cube02_3.jpg", "cube02_4.jpg", "cube02_5.jpg", "cube03_0.png", "cube03_1.png",
			"cube03_2.png", "cube03_3.png", "cube03_4.png", "cube03_5.png", "cube04_0.png", "cube04_1.png", "cube04_2.png", "cube04_3.png",
			"cube04_4.png", "cube04_5.png", "cube05_0.png", "cube05_1.png", "cube05_2.png", "cube05_3.png", "cube05_4.png", "cube05_5.png"}

		imgui.PushItemWidth(-1)
		imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{X: 0, Y: 6})
		imgui.PushStyleVarVec2(imgui.StyleVarWindowMinSize, imgui.Vec2{X: 0, Y: 100})

		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: .0, Z: .0, W: 1})
		imgui.Text("Channel #0")
		imgui.PopStyleColor()
		imgui.Checkbox("Cubemap?##001", &comp.channel0Cube)
		if comp.channel0Cube {
			imgui.ListBox("##cubemap0", &comp.cubemapImage0, cubemapImages)
		} else {
			imgui.ListBox("##texImage0", &comp.texImage0, textureImages)
		}

		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: .0, Z: .0, W: 1})
		imgui.Text("Channel #1")
		imgui.PopStyleColor()
		imgui.Checkbox("Cubemap?##002", &comp.channel1Cube)
		if comp.channel1Cube {
			imgui.ListBox("##cubemap1", &comp.cubemapImage1, cubemapImages)
		} else {
			imgui.ListBox("##texImage1", &comp.texImage1, textureImages)
		}

		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: .0, Z: .0, W: 1})
		imgui.Text("Channel #2")
		imgui.PopStyleColor()
		imgui.Checkbox("Cubemap?##003", &comp.channel2Cube)
		if comp.channel2Cube {
			imgui.ListBox("##cubemap2", &comp.cubemapImage2, cubemapImages)
		} else {
			imgui.ListBox("##texImage2", &comp.texImage2, textureImages)
		}

		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1, Y: .0, Z: .0, W: 1})
		imgui.Text("Channel #3")
		imgui.PopStyleColor()
		imgui.Checkbox("Cubemap?##004", &comp.channel3Cube)
		if comp.channel3Cube {
			imgui.ListBox("##cubemap3", &comp.cubemapImage3, cubemapImages)
		} else {
			imgui.ListBox("##texImage3", &comp.texImage3, textureImages)
		}

		imgui.PopStyleVarV(2)

		if comp.texImage0 > 0 {
			comp.engineShadertoy.IChannel0CubeImage = ""
			if comp.texImage0 < 10 {
				comp.engineShadertoy.IChannel0Image = sett.App.AppFolder + "tex0" + string(comp.texImage0) + ".jpg"
			} else {
				comp.engineShadertoy.IChannel0Image = sett.App.AppFolder + "tex" + string(comp.texImage0) + ".jpg"
			}
			comp.texImage0 = 0
		}
		if comp.cubemapImage0 > 0 {
			comp.engineShadertoy.IChannel0Image = ""
			comp.engineShadertoy.IChannel0CubeImage = sett.App.AppFolder + cubemapImages[comp.cubemapImage0]
			comp.cubemapImage0 = 0
		}

		if comp.texImage1 > 0 {
			comp.engineShadertoy.IChannel1CubeImage = ""
			if comp.texImage1 < 10 {
				comp.engineShadertoy.IChannel1Image = sett.App.AppFolder + "tex0" + string(comp.texImage1) + ".jpg"
			} else {
				comp.engineShadertoy.IChannel1Image = sett.App.AppFolder + "tex" + string(comp.texImage1) + ".jpg"
			}
			comp.texImage1 = 0
		}
		if comp.cubemapImage1 > 0 {
			comp.engineShadertoy.IChannel1Image = ""
			comp.engineShadertoy.IChannel1CubeImage = sett.App.AppFolder + cubemapImages[comp.cubemapImage1]
			comp.cubemapImage1 = 0
		}

		if comp.texImage2 > 0 {
			comp.engineShadertoy.IChannel2CubeImage = ""
			if comp.texImage2 < 10 {
				comp.engineShadertoy.IChannel2Image = sett.App.AppFolder + "tex0" + string(comp.texImage2) + ".jpg"
			} else {
				comp.engineShadertoy.IChannel2Image = sett.App.AppFolder + "tex" + string(comp.texImage2) + ".jpg"
			}
			comp.texImage2 = 0
		}
		if comp.cubemapImage2 > 0 {
			comp.engineShadertoy.IChannel2Image = ""
			comp.engineShadertoy.IChannel2CubeImage = sett.App.AppFolder + cubemapImages[comp.cubemapImage2]
			comp.cubemapImage2 = 0
		}

		if comp.texImage3 > 0 {
			comp.engineShadertoy.IChannel3CubeImage = ""
			if comp.texImage3 < 10 {
				comp.engineShadertoy.IChannel3Image = sett.App.AppFolder + "tex0" + string(comp.texImage3) + ".jpg"
			} else {
				comp.engineShadertoy.IChannel3Image = sett.App.AppFolder + "tex" + string(comp.texImage3) + ".jpg"
			}
			comp.texImage3 = 0
		}
		if comp.cubemapImage3 > 0 {
			comp.engineShadertoy.IChannel3Image = ""
			comp.engineShadertoy.IChannel3CubeImage = sett.App.AppFolder + cubemapImages[comp.cubemapImage3]
			comp.cubemapImage3 = 0
		}

		imgui.EndChild()

		imgui.SameLine()
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 89.0 / 255.0, Y: 91.0 / 255.0, Z: 94 / 255.0, W: 1})
		imgui.PushStyleColor(imgui.StyleColorButtonHovered, imgui.Vec4{X: 119.0 / 255.0, Y: 122.0 / 255.0, Z: 124.0 / 255.0, W: 1})
		imgui.PushStyleColor(imgui.StyleColorButtonActive, imgui.Vec4{X: .0, Y: .0, Z: .0, W: 1})
		imgui.ButtonV("###splitterGUI2", imgui.Vec2{X: 4.0, Y: -1})
		imgui.PopStyleColorV(3)
		if imgui.IsItemActive() {
			comp.widthTexturesPanel += imgui.CurrentIO().MouseDelta().X
		}
		if imgui.IsItemHovered() {
			imgui.SetMouseCursor(imgui.MouseCursorResizeNS)
		} else {
			imgui.SetMouseCursor(imgui.MouseCursorNone)
		}
		imgui.SameLine()

		// BEGIN IDE
		imgui.BeginChildV("IDE", imgui.Vec2{X: 0.0, Y: 0.0}, false, 0)
		lines := (imgui.WindowHeight() - 4.0) / imgui.TextLineHeight()
		imgui.InputTextMultilineV("##source", &comp.shadertoyEditorText, imgui.Vec2{X: -1.0, Y: imgui.TextLineHeight() * lines}, 0, nil)
		imgui.EndChild()

		// END editor
		imgui.EndChild()

		if imgui.IsWindowHovered() { // TODO: && !imgui.IsAnyItemActive() && imgui.IsMouseDragging(2, 0.0) {
			comp.scrolling = comp.scrolling.Sub(mgl32.Vec2{imgui.CurrentIO().MouseDelta().X, imgui.CurrentIO().MouseDelta().Y})
		}

		imgui.End()
	}
}

func (comp *ComponentShadertoy) compileShader() {
	comp.engineShadertoy = objects.InitShadertoy(comp.window)
	comp.engineShadertoy.InitShaderProgram(comp.shadertoyEditorText)
	comp.engineShadertoy.InitBuffers()
	comp.engineShadertoy.InitFBO(int32(comp.windowWidth), int32(comp.heightTopPanel), &comp.vboTexture)
}

func (comp *ComponentShadertoy) openExample(fileName string) {
	comp.shadertoyEditorText = settings.ReadFile(fileName, false)
}

func (comp *ComponentShadertoy) getFromClipboard() {
	comp.shadertoyEditorText, _ = comp.window.ClipboardText()
}
