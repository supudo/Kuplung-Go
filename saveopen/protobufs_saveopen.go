package saveopen

import (
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/types"
)

// ProtoBufsSaveOpen ...
type ProtoBufsSaveOpen struct {
	doProgress func(float32)
}

// NewProtoBufsSaveOpen ...
func NewProtoBufsSaveOpen(doProgress func(float32)) *ProtoBufsSaveOpen {
	pm := &ProtoBufsSaveOpen{}
	pm.doProgress = doProgress
	return pm
}

// Save ...
func (pm *ProtoBufsSaveOpen) Save(file *types.FBEntity, meshes []*meshes.ModelFace) {
}

// Open ...
func (pm *ProtoBufsSaveOpen) Open(file *types.FBEntity) []*meshes.ModelFace {
	meshes := []*meshes.ModelFace{}
	return meshes
}
