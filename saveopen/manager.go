package saveopen

import (
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/types"
)

// SaveOpenManager ...
type SaveOpenManager struct {
	soProtobufs *ProtoBufsSaveOpen

	doProgress func(float32)
}

// NewSaveOpenManager ...
func NewSaveOpenManager(doProgress func(float32)) *SaveOpenManager {
	som := &SaveOpenManager{}
	som.doProgress = doProgress
	som.initProtobufs()
	return som
}

// Save ...
func (som *SaveOpenManager) Save(file *types.FBEntity, meshes []*meshes.ModelFace) {
	som.soProtobufs.Save(file, meshes)
}

// Open ...
func (som *SaveOpenManager) Open(file *types.FBEntity) []*meshes.ModelFace {
	return som.soProtobufs.Open(file)
}

func (som *SaveOpenManager) initProtobufs() {
	som.soProtobufs = NewProtoBufsSaveOpen(som.doProgress)
}
