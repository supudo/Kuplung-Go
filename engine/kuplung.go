package engine

import (
	"github.com/supudo/Kuplung-Go/gui"
	"github.com/veandco/go-sdl2/sdl"
)

// Kuplung is the main application
type Kuplung struct {
	window *KuplungWindow
	gui    *gui.GUI
}

// NewKuplung ...
func (k *Kuplung) NewKuplung(window *KuplungWindow) {
	k.window = window
	k.gui = gui.NewGUI(k.window.Window)

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
