package engine

import "github.com/supudo/Kuplung-Go/interfaces"

// ClipboardAdapter ...
type ClipboardAdapter struct {
	Window interfaces.Window
}

// ClipboardText ...
func (adapter ClipboardAdapter) ClipboardText() (string, error) {
	return adapter.Window.ClipboardText()
}

// SetClipboardText ...
func (adapter ClipboardAdapter) SetClipboardText(value string) {
	adapter.Window.SetClipboardText(value)
}
