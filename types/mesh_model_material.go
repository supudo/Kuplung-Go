package types

import "github.com/go-gl/mathgl/mgl32"

// MeshModelMaterial ...
type MeshModelMaterial struct {
	MaterialID    uint32
	MaterialTitle string

	SpecularExp float32

	AmbientColor  mgl32.Vec3
	DiffuseColor  mgl32.Vec3
	SpecularColor mgl32.Vec3
	EmissionColor mgl32.Vec3

	Transparency     float32
	IlluminationMode uint32
	OpticalDensity   float32

	TextureAmbient      MeshMaterialTextureImage
	TextureDiffuse      MeshMaterialTextureImage
	TextureSpecular     MeshMaterialTextureImage
	TextureSpecularExp  MeshMaterialTextureImage
	TextureDissolve     MeshMaterialTextureImage
	TextureBump         MeshMaterialTextureImage
	TextureDisplacement MeshMaterialTextureImage
}
