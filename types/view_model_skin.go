package types

// ViewModelSkin ...
type ViewModelSkin uint32

// ViewModelSkins ...
const (
	ViewModelSkinSolid ViewModelSkin = 0 + iota
	ViewModelSkinMaterial
	ViewModelSkinTexture
	ViewModelSkinWireframe
	ViewModelSkinRendered
)
