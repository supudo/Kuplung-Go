package types

// LightSourceType ...
type LightSourceType uint32

// LightSourceType ...
const (
	LightSourceTypeDirectional LightSourceType = 0 + iota
	LightSourceTypePoint
	LightSourceTypeSpot
)
