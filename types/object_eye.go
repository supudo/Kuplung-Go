package types

import "github.com/go-gl/mathgl/mgl32"

// ObjectEye ...
type ObjectEye struct {
	ViewEye    mgl32.Vec3
	ViewCenter mgl32.Vec3
	ViewUp     mgl32.Vec3
}
