package types

import "github.com/go-gl/mathgl/mgl32"

// PackedVertex ...
type PackedVertex struct {
	Position mgl32.Vec3
	UV       mgl32.Vec2
	Normal   mgl32.Vec3
}
