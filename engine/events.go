package engine

import "github.com/supudo/Kuplung-Go/interfaces"

// WindowEventDispatcher ...
type WindowEventDispatcher struct {
	CallClosed interfaces.ClosedCallback
	CallRender interfaces.RenderCallback
}

// OnClosed implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnClosed(callback interfaces.ClosedCallback) {
	window.CallClosed = callback
}

// OnRender implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnRender(callback interfaces.RenderCallback) {
	window.CallRender = callback
}
