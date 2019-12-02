package engine

import (
	"time"

	"github.com/inkyblackness/imgui-go"
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
}

// NewKuplungWindow ...
func NewKuplungWindow(title string) *KuplungWindow {
	var sett = settings.GetSettings()
	fps := sett.Rendering.FramesPerSecond
	w, g := initSDL()
	window := &KuplungWindow{
		WindowEventDispatcher: NullWindowEventDispatcher(),
		sdlWindow:             w,
		glContext:             g,
		glWrapper:             NewOpenGL(),
		framesPerSecond:       fps,
		frameTime:             time.Duration(int64(float64(time.Second) / fps)),
		nextRenderTick:        time.Now(),
	}
	//window.SetKeyMapping()
	return window
}

func initSDL() (sdlWindow *sdl.Window, glContext sdl.GLContext) {
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
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

	sett := settings.GetSettings()
	window, err := sdl.CreateWindow("Kuplung "+sett.App.ApplicationVersion, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight), sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN)
	if err != nil {
		sdl.Quit()
		settings.LogError("[initSDL] Failed to create window: %v", err)
	}

	glContext, err = window.GLCreateContext()
	if err != nil {
		window.Destroy()
		sdl.Quit()
		settings.LogError("[initSDL] Failed to create OpenGL context: %v", err)
	}

	err = window.GLMakeCurrent(glContext)
	if err != nil {
		sdl.GLDeleteContext(glContext)
		window.Destroy()
		sdl.Quit()
		settings.LogError("[initSDL] Failed to set current OpenGL context: %v", err)
	}

	err = sdl.GLSetSwapInterval(1)
	if err != nil {
		sdl.GLDeleteContext(glContext)
		window.Destroy()
		sdl.Quit()
		settings.LogError("[initSDL] Failed to set swap interval: %v", err)
	}

	return sdlWindow, glContext
}

func (window *KuplungWindow) processEvent(event sdl.Event) {
	var sett = settings.GetSettings()
	io := imgui.CurrentIO()
	switch event.GetType() {
	case sdl.QUIT:
		sett.MemSettings.QuitApplication = true
	case sdl.MOUSEWHEEL:
		wheelEvent := event.(*sdl.MouseWheelEvent)
		var deltaX, deltaY float32
		if wheelEvent.X > 0 {
			deltaX++
		} else if wheelEvent.X < 0 {
			deltaX--
		}
		if wheelEvent.Y > 0 {
			deltaY++
		} else if wheelEvent.Y < 0 {
			deltaY--
		}
		io.AddMouseWheelDelta(deltaX, deltaY)
	case sdl.MOUSEBUTTONDOWN:
		buttonEvent := event.(*sdl.MouseButtonEvent)
		switch buttonEvent.Button {
		case sdl.BUTTON_LEFT:
			//platform.buttonsDown[0] = true
		case sdl.BUTTON_RIGHT:
			//platform.buttonsDown[1] = true
		case sdl.BUTTON_MIDDLE:
			//platform.buttonsDown[2] = true
		}
	case sdl.TEXTINPUT:
		inputEvent := event.(*sdl.TextInputEvent)
		io.AddInputCharacters(string(inputEvent.Text[:]))
	case sdl.KEYDOWN:
		keyEvent := event.(*sdl.KeyboardEvent)
		io.KeyPress(int(keyEvent.Keysym.Scancode))
		window.updateKeyModifier()
	case sdl.KEYUP:
		keyEvent := event.(*sdl.KeyboardEvent)
		io.KeyRelease(int(keyEvent.Keysym.Scancode))
		window.updateKeyModifier()
	}
}

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
	settings.LogWarn("%v ---- %v", w, h)
	return int(w), int(h)
}

// SetFullScreen toggles the windowed mode.
func (window *KuplungWindow) SetFullScreen(on bool) {
	if on {
		_ = window.sdlWindow.SetFullscreen(sdl.WINDOW_FULLSCREEN)
	} else {
		var sett = settings.GetSettings()
		window.sdlWindow.SetSize(sett.AppWindow.SDLWindowWidth, sett.AppWindow.SDLWindowHeight)
	}
}
