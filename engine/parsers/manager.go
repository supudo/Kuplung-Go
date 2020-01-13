package parsers

import "github.com/supudo/Kuplung-Go/types"

// ParserManager ...
type ParserManager struct {
	objParser *ObjParser

	doProgress func(float32)
}

// NewParserManager ...
func NewParserManager(doProgress func(float32)) *ParserManager {
	pm := &ParserManager{}
	pm.doProgress = doProgress
	pm.initObjParser()
	return pm
}

// Parse ...
func (pm *ParserManager) Parse(filename string, psettings []string, itype types.ImportExportFormat) []types.MeshModel {
	switch itype {
	case types.ImportExportFormatOBJ:
		return pm.objParser.Parse(filename, psettings)
	}
	return nil
}

func (pm *ParserManager) initObjParser() {
	pm.objParser = NewObjParser(pm.doProgress)
}
