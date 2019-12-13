package parsers

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// ObjParser ...
type ObjParser struct {
	objFileLinesCount int32
	filename          string
	doProgress        func(float32)

	models                        []types.MeshModel
	materials                     map[string]*types.MeshModelMaterial
	vectorVertices, vectorNormals []mgl32.Vec3
	vectorTextureCoordinates      []mgl32.Vec2
	vectorIndices                 []uint32

	// current object name
	id_objTitle string
	// vertex coordinates
	id_geometricVertices string
	// texture coordinates
	id_textureCoordinates string
	// normals
	id_vertexNormals string
	// space vertices
	id_spaceVertices string
	// polygon faces
	id_face string
	// material file
	id_materialFile string
	// material name for the current object
	id_useMaterial string

	// material
	id_materialNew string

	// To specify the ambient reflectivity of the current material, you can use the "Ka" statement,
	// the "Ka spectral" statement, or the "Ka xyz" statement.
	id_materialAmbientColor string
	// To specify the diffuse reflectivity of the current material, you can use the "Kd" statement,
	// the "Kd spectral" statement, or the "Kd xyz" statement.
	id_materialDiffuseColor string
	// To specify the specular reflectivity of the current material, you can use the "Ks" statement,
	// the "Ks spectral" statement, or the "Ks xyz" statement.
	id_materialSpecularColor string
	// The emission constant color of the material
	id_materialEmissionColor string
	// Specifies the specular exponent for the current material. This defines the focus of the specular highlight.
	id_materialSpecularExp string
	// Specifies the dissolve for the current material.  Tr or d, depending on the formats. Transperancy
	id_materialTransperant1 string
	id_materialTransperant2 string
	// Specifies the optical density for the surface. This is also known as index of refraction.
	id_materialOpticalDensity string
	// The "illum" statement specifies the illumination model to use in the material.
	// Illumination models are mathematical equations that represent various material lighting and shading effects.
	id_materialIllumination string
	// Specifies that a color texture file or a color procedural texture file is applied to the ambient reflectivity of the material.
	// During rendering, the "map_Ka" value is multiplied by the "Ka" value.
	id_materialTextureAmbient string
	// Specifies that a color texture file or color procedural texture file is linked to the diffuse reflectivity of the material.
	// During rendering, the map_Kd value is multiplied by the Kd value.
	id_materialTextureDiffuse string
	// Bump map
	id_materialTextureBump string
	// Displacement map
	id_materialTextureDisplacement string
	// Specifies that a color texture file or color procedural texture file is linked to the specular reflectivity of the material.
	// During rendering, the map_Ks value is multiplied by the Ks value.
	id_materialTextureSpecular string
	// Specifies that a scalar texture file or scalar procedural texture file is linked to the specular exponent of the material.
	// During rendering, the map_Ns value is multiplied by the Ns value.
	id_materialTextureSpecularExp string
	// Specifies that a scalar texture file or scalar procedural texture file is linked to the dissolve of the material.
	// During rendering, the map_d value is multiplied by the d value.
	id_materialTextureDissolve string
}

// NewObjParser ...
func NewObjParser(doProgress func(float32)) *ObjParser {
	objp := &ObjParser{}

	objp.doProgress = doProgress
	objp.objFileLinesCount = 0

	objp.id_objTitle = "o "
	objp.id_geometricVertices = "v "
	objp.id_textureCoordinates = "vt "
	objp.id_vertexNormals = "vn "
	objp.id_spaceVertices = "vp "
	objp.id_face = "f "
	objp.id_materialFile = "mtllib "
	objp.id_useMaterial = "usemtl "
	objp.id_materialNew = "newmtl "

	objp.id_materialAmbientColor = "Ka "
	objp.id_materialDiffuseColor = "Kd "
	objp.id_materialSpecularColor = "Ks "
	objp.id_materialEmissionColor = "Ke "
	objp.id_materialSpecularExp = "Ns "
	objp.id_materialTransperant1 = "Tr "
	objp.id_materialTransperant2 = "d "
	objp.id_materialOpticalDensity = "Ni "
	objp.id_materialIllumination = "illum "
	objp.id_materialTextureAmbient = "map_Ka "
	objp.id_materialTextureDiffuse = "map_Kd "
	objp.id_materialTextureBump = "map_Bump "
	objp.id_materialTextureDisplacement = "disp "
	objp.id_materialTextureSpecular = "map_Ks "
	objp.id_materialTextureSpecularExp = "map_Ns "
	objp.id_materialTextureDissolve = "map_d "

	return objp
}

// Parse ...
func (objp *ObjParser) Parse(filename string, psettings []string) []types.MeshModel {
	objp.resetSettings()

	objp.filename = filename

	file, err := os.Open(objp.filename)
	if err != nil {
		settings.LogWarn("[OBJ Parser] Can't open obj file (%v): %v", objp.filename, err)
	}
	defer file.Close()

	var indexModels, indexVertices, indexTexture, indexNormals []uint32
	var vVertices, vNormals []mgl32.Vec3
	var vTextureCoordinates []mgl32.Vec2

	modelCounter, currentModelID, progressStageCounter := uint32(0), uint32(0), uint32(0)

	var singleLine string
	var x, y, z float32
	scanner := bufio.NewScanner(file)
	progressStageTotal := objp.getNumberOfLines(objp.filename)
	for scanner.Scan() {
		singleLine = scanner.Text()

		if strings.HasPrefix(singleLine, objp.id_materialFile) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialFile, "")
			objp.loadMaterialFile(singleLine)
		} else if strings.HasPrefix(singleLine, objp.id_objTitle) {
			currentModelID = modelCounter
			entityModel := types.MeshModel{
				File:                    objp.filename,
				ID:                      currentModelID,
				ModelTitle:              singleLine,
				CountVertices:           0,
				CountTextureCoordinates: 0,
				CountNormals:            0,
				CountIndices:            0,
			}
			entityModel.ModelTitle = strings.ReplaceAll(entityModel.ModelTitle, objp.id_objTitle, "")
			modelCounter++
			objp.models = append(objp.models, entityModel)
		} else if strings.HasPrefix(singleLine, objp.id_geometricVertices) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_geometricVertices, "")
			fmt.Sscanf(singleLine, "%f %f %f", &x, &y, &z)
			vVertices = append(vVertices, mgl32.Vec3{x, y, z})
		} else if strings.HasPrefix(singleLine, objp.id_textureCoordinates) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_textureCoordinates, "")
			fmt.Sscanf(singleLine, "%f %f", &x, &y)
			vTextureCoordinates = append(vTextureCoordinates, mgl32.Vec2{x, y})
		} else if strings.HasPrefix(singleLine, objp.id_vertexNormals) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_vertexNormals, "")
			fmt.Sscanf(singleLine, "%f %f %f", &x, &y, &z)
			vNormals = append(vNormals, mgl32.Vec3{x, y, z})
		} else if strings.HasPrefix(singleLine, objp.id_useMaterial) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_useMaterial, "")
			objp.models[currentModelID].ModelMaterial = *objp.materials[singleLine]
			objp.models[currentModelID].MaterialTitle = objp.models[currentModelID].ModelMaterial.MaterialTitle
		} else if strings.HasPrefix(singleLine, objp.id_face) {
			var ft []string
			ft = strings.Split(singleLine, " ")
			if len(ft) == 5 {
				var tri_vertexIndex, tri_uvIndex, tri_normalIndex [4]uint32
				face := objp.id_face + "%v/%v/%v %v/%v/%v %v/%v/%v %v/%v/%v"
				matches, _ := fmt.Sscanf(singleLine, face,
					&tri_vertexIndex[0], &tri_uvIndex[0], &tri_normalIndex[0],
					&tri_vertexIndex[1], &tri_uvIndex[1], &tri_normalIndex[1],
					&tri_vertexIndex[2], &tri_uvIndex[2], &tri_normalIndex[2],
					&tri_vertexIndex[3], &tri_uvIndex[3], &tri_normalIndex[3])
				if matches != 12 {
					face = objp.id_face + "%v//%v %v//%v %v//%v %v//%v"
					matches, _ := fmt.Sscanf(singleLine, face,
						&tri_vertexIndex[0], &tri_normalIndex[0],
						&tri_vertexIndex[1], &tri_normalIndex[1],
						&tri_vertexIndex[2], &tri_normalIndex[2],
						&tri_vertexIndex[3], &tri_normalIndex[3])
					if matches != 8 {
						settings.LogWarn("[OBJ Parser] OBJ file is in wrong format: %v", objp.filename)
						return objp.models
					}
				}
				indexModels = append(indexModels, currentModelID)
				indexModels = append(indexModels, currentModelID)
				indexModels = append(indexModels, currentModelID)
				indexVertices = append(indexVertices, tri_vertexIndex[0])
				indexVertices = append(indexVertices, tri_vertexIndex[1])
				indexVertices = append(indexVertices, tri_vertexIndex[2])
				indexTexture = append(indexTexture, tri_uvIndex[0])
				indexTexture = append(indexTexture, tri_uvIndex[1])
				indexTexture = append(indexTexture, tri_uvIndex[2])
				indexNormals = append(indexNormals, tri_normalIndex[0])
				indexNormals = append(indexNormals, tri_normalIndex[1])
				indexNormals = append(indexNormals, tri_normalIndex[2])

				indexModels = append(indexModels, currentModelID)
				indexModels = append(indexModels, currentModelID)
				indexModels = append(indexModels, currentModelID)
				indexVertices = append(indexVertices, tri_vertexIndex[2])
				indexVertices = append(indexVertices, tri_vertexIndex[3])
				indexVertices = append(indexVertices, tri_vertexIndex[0])
				indexTexture = append(indexTexture, tri_uvIndex[2])
				indexTexture = append(indexTexture, tri_uvIndex[3])
				indexTexture = append(indexTexture, tri_uvIndex[0])
				indexNormals = append(indexNormals, tri_normalIndex[2])
				indexNormals = append(indexNormals, tri_normalIndex[3])
				indexNormals = append(indexNormals, tri_normalIndex[0])
			} else {
				var vertexIndex, uvIndex, normalIndex [3]uint32
				face := objp.id_face + "%v/%v/%v %v/%v/%v %v/%v/%v"
				matches, _ := fmt.Sscanf(singleLine, face,
					&vertexIndex[0], &uvIndex[0], &normalIndex[0],
					&vertexIndex[1], &uvIndex[1], &normalIndex[1],
					&vertexIndex[2], &uvIndex[2], &normalIndex[2])
				if matches != 9 {
					face = objp.id_face + "%v//%v %v//%v %v//%v"
					matches, _ := fmt.Sscanf(singleLine, face,
						&vertexIndex[0], &normalIndex[0],
						&vertexIndex[1], &normalIndex[1],
						&vertexIndex[2], &normalIndex[2])
					if matches != 6 {
						settings.LogWarn("[OBJ Parser] OBJ file is in wrong format: %v", objp.filename)
						return objp.models
					}
				}
				indexModels = append(indexModels, currentModelID)
				indexModels = append(indexModels, currentModelID)
				indexModels = append(indexModels, currentModelID)
				indexVertices = append(indexVertices, vertexIndex[0])
				indexVertices = append(indexVertices, vertexIndex[1])
				indexVertices = append(indexVertices, vertexIndex[2])
				indexTexture = append(indexTexture, uvIndex[0])
				indexTexture = append(indexTexture, uvIndex[1])
				indexTexture = append(indexTexture, uvIndex[2])
				indexNormals = append(indexNormals, normalIndex[0])
				indexNormals = append(indexNormals, normalIndex[1])
				indexNormals = append(indexNormals, normalIndex[2])
			}
		}

		progressStageCounter++
		_ = (float32(progressStageCounter) / float32(progressStageTotal)) * 100.0
	}

	if err := scanner.Err(); err != nil {
		settings.LogWarn("[OBJ Parser] Scanner error: %v", err)
	}

	Setting_Axis_Forward := int32(4)
	if len(psettings) > 0 && len(psettings[0]) != 0 {
		i64, _ := strconv.ParseUint(psettings[0], 10, 32)
		Setting_Axis_Forward = int32(i64)
	}
	Setting_Axis_Up := int32(5)
	if len(psettings) > 1 && len(psettings[1]) != 0 {
		i64, _ := strconv.ParseUint(psettings[1], 10, 32)
		Setting_Axis_Up = int32(i64)
	}

	if len(objp.models) > 0 {
		progressStageCounter := 0
		progressStageTotal := len(indexVertices)
		objp.doProgress(0)
		for i := 0; i < len(indexVertices); i++ {
			modelIndex := indexModels[i]
			vertexIndex := indexVertices[i]
			normalIndex := indexNormals[i]

			vertex := FixVectorAxis(vVertices[vertexIndex-1], Setting_Axis_Forward, Setting_Axis_Up)
			normal := FixVectorAxis(vNormals[normalIndex-1], Setting_Axis_Forward, Setting_Axis_Up)

			objp.models[modelIndex].Vertices = append(objp.models[modelIndex].Vertices, vertex)
			objp.models[modelIndex].CountVertices++
			objp.models[modelIndex].Normals = append(objp.models[modelIndex].Normals, normal)
			objp.models[modelIndex].CountNormals++

			if len(vTextureCoordinates) > 0 {
				uvIndex := indexTexture[i]
				uv := vTextureCoordinates[uvIndex-1]
				objp.models[modelIndex].TextureCoordinates = append(objp.models[modelIndex].TextureCoordinates, uv)
				objp.models[modelIndex].CountTextureCoordinates++
			} else {
				objp.models[modelIndex].CountTextureCoordinates = 0
			}

			progressStageCounter++
			progress := (float32(progressStageCounter) / float32(progressStageTotal)) * 100.0
			objp.doProgress(progress)
		}

		progressStageCounter = 0
		progressStageTotal = len(objp.models)
		objp.doProgress(0.0)
		vertexToOutIndex := make(map[types.PackedVertex]uint32)
		for i := 0; i < len(objp.models); i++ {
			m := objp.models[i]
			var outVertices, outNormals []mgl32.Vec3
			var outTextureCoordinates []mgl32.Vec2
			for j := 0; j < len(m.Vertices); j++ {
				var packed types.PackedVertex
				if len(m.TextureCoordinates) > 0 {
					packed = types.PackedVertex{m.Vertices[j], m.TextureCoordinates[j], m.Normals[j]}
				} else {
					packed = types.PackedVertex{m.Vertices[j], mgl32.Vec2{0, 0}, m.Normals[j]}
				}

				index, found := objp.getSimilarVertexIndex(packed, vertexToOutIndex)
				if found {
					m.Indices = append(m.Indices, index)
				} else {
					outVertices = append(outVertices, m.Vertices[j])
					if len(m.TextureCoordinates) > 0 {
						outTextureCoordinates = append(outTextureCoordinates, m.TextureCoordinates[j])
					}
					outNormals = append(outNormals, m.Normals[j])
					newIndex := uint32(len(outVertices) - 1)
					m.Indices = append(m.Indices, newIndex)
					vertexToOutIndex[packed] = newIndex
				}
			}
			objp.models[i].Vertices = outVertices
			objp.models[i].TextureCoordinates = outTextureCoordinates
			objp.models[i].Normals = outNormals
			objp.models[i].Indices = m.Indices
			objp.models[i].CountIndices = int32(len(m.Indices))

			progressStageCounter++
			progress := (float32(progressStageCounter) / float32(progressStageTotal)) * 100.0
			objp.doProgress(progress)
		}
	}

	// if objp.models[0].ModelTitle == "XPlus" {
	// 	types.KuplungPrintObjModels(objp.models, false, true)
	// }

	return objp.models
}

func (objp *ObjParser) loadMaterialFile(materialFile string) {
	objp.materials = make(map[string]*types.MeshModelMaterial)

	materialPath := filepath.Dir(objp.filename) + "/" + materialFile

	file, err := os.Open(materialPath)
	if err != nil {
		settings.LogWarn("[OBJ Parser] Can't open .mtl file (%v): %v", materialPath, err)
	}
	defer file.Close()

	MaterialID := uint32(0)
	var singleLine, currentMaterialTitle string
	var r, g, b float32
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		singleLine = scanner.Text()
		if strings.HasPrefix(singleLine, objp.id_materialNew) {
			currentMaterialTitle = singleLine
			currentMaterialTitle = strings.ReplaceAll(currentMaterialTitle, objp.id_materialNew, "")
			entityMaterial := types.MeshModelMaterial{
				MaterialID:       MaterialID,
				MaterialTitle:    currentMaterialTitle,
				SpecularExp:      1.0,
				Transparency:     1.0,
				IlluminationMode: 2,
				OpticalDensity:   1.0,
				AmbientColor:     mgl32.Vec3{0, 0, 0},
				DiffuseColor:     mgl32.Vec3{0, 0, 0},
				SpecularColor:    mgl32.Vec3{0, 0, 0},
				EmissionColor:    mgl32.Vec3{0, 0, 0}}
			MaterialID++
			objp.materials[currentMaterialTitle] = &entityMaterial
		} else if strings.HasPrefix(singleLine, objp.id_materialAmbientColor) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialAmbientColor, "")
			fmt.Sscanf(singleLine, "%f %f %f", &r, &g, &b)
			mat := objp.materials[currentMaterialTitle]
			mat.AmbientColor = mgl32.Vec3{r, g, b}
			objp.materials[currentMaterialTitle].AmbientColor = mgl32.Vec3{r, g, b}
		} else if strings.HasPrefix(singleLine, objp.id_materialDiffuseColor) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialDiffuseColor, "")
			fmt.Sscanf(singleLine, "%f %f %f", &r, &g, &b)
			objp.materials[currentMaterialTitle].DiffuseColor = mgl32.Vec3{r, g, b}
		} else if strings.HasPrefix(singleLine, objp.id_materialSpecularColor) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialSpecularColor, "")
			fmt.Sscanf(singleLine, "%f %f %f", &r, &g, &b)
			objp.materials[currentMaterialTitle].SpecularColor = mgl32.Vec3{r, g, b}
		} else if strings.HasPrefix(singleLine, objp.id_materialEmissionColor) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialEmissionColor, "")
			fmt.Sscanf(singleLine, "%f %f %f", &r, &g, &b)
			objp.materials[currentMaterialTitle].EmissionColor = mgl32.Vec3{r, g, b}
		} else if strings.HasPrefix(singleLine, objp.id_materialSpecularExp) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialSpecularExp, "")
			f64, _ := strconv.ParseFloat(singleLine, 32)
			objp.materials[currentMaterialTitle].SpecularExp = float32(f64)
		} else if strings.HasPrefix(singleLine, objp.id_materialTransperant1) || strings.HasPrefix(singleLine, objp.id_materialTransperant2) {
			if strings.HasPrefix(singleLine, objp.id_materialTransperant1) {
				singleLine = strings.ReplaceAll(singleLine, objp.id_materialTransperant1, "")
			} else {
				singleLine = strings.ReplaceAll(singleLine, objp.id_materialTransperant2, "")
			}
			f64, _ := strconv.ParseFloat(singleLine, 32)
			objp.materials[currentMaterialTitle].Transparency = float32(f64)
		} else if strings.HasPrefix(singleLine, objp.id_materialOpticalDensity) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialOpticalDensity, "")
			f64, _ := strconv.ParseFloat(singleLine, 32)
			objp.materials[currentMaterialTitle].OpticalDensity = float32(f64)
		} else if strings.HasPrefix(singleLine, objp.id_materialIllumination) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialIllumination, "")
			i64, _ := strconv.ParseUint(singleLine, 10, 32)
			objp.materials[currentMaterialTitle].IlluminationMode = uint32(i64)
		} else if strings.HasPrefix(singleLine, objp.id_materialTextureAmbient) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialTextureAmbient, "")
			objp.materials[currentMaterialTitle].TextureAmbient = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.id_materialTextureBump) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialTextureBump, "")
			objp.materials[currentMaterialTitle].TextureBump = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.id_materialTextureDiffuse) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialTextureDiffuse, "")
			objp.materials[currentMaterialTitle].TextureDiffuse = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.id_materialTextureDisplacement) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialTextureDisplacement, "")
			objp.materials[currentMaterialTitle].TextureDisplacement = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.id_materialTextureDissolve) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialTextureDissolve, "")
			objp.materials[currentMaterialTitle].TextureDissolve = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.id_materialTextureSpecular) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialTextureSpecular, "")
			objp.materials[currentMaterialTitle].TextureSpecular = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.id_materialTextureSpecularExp) {
			singleLine = strings.ReplaceAll(singleLine, objp.id_materialTextureSpecularExp, "")
			objp.materials[currentMaterialTitle].TextureSpecularExp = objp.parseTextureImage(singleLine)
		}
	}
}

func (objp *ObjParser) parseTextureImage(textureLine string) types.MeshMaterialTextureImage {
	var materialImage types.MeshMaterialTextureImage

	materialImage.Height = 0
	materialImage.Width = 0
	materialImage.UseTexture = true

	if strings.Contains(textureLine, "-") {
		lineElements := strings.Split(textureLine, "-")

		if lineElements[0] == "" {
			lineElements = append(lineElements[:0], lineElements[1:]...)
		}

		lastElements := strings.Split(lineElements[len(lineElements)-1], " ")
		materialImage.Image = lastElements[len(lastElements)-1]
		lastElements = lastElements[:len(lastElements)-1]

		// TODO: commands
		materialImage.Commands = []string{""}
	} else {
		materialImage.Image = textureLine
	}
	strings.TrimSpace(materialImage.Image)

	folderPath := objp.filename
	folderPath = strings.ReplaceAll(folderPath, objp.filename, "")

	if !objp.fileExists(materialImage.Image) {
		materialImage.Image = folderPath + materialImage.Image
	}

	fileElements := strings.Split(materialImage.Image, "/")
	materialImage.Filename = fileElements[len(fileElements)-1]
	return materialImage
}

func (objp *ObjParser) resetSettings() {
	objp.objFileLinesCount = 0
	objp.filename = ""
	objp.models = nil
	objp.materials = nil
	objp.vectorVertices = nil
	objp.vectorNormals = nil
	objp.vectorTextureCoordinates = nil
	objp.vectorIndices = nil
}

func (objp *ObjParser) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (objp *ObjParser) getNumberOfLines(filename string) uint32 {
	f, _ := os.Open(filename)
	scanner := bufio.NewScanner(f)
	lineCounter := uint32(0)
	for scanner.Scan() {
		lineCounter++
	}
	return lineCounter
}

func (objp *ObjParser) getSimilarVertexIndex(packed types.PackedVertex, vertexToOutIndex map[types.PackedVertex]uint32) (uint32, bool) {
	index, found := vertexToOutIndex[packed]
	return index, found
}
