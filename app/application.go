package app

import (
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/input"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/gui"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/rendering"
	"github.com/supudo/Kuplung-Go/settings"
)

// KuplungApp ...
type KuplungApp struct {
	Version string

	window     interfaces.Window
	gl         interfaces.OpenGL
	clipboard  engine.ClipboardAdapter
	guiContext *gui.Context

	lastModifier input.Modifier
	lastMouseX   float32
	lastMouseY   float32

	renderManager *rendering.RenderManager

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
	kapp.initRenderingManager()
}

func (kapp *KuplungApp) initWindowCallbacks() {
	kapp.window.OnClosed(kapp.onWindowClosed)
	kapp.window.OnRender(kapp.render)

	kapp.window.OnMouseMove(kapp.onMouseMove)
	kapp.window.OnMouseScroll(kapp.onMouseScroll)
	kapp.window.OnMouseButtonDown(kapp.onMouseButtonDown)
	kapp.window.OnMouseButtonUp(kapp.onMouseButtonUp)

	kapp.window.OnKey(kapp.onKey)
	kapp.window.OnModifier(kapp.onModifier)
	settings.LogInfo("[Application] Callbacks initialized.")
}

func (kapp *KuplungApp) render() {
	sett := settings.GetSettings()
	if !sett.MemSettings.QuitApplication {
		currentTime := kapp.window.GetTicks()
		sett.MemSettings.NbFrames++
		if currentTime-sett.MemSettings.NbLastTime >= 1000.0 {
			sett.MemSettings.NbResult = 1000.0 / float32(sett.MemSettings.NbFrames)
			sett.MemSettings.NbFrames = 0
			sett.MemSettings.NbLastTime += 1000
		}

		kapp.guiContext.NewFrame()
		kapp.gl.Clear(oglconsts.COLOR_BUFFER_BIT | oglconsts.DEPTH_BUFFER_BIT | oglconsts.STENCIL_BUFFER_BIT)
		kapp.guiContext.DrawGUI(true, kapp.renderManager)

		kapp.renderManager.Render()

		kapp.guiContext.Render()
	}
}

func (kapp *KuplungApp) initOpenGL() {
	sett := settings.GetSettings()
	kapp.gl.Enable(oglconsts.DEPTH_TEST)
	kapp.gl.Enable(oglconsts.BLEND)
	kapp.gl.BlendFunc(oglconsts.SRC_ALPHA, oglconsts.ONE_MINUS_SRC_ALPHA)
	kapp.gl.ClearColor(sett.AppGui.GUIClearColor[0], sett.AppGui.GUIClearColor[1], sett.AppGui.GUIClearColor[2], sett.AppGui.GUIClearColor[3])
}

func (kapp *KuplungApp) initGui() {
	kapp.initGuiSizes()
	kapp.guiContext = gui.NewContext(kapp.window)
	kapp.initGuiStyle()
	settings.LogInfo("[Application] GUI initialized.")
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
}

func (kapp *KuplungApp) initRenderingManager() {
	kapp.renderManager = rendering.NewRenderManager(kapp.window, kapp.DoProgress)
	settings.LogInfo("[Application] Rendering Manager initialized.")
}

// DoProgress ...
func (kapp *KuplungApp) DoProgress(progress float32) {
	//settings.LogWarn("[Progress] %v", math.Round(float64(progress)*100)/100)
	kapp.guiContext.GuiVars.ParsingPercentage = progress
}

func (kapp *KuplungApp) onWindowClosed() {
	if kapp.guiContext != nil {
		kapp.guiContext.Destroy()
		kapp.guiContext = nil
	}
	if kapp.renderManager != nil {
		kapp.renderManager.Dispose()
	}
}

func (kapp *KuplungApp) onKey(key input.Key, modifier input.Modifier) {
}

func (kapp *KuplungApp) onModifier(modifier input.Modifier) {
}

func (kapp *KuplungApp) onMouseMove(x, y float32) {
	kapp.lastMouseX = x
	kapp.lastMouseY = y
	if !kapp.guiContext.IsUsingMouse() {
	}
}

func (kapp *KuplungApp) onMouseScroll(dx, dy float32) {
	if !kapp.guiContext.IsUsingMouse() {
	}
	kapp.guiContext.MouseScroll(dx, dy)
}

func (kapp *KuplungApp) onMouseButtonDown(btn uint32, modifier input.Modifier) {
	if !kapp.guiContext.IsUsingMouse() {
	}
}

func (kapp *KuplungApp) onMouseButtonUp(btn uint32, modifier input.Modifier) {
	if !kapp.guiContext.IsUsingMouse() {
	}
}
