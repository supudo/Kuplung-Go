package export

import "github.com/supudo/Kuplung-Go/types"

// ExporterManager ...
type ExporterManager struct {
	exporterObj *ExporterObj

	doProgress func(float32)
}

// NewExportManager ...
func NewExportManager(doProgress func(float32)) *ExporterManager {
	pm := &ExporterManager{}
	pm.doProgress = doProgress
	pm.initExporterObj()
	return pm
}

// Export ...
func (pm *ExporterManager) Export(file types.FBEntity, psettings []string, itype types.ImportExportFormat) {
	switch itype {
	case types.ImportExportFormatOBJ:
		pm.exporterObj.Export(file, psettings)
	}
}

func (pm *ExporterManager) initExporterObj() {
	pm.exporterObj = NewExporterObj(pm.doProgress)
}
