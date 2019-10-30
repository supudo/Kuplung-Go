package gui

import (
	"fmt"
	"os"

	"github.com/inkyblackness/imgui-go"
)

// RenderGUIStart handles GUI rendering
func RenderGUIStart() {
	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	platform, err := NewSDL(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer platform.Dispose()

	renderer, err := NewOpenGL3(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer renderer.Dispose()

	Run(platform, renderer)
}
