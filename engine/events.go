package engine

import (
	"github.com/supudo/Kuplung-Go/engine/input"
	"github.com/supudo/Kuplung-Go/interfaces"
)

// WindowEventDispatcher ...
type WindowEventDispatcher struct {
	CallClosed            interfaces.ClosedCallback
	CallRender            interfaces.RenderCallback
	CallOnMouseMove       interfaces.MouseMoveCallback
	CallOnMouseButtonUp   interfaces.MouseButtonCallback
	CallOnMouseButtonDown interfaces.MouseButtonCallback
	CallOnMouseScroll     interfaces.MouseScrollCallback
	CallModifier          interfaces.ModifierCallback
	CallKey               interfaces.KeyCallback
	CallCharCallback      interfaces.CharCallback
}

// NullWindowEventDispatcher returns an initialized instance with empty callbacks.
func NullWindowEventDispatcher() WindowEventDispatcher {
	return WindowEventDispatcher{
		CallClosed:            func() {},
		CallRender:            func() {},
		CallOnMouseMove:       func(float32, float32) {},
		CallOnMouseButtonUp:   func(uint32, input.Modifier) {},
		CallOnMouseButtonDown: func(uint32, input.Modifier) {},
		CallOnMouseScroll:     func(float32, float32) {},
		CallKey:               func(input.Key, input.Modifier) {},
		CallModifier:          func(input.Modifier) {},
		CallCharCallback:      func(rune) {},
	}
}

// OnClosed implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnClosed(callback interfaces.ClosedCallback) {
	window.CallClosed = callback
}

// OnRender implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnRender(callback interfaces.RenderCallback) {
	window.CallRender = callback
}
