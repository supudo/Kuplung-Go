package export

import "github.com/supudo/Kuplung-Go/types"

// ExporterObj ...
type ExporterObj struct {
}

// NewExporterObj ...
func NewExporterObj(doProgress func(float32)) *ExporterObj {
	eobj := &ExporterObj{}
	return eobj
}

// Export ...
func (eobj *ExporterObj) Export(file types.FBEntity, psettings []string) {
}
