package types

import "github.com/go-gl/mathgl/mgl32"

// MeshModel ...
type MeshModel struct {
	ID       uint32
	File     string
	FilePath string

	ModelTitle    string
	MaterialTitle string

	CountVertices           int32
	CountColors             int32
	CountTextureCoordinates int32
	CountNormals            int32
	CountIndices            int32

	Vertices           []mgl32.Vec3
	Colors             []mgl32.Vec3
	TextureCoordinates []mgl32.Vec2
	Normals            []mgl32.Vec3
	Indices            []uint32

	ModelMaterial MeshModelMaterial
}
