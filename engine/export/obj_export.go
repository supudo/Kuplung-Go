package export

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// ExporterObj ...
type ExporterObj struct {
	funcProgress func(float32)

	uniqueVertices           []mgl32.Vec3
	uniqueTextureCoordinates []mgl32.Vec2
	uniqueNormals            []mgl32.Vec3
	vCounter                 int32
	vtCounter                int32
	vnCounter                int32

	addSuffix   bool
	objSettings []string
	exportFile  types.FBEntity
	nlDelimiter string
}

// NewExporterObj ...
func NewExporterObj(doProgress func(float32)) *ExporterObj {
	eobj := &ExporterObj{
		vCounter:     1,
		vtCounter:    1,
		vnCounter:    1,
		addSuffix:    false,
		funcProgress: doProgress,
	}
	if runtime.GOOS == "darwin" {
		eobj.nlDelimiter = "\n"
	}
	return eobj
}

// Export ...
func (eobj *ExporterObj) Export(faces []*meshes.ModelFace, file types.FBEntity, psettings []string) {
	eobj.objSettings = psettings
	eobj.addSuffix = false
	eobj.exportFile = file
	eobj.exportGeometry(faces)
	eobj.exportMaterials(faces)
}

func (eobj *ExporterObj) exportGeometry(faces []*meshes.ModelFace) {
	fileContents := "# Kuplung v1.0 OBJ File Export" + eobj.nlDelimiter
	fileContents += "# http://www.github.com/supudo/kuplung/" + eobj.nlDelimiter
	fn := eobj.exportFile.Title
	fn = strings.TrimSuffix(fn, ".obj")
	fileContents += "mtllib " + fn + ".mtl" + eobj.nlDelimiter

	eobj.uniqueVertices = nil
	eobj.uniqueTextureCoordinates = nil
	eobj.uniqueNormals = nil
	eobj.vCounter = 1
	eobj.vtCounter = 1
	eobj.vnCounter = 1

	for i := range faces {
		fileContents += eobj.exportMesh(*faces[i])
	}
	fileContents += eobj.nlDelimiter

	if len(fileContents) > 0 {
		fileSuffix := "_" + strconv.FormatInt(time.Now().Unix(), 10)
		if !eobj.addSuffix {
			fileSuffix = ""
		}
		filePath := filepath.Dir(eobj.exportFile.Path)
		fileName := eobj.exportFile.Title
		fileName = strings.TrimSuffix(fileName, ".obj")
		settings.SaveStringToFile(fileContents, filePath+"/"+fileName+fileSuffix+".obj", "ExporterOBJ")
	}
}

func (eobj *ExporterObj) exportMaterials(faces []*meshes.ModelFace) {
	materials := make(map[string]string)
	for i := 0; i < len(faces); i++ {
		mat := faces[i].MeshModel.ModelMaterial
		if len(materials[mat.MaterialTitle]) == 0 {
			materials[mat.MaterialTitle] = eobj.nlDelimiter
			materials[mat.MaterialTitle] += "newmtl " + mat.MaterialTitle + eobj.nlDelimiter
			materials[mat.MaterialTitle] += fmt.Sprintf("Ns %g", mat.SpecularExp) + eobj.nlDelimiter
			materials[mat.MaterialTitle] += fmt.Sprintf("Ka %g %g %g", mat.AmbientColor.X(), mat.AmbientColor.Y(), mat.AmbientColor.Z()) + eobj.nlDelimiter
			materials[mat.MaterialTitle] += fmt.Sprintf("Kd %g %g %g", mat.DiffuseColor.X(), mat.DiffuseColor.Y(), mat.DiffuseColor.Z()) + eobj.nlDelimiter
			materials[mat.MaterialTitle] += fmt.Sprintf("Ks %g %g %g", mat.SpecularColor.X(), mat.SpecularColor.Y(), mat.SpecularColor.Z()) + eobj.nlDelimiter
			materials[mat.MaterialTitle] += fmt.Sprintf("Ke %g %g %g", mat.EmissionColor.X(), mat.EmissionColor.Y(), mat.EmissionColor.Z()) + eobj.nlDelimiter
			if mat.OpticalDensity >= 0.0 {
				materials[mat.MaterialTitle] += fmt.Sprintf("Ni %g", mat.OpticalDensity) + eobj.nlDelimiter
			}
			materials[mat.MaterialTitle] += fmt.Sprintf("d %g", mat.Transparency) + eobj.nlDelimiter
			materials[mat.MaterialTitle] += fmt.Sprintf("illum %d", mat.IlluminationMode) + eobj.nlDelimiter

			if len(mat.TextureAmbient.Image) > 0 {
				materials[mat.MaterialTitle] += "map_Ka " + mat.TextureAmbient.Image + eobj.nlDelimiter
			}
			if len(mat.TextureDiffuse.Image) > 0 {
				materials[mat.MaterialTitle] += "map_Kd " + mat.TextureDiffuse.Image + eobj.nlDelimiter
			}
			if len(mat.TextureDissolve.Image) > 0 {
				materials[mat.MaterialTitle] += "map_d " + mat.TextureDissolve.Image + eobj.nlDelimiter
			}
			if len(mat.TextureBump.Image) > 0 {
				materials[mat.MaterialTitle] += "map_Bump " + mat.TextureBump.Image + eobj.nlDelimiter
			}
			if len(mat.TextureDisplacement.Image) > 0 {
				materials[mat.MaterialTitle] += "disp " + mat.TextureDisplacement.Image + eobj.nlDelimiter
			}
			if len(mat.TextureSpecular.Image) > 0 {
				materials[mat.MaterialTitle] += "map_Ks " + mat.TextureSpecular.Image + eobj.nlDelimiter
			}
			if len(mat.TextureSpecularExp.Image) > 0 {
				materials[mat.MaterialTitle] += "map_Ns " + mat.TextureSpecularExp.Image + eobj.nlDelimiter
			}
		}
	}

	fileContents := "# Kuplung MTL File" + eobj.nlDelimiter
	fileContents += fmt.Sprintf("# Material Count: %d", len(materials)) + eobj.nlDelimiter
	fileContents += "# http://www.github.com/supudo/kuplung/" + eobj.nlDelimiter

	for _, line := range materials {
		fileContents += line
	}

	fileContents += eobj.nlDelimiter

	if len(fileContents) > 0 {
		fileSuffix := "_" + strconv.FormatInt(time.Now().Unix(), 10)
		if !eobj.addSuffix {
			fileSuffix = ""
		}
		filePath := filepath.Dir(eobj.exportFile.Path)
		fileName := eobj.exportFile.Title
		fileName = strings.TrimSuffix(fileName, ".obj")
		settings.SaveStringToFile(fileContents, filePath+"/"+fileName+fileSuffix+".mtl", "ExporterOBJ")
	}
}

func (eobj *ExporterObj) exportMesh(face meshes.ModelFace) string {
	model := face.MeshModel
	meshData := ""
	v := ""
	vt := ""
	vn := ""
	f := ""

	eobj.funcProgress(0.0)
	totalProgress := int32(len(model.Indices)) * 2
	progressCounter := float32(0.0)

	meshData += eobj.nlDelimiter
	meshData += "o " + model.ModelTitle + eobj.nlDelimiter
	for j := 0; j < len(model.Indices); j++ {
		idx := model.Indices[j]
		vertex := model.Vertices[idx]

		vertex = vertex.Add(mgl32.Vec3{face.PositionX.Point, face.PositionY.Point, face.PositionZ.Point})
		// TODO: rotation & scale vector
		// vertex = vertex.Cross(rotation)
		// vertex = vertex.Mul(scale)

		textureCoordinate := mgl32.Vec2{0}
		if len(model.TextureCoordinates) > 0 {
			textureCoordinate = model.TextureCoordinates[idx]
		}
		normal := model.Normals[idx]

		if !eobj.hasVec3(eobj.uniqueVertices, vertex) {
			eobj.uniqueVertices = append(eobj.uniqueVertices, vertex)
			v += fmt.Sprintf("v %.6f %.6f %.6f", vertex.X(), vertex.Y(), vertex.Z()) + eobj.nlDelimiter
		}

		if len(face.MeshModel.TextureCoordinates) > 0 && !eobj.hasVec2(eobj.uniqueTextureCoordinates, textureCoordinate) {
			eobj.uniqueTextureCoordinates = append(eobj.uniqueTextureCoordinates, textureCoordinate)
			vt += fmt.Sprintf("vt %.6f %.6f", textureCoordinate.X(), textureCoordinate.Y()) + eobj.nlDelimiter
		}

		if !eobj.hasVec3(eobj.uniqueNormals, normal) {
			eobj.uniqueNormals = append(eobj.uniqueNormals, normal)
			vn += fmt.Sprintf("vn %.6f %.6f %.6f", normal.X(), normal.Y(), normal.Z()) + eobj.nlDelimiter
		}

		progressCounter++
		eobj.funcProgress((progressCounter / float32(totalProgress)) * 100.0)
	}

	meshData += v
	meshData += vt
	meshData += vn

	triangleFace := ""
	vCounterT := 0
	for k := 0; k < len(model.Indices); k++ {
		vCounterT++
		j := model.Indices[k]

		vertex := model.Vertices[j]
		vertex = vertex.Add(mgl32.Vec3{face.PositionX.Point, face.PositionY.Point, face.PositionZ.Point})
		// TODO: rotation & scale vector
		// vertex = vertex.Mul(rotation)
		// vertex = vertex.Mul(scale)

		v := eobj.findIndexVec3(eobj.uniqueVertices, vertex) + 1
		vn := eobj.findIndexVec3(eobj.uniqueNormals, model.Normals[j]) + 1

		vt := int32(-1)
		if len(model.TextureCoordinates) > 0 {
			vt = eobj.findIndexVec2(eobj.uniqueTextureCoordinates, model.TextureCoordinates[j])
		}

		if vt > -1 {
			triangleFace += fmt.Sprintf(" %d/%d/%d", v, (vt + 1), vn)
		} else {
			triangleFace += fmt.Sprintf(" %d//%d", v, vn)
		}

		if vCounterT%3 == 0 {
			f += "f" + triangleFace + eobj.nlDelimiter
			triangleFace = ""
		}

		progressCounter++
		eobj.funcProgress((progressCounter / float32(totalProgress)) * 100.0)
	}

	meshData += "usemtl " + model.MaterialTitle + eobj.nlDelimiter
	meshData += "s off" + eobj.nlDelimiter
	meshData += f

	return meshData
}

func (eobj *ExporterObj) hasVec2(haystack []mgl32.Vec2, needle mgl32.Vec2) bool {
	for i := 0; i < len(haystack); i++ {
		v := haystack[i]
		if v.X() == needle.X() && v.Y() == needle.Y() {
			return true
		}
	}
	return false
}

func (eobj *ExporterObj) hasVec3(haystack []mgl32.Vec3, needle mgl32.Vec3) bool {
	for i := 0; i < len(haystack); i++ {
		v := haystack[i]
		if v.X() == needle.X() && v.Y() == needle.Y() && v.Z() == needle.Z() {
			return true
		}
	}
	return false
}

func (eobj *ExporterObj) findIndexVec2(haystack []mgl32.Vec2, needle mgl32.Vec2) int32 {
	for i := int32(0); i < int32(len(haystack)); i++ {
		v := haystack[i]
		if v.X() == needle.X() && v.Y() == needle.Y() {
			return i
		}
	}
	return 0
}

func (eobj *ExporterObj) findIndexVec3(haystack []mgl32.Vec3, needle mgl32.Vec3) int32 {
	for i := int32(0); i < int32(len(haystack)); i++ {
		v := haystack[i]
		if v.X() == needle.X() && v.Y() == needle.Y() && v.Z() == needle.Z() {
			return i
		}
	}
	return 0
}
