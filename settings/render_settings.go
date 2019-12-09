package settings

import "sync"

// RenderingSettings ...
type RenderingSettings struct {
	ShowCube bool

	ZoomFactor float32

	Fov         float32
	RatioWidth  float32
	RatioHeight float32
	PlaneClose  float32
	PlaneFar    float32

	GammaCoeficient float32

	ShowAxisHelpers bool
	ShowZAxis       bool

	WorldGridSizeSquares    int32
	WorldGridFixedWithWorld bool
	UseWorldGrid            bool
	ShowGrid                bool
	ActAsMirror             bool
}

var instantiatedRendering *RenderingSettings
var onceRendering sync.Once

// GetRenderingSettings singleton for our application settings
func GetRenderingSettings() *RenderingSettings {
	onceRendering.Do(func() {
		as := InitRenderingSettings()
		instantiatedRendering = &as
	})
	return instantiatedRendering
}

// InitRenderingSettings will initialize application settings
func InitRenderingSettings() RenderingSettings {
	var rSettings RenderingSettings

	rSettings.ShowCube = false

	rSettings.ZoomFactor = 1.0

	rSettings.Fov = 45.0
	rSettings.RatioWidth = 4.0
	rSettings.RatioHeight = 3.0
	rSettings.PlaneClose = 1.0
	rSettings.PlaneFar = 1000.0

	rSettings.GammaCoeficient = 1.0

	rSettings.ShowAxisHelpers = true
	rSettings.ShowZAxis = true

	rSettings.WorldGridSizeSquares = 10
	rSettings.WorldGridFixedWithWorld = true
	rSettings.UseWorldGrid = true
	rSettings.ShowGrid = true
	rSettings.ActAsMirror = false

	return rSettings
}
