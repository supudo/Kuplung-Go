package saveopen

import (
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/types"
)

// SOManager ...
type SOManager struct {
	soProtobufs *ProtoBufsSaveOpen

	doProgress func(float32)
}

// NewSaveOpenManager ...
func NewSaveOpenManager(doProgress func(float32)) *SOManager {
	som := &SOManager{}
	som.doProgress = doProgress
	som.initProtobufs()
	return som
}

// Save ...
func (som *SOManager) Save(file *types.FBEntity, meshes []*meshes.ModelFace) {
	som.soProtobufs.Save(file, meshes)
}

// Open ...
func (som *SOManager) Open(file *types.FBEntity) []*meshes.ModelFace {
	return som.soProtobufs.Open(file)
}

func (som *SOManager) initProtobufs() {
	som.soProtobufs = NewProtoBufsSaveOpen(som.doProgress)
}
