package types

// MaterialTextureType ...
const (
	MaterialTextureTypeUndefined uint16 = iota + 1
	MaterialTextureTypeAmbient
	MaterialTextureTypeDiffuse
	MaterialTextureTypeDissolve
	MaterialTextureTypeBump
	MaterialTextureTypeSpecular
	MaterialTextureTypeSpecularExp
	MaterialTextureTypeDisplacement
)

// GetMaterialTextureName ...
func GetMaterialTextureName(texType uint16) string {
	switch texType {
	case MaterialTextureTypeAmbient:
		return "Ambient"
	case MaterialTextureTypeDiffuse:
		return "Diffuse"
	case MaterialTextureTypeDissolve:
		return "Dissolve"
	case MaterialTextureTypeBump:
		return "Normal"
	case MaterialTextureTypeSpecular:
		return "Specular"
	case MaterialTextureTypeSpecularExp:
		return "Specular Exp"
	case MaterialTextureTypeDisplacement:
		return "Displacement"
	default:
		return ""
	}
}
