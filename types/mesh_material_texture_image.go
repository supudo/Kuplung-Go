package types

// MeshMaterialTextureImage ...
type MeshMaterialTextureImage struct {
	Width  int32
	Height int32

	UseTexture bool

	Filename string
	Image    string

	Commands []string
}
