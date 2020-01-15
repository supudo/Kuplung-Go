package saveopen

import (
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/objects"
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
func (som *SOManager) Save(file *types.FBEntity, meshes []*meshes.ModelFace, lights []*objects.Light, rprops types.RenderProperties, cam *objects.Camera, grid *objects.WorldGrid) {
	som.soProtobufs.Save(file, meshes, lights, rprops, cam, grid)
}

// Open ...
func (som *SOManager) Open(file *types.FBEntity, window interfaces.Window, systemModels map[string]types.MeshModel, faces *[]*meshes.ModelFace, lights *[]*objects.Light, rprops *types.RenderProperties, cam *objects.Camera, grid *objects.WorldGrid) {
	som.soProtobufs.Open(file, window, systemModels, faces, lights, rprops, cam, grid)
}

func (som *SOManager) initProtobufs() {
	som.soProtobufs = NewProtoBufsSaveOpen(som.doProgress)
}
