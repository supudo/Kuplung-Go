package rendering

import (
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/objects"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// RenderManager is the main structure for rendering
type RenderManager struct {
	window interfaces.Window

	cube *objects.Cube

	Settings types.RenderSettings
}

// NewRenderManager will return an instance of the rendering manager
func NewRenderManager(window interfaces.Window) *RenderManager {
	rm := &RenderManager{}
	rm.window = window
	sett := settings.GetSettings()

	rm.Settings.GLSLVersion = "#version 410"

	rm.Settings.Fov = sett.MemSettings.ZoomFactor
	rm.Settings.RatioWidth = 4.0
	rm.Settings.RatioHeight = 3.0
	rm.Settings.PlaneClose = 1.0
	rm.Settings.PlaneFar = 1000.0

	rm.initCube()

	return rm
}

// Render handles rendering of all scene objects
func (rm *RenderManager) Render(gvars types.ObjectVariables) {
	sett := settings.GetSettings()
	rm.Settings.Fov = sett.MemSettings.ZoomFactor
	if gvars.ShowCube {
		rm.cube.Render(rm.Settings)
	}
}

// Dispose will cleanup everything
func (rm *RenderManager) Dispose() {
	rm.cube.Dispose()
}

func (rm *RenderManager) initCube() {
	rm.cube = objects.CubeInit(rm.window, rm.Settings)
}
