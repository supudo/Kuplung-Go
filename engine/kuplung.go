package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Kuplung is the main application
type Kuplung struct {
	window *KuplungWindow
	gui    *Gui
}

// NewKuplung ...
func (k *Kuplung) NewKuplung(window *KuplungWindow) {
	k.window = window
	k.gui = NewGUI(k)

	cube := CubeInit()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}

		k.gui.UIRenderStart()

		cube.CubeRender()

		k.gui.UIRenderEnd()
	}
}
