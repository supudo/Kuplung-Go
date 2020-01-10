package settings

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/types"
	"gopkg.in/yaml.v2"
)

// RenderingSettings ...
type RenderingSettings struct {
	MatrixProjection mgl32.Mat4
	MatrixCamera     mgl32.Mat4

	General struct {
		ShowCube        bool    `yaml:"ShowCube"`
		Fov             float32 `yaml:"Fov"`
		RatioWidth      float32 `yaml:"RatioWidth"`
		RatioHeight     float32 `yaml:"RatioHeight"`
		PlaneClose      float32 `yaml:"PlaneClose"`
		PlaneFar        float32 `yaml:"PlaneFar"`
		GammaCoeficient float32 `yaml:"GammaCoeficient"`

		ShowPickRays       bool    `yaml:"ShowPickRays"`
		ShowPickRaysSingle bool    `yaml:"ShowPickRaysSingle"`
		RayAnimate         bool    `yaml:"RayAnimate"`
		RayOriginX         float32 `yaml:"RayOriginX"`
		RayOriginY         float32 `yaml:"RayOriginY"`
		RayOriginZ         float32 `yaml:"RayOriginZ"`
		RayOriginXS        string
		RayOriginYS        string
		RayOriginZS        string
		RayDraw            bool    `yaml:"RayDraw"`
		RayDirectionX      float32 `yaml:"RayDirectionX"`
		RayDirectionY      float32 `yaml:"RayDirectionY"`
		RayDirectionZ      float32 `yaml:"RayDirectionZ"`
		RayDirectionXS     string
		RayDirectionYS     string
		RayDirectionZS     string

		OcclusionCulling       bool `yaml:"OcclusionCulling"`
		RenderingDepth         bool
		SelectedViewModelSkin  types.ViewModelSkin
		ShowBoundingBox        bool
		BoundingBoxRefresh     bool
		BoundingBoxPadding     float32
		OutlineColor           mgl32.Vec4
		OutlineColorPickerOpen bool
		OutlineThickness       float32

		VertexSphereVisible         bool
		VertexSphereColorPickerOpen bool
		VertexSphereIsSphere        bool
		VertexSphereShowWireframes  bool
		VertexSphereRadius          float32
		VertexSphereSegments        int32
		VertexSphereColor           mgl32.Vec4

		ShowAllVisualArtefacts bool
	} `yaml:"General"`

	Axis struct {
		ShowAxisHelpers bool `yaml:"ShowAxisHelpers"`
		ShowZAxis       bool `yaml:"ShowZAxis"`
	} `yaml:"Axis"`

	Grid struct {
		WorldGridSizeSquares    int32 `yaml:"WorldGridSizeSquares"`
		WorldGridFixedWithWorld bool  `yaml:"WorldGridFixedWithWorld"`
		ShowGrid                bool  `yaml:"ShowGrid"`
		ActAsMirror             bool  `yaml:"ActAsMirror"`
	} `yaml:"Grid"`

	SkyBox struct {
		SkyboxSelectedItem int32
	}
}

var instantiatedRendering *RenderingSettings
var onceRendering sync.Once

// GetRenderingSettings singleton for our application settings
func GetRenderingSettings() *RenderingSettings {
	onceRendering.Do(func() {
		as := InitRenderingSettings()
		instantiatedRendering = &as
	})
	return instantiatedRendering
}

// InitRenderingSettings will initialize application settings
func InitRenderingSettings() RenderingSettings {
	var rSettings RenderingSettings

	dir, err := os.Getwd()
	if err != nil {
		LogError("Rendering Settings error: %v", err)
	}

	if runtime.GOOS == "darwin" {
		dir += "/../Resources/resources/"
	} else if runtime.GOOS == "windows" {
		dir += "./"
	} else {
		// TODO: other platforms
	}

	appConfig, err := ioutil.ReadFile(dir + "Kuplung_RenderingSettings.yaml")
	if err != nil {
		LogError("Rendering Settings error: %v", err)
	}

	err = yaml.Unmarshal(appConfig, &rSettings)
	if err != nil {
		LogError("Rendering Settings error: %v", err)
	}

	rSettings.General.RayOriginXS = fmt.Sprintf("%f", rSettings.General.RayOriginX)
	rSettings.General.RayOriginYS = fmt.Sprintf("%f", rSettings.General.RayOriginY)
	rSettings.General.RayOriginZS = fmt.Sprintf("%f", rSettings.General.RayOriginZ)
	rSettings.General.RayDirectionXS = fmt.Sprintf("%f", rSettings.General.RayDirectionX)
	rSettings.General.RayDirectionYS = fmt.Sprintf("%f", rSettings.General.RayDirectionY)
	rSettings.General.RayDirectionZS = fmt.Sprintf("%f", rSettings.General.RayDirectionZ)
	rSettings.General.RenderingDepth = false

	rSettings.MatrixProjection = mgl32.Perspective(mgl32.DegToRad(rSettings.General.Fov), rSettings.General.RatioWidth/rSettings.General.RatioHeight, rSettings.General.PlaneClose, rSettings.General.PlaneFar)
	rSettings.MatrixCamera = mgl32.Ident4()

	rSettings.General.SelectedViewModelSkin = types.ViewModelSkinRendered

	return rSettings
}

// ResetRenderSettings ...
func ResetRenderSettings() {
	rSettings := GetRenderingSettings()

	rSettings.General.ShowCube = false

	rSettings.General.Fov = 45.0
	rSettings.General.RatioWidth = 4.0
	rSettings.General.RatioHeight = 3.0
	rSettings.General.PlaneClose = 1.0
	rSettings.General.PlaneFar = 1000.0

	rSettings.General.GammaCoeficient = 1.0

	rSettings.Axis.ShowAxisHelpers = true
	rSettings.Axis.ShowZAxis = true

	rSettings.Grid.WorldGridSizeSquares = 30
	rSettings.Grid.WorldGridFixedWithWorld = true
	rSettings.Grid.ShowGrid = true
	rSettings.Grid.ActAsMirror = false

	rSettings.General.ShowPickRays = false
	rSettings.General.ShowPickRays = false
	rSettings.General.ShowPickRaysSingle = true
	rSettings.General.RayAnimate = false
	rSettings.General.RayOriginX = 0.0
	rSettings.General.RayOriginY = 0.0
	rSettings.General.RayOriginZ = 0.0
	rSettings.General.RayOriginXS = "0.0"
	rSettings.General.RayOriginYS = "0.0"
	rSettings.General.RayOriginZS = "0.0"
	rSettings.General.RayDraw = true
	rSettings.General.RayDirectionX = 0.0
	rSettings.General.RayDirectionY = 0.0
	rSettings.General.RayDirectionZ = 0.0
	rSettings.General.RayDirectionXS = "0.0"
	rSettings.General.RayDirectionYS = "0.0"
	rSettings.General.RayDirectionZS = "0.0"

	rSettings.General.RenderingDepth = false

	rSettings.MatrixProjection = mgl32.Perspective(mgl32.DegToRad(rSettings.General.Fov), rSettings.General.RatioWidth/rSettings.General.RatioHeight, rSettings.General.PlaneClose, rSettings.General.PlaneFar)
	rSettings.MatrixCamera = mgl32.Ident4()

	rSettings.General.OcclusionCulling = false

	rSettings.General.OutlineColor = mgl32.Vec4{0.0, 0.0, 0.0, 0.0}
	rSettings.General.OutlineColorPickerOpen = false
	rSettings.General.OutlineThickness = 1.01
	rSettings.General.ShowBoundingBox = false
	rSettings.General.BoundingBoxRefresh = false
	rSettings.General.BoundingBoxPadding = 0.01
	rSettings.General.ShowAllVisualArtefacts = true
}

// SaveRenderingSettings will save the settings back to yaml file
func SaveRenderingSettings() {
	rsett := GetRenderingSettings()

	dir, err := os.Getwd()
	if err != nil {
		LogError("Rendering Settings error: %v", err)
	}

	if runtime.GOOS == "darwin" {
		dir += "/../Resources/resources/"
	} else if runtime.GOOS == "windows" {
		dir += "./"
	} else {
		// TODO: other platforms
	}

	data, err := yaml.Marshal(&rsett)
	if err != nil {
		LogError("Rendering Settings save error: %v", err)
	}

	err = ioutil.WriteFile(dir+"Kuplung_RenderingSettings.yaml", data, 0644)
	if err != nil {
		LogError("Rendering Settings save error: %v", err)
	}
}

// KuplungPrintObjModels ..
func KuplungPrintObjModels(models []types.MeshModel, byIndices, shouldExit bool) {
	for i := 0; i < len(models); i++ {
		m := models[i]
		LogWarn("model.ID = %v", m.ID)
		LogWarn("model.countIndices = %v (%v)", m.CountIndices, MaxElement(m.Indices))
		LogWarn("model.countNormals = %v (%v)", m.CountNormals, (m.CountNormals * 3))
		LogWarn("model.countTextureCoordinates = %v (%v)", m.CountTextureCoordinates, (m.CountTextureCoordinates * 2))
		LogWarn("model.countVertices = %v (%v)", m.CountVertices, (m.CountVertices * 3))
		LogWarn("model.MaterialTitle = %v", m.MaterialTitle)
		LogWarn("model.ModelTitle = %v", m.ModelTitle)

		if byIndices {
			LogWarn("m.geometry :")
			for j := 0; j < len(m.Indices); j++ {
				idx := m.Indices[j]
				geom := fmt.Sprintf("index = %v ---> ", idx)
				vert := m.Vertices[idx]
				tc := m.TextureCoordinates[idx]
				n := m.Normals[idx]
				geom += fmt.Sprintf("vertex = [%v, %v, %v]", vert.X(), vert.Y(), vert.Z())
				geom += fmt.Sprintf(", uv = [%v, %v]", tc.X(), tc.Y())
				geom += fmt.Sprintf(", normal = [%v, %v, %v]", n.X(), n.Y(), n.Z())
				//LogWarn("%v", geom)
			}
		} else {
			var verts string
			for j := 0; j < len(m.Vertices); j++ {
				v := m.Vertices[j]
				verts += fmt.Sprintf("[%v, %v, %v], ", v.X(), v.Y(), v.Z())
			}
			LogWarn("model.vertices : %v", verts)

			var uvs string
			for j := 0; j < len(m.TextureCoordinates); j++ {
				uvs += fmt.Sprintf("[%v, %v], ", m.TextureCoordinates[j].X(), m.TextureCoordinates[j].Y())
			}
			LogWarn("model.texture_coordinates : %v", uvs)

			var normals string
			for j := 0; j < len(m.Normals); j++ {
				n := m.Normals[j]
				normals += fmt.Sprintf("[%f, %f, %f], ", n.X(), n.Y(), n.Z())
			}
			LogWarn("model.normals : %v", normals)

			var indices string
			for j := 0; j < len(m.Indices); j++ {
				indices += fmt.Sprintf("%v, ", m.Indices[j])
			}
			LogWarn("model.indices : %v", indices)
		}

		LogWarn("model.ModelMaterial.MaterialID = %v", m.ModelMaterial.MaterialID)
		LogWarn("model.ModelMaterial.MaterialTitle = %v", m.ModelMaterial.MaterialTitle)

		LogWarn("model.ModelMaterial.AmbientColor = %v, %v, %v", m.ModelMaterial.AmbientColor.X(), m.ModelMaterial.AmbientColor.Y(), m.ModelMaterial.AmbientColor.Z())
		LogWarn("model.ModelMaterial.DiffuseColor = %v, %v, %v", m.ModelMaterial.DiffuseColor.X(), m.ModelMaterial.DiffuseColor.Y(), m.ModelMaterial.DiffuseColor.Z())
		LogWarn("model.ModelMaterial.SpecularColor = %v, %v, %v", m.ModelMaterial.SpecularColor.X(), m.ModelMaterial.SpecularColor.Y(), m.ModelMaterial.SpecularColor.Z())
		LogWarn("model.ModelMaterial.EmissionColor = %v, %v, %v", m.ModelMaterial.EmissionColor.X(), m.ModelMaterial.EmissionColor.Y(), m.ModelMaterial.EmissionColor.Z())

		LogWarn("model.ModelMaterial.SpecularExp = %v", m.ModelMaterial.SpecularExp)
		LogWarn("model.ModelMaterial.Transparency = %v", m.ModelMaterial.Transparency)
		LogWarn("model.ModelMaterial.OpticalDensity = %v", m.ModelMaterial.OpticalDensity)
		LogWarn("model.ModelMaterial.IlluminationMode = %v", m.ModelMaterial.IlluminationMode)

		LogWarn("model.ModelMaterial.textures_ambient.Filename = %v", m.ModelMaterial.TextureAmbient.Filename)
		LogWarn("model.ModelMaterial.textures_diffuse.Filename = %v", m.ModelMaterial.TextureDiffuse.Filename)
		LogWarn("model.ModelMaterial.textures_specular.Filename = %v", m.ModelMaterial.TextureSpecular.Filename)
		LogWarn("model.ModelMaterial.textures_specularExp.Filename = %v", m.ModelMaterial.TextureSpecularExp.Filename)
		LogWarn("model.ModelMaterial.textures_dissolve.Filename = %v", m.ModelMaterial.TextureDissolve.Filename)
		LogWarn("model.ModelMaterial.textures_bump.Filename = %v", m.ModelMaterial.TextureBump.Filename)
		LogWarn("model.ModelMaterial.textures_displacement.Filename = %v", m.ModelMaterial.TextureDisplacement.Filename)
	}
	if shouldExit {
		LogError("--------")
	} else {
		LogWarn("--------")
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

// MinMax ...
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
