package app

import (
	"time"

	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/gui"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
)

// KuplungApp ...
type KuplungApp struct {
	Version string

	window     interfaces.Window
	gl         interfaces.OpenGL
	clipboard  engine.ClipboardAdapter
	guiContext *gui.Context

	cube engine.Cube

	FontFile string
	FontSize float32
	GuiScale float32
}

// InitializeKuplungWindow ...
func (kapp *KuplungApp) InitializeKuplungWindow(window interfaces.Window) {
	kapp.window = window
	kapp.clipboard.Window = window
	kapp.gl = window.OpenGL()

	kapp.initWindowCallbacks()
	kapp.initOpenGL()
	kapp.initGui()
	kapp.initCube()
}

func (kapp *KuplungApp) initWindowCallbacks() {
	kapp.window.OnClosed(kapp.onWindowClosed)
	kapp.window.OnRender(kapp.render)
}

func (kapp *KuplungApp) onWindowClosed() {
	if kapp.guiContext != nil {
		kapp.guiContext.Destroy()
		kapp.guiContext = nil
	}
}

func (kapp *KuplungApp) render() {
	kapp.guiContext.NewFrame()
	kapp.gl.Clear(engine.COLOR_BUFFER_BIT)
	kapp.guiContext.DrawMainMenu()

	//kapp.cube.CubeRender()

	kapp.guiContext.Render()
	// sleep to avoid 100% CPU usage for this demo
	<-time.After(time.Millisecond * 25)
}

func (kapp *KuplungApp) initCube() {
	kapp.cube = *engine.CubeInit()
}

func (kapp *KuplungApp) initOpenGL() {
	kapp.gl.Enable(engine.DEPTH_TEST)
	kapp.gl.Enable(engine.BLEND)
	kapp.gl.BlendFunc(engine.SRC_ALPHA, engine.ONE_MINUS_SRC_ALPHA)
	sett := settings.GetSettings()
	kapp.gl.ClearColor(sett.AppGui.GUIClearColor[0], sett.AppGui.GUIClearColor[1], sett.AppGui.GUIClearColor[2], sett.AppGui.GUIClearColor[3])
}

func (kapp *KuplungApp) initGui() {
	kapp.initGuiSizes()
	param := gui.ContextParameters{
		FontFile: kapp.FontFile,
		FontSize: kapp.FontSize,
	}
	kapp.guiContext = gui.NewContext(kapp.window, param)
	kapp.initGuiStyle()
}

func (kapp *KuplungApp) initGuiSizes() {
	if kapp.GuiScale < 0.5 {
		kapp.GuiScale = 1.0
	} else if kapp.GuiScale > 10.0 {
		kapp.GuiScale = 10.0
	}

	if kapp.FontSize <= 0 {
		kapp.FontSize = 16.0
	}

	kapp.FontSize *= kapp.GuiScale
}

func (kapp *KuplungApp) initGuiStyle() {
	// if len(kapp.FontFile) == 0 {
	// 	imgui.CurrentIO().SetFontGlobalScale(kapp.GuiScale)
	// }
	// imgui.CurrentStyle().ScaleAllSizes(kapp.GuiScale)

	// color := func(r, g, b byte, alpha float32) imgui.Vec4 {
	// 	return imgui.Vec4{X: float32(r) / 255.0, Y: float32(g) / 255.0, Z: float32(b) / 255.0, W: alpha}
	// }
	// colorDoubleFull := func(alpha float32) imgui.Vec4 { return color(0xC4, 0x38, 0x9F, alpha) }
	// colorDoubleDark := func(alpha float32) imgui.Vec4 { return color(0x31, 0x01, 0x38, alpha) }

	// colorTripleFull := func(alpha float32) imgui.Vec4 { return color(0x21, 0xFF, 0x43, alpha) }
	// colorTripleDark := func(alpha float32) imgui.Vec4 { return color(0x06, 0xCC, 0x94, alpha) }
	// colorTripleLight := func(alpha float32) imgui.Vec4 { return color(0x51, 0x99, 0x58, alpha) }

	// style := imgui.CurrentStyle()
	// style.SetColor(imgui.StyleColorText, colorTripleFull(1.0))
	// style.SetColor(imgui.StyleColorTextDisabled, colorTripleDark(1.0))

	// style.SetColor(imgui.StyleColorWindowBg, colorDoubleDark(0.80))
	// style.SetColor(imgui.StyleColorPopupBg, colorDoubleDark(0.75))

	// style.SetColor(imgui.StyleColorTitleBgActive, colorTripleLight(1.0))
	// style.SetColor(imgui.StyleColorFrameBg, colorTripleLight(0.54))

	// style.SetColor(imgui.StyleColorFrameBgHovered, colorTripleDark(0.4))
	// style.SetColor(imgui.StyleColorFrameBgActive, colorTripleDark(0.67))
	// style.SetColor(imgui.StyleColorCheckMark, colorTripleDark(1.0))
	// style.SetColor(imgui.StyleColorSliderGrabActive, colorTripleDark(1.0))
	// style.SetColor(imgui.StyleColorButton, colorTripleDark(0.4))
	// style.SetColor(imgui.StyleColorButtonHovered, colorTripleDark(1.0))
	// style.SetColor(imgui.StyleColorHeader, colorTripleLight(0.70))
	// style.SetColor(imgui.StyleColorHeaderHovered, colorTripleDark(0.8))
	// style.SetColor(imgui.StyleColorHeaderActive, colorTripleDark(1.0))
	// style.SetColor(imgui.StyleColorResizeGrip, colorTripleDark(0.25))
	// style.SetColor(imgui.StyleColorResizeGripHovered, colorTripleDark(0.67))
	// style.SetColor(imgui.StyleColorResizeGripActive, colorTripleDark(0.95))
	// style.SetColor(imgui.StyleColorTextSelectedBg, colorTripleDark(0.35))

	// style.SetColor(imgui.StyleColorSliderGrab, colorDoubleFull(1.0))
	// style.SetColor(imgui.StyleColorButtonActive, colorDoubleFull(1.0))
	// style.SetColor(imgui.StyleColorSeparatorHovered, colorDoubleFull(0.78))
	// style.SetColor(imgui.StyleColorSeparatorActive, colorTripleLight(1.0))
}
