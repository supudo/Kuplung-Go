package settings

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

// RenderingSettings ...
type RenderingSettings struct {
	MatrixProjection mgl32.Mat4
	MatrixCamera     mgl32.Mat4

	ShowCube bool

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

	rSettings.MatrixProjection = mgl32.Perspective(mgl32.DegToRad(rSettings.Fov), rSettings.RatioWidth/rSettings.RatioHeight, rSettings.PlaneClose, rSettings.PlaneFar)
	rSettings.MatrixCamera = mgl32.Ident4()

	return rSettings
}
