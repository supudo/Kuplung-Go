package types

import (
	"fmt"
	"github.com/supudo/Kuplung-Go/settings"
)

// KuplungPrintObjModels ..
func KuplungPrintObjModels(models []MeshModel, byIndices, shouldExit bool) {
	for i := 0; i < len(models); i++ {
		m := models[i]
		settings.LogWarn("model.ID = %v", m.ID)
		settings.LogWarn("model.countIndices = %v (%v)", m.CountIndices, MaxElement(m.Indices))
		settings.LogWarn("model.countNormals = %v (%v)", m.CountNormals, (m.CountNormals * 3))
		settings.LogWarn("model.countTextureCoordinates = %v (%v)", m.CountTextureCoordinates, (m.CountTextureCoordinates * 2))
		settings.LogWarn("model.countVertices = %v (%v)", m.CountVertices, (m.CountVertices * 3))
		settings.LogWarn("model.MaterialTitle = %v", m.MaterialTitle)
		settings.LogWarn("model.ModelTitle = %v", m.ModelTitle)

		if byIndices {
			settings.LogWarn("m.geometry :")
			for j := 0; j < len(m.Indices); j++ {
				idx := m.Indices[j]
				geom := fmt.Sprintf("index = %v ---> ", idx)
				vert := m.Vertices[idx]
				tc := m.TextureCoordinates[idx]
				n := m.Normals[idx]
				geom += fmt.Sprintf("vertex = [%v, %v, %v]", vert.X(), vert.Y(), vert.Z())
				geom += fmt.Sprintf(", uv = [%v, %v]", tc.X(), tc.Y())
				geom += fmt.Sprintf(", normal = [%v, %v, %v]", n.X(), n.Y(), n.Z())
				//settings.LogWarn("%v", geom)
			}
		} else {
			var verts string
			for j := 0; j < len(m.Vertices); j++ {
				v := m.Vertices[j]
				verts += fmt.Sprintf("[%v, %v, %v], ", v.X(), v.Y(), v.Z())
			}
			settings.LogWarn("m.vertices : %v", verts)

			var uvs string
			for j := 0; j < len(m.TextureCoordinates); j++ {
				uvs += fmt.Sprintf("[%v, %v], ", m.TextureCoordinates[j].X(), m.TextureCoordinates[j].Y())
			}
			settings.LogWarn("m.texture_coordinates : %v", uvs)

			var normals string
			for j := 0; j < len(m.Normals); j++ {
				n := m.Normals[j]
				normals += fmt.Sprintf("[%f, %f, %f], ", n.X(), n.Y(), n.Z())
			}
			settings.LogWarn("m.normals : %v", normals)

			var indices string
			for j := 0; j < len(m.Indices); j++ {
				indices += fmt.Sprintf("%v, ", m.Indices[j])
			}
			settings.LogWarn("m.indices : %v", indices)
		}

		settings.LogWarn("model.ModelMaterial.MaterialID = %v", m.ModelMaterial.MaterialID)
		settings.LogWarn("model.ModelMaterial.MaterialTitle = %v", m.ModelMaterial.MaterialTitle)

		settings.LogWarn("model.ModelMaterial.AmbientColor = %v, %v, %v", m.ModelMaterial.AmbientColor.X(), m.ModelMaterial.AmbientColor.Y(), m.ModelMaterial.AmbientColor.Z())
		settings.LogWarn("model.ModelMaterial.DiffuseColor = %v, %v, %v", m.ModelMaterial.DiffuseColor.X(), m.ModelMaterial.DiffuseColor.Y(), m.ModelMaterial.DiffuseColor.Z())
		settings.LogWarn("model.ModelMaterial.SpecularColor = %v, %v, %v", m.ModelMaterial.SpecularColor.X(), m.ModelMaterial.SpecularColor.Y(), m.ModelMaterial.SpecularColor.Z())
		settings.LogWarn("model.ModelMaterial.EmissionColor = %v, %v, %v", m.ModelMaterial.EmissionColor.X(), m.ModelMaterial.EmissionColor.Y(), m.ModelMaterial.EmissionColor.Z())

		settings.LogWarn("model.ModelMaterial.SpecularExp = %v", m.ModelMaterial.SpecularExp)
		settings.LogWarn("model.ModelMaterial.Transparency = %v", m.ModelMaterial.Transparency)
		settings.LogWarn("model.ModelMaterial.OpticalDensity = %v", m.ModelMaterial.OpticalDensity)
		settings.LogWarn("model.ModelMaterial.IlluminationMode = %v", m.ModelMaterial.IlluminationMode)

		settings.LogWarn("model.ModelMaterial.textures_ambient.Filename = %v", m.ModelMaterial.TextureAmbient.Filename)
		settings.LogWarn("model.ModelMaterial.textures_diffuse.Filename = %v", m.ModelMaterial.TextureDiffuse.Filename)
		settings.LogWarn("model.ModelMaterial.textures_specular.Filename = %v", m.ModelMaterial.TextureSpecular.Filename)
		settings.LogWarn("model.ModelMaterial.textures_specularExp.Filename = %v", m.ModelMaterial.TextureSpecularExp.Filename)
		settings.LogWarn("model.ModelMaterial.textures_dissolve.Filename = %v", m.ModelMaterial.TextureDissolve.Filename)
		settings.LogWarn("model.ModelMaterial.textures_bump.Filename = %v", m.ModelMaterial.TextureBump.Filename)
		settings.LogWarn("model.ModelMaterial.textures_displacement.Filename = %v", m.ModelMaterial.TextureDisplacement.Filename)
	}
	if shouldExit {
		settings.LogError("--------")
	} else {
		settings.LogWarn("--------")
	}
}

// MaxElement ...
func MaxElement(array []uint32) uint32 {
	var max uint32 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}

func MinMax(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
