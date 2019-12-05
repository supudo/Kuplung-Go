package engine

import (
	"time"

	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/engine/input"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/veandco/go-sdl2/sdl"
)

// KuplungWindow ...
type KuplungWindow struct {
	WindowEventDispatcher

	sdlWindow *sdl.Window
	glContext sdl.GLContext
	glWrapper *OpenGL

	framesPerSecond float64
	frameTime       time.Duration
	nextRenderTick  time.Time

	buttonsDown [3]bool
}

// NewKuplungWindow ...
func NewKuplungWindow(title string) (window *KuplungWindow) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		settings.LogError("[initSDL] Failed to create window: %v", err)
	}

	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 4)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG)
	_ = sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	_ = sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	_ = sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)
	_ = sdl.GLSetAttribute(sdl.GL_STENCIL_SIZE, 8)
	_ = sdl.GLSetAttribute(sdl.GL_ACCELERATED_VISUAL, 1)
	_ = sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
	_ = sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 4)

	_ = sdl.SetHint(sdl.HINT_MAC_CTRL_CLICK_EMULATE_RIGHT_CLICK, "1")
	_ = sdl.SetHint(sdl.HINT_VIDEO_HIGHDPI_DISABLED, "0")

	var sett = settings.GetSettings()
	wWidth, wHeight := int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight)
	win, err := sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, wWidth, wHeight, sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI|sdl.WINDOW_RESIZABLE)
	if err != nil {
		sdl.Quit()
		settings.LogError("[initSDL] Failed to create window: %v", err)
	}

	glContext, err := win.GLCreateContext()
	if err != nil {
		settings.LogError("[initSDL] Failed to create OpenGL context: %v", err)
	}

	err = win.GLMakeCurrent(glContext)
	if err != nil {
		settings.LogError("[initSDL] Failed to set current OpenGL context: %v", err)
	}

	err = sdl.GLSetSwapInterval(1)
	if err != nil {
		settings.LogError("[initSDL] Failed to set swap interval: %v", err)
	}

	fps := sett.Rendering.FramesPerSecond
	window = &KuplungWindow{
		WindowEventDispatcher: NullWindowEventDispatcher(),
		sdlWindow:             win,
		glContext:             glContext,
		glWrapper:             NewOpenGL(),
		framesPerSecond:       fps,
		frameTime:             time.Duration(int64(float64(time.Second) / fps)),
		nextRenderTick:        time.Now(),
	}

	window.OnClosed(window.onClosed)

	return
}

func (window *KuplungWindow) processEvent(event sdl.Event) {
	sett := settings.GetSettings()
	if sett.MemSettings.QuitApplication {
		return
	}
	io := imgui.CurrentIO()

	x, y, state := sdl.GetMouseState()
	io.SetMousePosition(imgui.Vec2{X: float32(x), Y: float32(y)})
	for i, button := range []uint32{sdl.BUTTON_LEFT, sdl.BUTTON_RIGHT, sdl.BUTTON_MIDDLE} {
		io.SetMouseButtonDown(i, window.buttonsDown[i] || (state&sdl.Button(button)) != 0)
		window.buttonsDown[i] = false
		window.CallOnMouseButtonDown(uint32(i), input.ModNone)
	}

	switch ev := event.(type) {
	case *sdl.QuitEvent:
		sett.MemSettings.QuitApplication = true
		window.CallClosed()
	case *sdl.MouseWheelEvent:
		var deltaX, deltaY float32
		if ev.X > 0 {
			deltaX++
		} else if ev.X < 0 {
			deltaX--
		}
		if ev.Y > 0 {
			deltaY++
			sett.MemSettings.ZoomFactor -= 4.0
		} else if ev.Y < 0 {
			deltaY--
			sett.MemSettings.ZoomFactor += 4.0
		}
		io.AddMouseWheelDelta(deltaX, deltaY)
	case *sdl.TextInputEvent:
		io.AddInputCharacters(string(ev.Text[:]))
	case *sdl.KeyboardEvent:
		switch ev.Type {
		case sdl.KEYDOWN:
			io.KeyPress(int(ev.Keysym.Scancode))
			window.updateKeyModifier()
		case sdl.KEYUP:
			io.KeyRelease(int(ev.Keysym.Scancode))
			window.updateKeyModifier()
		}
	case *sdl.WindowEvent:
		switch ev.Event {
		case sdl.WINDOWEVENT_RESIZED:
			width, height := ev.Data1, ev.Data2
			sett.AppWindow.SDLWindowWidth = float32(width)
			sett.AppWindow.SDLWindowHeight = float32(height)
			io.SetDisplaySize(imgui.Vec2{X: float32(width), Y: float32(height)})
			window.CallResize(int(width), int(height))
		case sdl.WINDOWEVENT_CLOSE:
			sett.MemSettings.QuitApplication = true
			window.CallClosed()
		}
	}
}

// func (window *KuplungWindow) processEvent2(event sdl.Event) {
// 	var sett = settings.GetSettings()

// 	if sett.MemSettings.QuitApplication {
// 		return
// 	}

// 	io := imgui.CurrentIO()

// 	x, y, state := sdl.GetMouseState()
// 	io.SetMousePosition(imgui.Vec2{X: float32(x), Y: float32(y)})
// 	for i, button := range []uint32{sdl.BUTTON_LEFT, sdl.BUTTON_RIGHT, sdl.BUTTON_MIDDLE} {
// 		io.SetMouseButtonDown(i, window.buttonsDown[i] || (state&sdl.Button(button)) != 0)
// 		window.buttonsDown[i] = false
// 		window.CallOnMouseButtonDown(uint32(i), input.ModNone)
// 	}

// 	switch event.GetType() {
// 	case sdl.QUIT:
// 		sett.MemSettings.QuitApplication = true
// 		window.CallClosed()
// 	case sdl.MOUSEWHEEL:
// 		wheelEvent := event.(*sdl.MouseWheelEvent)
// 		var deltaX, deltaY float32
// 		if wheelEvent.X > 0 {
// 			deltaX++
// 		} else if wheelEvent.X < 0 {
// 			deltaX--
// 		}
// 		if wheelEvent.Y > 0 {
// 			deltaY++
// 		} else if wheelEvent.Y < 0 {
// 			deltaY--
// 			sett.MemSettings.ZoomFactor += 4.0
// 		}
// 		io.AddMouseWheelDelta(deltaX, deltaY)
// 	case sdl.TEXTINPUT:
// 		inputEvent := event.(*sdl.TextInputEvent)
// 		io.AddInputCharacters(string(inputEvent.Text[:]))
// 	case sdl.KEYDOWN:
// 		keyEvent := event.(*sdl.KeyboardEvent)
// 		io.KeyPress(int(keyEvent.Keysym.Scancode))
// 		window.updateKeyModifier()
// 	case sdl.KEYUP:
// 		keyEvent := event.(*sdl.KeyboardEvent)
// 		io.KeyRelease(int(keyEvent.Keysym.Scancode))
// 		window.updateKeyModifier()
// 	case sdl.WINDOWEVENT_RESIZED:
// 		windowEvent := event.(*sdl.WindowEvent)
// 		width, height := windowEvent.Data1, windowEvent.Data2
// 		sett.AppWindow.SDLWindowWidth = width
// 		sett.AppWindow.SDLWindowHeight = height
// 		io.SetDisplaySize(imgui.Vec2{X: float32(width), Y: float32(height)})
// 		window.CallResize(int(width), int(height))
// 	case sdl.WINDOWEVENT_CLOSE:
// 		sett.MemSettings.QuitApplication = true
// 		window.CallClosed()
// 	}
// }

func (window *KuplungWindow) updateKeyModifier() {
	modState := sdl.GetModState()
	mapModifier := func(lMask sdl.Keymod, lKey int, rMask sdl.Keymod, rKey int) (lResult int, rResult int) {
		if (modState & lMask) != 0 {
			lResult = lKey
		}
		if (modState & rMask) != 0 {
			rResult = rKey
		}
		return
	}
	io := imgui.CurrentIO()
	io.KeyShift(mapModifier(sdl.KMOD_LSHIFT, sdl.SCANCODE_LSHIFT, sdl.KMOD_RSHIFT, sdl.SCANCODE_RSHIFT))
	io.KeyCtrl(mapModifier(sdl.KMOD_LCTRL, sdl.SCANCODE_LCTRL, sdl.KMOD_RCTRL, sdl.SCANCODE_RCTRL))
	io.KeyAlt(mapModifier(sdl.KMOD_LALT, sdl.SCANCODE_LALT, sdl.KMOD_RALT, sdl.SCANCODE_RALT))
}

// ShouldClose returns true if the user requested the window to close.
func (window *KuplungWindow) ShouldClose() bool {
	var sett = settings.GetSettings()
	return sett.MemSettings.QuitApplication
}

// Close closes the window and releases its resources.
func (window *KuplungWindow) Close() {
	window.sdlWindow.Destroy()
	sdl.Quit()
}

// ClipboardText returns the current value of the clipboard, if it is compatible with UTF-8.
func (window KuplungWindow) ClipboardText() (string, error) {
	return sdl.GetClipboardText()
}

// SetClipboardText sets the current value of the clipboard as UTF-8 string.
func (window KuplungWindow) SetClipboardText(value string) {
	sdl.SetClipboardText(value)
}

// Update must be called from within the main thread as often as possible.
func (window *KuplungWindow) Update() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		window.processEvent(event)
	}

	now := time.Now()
	delta := now.Sub(window.nextRenderTick)
	if delta.Nanoseconds() < 0 {
		// detected a change of wallclock time into the past; realign
		delta = window.frameTime
		window.nextRenderTick = now
	}

	if delta.Nanoseconds() >= window.frameTime.Nanoseconds() {
		window.sdlWindow.GLMakeCurrent(window.glContext)
		window.CallRender()
		window.sdlWindow.GLSwap()
		framesCovered := delta.Nanoseconds() / window.frameTime.Nanoseconds()
		window.nextRenderTick = window.nextRenderTick.Add(time.Duration(framesCovered * window.frameTime.Nanoseconds()))
	}
}

// OpenGL returns the OpenGL API.
func (window *KuplungWindow) OpenGL() interfaces.OpenGL {
	return window.glWrapper
}

// Size returns the dimension of the frame buffer of this window.
func (window *KuplungWindow) Size() (width int, height int) {
	w, h := window.sdlWindow.GLGetDrawableSize()
	if w == 0 || h == 0 {
		settings.LogError("[Size] Window size is 0 : %v x %v", w, h)
	}
	return int(w), int(h)
}

// SetFullScreen toggles the windowed mode.
func (window *KuplungWindow) SetFullScreen(on bool) {
	if on {
		_ = window.sdlWindow.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	} else {
		var sett = settings.GetSettings()
		window.sdlWindow.SetSize(int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight))
	}
}

func (window *KuplungWindow) onClosed() {
	window.CallClosed()
}
