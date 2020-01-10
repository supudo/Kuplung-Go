package types

// ImportExportFormat ...
type ImportExportFormat uint32

// Import/Export formats
const (
	ImportExportFormatUNDEFINED ImportExportFormat = iota
	ImportExportFormatOBJ
	ImportExportFormatGLTF
	ImportExportFormatSTL
	ImportExportFormatPLY
)
