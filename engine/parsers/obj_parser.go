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
	idObjTitle string
	// vertex coordinates
	idGeometricVertices string
	// texture coordinates
	idTextureCoordinates string
	// normals
	idVertexNormals string
	// space vertices
	idSpaceVertices string
	// polygon faces
	idFace string
	// material file
	idMaterialFile string
	// material name for the current object
	idUseMaterial string

	// material
	idMaterialNew string

	// To specify the ambient reflectivity of the current material, you can use the "Ka" statement,
	// the "Ka spectral" statement, or the "Ka xyz" statement.
	idMaterialAmbientColor string
	// To specify the diffuse reflectivity of the current material, you can use the "Kd" statement,
	// the "Kd spectral" statement, or the "Kd xyz" statement.
	idMaterialDiffuseColor string
	// To specify the specular reflectivity of the current material, you can use the "Ks" statement,
	// the "Ks spectral" statement, or the "Ks xyz" statement.
	idMaterialSpecularColor string
	// The emission constant color of the material
	idMaterialEmissionColor string
	// Specifies the specular exponent for the current material. This defines the focus of the specular highlight.
	idMaterialSpecularExp string
	// Specifies the dissolve for the current material.  Tr or d, depending on the formats. Transperancy
	idMaterialTransperant1 string
	idMaterialTransperant2 string
	// Specifies the optical density for the surface. This is also known as index of refraction.
	idMaterialOpticalDensity string
	// The "illum" statement specifies the illumination model to use in the material.
	// Illumination models are mathematical equations that represent various material lighting and shading effects.
	idMaterialIllumination string
	// Specifies that a color texture file or a color procedural texture file is applied to the ambient reflectivity of the material.
	// During rendering, the "map_Ka" value is multiplied by the "Ka" value.
	idMaterialTextureAmbient string
	// Specifies that a color texture file or color procedural texture file is linked to the diffuse reflectivity of the material.
	// During rendering, the map_Kd value is multiplied by the Kd value.
	idMaterialTextureDiffuse string
	// Bump map
	idMaterialTextureBump string
	// Displacement map
	idMaterialTextureDisplacement string
	// Specifies that a color texture file or color procedural texture file is linked to the specular reflectivity of the material.
	// During rendering, the map_Ks value is multiplied by the Ks value.
	idMaterialTextureSpecular string
	// Specifies that a scalar texture file or scalar procedural texture file is linked to the specular exponent of the material.
	// During rendering, the map_Ns value is multiplied by the Ns value.
	idMaterialTextureSpecularExp string
	// Specifies that a scalar texture file or scalar procedural texture file is linked to the dissolve of the material.
	// During rendering, the map_d value is multiplied by the d value.
	idMaterialTextureDissolve string
}

// NewObjParser ...
func NewObjParser(doProgress func(float32)) *ObjParser {
	objp := &ObjParser{}

	objp.doProgress = doProgress
	objp.objFileLinesCount = 0

	objp.idObjTitle = "o "
	objp.idGeometricVertices = "v "
	objp.idTextureCoordinates = "vt "
	objp.idVertexNormals = "vn "
	objp.idSpaceVertices = "vp "
	objp.idFace = "f "
	objp.idMaterialFile = "mtllib "
	objp.idUseMaterial = "usemtl "
	objp.idMaterialNew = "newmtl "

	objp.idMaterialAmbientColor = "Ka "
	objp.idMaterialDiffuseColor = "Kd "
	objp.idMaterialSpecularColor = "Ks "
	objp.idMaterialEmissionColor = "Ke "
	objp.idMaterialSpecularExp = "Ns "
	objp.idMaterialTransperant1 = "Tr "
	objp.idMaterialTransperant2 = "d "
	objp.idMaterialOpticalDensity = "Ni "
	objp.idMaterialIllumination = "illum "
	objp.idMaterialTextureAmbient = "map_Ka "
	objp.idMaterialTextureDiffuse = "map_Kd "
	objp.idMaterialTextureBump = "map_Bump "
	objp.idMaterialTextureDisplacement = "disp "
	objp.idMaterialTextureSpecular = "map_Ks "
	objp.idMaterialTextureSpecularExp = "map_Ns "
	objp.idMaterialTextureDissolve = "map_d "

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

		if strings.HasPrefix(singleLine, objp.idMaterialFile) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialFile, "")
			objp.loadMaterialFile(singleLine)
		} else if strings.HasPrefix(singleLine, objp.idObjTitle) {
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
			entityModel.ModelTitle = strings.ReplaceAll(entityModel.ModelTitle, objp.idObjTitle, "")
			modelCounter++
			objp.models = append(objp.models, entityModel)
		} else if strings.HasPrefix(singleLine, objp.idGeometricVertices) {
			singleLine = strings.ReplaceAll(singleLine, objp.idGeometricVertices, "")
			fmt.Sscanf(singleLine, "%f %f %f", &x, &y, &z)
			vVertices = append(vVertices, mgl32.Vec3{x, y, z})
		} else if strings.HasPrefix(singleLine, objp.idTextureCoordinates) {
			singleLine = strings.ReplaceAll(singleLine, objp.idTextureCoordinates, "")
			fmt.Sscanf(singleLine, "%f %f", &x, &y)
			vTextureCoordinates = append(vTextureCoordinates, mgl32.Vec2{x, y})
		} else if strings.HasPrefix(singleLine, objp.idVertexNormals) {
			singleLine = strings.ReplaceAll(singleLine, objp.idVertexNormals, "")
			fmt.Sscanf(singleLine, "%f %f %f", &x, &y, &z)
			vNormals = append(vNormals, mgl32.Vec3{x, y, z})
		} else if strings.HasPrefix(singleLine, objp.idUseMaterial) {
			singleLine = strings.ReplaceAll(singleLine, objp.idUseMaterial, "")
			objp.models[currentModelID].ModelMaterial = *objp.materials[singleLine]
			objp.models[currentModelID].MaterialTitle = objp.models[currentModelID].ModelMaterial.MaterialTitle
		} else if strings.HasPrefix(singleLine, objp.idFace) {
			var ft []string
			ft = strings.Split(singleLine, " ")
			if len(ft) == 5 {
				var triVertexIndex, triUvIndex, triNormalIndex [4]uint32
				face := objp.idFace + "%v/%v/%v %v/%v/%v %v/%v/%v %v/%v/%v"
				matches, _ := fmt.Sscanf(singleLine, face,
					&triVertexIndex[0], &triUvIndex[0], &triNormalIndex[0],
					&triVertexIndex[1], &triUvIndex[1], &triNormalIndex[1],
					&triVertexIndex[2], &triUvIndex[2], &triNormalIndex[2],
					&triVertexIndex[3], &triUvIndex[3], &triNormalIndex[3])
				if matches != 12 {
					face = objp.idFace + "%v//%v %v//%v %v//%v %v//%v"
					matches, _ := fmt.Sscanf(singleLine, face,
						&triVertexIndex[0], &triNormalIndex[0],
						&triVertexIndex[1], &triNormalIndex[1],
						&triVertexIndex[2], &triNormalIndex[2],
						&triVertexIndex[3], &triNormalIndex[3])
					if matches != 8 {
						settings.LogWarn("[OBJ Parser] OBJ file is in wrong format: %v", objp.filename)
						return objp.models
					}
				}
				indexModels = append(indexModels, currentModelID)
				indexModels = append(indexModels, currentModelID)
				indexModels = append(indexModels, currentModelID)
				indexVertices = append(indexVertices, triVertexIndex[0])
				indexVertices = append(indexVertices, triVertexIndex[1])
				indexVertices = append(indexVertices, triVertexIndex[2])
				indexTexture = append(indexTexture, triUvIndex[0])
				indexTexture = append(indexTexture, triUvIndex[1])
				indexTexture = append(indexTexture, triUvIndex[2])
				indexNormals = append(indexNormals, triNormalIndex[0])
				indexNormals = append(indexNormals, triNormalIndex[1])
				indexNormals = append(indexNormals, triNormalIndex[2])

				indexModels = append(indexModels, currentModelID)
				indexModels = append(indexModels, currentModelID)
				indexModels = append(indexModels, currentModelID)
				indexVertices = append(indexVertices, triVertexIndex[2])
				indexVertices = append(indexVertices, triVertexIndex[3])
				indexVertices = append(indexVertices, triVertexIndex[0])
				indexTexture = append(indexTexture, triUvIndex[2])
				indexTexture = append(indexTexture, triUvIndex[3])
				indexTexture = append(indexTexture, triUvIndex[0])
				indexNormals = append(indexNormals, triNormalIndex[2])
				indexNormals = append(indexNormals, triNormalIndex[3])
				indexNormals = append(indexNormals, triNormalIndex[0])
			} else {
				var vertexIndex, uvIndex, normalIndex [3]uint32
				face := objp.idFace + "%v/%v/%v %v/%v/%v %v/%v/%v"
				matches, _ := fmt.Sscanf(singleLine, face,
					&vertexIndex[0], &uvIndex[0], &normalIndex[0],
					&vertexIndex[1], &uvIndex[1], &normalIndex[1],
					&vertexIndex[2], &uvIndex[2], &normalIndex[2])
				if matches != 9 {
					face = objp.idFace + "%v//%v %v//%v %v//%v"
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

	SettingAxisForward := int32(4)
	if len(psettings) > 0 && len(psettings[0]) != 0 {
		i64, _ := strconv.ParseUint(psettings[0], 10, 32)
		SettingAxisForward = int32(i64)
	}
	SettingAxisUp := int32(5)
	if len(psettings) > 1 && len(psettings[1]) != 0 {
		i64, _ := strconv.ParseUint(psettings[1], 10, 32)
		SettingAxisUp = int32(i64)
	}

	if len(objp.models) > 0 {
		progressStageCounter := 0
		progressStageTotal := len(indexVertices)
		objp.doProgress(0)
		for i := 0; i < len(indexVertices); i++ {
			modelIndex := indexModels[i]
			vertexIndex := indexVertices[i]
			normalIndex := indexNormals[i]

			vertex := FixVectorAxis(vVertices[vertexIndex-1], SettingAxisForward, SettingAxisUp)
			normal := FixVectorAxis(vNormals[normalIndex-1], SettingAxisForward, SettingAxisUp)

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
					packed = types.PackedVertex{Position: m.Vertices[j], UV: m.TextureCoordinates[j], Normal: m.Normals[j]}
				} else {
					packed = types.PackedVertex{Position: m.Vertices[j], UV: mgl32.Vec2{0, 0}, Normal: m.Normals[j]}
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
		if strings.HasPrefix(singleLine, objp.idMaterialNew) {
			currentMaterialTitle = singleLine
			currentMaterialTitle = strings.ReplaceAll(currentMaterialTitle, objp.idMaterialNew, "")
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
		} else if strings.HasPrefix(singleLine, objp.idMaterialAmbientColor) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialAmbientColor, "")
			fmt.Sscanf(singleLine, "%f %f %f", &r, &g, &b)
			mat := objp.materials[currentMaterialTitle]
			mat.AmbientColor = mgl32.Vec3{r, g, b}
			objp.materials[currentMaterialTitle].AmbientColor = mgl32.Vec3{r, g, b}
		} else if strings.HasPrefix(singleLine, objp.idMaterialDiffuseColor) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialDiffuseColor, "")
			fmt.Sscanf(singleLine, "%f %f %f", &r, &g, &b)
			objp.materials[currentMaterialTitle].DiffuseColor = mgl32.Vec3{r, g, b}
		} else if strings.HasPrefix(singleLine, objp.idMaterialSpecularColor) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialSpecularColor, "")
			fmt.Sscanf(singleLine, "%f %f %f", &r, &g, &b)
			objp.materials[currentMaterialTitle].SpecularColor = mgl32.Vec3{r, g, b}
		} else if strings.HasPrefix(singleLine, objp.idMaterialEmissionColor) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialEmissionColor, "")
			fmt.Sscanf(singleLine, "%f %f %f", &r, &g, &b)
			objp.materials[currentMaterialTitle].EmissionColor = mgl32.Vec3{r, g, b}
		} else if strings.HasPrefix(singleLine, objp.idMaterialSpecularExp) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialSpecularExp, "")
			f64, _ := strconv.ParseFloat(singleLine, 32)
			objp.materials[currentMaterialTitle].SpecularExp = float32(f64)
		} else if strings.HasPrefix(singleLine, objp.idMaterialTransperant1) || strings.HasPrefix(singleLine, objp.idMaterialTransperant2) {
			if strings.HasPrefix(singleLine, objp.idMaterialTransperant1) {
				singleLine = strings.ReplaceAll(singleLine, objp.idMaterialTransperant1, "")
			} else {
				singleLine = strings.ReplaceAll(singleLine, objp.idMaterialTransperant2, "")
			}
			f64, _ := strconv.ParseFloat(singleLine, 32)
			objp.materials[currentMaterialTitle].Transparency = float32(f64)
		} else if strings.HasPrefix(singleLine, objp.idMaterialOpticalDensity) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialOpticalDensity, "")
			f64, _ := strconv.ParseFloat(singleLine, 32)
			objp.materials[currentMaterialTitle].OpticalDensity = float32(f64)
		} else if strings.HasPrefix(singleLine, objp.idMaterialIllumination) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialIllumination, "")
			i64, _ := strconv.ParseUint(singleLine, 10, 32)
			objp.materials[currentMaterialTitle].IlluminationMode = uint32(i64)
		} else if strings.HasPrefix(singleLine, objp.idMaterialTextureAmbient) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialTextureAmbient, "")
			objp.materials[currentMaterialTitle].TextureAmbient = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.idMaterialTextureBump) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialTextureBump, "")
			objp.materials[currentMaterialTitle].TextureBump = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.idMaterialTextureDiffuse) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialTextureDiffuse, "")
			objp.materials[currentMaterialTitle].TextureDiffuse = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.idMaterialTextureDisplacement) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialTextureDisplacement, "")
			objp.materials[currentMaterialTitle].TextureDisplacement = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.idMaterialTextureDissolve) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialTextureDissolve, "")
			objp.materials[currentMaterialTitle].TextureDissolve = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.idMaterialTextureSpecular) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialTextureSpecular, "")
			objp.materials[currentMaterialTitle].TextureSpecular = objp.parseTextureImage(singleLine)
		} else if strings.HasPrefix(singleLine, objp.idMaterialTextureSpecularExp) {
			singleLine = strings.ReplaceAll(singleLine, objp.idMaterialTextureSpecularExp, "")
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
