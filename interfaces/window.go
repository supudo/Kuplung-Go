package interfaces

// ClosedCallback is the function to clean up resources when the window is being closed.
type ClosedCallback func()

// RenderCallback is the function to receive render events. When the callback
// returns, the window will swap the internal buffer.
type RenderCallback func()

// Window represents an OpenGL render surface.
type Window interface {
	// ClipboardText returns the current value of the clipboard, if it is compatible with UTF-8.
	ClipboardText() (string, error)
	// SetClipboardText sets the current value of the clipboard as UTF-8 string.
	SetClipboardText(value string)

	// OnClosed registers a callback function which shall be called when the window is being closed.
	OnClosed(callback ClosedCallback)

	// OpenGL returns the OpenGL API wrapper for this window.
	OpenGL() OpenGL
	// OnRender registers a callback function which shall be called to update the scene.
	OnRender(callback RenderCallback)

	// Size returns the dimensions of the window display area in pixel.
	Size() (width int32, height int32)
	// SetFullScreen sets the full screen state of the window.
	SetFullScreen(on bool)
}
