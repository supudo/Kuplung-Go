package components

import (
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// ComponentExport ...
type ComponentExport struct {
	positionX, positionY int32
	width, height        int32
	panelWidthOptions    float32
	panelWidthOptionsMin float32

	SettingForward, SettingUp int32

	currentFolder string
}

// NewComponentExport ...
func NewComponentExport() *ComponentExport {
	sett := settings.GetSettings()
	cide := &ComponentExport{}
	cide.panelWidthOptions = 200.0
	cide.panelWidthOptionsMin = 200.0
	cide.SettingForward = 2
	cide.SettingUp = 4
	cide.currentFolder = sett.App.CurrentPath
	return cide
}

// Render ...
func (comp *ComponentExport) Render(open *bool, dialogExportType *types.ImportExportFormat) {
}
