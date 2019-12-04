package engine

import (
	"github.com/supudo/Kuplung-Go/engine/input"
	"github.com/supudo/Kuplung-Go/interfaces"
)

// WindowEventDispatcher ...
type WindowEventDispatcher struct {
	CallClosed            interfaces.ClosedCallback
	CallRender            interfaces.RenderCallback
	CallResize            interfaces.ResizeCallback
	CallOnMouseMove       interfaces.MouseMoveCallback
	CallOnMouseButtonUp   interfaces.MouseButtonCallback
	CallOnMouseButtonDown interfaces.MouseButtonCallback
	CallOnMouseScroll     interfaces.MouseScrollCallback
	CallModifier          interfaces.ModifierCallback
	CallKey               interfaces.KeyCallback
}

// NullWindowEventDispatcher returns an initialized instance with empty callbacks.
func NullWindowEventDispatcher() WindowEventDispatcher {
	return WindowEventDispatcher{
		CallClosed:            func() {},
		CallRender:            func() {},
		CallResize:            func(int, int) {},
		CallOnMouseMove:       func(float32, float32) {},
		CallOnMouseButtonUp:   func(uint32, input.Modifier) {},
		CallOnMouseButtonDown: func(uint32, input.Modifier) {},
		CallOnMouseScroll:     func(float32, float32) {},
		CallKey:               func(input.Key, input.Modifier) {},
		CallModifier:          func(input.Modifier) {},
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

// OnResize implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnResize(callback interfaces.ResizeCallback) {
	window.CallResize = callback
}

// OnMouseMove implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnMouseMove(callback interfaces.MouseMoveCallback) {
	window.CallOnMouseMove = callback
}

// OnMouseButtonDown implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnMouseButtonDown(callback interfaces.MouseButtonCallback) {
	window.CallOnMouseButtonDown = callback
}

// OnMouseButtonUp implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnMouseButtonUp(callback interfaces.MouseButtonCallback) {
	window.CallOnMouseButtonUp = callback
}

// OnMouseScroll implements the WindowEventDispatcher interface.
func (window *WindowEventDispatcher) OnMouseScroll(callback interfaces.MouseScrollCallback) {
	window.CallOnMouseScroll = callback
}

// OnKey implements the WindowEventDispatcher interface
func (window *WindowEventDispatcher) OnKey(callback interfaces.KeyCallback) {
	window.CallKey = callback
}

// OnModifier implements the WindowEventDispatcher interface
func (window *WindowEventDispatcher) OnModifier(callback interfaces.ModifierCallback) {
	window.CallModifier = callback
}
