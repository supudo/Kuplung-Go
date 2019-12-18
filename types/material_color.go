package types

import "github.com/go-gl/mathgl/mgl32"

// MaterialColor ...
type MaterialColor struct {
	ColorPickerOpen bool
	Animate         bool
	Strength        float32
	Color           mgl32.Vec3
}
