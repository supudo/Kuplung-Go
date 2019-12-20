package types

import "github.com/go-gl/mathgl/mgl32"

// RenderProperties ...
type RenderProperties struct {
	UIAmbientLightX, UIAmbientLightY, UIAmbientLightZ float32

	SolidLightDirectionX, SolidLightDirectionY, SolidLightDirectionZ float32

	SolidLightMaterialColor, SolidLightAmbient mgl32.Vec3
	SolidLightDiffuse, SolidLightSpecular      mgl32.Vec3

	SolidLightAmbientStrength, SolidLightDiffuseStrength, SolidLightSpecularStrength float32

	SolidLightMaterialColorColorPicker, SolidLightAmbientColorPicker bool
	SolidLightDiffuseColorPicker, SolidLightSpecularColorPicker      bool
}
