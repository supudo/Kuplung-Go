package types

// ModelFaceLightSourceDirectional ...
type ModelFaceLightSourceDirectional struct {
	InUse                                              int32
	Direction                                          int32
	Ambient, Diffuse, Specular                         int32
	StrengthAmbient, StrengthDiffuse, StrengthSpecular int32
}

// ModelFaceLightSourcePoint ...
type ModelFaceLightSourcePoint struct {
	InUse                                              int32
	Position                                           int32
	Constant, Linear, Quadratic                        int32
	Ambient, Diffuse, Specular                         int32
	StrengthAmbient, StrengthDiffuse, StrengthSpecular int32
}

// ModelFaceLightSourceSpot ...
type ModelFaceLightSourceSpot struct {
	InUse                                              int32
	Position, Direction                                int32
	CutOff, OuterCutOff                                int32
	Constant, Linear, Quadratic                        int32
	Ambient, Diffuse, Specular                         int32
	StrengthAmbient, StrengthDiffuse, StrengthSpecular int32
}
