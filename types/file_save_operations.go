package types

// FileSaverOperation ...
type FileSaverOperation uint32

// FileSaverOperation ...
const (
	FileSaverOperationSaveScene FileSaverOperation = 0 + iota
	FileSaverOperationOpenScene
	FileSaverOperationRenderer
)
