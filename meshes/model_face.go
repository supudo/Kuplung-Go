package meshes

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
	"github.com/supudo/Kuplung-Go/utilities"
)

// ModelFace ...
type ModelFace struct {
	window interfaces.Window

	shaderProgram uint32
	GLVAO         uint32

	VertexSphereVisible, VertexSphereIsSphere, VertexSphereShowWireframes bool

	VertexSphereRadius   float32
	VertexSphereSegments int32
	VertexSphereColor    mgl32.Vec4

	VboTextureDiffuse                         uint32
	VboTextureAmbient, VboTextureSpecular     uint32
	VboTextureSpecularExp, VboTextureDissolve uint32
	VboTextureBump, VboTextureDisplacement    uint32

	HasTextureDiffuse                         bool
	HasTextureAmbient, HasTextureSpecular     bool
	HasTextureSpecularExp, HasTextureDissolve bool
	HasTextureBump, HasTextureDisplacement    bool

	OccQuery uint32

	AssetsFolder string

	ModelViewSkin types.ViewModelSkin

	MeshModel   types.MeshModel
	MatrixModel mgl32.Mat4

	SettingsDeferredRender bool
	SettingCelShading      bool
	SettingWireframe       bool
	SettingUseTessellation bool
	SettingUseCullFace     bool
	ShowMaterialEditor     bool
	IsModelSelected        bool

	SettingAlpha                    float32
	SettingTessellationSubdivision  uint32
	PositionX, PositionY, PositionZ types.ObjectCoordinate
	ScaleX, ScaleY, ScaleZ          types.ObjectCoordinate
	RotateX, RotateY, RotateZ       types.ObjectCoordinate
	DisplaceX, DisplaceY, DisplaceZ types.ObjectCoordinate

	SettingMaterialRefraction  types.ObjectCoordinate
	SettingMaterialSpecularExp types.ObjectCoordinate

	SettingLightPosition  mgl32.Vec3
	SettingLightDirection mgl32.Vec3
	SettingLightAmbient   mgl32.Vec3
	SettingLightDiffuse   mgl32.Vec3
	SettingLightSpecular  mgl32.Vec3

	OutlineColor     mgl32.Vec4
	OutlineThickness float32

	SettingLightStrengthAmbient  float32
	SettingLightStrengthDiffuse  float32
	SettingLightStrengthSpecular float32
	SettingLightingPassDrawMode  uint32

	MaterialIlluminationModel uint32
	SettingParallaxMapping    bool

	MaterialAmbient         types.MaterialColor
	MaterialDiffuse         types.MaterialColor
	MaterialSpecular        types.MaterialColor
	MaterialEmission        types.MaterialColor
	DisplacementHeightScale types.ObjectCoordinate

	EffectGBlurMode   int32
	EffectGBlurRadius types.ObjectCoordinate
	EffectGBlurWidth  types.ObjectCoordinate

	EffectBloomDoBloom     bool
	EffectBloomWeightA     float32
	EffectBloomWeightB     float32
	EffectBloomWeightC     float32
	EffectBloomWeightD     float32
	EffectBloomVignette    float32
	EffectBloomVignetteAtt float32

	EffectToneMappingACESFilmRec2020 bool
	EffectHDRTonemapping             bool

	SettingShowShadows bool

	SettingRenderingPBR          bool
	SettingRenderingPBRMetallic  float32
	SettingRenderingPBRRoughness float32
	SettingRenderingPBRAO        float32

	// view skin
	SettingModelViewSkin        types.ViewModelSkin
	SolidLightSkinMaterialColor mgl32.Vec3
	SolidLightSkinAmbient       mgl32.Vec3
	SolidLightSkinDiffuse       mgl32.Vec3
	SolidLightSkinSpecular      mgl32.Vec3

	SolidLightSkinAmbientStrength float32
	SolidLightSkinDiffuseStrength float32
	SlidLightSkinSpecularStrength float32
}

// NewModelFace ...
func NewModelFace(window interfaces.Window, model types.MeshModel) *ModelFace {
	mesh := &ModelFace{}
	mesh.window = window
	mesh.MeshModel = model
	mesh.InitProperties()
	return mesh
}

// InitProperties ...
func (mesh *ModelFace) InitProperties() {
	mesh.AssetsFolder = ""

	mesh.SettingCelShading = false
	mesh.SettingWireframe = false
	mesh.SettingAlpha = 1.0
	mesh.ShowMaterialEditor = false
	mesh.SettingsDeferredRender = false
	mesh.SettingUseCullFace = false

	mesh.HasTextureAmbient = false
	mesh.HasTextureSpecular = false
	mesh.HasTextureSpecularExp = false
	mesh.HasTextureDissolve = false
	mesh.HasTextureBump = false
	mesh.HasTextureDisplacement = false

	mesh.PositionX = types.ObjectCoordinate{Animate: false, Point: 0.0}
	mesh.PositionY = types.ObjectCoordinate{Animate: false, Point: 0.0}
	mesh.PositionZ = types.ObjectCoordinate{Animate: false, Point: 0.0}

	mesh.ScaleX = types.ObjectCoordinate{Animate: false, Point: 1.0}
	mesh.ScaleY = types.ObjectCoordinate{Animate: false, Point: 1.0}
	mesh.ScaleZ = types.ObjectCoordinate{Animate: false, Point: 1.0}

	mesh.RotateX = types.ObjectCoordinate{Animate: false, Point: 0.0}
	mesh.RotateY = types.ObjectCoordinate{Animate: false, Point: 0.0}
	mesh.RotateZ = types.ObjectCoordinate{Animate: false, Point: 0.0}

	mesh.DisplaceX = types.ObjectCoordinate{Animate: false, Point: 0.0}
	mesh.DisplaceY = types.ObjectCoordinate{Animate: false, Point: 0.0}
	mesh.DisplaceZ = types.ObjectCoordinate{Animate: false, Point: 0.0}

	mesh.MatrixModel = mgl32.Ident4()

	mesh.SettingMaterialRefraction = types.ObjectCoordinate{Animate: false, Point: mesh.MeshModel.ModelMaterial.OpticalDensity}
	mesh.SettingMaterialSpecularExp = types.ObjectCoordinate{Animate: false, Point: mesh.MeshModel.ModelMaterial.SpecularExp}

	mesh.SettingLightPosition = mgl32.Vec3{0.0, 0.0, 0.0}
	mesh.SettingLightDirection = mgl32.Vec3{0.0, 0.0, 0.0}
	mesh.SettingLightAmbient = mgl32.Vec3{0.0, 0.0, 0.0}
	mesh.SettingLightDiffuse = mgl32.Vec3{0.0, 0.0, 0.0}
	mesh.SettingLightSpecular = mgl32.Vec3{0.0, 0.0, 0.0}
	mesh.SettingLightStrengthAmbient = 1.0
	mesh.SettingLightStrengthDiffuse = 1.0
	mesh.SettingLightStrengthSpecular = 1.0
	mesh.SettingTessellationSubdivision = 1
	mesh.SettingLightingPassDrawMode = 1
	mesh.OutlineColor = mgl32.Vec4{0.0, 0.0, 0.0, 0.0}
	mesh.OutlineThickness = 1.0

	mesh.MaterialIlluminationModel = mesh.MeshModel.ModelMaterial.IlluminationMode
	mesh.SettingParallaxMapping = false

	mesh.MaterialAmbient = types.MaterialColor{ColorPickerOpen: false, Animate: false, Strength: 1.0, Color: mesh.MeshModel.ModelMaterial.AmbientColor}
	mesh.MaterialDiffuse = types.MaterialColor{ColorPickerOpen: false, Animate: false, Strength: 1.0, Color: mesh.MeshModel.ModelMaterial.DiffuseColor}
	mesh.MaterialSpecular = types.MaterialColor{ColorPickerOpen: false, Animate: false, Strength: 1.0, Color: mesh.MeshModel.ModelMaterial.SpecularColor}
	mesh.MaterialEmission = types.MaterialColor{ColorPickerOpen: false, Animate: false, Strength: 1.0, Color: mesh.MeshModel.ModelMaterial.EmissionColor}
	mesh.DisplacementHeightScale = types.ObjectCoordinate{Animate: false, Point: 0.0}

	mesh.EffectGBlurMode = -1
	mesh.EffectGBlurRadius = types.ObjectCoordinate{Animate: false, Point: 0.0}
	mesh.EffectGBlurWidth = types.ObjectCoordinate{Animate: false, Point: 0.0}

	mesh.EffectBloomDoBloom = false
	mesh.EffectBloomWeightA = 0.0
	mesh.EffectBloomWeightB = 0.0
	mesh.EffectBloomWeightC = 0.0
	mesh.EffectBloomWeightD = 0.0
	mesh.EffectBloomVignette = 0.0
	mesh.EffectBloomVignetteAtt = 0.0

	mesh.EffectToneMappingACESFilmRec2020 = false
	mesh.EffectHDRTonemapping = false

	mesh.SettingShowShadows = true

	mesh.SettingRenderingPBR = true
	mesh.SettingRenderingPBRMetallic = 0.1
	mesh.SettingRenderingPBRRoughness = 0.1
	mesh.SettingRenderingPBRAO = 0.1
}

// InitBuffers ...
func (mesh *ModelFace) InitBuffers() {
	gl := mesh.window.OpenGL()

	mesh.GLVAO = gl.GenVertexArrays(1)[0]

	gl.BindVertexArray(mesh.GLVAO)

	vboVertices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(mesh.MeshModel.Vertices)*3*4, gl.Ptr(mesh.MeshModel.Vertices), oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

	vboNormals := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboNormals)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(mesh.MeshModel.Normals)*3*4, gl.Ptr(mesh.MeshModel.Normals), oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

	if len(mesh.MeshModel.TextureCoordinates) > 0 {
		vboTextureCoordinates := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboTextureCoordinates)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(mesh.MeshModel.TextureCoordinates)*2*4, gl.Ptr(mesh.MeshModel.TextureCoordinates), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 2, oglconsts.FLOAT, false, 2*4, gl.PtrOffset(0))

		if len(mesh.MeshModel.ModelMaterial.TextureAmbient.Image) > 0 {
			mesh.VboTextureAmbient = engine.LoadTextureRepeat(mesh.window.OpenGL(), mesh.MeshModel.ModelMaterial.TextureAmbient.Image)
			mesh.HasTextureAmbient = true
		}
		if len(mesh.MeshModel.ModelMaterial.TextureDiffuse.Image) > 0 {
			mesh.VboTextureDiffuse = engine.LoadTextureRepeat(mesh.window.OpenGL(), mesh.MeshModel.ModelMaterial.TextureDiffuse.Image)
			mesh.HasTextureDiffuse = true
		}
		if len(mesh.MeshModel.ModelMaterial.TextureSpecular.Image) > 0 {
			mesh.VboTextureSpecular = engine.LoadTextureRepeat(mesh.window.OpenGL(), mesh.MeshModel.ModelMaterial.TextureSpecular.Image)
			mesh.HasTextureSpecular = true
		}
		if len(mesh.MeshModel.ModelMaterial.TextureSpecularExp.Image) > 0 {
			mesh.VboTextureSpecularExp = engine.LoadTextureRepeat(mesh.window.OpenGL(), mesh.MeshModel.ModelMaterial.TextureSpecularExp.Image)
			mesh.HasTextureSpecularExp = true
		}
		if len(mesh.MeshModel.ModelMaterial.TextureDissolve.Image) > 0 {
			mesh.VboTextureDissolve = engine.LoadTextureRepeat(mesh.window.OpenGL(), mesh.MeshModel.ModelMaterial.TextureDissolve.Image)
			mesh.HasTextureDissolve = true
		}
		if len(mesh.MeshModel.ModelMaterial.TextureBump.Image) > 0 {
			mesh.VboTextureBump = engine.LoadTextureRepeat(mesh.window.OpenGL(), mesh.MeshModel.ModelMaterial.TextureBump.Image)
			mesh.HasTextureBump = true
		}
		if len(mesh.MeshModel.ModelMaterial.TextureDisplacement.Image) > 0 {
			mesh.VboTextureDisplacement = engine.LoadTextureRepeat(mesh.window.OpenGL(), mesh.MeshModel.ModelMaterial.TextureDisplacement.Image)
			mesh.HasTextureDisplacement = true
		}
	}

	vboIndices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, vboIndices)
	gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, int(mesh.MeshModel.CountIndices)*4, gl.Ptr(mesh.MeshModel.Indices), oglconsts.STATIC_DRAW)

	if len(mesh.MeshModel.ModelMaterial.TextureBump.Image) > 0 && len(mesh.MeshModel.Vertices) > 0 && len(mesh.MeshModel.TextureCoordinates) > 0 && len(mesh.MeshModel.Normals) > 0 {
		tangents, bitangents := utilities.ComputeTangentBasis(mesh.MeshModel.TextureCoordinates, mesh.MeshModel.Vertices, mesh.MeshModel.Normals)

		// tangents
		vboTangents := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboTangents)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(tangents)*3*4, gl.Ptr(tangents), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(3)
		gl.VertexAttribPointer(3, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		// bitangents
		vboBitangents := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboBitangents)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(bitangents)*3*4, gl.Ptr(bitangents), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(4)
		gl.VertexAttribPointer(4, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))
	}

	mesh.OccQuery = gl.GenQueries(1)[0]

	gl.BindVertexArray(0)

	gl.DeleteBuffers([]uint32{vboVertices, vboNormals, vboIndices})
}

// Render ...
func (mesh *ModelFace) Render(useTessellation bool) {
	gl := mesh.window.OpenGL()
	rsett := settings.GetRenderingSettings()

	mesh.SettingUseTessellation = useTessellation

	if mesh.SettingWireframe {
		gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.LINE)
	}

	if rsett.General.OcclusionCulling {
		gl.BeginConditionalRender(mesh.OccQuery, oglconsts.QUERY_BY_REGION_WAIT)
	}

	gl.BindVertexArray(mesh.GLVAO)

	if useTessellation {
		gl.DrawElements(oglconsts.PATCHES, mesh.MeshModel.CountIndices, oglconsts.UNSIGNED_INT, 0)
	} else {
		gl.DrawElements(oglconsts.TRIANGLES, mesh.MeshModel.CountIndices, oglconsts.UNSIGNED_INT, 0)
	}

	if rsett.General.OcclusionCulling {
		gl.EndConditionalRender()
	}

	gl.BindVertexArray(0)

	if mesh.SettingWireframe {
		gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.FILL)
	}

	// matrixBB := rsett.MatrixProjection.Mul(rsett.MatrixCamera).Mul(mgl32.Ident4())

	if rsett.General.ShowBoundingBox && mesh.IsModelSelected {
		// mesh.boundingBox.Render(matrixBB, mesh.SettingOutlineColor)
	}

	// if (this->vertexSphereVisible) {
	//   this->vertexSphere->isSphere = this->vertexSphereIsSphere;
	//   this->vertexSphere->showWireframes = this->vertexSphereShowWireframes;
	//   this->vertexSphere->initBuffers(this->meshModel, this->vertexSphereSegments, this->vertexSphereRadius);
	//   this->vertexSphere->render(matrixBB, vertexSphereColor);
	// }

	// if (this->getOptionsSelected() && (this->Setting_Gizmo_Rotate || this->Setting_Gizmo_Translate || this->Setting_Gizmo_Scale)) {
	//   ImGuizmo::Enable(true);
	//   ImGuizmo::OPERATION gizmo_operation = ImGuizmo::TRANSLATE;
	//   if (this->Setting_Gizmo_Rotate)
	//     gizmo_operation = ImGuizmo::ROTATE;
	//   else if (this->Setting_Gizmo_Scale)
	//     gizmo_operation = ImGuizmo::SCALE;
	//   glm::mat4 mtx = glm::mat4(1.0);
	//   ImGuiIO& io = ImGui::GetIO();
	//   ImGuizmo::SetRect(0, 0, io.DisplaySize.x, io.DisplaySize.y);
	//   ImGuizmo::Manipulate(glm::value_ptr(this->matrixCamera), glm::value_ptr(this->matrixProjection), gizmo_operation, ImGuizmo::LOCAL, glm::value_ptr(this->matrixModel), glm::value_ptr(mtx));

	//   glm::vec3 scale;
	//   glm::quat rotation;
	//   glm::vec3 translation;
	//   glm::vec3 skew;
	//   glm::vec4 perspective;
	//   glm::decompose(mtx, scale, rotation, translation, skew, perspective);

	//   if (this->Setting_Gizmo_Translate) {
	//     this->positionX->point += translation.x;
	//     this->positionY->point += translation.y;
	//     this->positionZ->point += translation.z;
	//   }

	//   if (this->Setting_Gizmo_Rotate) {
	//     this->rotateX->point += glm::degrees(rotation.x);
	//     this->rotateY->point += glm::degrees(rotation.y);
	//     this->rotateZ->point += glm::degrees(rotation.z);
	//   }

	//   if (this->Setting_Gizmo_Scale) {
	//     this->scaleX->point *= scale.x;
	//     this->scaleY->point *= scale.y;
	//     this->scaleZ->point *= scale.z;
	//   }
	// }
}

// Dispose ...
func (mesh *ModelFace) Dispose() {
	gl := mesh.window.OpenGL()
	gl.DeleteProgram(mesh.GLVAO)
}
