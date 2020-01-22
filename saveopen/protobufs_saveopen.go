package saveopen

import (
	fmt "fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	proto "github.com/golang/protobuf/proto"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/objects"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
	"github.com/supudo/Kuplung-Go/utilities"
)

// ProtoBufsSaveOpen ...
type ProtoBufsSaveOpen struct {
	doProgress func(float32)

	fileNameSettings string
	fileNameScene    string
}

// NewProtoBufsSaveOpen ...
func NewProtoBufsSaveOpen(doProgress func(float32)) *ProtoBufsSaveOpen {
	pm := &ProtoBufsSaveOpen{}
	pm.doProgress = doProgress
	return pm
}

// Save ...
func (pm *ProtoBufsSaveOpen) Save(file *types.FBEntity, meshModelFaces []*meshes.ModelFace, lights []*objects.Light, rprops types.RenderProperties, cam *objects.Camera, grid *objects.WorldGrid) {
	fileName := file.Path
	if !strings.HasSuffix(fileName, ".kuplung") {
		fileName += ".kuplung"
	}
	pm.fileNameSettings = fileName + ".settings"
	pm.fileNameScene = fileName + ".scene"

	pm.storeRenderingSettings(lights, rprops, cam, grid)
	pm.storeObjects(meshModelFaces)

	zfiles := []string{pm.fileNameSettings, pm.fileNameScene}
	utilities.ZipFiles(fileName, zfiles)

	if err := os.Remove(pm.fileNameSettings); err != nil {
		settings.LogWarn("[SaveOpen-ProtoBufs] [Save] Can't delete temp file : %v!", pm.fileNameSettings)
	}
	if err := os.Remove(pm.fileNameScene); err != nil {
		settings.LogWarn("[SaveOpen-ProtoBufs] [Save] Can't delete temp file : %v!", pm.fileNameScene)
	}
}

// Open ...
func (pm *ProtoBufsSaveOpen) Open(file *types.FBEntity, window interfaces.Window, systemModels map[string]types.MeshModel, faces *[]*meshes.ModelFace, lights *[]*objects.Light, rprops *types.RenderProperties, cam *objects.Camera, grid *objects.WorldGrid) {
	pfiles := utilities.UnzipFiles(file.Path, filepath.Dir(file.Path))
	if len(pfiles) == 2 {
		pm.openRenderingSettings(pfiles[0], window, systemModels, lights, rprops, cam, grid)
		pm.readObjects(pfiles[1], window, faces)
		if err := os.Remove(pfiles[0]); err != nil {
			settings.LogWarn("[SaveOpen-ProtoBufs] [Open] Can't delete temp file : %v!", pfiles[0])
		}
		if err := os.Remove(pfiles[1]); err != nil {
			settings.LogWarn("[SaveOpen-ProtoBufs] [Open] Can't delete temp file : %v!", pfiles[1])
		}
	}
}

func (pm *ProtoBufsSaveOpen) openRenderingSettings(filename string, window interfaces.Window, systemModels map[string]types.MeshModel, lights *[]*objects.Light, rprops *types.RenderProperties, cam *objects.Camera, grid *objects.WorldGrid) {
	fSettingsHandle, err := ioutil.ReadFile(filename)
	if err != nil {
		settings.LogWarn("[SaveOpen-ProtoBufs] [Open] Error reading file %v : %v", filename, err)
		return
	}
	gs := &GUISettings{}
	if err := proto.Unmarshal(fSettingsHandle, gs); err != nil {
		settings.LogWarn("[SaveOpen-ProtoBufs] [Open] Can't decode general settings: %v", err)
		return
	}

	rsett := settings.GetRenderingSettings()

	// Render Settings
	rsett.General.ShowCube = *gs.ShowCube
	rsett.General.Fov = *gs.Fov
	rsett.General.RatioWidth = *gs.RatioWidth
	rsett.General.RatioHeight = *gs.RatioHeight
	rsett.General.PlaneClose = *gs.PlaneClose
	rsett.General.PlaneFar = *gs.PlaneFar
	rsett.General.GammaCoeficient = *gs.GammaCoeficient

	rsett.General.ShowPickRays = *gs.ShowPickRays
	rsett.General.ShowPickRaysSingle = *gs.ShowPickRaysSingle
	rsett.General.RayAnimate = *gs.RayAnimate
	rsett.General.RayOriginX = *gs.RayOriginX
	rsett.General.RayOriginX = *gs.RayOriginY
	rsett.General.RayOriginZ = *gs.RayOriginZ
	rsett.General.RayOriginXS = *gs.RayOriginXS
	rsett.General.RayOriginYS = *gs.RayOriginYS
	rsett.General.RayOriginZS = *gs.RayOriginZS
	rsett.General.RayDraw = *gs.RayDraw
	rsett.General.RayDirectionX = *gs.RayDirectionX
	rsett.General.RayDirectionY = *gs.RayDirectionY
	rsett.General.RayDirectionZ = *gs.RayDirectionZ
	rsett.General.RayDirectionXS = *gs.RayDirectionXS
	rsett.General.RayDirectionYS = *gs.RayDirectionYS
	rsett.General.RayDirectionZS = *gs.RayDirectionZS

	rsett.General.OcclusionCulling = *gs.OcclusionCulling
	rsett.General.RenderingDepth = *gs.RenderingDepth
	rsett.General.SelectedViewModelSkin = types.ViewModelSkin(*gs.SelectedViewModelSkin)
	rsett.General.ShowBoundingBox = *gs.ShowBoundingBox
	rsett.General.BoundingBoxRefresh = *gs.BoundingBoxRefresh
	rsett.General.BoundingBoxPadding = *gs.BoundingBoxPadding
	rsett.General.OutlineColor = mgl32.Vec4{*gs.OutlineColor.X, *gs.OutlineColor.Y, *gs.OutlineColor.Z, 1}
	rsett.General.OutlineColorPickerOpen = *gs.OutlineColorPickerOpen
	rsett.General.OutlineThickness = *gs.OutlineThickness

	rsett.General.VertexSphereVisible = *gs.VertexSphereVisible
	rsett.General.VertexSphereColorPickerOpen = *gs.VertexSphereColorPickerOpen
	rsett.General.VertexSphereIsSphere = *gs.VertexSphereIsSphere
	rsett.General.VertexSphereShowWireframes = *gs.VertexSphereShowWireframes
	rsett.General.VertexSphereRadius = *gs.VertexSphereRadius
	rsett.General.VertexSphereSegments = *gs.VertexSphereSegments
	rsett.General.VertexSphereColor = mgl32.Vec4{*gs.VertexSphereColor.X, *gs.VertexSphereColor.Y, *gs.VertexSphereColor.Z, 1}

	rsett.General.ShowAllVisualArtefacts = *gs.ShowAllVisualArtefacts

	rsett.Axis.ShowZAxis = *gs.ShowZAxis

	rsett.Grid.WorldGridSizeSquares = *gs.WorldGridSizeSquares
	rsett.Grid.WorldGridFixedWithWorld = *gs.WorldGridFixedWithWorld
	rsett.Grid.ShowGrid = *gs.ShowGrid
	rsett.Grid.ActAsMirror = *gs.ActAsMirror

	rsett.SkyBox.SkyboxSelectedItem = *gs.SkyboxSelectedItem

	rsett.Defered.DeferredTestMode = *gs.DeferredTestMode
	rsett.Defered.DeferredTestLights = *gs.DeferredTestLights
	rsett.Defered.DeferredRandomizeLightPositions = *gs.DeferredRandomizeLightPositions
	rsett.Defered.LightingPassDrawMode = *gs.LightingPassDrawMode
	rsett.Defered.DeferredTestLightsNumber = *gs.DeferredTestLightsNumber
	rsett.Defered.DeferredAmbientStrength = *gs.DeferredAmbientStrength

	rsett.General.DebugShadowTexture = *gs.DebugShadowTexture

	// Render Properties
	rprops.UIAmbientLightX = *gs.UIAmbientLightX
	rprops.UIAmbientLightY = *gs.UIAmbientLightY
	rprops.UIAmbientLightZ = *gs.UIAmbientLightZ

	rprops.SolidLightDirectionX = *gs.SolidLightDirectionX
	rprops.SolidLightDirectionY = *gs.SolidLightDirectionY
	rprops.SolidLightDirectionZ = *gs.SolidLightDirectionZ

	rprops.SolidLightMaterialColor = mgl32.Vec3{*gs.SolidLightMaterialColor.X, *gs.SolidLightMaterialColor.Y, *gs.SolidLightMaterialColor.Z}
	rprops.SolidLightAmbient = mgl32.Vec3{*gs.SolidLightAmbient.X, *gs.SolidLightAmbient.Y, *gs.SolidLightAmbient.Z}
	rprops.SolidLightDiffuse = mgl32.Vec3{*gs.SolidLightDiffuse.X, *gs.SolidLightDiffuse.Y, *gs.SolidLightDiffuse.Z}
	rprops.SolidLightSpecular = mgl32.Vec3{*gs.SolidLightSpecular.X, *gs.SolidLightSpecular.Y, *gs.SolidLightSpecular.Z}

	rprops.SolidLightAmbientStrength = *gs.SolidLightAmbientStrength
	rprops.SolidLightDiffuseStrength = *gs.SolidLightDiffuseStrength
	rprops.SolidLightSpecularStrength = *gs.SolidLightSpecularStrength

	rprops.SolidLightMaterialColorColorPicker = *gs.SolidLightMaterialColorColorPicker
	rprops.SolidLightAmbientColorPicker = *gs.SolidLightAmbientColorPicker
	rprops.SolidLightDiffuseColorPicker = *gs.SolidLightDiffuseColorPicker
	rprops.SolidLightSpecularColorPicker = *gs.SolidLightSpecularColorPicker

	// Camera
	c := *gs.Camera
	cam.CameraPosition = mgl32.Vec3{*c.CameraPosition.X, *c.CameraPosition.Y, *c.CameraPosition.Z}
	cam.EyeSettings.ViewEye = mgl32.Vec3{*c.View_Eye.X, *c.View_Eye.Y, *c.View_Eye.Z}
	cam.EyeSettings.ViewCenter = mgl32.Vec3{*c.View_Center.X, *c.View_Center.Y, *c.View_Center.Z}
	cam.EyeSettings.ViewUp = mgl32.Vec3{*c.View_Up.X, *c.View_Up.Y, *c.View_Up.Z}
	cam.PositionX = types.ObjectCoordinate{Animate: *c.PositionX.Animate, Point: *c.PositionX.Point}
	cam.PositionY = types.ObjectCoordinate{Animate: *c.PositionY.Animate, Point: *c.PositionY.Point}
	cam.PositionZ = types.ObjectCoordinate{Animate: *c.PositionZ.Animate, Point: *c.PositionZ.Point}
	cam.RotateX = types.ObjectCoordinate{Animate: *c.RotateX.Animate, Point: *c.RotateX.Point}
	cam.RotateY = types.ObjectCoordinate{Animate: *c.RotateY.Animate, Point: *c.RotateY.Point}
	cam.RotateZ = types.ObjectCoordinate{Animate: *c.RotateZ.Animate, Point: *c.RotateZ.Point}
	cam.RotateCenterX = types.ObjectCoordinate{Animate: *c.RotateCenterX.Animate, Point: *c.RotateCenterX.Point}
	cam.RotateCenterY = types.ObjectCoordinate{Animate: *c.RotateCenterY.Animate, Point: *c.RotateCenterY.Point}
	cam.RotateCenterZ = types.ObjectCoordinate{Animate: *c.RotateCenterZ.Animate, Point: *c.RotateCenterZ.Point}

	// Grid
	g := *gs.Grid
	grid.ActAsMirror = *g.ActAsMirror
	grid.GridSize = *g.GridSize
	grid.Transparency = *g.Transparency
	grid.PositionX = types.ObjectCoordinate{Animate: *g.PositionX.Animate, Point: *g.PositionX.Point}
	grid.PositionY = types.ObjectCoordinate{Animate: *g.PositionY.Animate, Point: *g.PositionY.Point}
	grid.PositionZ = types.ObjectCoordinate{Animate: *g.PositionZ.Animate, Point: *g.PositionZ.Point}
	grid.RotateX = types.ObjectCoordinate{Animate: *g.RotateX.Animate, Point: *g.RotateX.Point}
	grid.RotateY = types.ObjectCoordinate{Animate: *g.RotateY.Animate, Point: *g.RotateY.Point}
	grid.RotateZ = types.ObjectCoordinate{Animate: *g.RotateZ.Animate, Point: *g.RotateZ.Point}
	grid.ScaleX = types.ObjectCoordinate{Animate: *g.ScaleX.Animate, Point: *g.ScaleX.Point}
	grid.ScaleY = types.ObjectCoordinate{Animate: *g.ScaleY.Animate, Point: *g.ScaleY.Point}
	grid.ScaleZ = types.ObjectCoordinate{Animate: *g.ScaleZ.Animate, Point: *g.ScaleZ.Point}

	// Lights
	*lights = []*objects.Light{}
	for i := 0; i < len(gs.Lights); i++ {
		l := gs.Lights[i]

		var lShape types.LightSourceType
		var lTitle, lDescription, lModel string
		switch *l.Type {
		case 0:
			lShape = types.LightSourceTypeDirectional
			lTitle = fmt.Sprintf("Directional %v", i)
			lDescription = "Directional area light source"
			lModel = "light_directional"
		case 1:
			lShape = types.LightSourceTypePoint
			lTitle = fmt.Sprintf("Point %v", i)
			lDescription = "Omnidirectional point light source"
			lModel = "light_point"
		case 2:
			lShape = types.LightSourceTypeSpot
			lTitle = fmt.Sprintf("Spot %v", i)
			lDescription = "Directional cone light source"
			lModel = "light_spot"
		}
		ll := objects.InitLight(window)
		ll.InitProperties(lShape)
		ll.Title = lTitle
		ll.Description = lDescription
		ll.SetModel(systemModels[lModel])
		ll.InitBuffers()

		*lights = append(*lights, ll)
	}
}

func (pm *ProtoBufsSaveOpen) readObjects(filename string, window interfaces.Window, faces *[]*meshes.ModelFace) {
	fModelsHandle, err := ioutil.ReadFile(filename)
	if err != nil {
		settings.LogError("[SaveOpen-ProtoBufs] [readObjects] Error reading file %v : %v", filename, err)
		return
	}
	gs := &Scene{}
	if err := proto.Unmarshal(fModelsHandle, gs); err != nil {
		settings.LogError("[SaveOpen-ProtoBufs] [readObjects] Can't decode general settings: %v", err)
		return
	}

	sett := settings.GetSettings()
	*faces = []*meshes.ModelFace{}
	var i int32
	for i = 0; i < int32(len(gs.Models)); i++ {
		gm := gs.Models[i]
		gmo := gm.MeshObject
		gmom := gm.MeshObject.ModelMaterial

		// MeshModel
		mm := types.MeshModel{}
		mm.ID = uint32(*gmo.ID)
		mm.File = *gmo.File
		mm.FilePath = *gmo.FilePath

		mm.ModelTitle = *gmo.ModelTitle
		mm.MaterialTitle = *gmo.MaterialTitle

		mm.CountVertices = *gmo.CountVertices
		mm.CountTextureCoordinates = *gmo.CountTextureCoordinates
		mm.CountNormals = *gmo.CountNormals
		mm.CountIndices = *gmo.CountIndices

		for j := 0; j < len(gmo.Vertices); j++ {
			mm.Vertices = append(mm.Vertices, mgl32.Vec3{*gmo.Vertices[j].X, *gmo.Vertices[j].Y, *gmo.Vertices[j].Z})
		}
		for j := 0; j < len(gmo.TextureCoordinates); j++ {
			mm.TextureCoordinates = append(mm.TextureCoordinates, mgl32.Vec2{*gmo.TextureCoordinates[j].X, *gmo.TextureCoordinates[j].Y})
		}
		for j := 0; j < len(gmo.Normals); j++ {
			mm.Normals = append(mm.Normals, mgl32.Vec3{*gmo.Normals[j].X, *gmo.Normals[j].Y, *gmo.Normals[j].Z})
		}
		mm.Indices = gmo.Indices

		// MeshModelMaterial
		mmm := types.MeshModelMaterial{}
		mmm.MaterialID = uint32(*gmom.MaterialID)
		mmm.MaterialTitle = *gmom.MaterialTitle

		mmm.SpecularExp = *gmom.SpecularExp

		mmm.AmbientColor = mgl32.Vec3{*gmom.AmbientColor.X, *gmom.AmbientColor.Y, *gmom.AmbientColor.Z}
		mmm.DiffuseColor = mgl32.Vec3{*gmom.DiffuseColor.X, *gmom.DiffuseColor.Y, *gmom.DiffuseColor.Z}
		mmm.SpecularColor = mgl32.Vec3{*gmom.SpecularColor.X, *gmom.SpecularColor.Y, *gmom.SpecularColor.Z}
		mmm.EmissionColor = mgl32.Vec3{*gmom.EmissionColor.X, *gmom.EmissionColor.Y, *gmom.EmissionColor.Z}

		mmm.Transparency = *gmom.Transparency
		mmm.IlluminationMode = *gmom.IlluminationMode
		mmm.OpticalDensity = *gmom.OpticalDensity

		mmmtia := types.MeshMaterialTextureImage{
			Filename:   *gmom.TextureAmbient.Filename,
			Image:      *gmom.TextureAmbient.Image,
			Width:      *gmom.TextureAmbient.Width,
			Height:     *gmom.TextureAmbient.Height,
			UseTexture: *gmom.TextureAmbient.UseTexture,
			Commands:   gmom.TextureAmbient.Commands}
		mmm.TextureAmbient = mmmtia
		mmmtid := types.MeshMaterialTextureImage{
			Filename:   *gmom.TextureDiffuse.Filename,
			Image:      *gmom.TextureDiffuse.Image,
			Width:      *gmom.TextureDiffuse.Width,
			Height:     *gmom.TextureDiffuse.Height,
			UseTexture: *gmom.TextureDiffuse.UseTexture,
			Commands:   gmom.TextureDiffuse.Commands}
		mmm.TextureDiffuse = mmmtid
		mmmtis := types.MeshMaterialTextureImage{
			Filename:   *gmom.TextureSpecular.Filename,
			Image:      *gmom.TextureSpecular.Image,
			Width:      *gmom.TextureSpecular.Width,
			Height:     *gmom.TextureSpecular.Height,
			UseTexture: *gmom.TextureSpecular.UseTexture,
			Commands:   gmom.TextureSpecular.Commands}
		mmm.TextureSpecular = mmmtis
		mmmtise := types.MeshMaterialTextureImage{
			Filename:   *gmom.TextureSpecularExp.Filename,
			Image:      *gmom.TextureSpecularExp.Image,
			Width:      *gmom.TextureSpecularExp.Width,
			Height:     *gmom.TextureSpecularExp.Height,
			UseTexture: *gmom.TextureSpecularExp.UseTexture,
			Commands:   gmom.TextureSpecularExp.Commands}
		mmm.TextureSpecularExp = mmmtise
		mmmtidi := types.MeshMaterialTextureImage{
			Filename:   *gmom.TextureDissolve.Filename,
			Image:      *gmom.TextureDissolve.Image,
			Width:      *gmom.TextureDissolve.Width,
			Height:     *gmom.TextureDissolve.Height,
			UseTexture: *gmom.TextureDissolve.UseTexture,
			Commands:   gmom.TextureDissolve.Commands}
		mmm.TextureDissolve = mmmtidi
		mmmtib := types.MeshMaterialTextureImage{
			Filename:   *gmom.TextureBump.Filename,
			Image:      *gmom.TextureBump.Image,
			Width:      *gmom.TextureBump.Width,
			Height:     *gmom.TextureBump.Height,
			UseTexture: *gmom.TextureBump.UseTexture,
			Commands:   gmom.TextureBump.Commands}
		mmm.TextureBump = mmmtib
		mmmtids := types.MeshMaterialTextureImage{
			Filename:   *gmom.TextureDisplacement.Filename,
			Image:      *gmom.TextureDisplacement.Image,
			Width:      *gmom.TextureDisplacement.Width,
			Height:     *gmom.TextureDisplacement.Height,
			UseTexture: *gmom.TextureDisplacement.UseTexture,
			Commands:   gmom.TextureDisplacement.Commands}
		mmm.TextureDisplacement = mmmtids
		mm.ModelMaterial = mmm

		mesh := meshes.NewModelFace(window, mm)
		mesh.InitProperties()
		mesh.InitBuffers()

		mesh.ModelID = *gm.ModelID
		mesh.ModelViewSkin = types.ViewModelSkin(*gm.Setting_ModelViewSkin)

		mesh.DeferredRender = *gm.Settings_DeferredRender
		mesh.CelShading = *gm.Setting_CelShading
		mesh.Wireframe = *gm.Setting_Wireframe
		mesh.UseTessellation = *gm.Setting_UseTessellation
		mesh.UseCullFace = *gm.Setting_UseCullFace
		mesh.ShowMaterialEditor = *gm.ShowMaterialEditor

		mesh.Alpha = *gm.Setting_Alpha
		mesh.TessellationSubdivision = *gm.Setting_TessellationSubdivision
		mesh.PositionX = types.ObjectCoordinate{Animate: *gm.PositionX.Animate, Point: *gm.PositionX.Point}
		mesh.PositionY = types.ObjectCoordinate{Animate: *gm.PositionY.Animate, Point: *gm.PositionY.Point}
		mesh.PositionZ = types.ObjectCoordinate{Animate: *gm.PositionZ.Animate, Point: *gm.PositionZ.Point}
		mesh.ScaleX = types.ObjectCoordinate{Animate: *gm.ScaleX.Animate, Point: *gm.ScaleX.Point}
		mesh.ScaleY = types.ObjectCoordinate{Animate: *gm.ScaleY.Animate, Point: *gm.ScaleY.Point}
		mesh.ScaleZ = types.ObjectCoordinate{Animate: *gm.ScaleZ.Animate, Point: *gm.ScaleZ.Point}
		mesh.RotateX = types.ObjectCoordinate{Animate: *gm.RotateX.Animate, Point: *gm.RotateX.Point}
		mesh.RotateY = types.ObjectCoordinate{Animate: *gm.RotateY.Animate, Point: *gm.RotateY.Point}
		mesh.RotateZ = types.ObjectCoordinate{Animate: *gm.RotateZ.Animate, Point: *gm.RotateZ.Point}
		mesh.DisplaceX = types.ObjectCoordinate{Animate: *gm.DisplaceX.Animate, Point: *gm.DisplaceX.Point}
		mesh.DisplaceY = types.ObjectCoordinate{Animate: *gm.DisplaceY.Animate, Point: *gm.DisplaceY.Point}
		mesh.DisplaceZ = types.ObjectCoordinate{Animate: *gm.DisplaceZ.Animate, Point: *gm.DisplaceZ.Point}

		mesh.MaterialRefraction = types.ObjectCoordinate{Animate: *gm.Setting_MaterialRefraction.Animate, Point: *gm.Setting_MaterialRefraction.Point}
		mesh.MaterialSpecularExp = types.ObjectCoordinate{Animate: *gm.Setting_MaterialSpecularExp.Animate, Point: *gm.Setting_MaterialSpecularExp.Point}

		mesh.LightPosition = mgl32.Vec3{*gm.Setting_LightPosition.X, *gm.Setting_LightPosition.Y, *gm.Setting_LightPosition.Z}
		mesh.LightDirection = mgl32.Vec3{*gm.Setting_LightDirection.X, *gm.Setting_LightDirection.Y, *gm.Setting_LightDirection.Z}
		mesh.LightAmbient = mgl32.Vec3{*gm.Setting_LightAmbient.X, *gm.Setting_LightAmbient.Y, *gm.Setting_LightAmbient.Z}
		mesh.LightDiffuse = mgl32.Vec3{*gm.Setting_LightDiffuse.X, *gm.Setting_LightDiffuse.Y, *gm.Setting_LightDiffuse.Z}
		mesh.LightSpecular = mgl32.Vec3{*gm.Setting_LightSpecular.X, *gm.Setting_LightSpecular.Y, *gm.Setting_LightSpecular.Z}

		mesh.LightStrengthAmbient = *gm.Setting_LightStrengthAmbient
		mesh.LightStrengthDiffuse = *gm.Setting_LightStrengthDiffuse
		mesh.LightStrengthSpecular = *gm.Setting_LightStrengthSpecular
		mesh.LightingPassDrawMode = uint32(*gm.Setting_LightingPass_DrawMode)

		mesh.MaterialIlluminationModel = uint32(*gm.MaterialIlluminationModel)
		mesh.ParallaxMapping = *gm.Setting_ParallaxMapping

		mesh.MaterialAmbient = types.MaterialColor{
			ColorPickerOpen: *gm.MaterialAmbient.ColorPickerOpen,
			Animate:         *gm.MaterialAmbient.Animate,
			Strength:        *gm.MaterialAmbient.Strength,
			Color:           mgl32.Vec3{*gm.MaterialAmbient.Color.X, *gm.MaterialAmbient.Color.Y, *gm.MaterialAmbient.Color.Z}}
		mesh.MaterialDiffuse = types.MaterialColor{
			ColorPickerOpen: *gm.MaterialDiffuse.ColorPickerOpen,
			Animate:         *gm.MaterialDiffuse.Animate,
			Strength:        *gm.MaterialDiffuse.Strength,
			Color:           mgl32.Vec3{*gm.MaterialDiffuse.Color.X, *gm.MaterialDiffuse.Color.Y, *gm.MaterialDiffuse.Color.Z}}
		mesh.MaterialSpecular = types.MaterialColor{
			ColorPickerOpen: *gm.MaterialSpecular.ColorPickerOpen,
			Animate:         *gm.MaterialSpecular.Animate,
			Strength:        *gm.MaterialSpecular.Strength,
			Color:           mgl32.Vec3{*gm.MaterialSpecular.Color.X, *gm.MaterialSpecular.Color.Y, *gm.MaterialSpecular.Color.Z}}
		mesh.MaterialEmission = types.MaterialColor{
			ColorPickerOpen: *gm.MaterialEmission.ColorPickerOpen,
			Animate:         *gm.MaterialEmission.Animate,
			Strength:        *gm.MaterialEmission.Strength,
			Color:           mgl32.Vec3{*gm.MaterialEmission.Color.X, *gm.MaterialEmission.Color.Y, *gm.MaterialEmission.Color.Z}}

		mesh.DisplacementHeightScale = types.ObjectCoordinate{Animate: *gm.DisplacementHeightScale.Animate, Point: *gm.DisplacementHeightScale.Point}

		mesh.EffectGBlurMode = *gm.Effect_GBlur_Mode
		mesh.EffectGBlurRadius = types.ObjectCoordinate{Animate: *gm.Effect_GBlur_Radius.Animate, Point: *gm.Effect_GBlur_Radius.Point}
		mesh.EffectGBlurWidth = types.ObjectCoordinate{Animate: *gm.Effect_GBlur_Width.Animate, Point: *gm.Effect_GBlur_Width.Point}

		mesh.EffectBloomDoBloom = *gm.Effect_Bloom_DoBloom
		mesh.EffectBloomWeightA = *gm.Effect_Bloom_WeightA
		mesh.EffectBloomWeightB = *gm.Effect_Bloom_WeightB
		mesh.EffectBloomWeightC = *gm.Effect_Bloom_WeightC
		mesh.EffectBloomWeightD = *gm.Effect_Bloom_WeightD
		mesh.EffectBloomVignette = *gm.Effect_Bloom_Vignette
		mesh.EffectBloomVignetteAtt = *gm.Effect_Bloom_VignetteAtt

		mesh.EffectToneMappingACESFilmRec2020 = *gm.EffectToneMappingACESFilmRec2020
		mesh.EffectHDRTonemapping = *gm.EffectHDRTonemapping

		mesh.ShowShadows = *gm.ShowShadows

		mesh.RenderingPBR = *gm.RenderingPBR
		mesh.RenderingPBRMetallic = *gm.RenderingPBRMetallic
		mesh.RenderingPBRRoughness = *gm.RenderingPBRRoughness
		mesh.RenderingPBRAO = *gm.RenderingPBRAO

		mesh.SolidLightSkinMaterialColor = mgl32.Vec3{*gm.SolidLightSkin_MaterialColor.X, *gm.SolidLightSkin_MaterialColor.Y, *gm.SolidLightSkin_MaterialColor.Z}
		mesh.SolidLightSkinAmbient = mgl32.Vec3{*gm.SolidLightSkin_Ambient.X, *gm.SolidLightSkin_Ambient.Y, *gm.SolidLightSkin_Ambient.Z}
		mesh.SolidLightSkinDiffuse = mgl32.Vec3{*gm.SolidLightSkin_Diffuse.X, *gm.SolidLightSkin_Diffuse.Y, *gm.SolidLightSkin_Diffuse.Z}
		mesh.SolidLightSkinSpecular = mgl32.Vec3{*gm.SolidLightSkin_Specular.X, *gm.SolidLightSkin_Specular.Y, *gm.SolidLightSkin_Specular.Z}

		mesh.SolidLightSkinAmbientStrength = *gm.SolidLightSkin_Ambient_Strength
		mesh.SolidLightSkinDiffuseStrength = *gm.SolidLightSkin_Diffuse_Strength
		mesh.SlidLightSkinSpecularStrength = *gm.SolidLightSkin_Specular_Strength

		*faces = append(*faces, mesh)

		sett.MemSettings.TotalVertices += mesh.MeshModel.CountVertices
		sett.MemSettings.TotalIndices += mesh.MeshModel.CountIndices
		sett.MemSettings.TotalTriangles += mesh.MeshModel.CountVertices / 3
		sett.MemSettings.TotalFaces += mesh.MeshModel.CountVertices / 6
		sett.MemSettings.TotalObjects++
	}
}

func (pm *ProtoBufsSaveOpen) storeRenderingSettings(lights []*objects.Light, rprops types.RenderProperties, cam *objects.Camera, grid *objects.WorldGrid) {
	gs := &GUISettings{}

	// Render Settings
	rsett := settings.GetRenderingSettings()
	gs.ShowCube = proto.Bool(rsett.General.ShowCube)
	gs.Fov = proto.Float32(rsett.General.Fov)
	gs.RatioWidth = proto.Float32(rsett.General.RatioWidth)
	gs.RatioHeight = proto.Float32(rsett.General.RatioHeight)
	gs.PlaneClose = proto.Float32(rsett.General.PlaneClose)
	gs.PlaneFar = proto.Float32(rsett.General.PlaneFar)
	gs.GammaCoeficient = proto.Float32(rsett.General.GammaCoeficient)

	gs.ShowPickRays = proto.Bool(rsett.General.ShowPickRays)
	gs.ShowPickRaysSingle = proto.Bool(rsett.General.ShowPickRaysSingle)
	gs.RayAnimate = proto.Bool(rsett.General.RayAnimate)
	gs.RayOriginX = proto.Float32(rsett.General.RayOriginX)
	gs.RayOriginY = proto.Float32(rsett.General.RayOriginY)
	gs.RayOriginZ = proto.Float32(rsett.General.RayOriginZ)
	gs.RayOriginXS = proto.String(rsett.General.RayOriginXS)
	gs.RayOriginYS = proto.String(rsett.General.RayOriginYS)
	gs.RayOriginZS = proto.String(rsett.General.RayOriginZS)
	gs.RayDraw = proto.Bool(rsett.General.RayDraw)
	gs.RayDirectionX = proto.Float32(rsett.General.RayDirectionX)
	gs.RayDirectionY = proto.Float32(rsett.General.RayDirectionY)
	gs.RayDirectionZ = proto.Float32(rsett.General.RayDirectionZ)
	gs.RayDirectionXS = proto.String(rsett.General.RayDirectionXS)
	gs.RayDirectionYS = proto.String(rsett.General.RayDirectionYS)
	gs.RayDirectionZS = proto.String(rsett.General.RayDirectionZS)

	gs.OcclusionCulling = proto.Bool(rsett.General.OcclusionCulling)
	gs.RenderingDepth = proto.Bool(rsett.General.RenderingDepth)
	gs.SelectedViewModelSkin = proto.Uint32(uint32(rsett.General.SelectedViewModelSkin))
	gs.ShowBoundingBox = proto.Bool(rsett.General.ShowBoundingBox)
	gs.BoundingBoxRefresh = proto.Bool(rsett.General.BoundingBoxRefresh)
	gs.BoundingBoxPadding = proto.Float32(rsett.General.BoundingBoxPadding)
	gs.OutlineColor = &Vec4{
		X: proto.Float32(rsett.General.OutlineColor.X()),
		Y: proto.Float32(rsett.General.OutlineColor.Y()),
		Z: proto.Float32(rsett.General.OutlineColor.Z()),
		W: proto.Float32(rsett.General.OutlineColor.W()),
	}
	gs.OutlineColorPickerOpen = proto.Bool(rsett.General.OutlineColorPickerOpen)
	gs.OutlineThickness = proto.Float32(rsett.General.OutlineThickness)

	gs.VertexSphereVisible = proto.Bool(rsett.General.VertexSphereVisible)
	gs.VertexSphereColorPickerOpen = proto.Bool(rsett.General.VertexSphereColorPickerOpen)
	gs.VertexSphereIsSphere = proto.Bool(rsett.General.VertexSphereIsSphere)
	gs.VertexSphereShowWireframes = proto.Bool(rsett.General.VertexSphereShowWireframes)
	gs.VertexSphereRadius = proto.Float32(rsett.General.VertexSphereRadius)
	gs.VertexSphereSegments = proto.Int32(rsett.General.VertexSphereSegments)
	gs.VertexSphereColor = &Vec4{
		X: proto.Float32(rsett.General.VertexSphereColor.X()),
		Y: proto.Float32(rsett.General.VertexSphereColor.Y()),
		Z: proto.Float32(rsett.General.VertexSphereColor.Z()),
		W: proto.Float32(rsett.General.VertexSphereColor.W()),
	}

	gs.ShowAllVisualArtefacts = proto.Bool(rsett.General.ShowAllVisualArtefacts)

	gs.ShowZAxis = proto.Bool(rsett.Axis.ShowZAxis)

	gs.WorldGridSizeSquares = proto.Int32(rsett.Grid.WorldGridSizeSquares)
	gs.WorldGridFixedWithWorld = proto.Bool(rsett.Grid.WorldGridFixedWithWorld)
	gs.ShowGrid = proto.Bool(rsett.Grid.ShowGrid)
	gs.ActAsMirror = proto.Bool(rsett.Grid.ActAsMirror)

	gs.SkyboxSelectedItem = proto.Int32(rsett.SkyBox.SkyboxSelectedItem)

	gs.DeferredTestMode = proto.Bool(rsett.Defered.DeferredTestMode)
	gs.DeferredTestLights = proto.Bool(rsett.Defered.DeferredTestLights)
	gs.DeferredRandomizeLightPositions = proto.Bool(rsett.Defered.DeferredRandomizeLightPositions)
	gs.LightingPassDrawMode = proto.Int32(rsett.Defered.LightingPassDrawMode)
	gs.DeferredTestLightsNumber = proto.Int32(rsett.Defered.DeferredTestLightsNumber)
	gs.DeferredAmbientStrength = proto.Float32(rsett.Defered.DeferredAmbientStrength)

	gs.DebugShadowTexture = proto.Bool(rsett.General.DebugShadowTexture)

	// Render Properties
	gs.UIAmbientLightX = proto.Float32(rprops.UIAmbientLightX)
	gs.UIAmbientLightY = proto.Float32(rprops.UIAmbientLightY)
	gs.UIAmbientLightZ = proto.Float32(rprops.UIAmbientLightZ)

	gs.SolidLightDirectionX = proto.Float32(rprops.SolidLightDirectionX)
	gs.SolidLightDirectionY = proto.Float32(rprops.SolidLightDirectionY)
	gs.SolidLightDirectionZ = proto.Float32(rprops.SolidLightDirectionZ)

	gs.SolidLightMaterialColor = &Vec3{
		X: proto.Float32(rprops.SolidLightMaterialColor.X()),
		Y: proto.Float32(rprops.SolidLightMaterialColor.Y()),
		Z: proto.Float32(rprops.SolidLightMaterialColor.Z()),
	}
	gs.SolidLightAmbient = &Vec3{
		X: proto.Float32(rprops.SolidLightAmbient.X()),
		Y: proto.Float32(rprops.SolidLightAmbient.Y()),
		Z: proto.Float32(rprops.SolidLightAmbient.Z()),
	}
	gs.SolidLightDiffuse = &Vec3{
		X: proto.Float32(rprops.SolidLightDiffuse.X()),
		Y: proto.Float32(rprops.SolidLightDiffuse.Y()),
		Z: proto.Float32(rprops.SolidLightDiffuse.Z()),
	}
	gs.SolidLightSpecular = &Vec3{
		X: proto.Float32(rprops.SolidLightSpecular.X()),
		Y: proto.Float32(rprops.SolidLightSpecular.Y()),
		Z: proto.Float32(rprops.SolidLightSpecular.Z()),
	}

	gs.SolidLightAmbientStrength = proto.Float32(rprops.SolidLightAmbientStrength)
	gs.SolidLightDiffuseStrength = proto.Float32(rprops.SolidLightDiffuseStrength)
	gs.SolidLightSpecularStrength = proto.Float32(rprops.SolidLightSpecularStrength)

	gs.SolidLightMaterialColorColorPicker = proto.Bool(rprops.SolidLightMaterialColorColorPicker)
	gs.SolidLightAmbientColorPicker = proto.Bool(rprops.SolidLightAmbientColorPicker)
	gs.SolidLightDiffuseColorPicker = proto.Bool(rprops.SolidLightDiffuseColorPicker)
	gs.SolidLightSpecularColorPicker = proto.Bool(rprops.SolidLightSpecularColorPicker)

	// Camera
	c := &CameraSettings{}
	c.CameraPosition = &Vec3{X: proto.Float32(cam.CameraPosition.X()), Y: proto.Float32(cam.CameraPosition.Y()), Z: proto.Float32(cam.CameraPosition.Z())}
	c.View_Eye = &Vec3{X: proto.Float32(cam.EyeSettings.ViewEye.X()), Y: proto.Float32(cam.EyeSettings.ViewEye.Y()), Z: proto.Float32(cam.EyeSettings.ViewEye.Z())}
	c.View_Center = &Vec3{X: proto.Float32(cam.EyeSettings.ViewCenter.X()), Y: proto.Float32(cam.EyeSettings.ViewCenter.Y()), Z: proto.Float32(cam.EyeSettings.ViewCenter.Z())}
	c.View_Up = &Vec3{X: proto.Float32(cam.EyeSettings.ViewUp.X()), Y: proto.Float32(cam.EyeSettings.ViewUp.Y()), Z: proto.Float32(cam.EyeSettings.ViewUp.Z())}
	c.PositionX = &ObjectCoordinate{Animate: proto.Bool(cam.PositionX.Animate), Point: proto.Float32(cam.PositionX.Point)}
	c.PositionY = &ObjectCoordinate{Animate: proto.Bool(cam.PositionY.Animate), Point: proto.Float32(cam.PositionY.Point)}
	c.PositionZ = &ObjectCoordinate{Animate: proto.Bool(cam.PositionZ.Animate), Point: proto.Float32(cam.PositionZ.Point)}
	c.RotateX = &ObjectCoordinate{Animate: proto.Bool(cam.RotateX.Animate), Point: proto.Float32(cam.RotateX.Point)}
	c.RotateY = &ObjectCoordinate{Animate: proto.Bool(cam.RotateY.Animate), Point: proto.Float32(cam.RotateY.Point)}
	c.RotateZ = &ObjectCoordinate{Animate: proto.Bool(cam.RotateZ.Animate), Point: proto.Float32(cam.RotateZ.Point)}
	c.RotateCenterX = &ObjectCoordinate{Animate: proto.Bool(cam.RotateCenterX.Animate), Point: proto.Float32(cam.RotateCenterX.Point)}
	c.RotateCenterY = &ObjectCoordinate{Animate: proto.Bool(cam.RotateCenterY.Animate), Point: proto.Float32(cam.RotateCenterY.Point)}
	c.RotateCenterZ = &ObjectCoordinate{Animate: proto.Bool(cam.RotateCenterZ.Animate), Point: proto.Float32(cam.RotateCenterZ.Point)}
	gs.Camera = c

	// Grid
	g := &GridSettings{}
	g.ActAsMirror = proto.Bool(grid.ActAsMirror)
	g.GridSize = proto.Int32(grid.GridSize)
	g.PositionX = &ObjectCoordinate{Animate: proto.Bool(grid.PositionX.Animate), Point: proto.Float32(grid.PositionX.Point)}
	g.PositionY = &ObjectCoordinate{Animate: proto.Bool(grid.PositionY.Animate), Point: proto.Float32(grid.PositionY.Point)}
	g.PositionZ = &ObjectCoordinate{Animate: proto.Bool(grid.PositionZ.Animate), Point: proto.Float32(grid.PositionZ.Point)}
	g.RotateX = &ObjectCoordinate{Animate: proto.Bool(grid.RotateX.Animate), Point: proto.Float32(grid.RotateX.Point)}
	g.RotateY = &ObjectCoordinate{Animate: proto.Bool(grid.RotateY.Animate), Point: proto.Float32(grid.RotateY.Point)}
	g.RotateZ = &ObjectCoordinate{Animate: proto.Bool(grid.RotateZ.Animate), Point: proto.Float32(grid.RotateZ.Point)}
	g.ScaleX = &ObjectCoordinate{Animate: proto.Bool(grid.ScaleX.Animate), Point: proto.Float32(grid.ScaleX.Point)}
	g.ScaleY = &ObjectCoordinate{Animate: proto.Bool(grid.ScaleY.Animate), Point: proto.Float32(grid.ScaleY.Point)}
	g.ScaleZ = &ObjectCoordinate{Animate: proto.Bool(grid.ScaleZ.Animate), Point: proto.Float32(grid.ScaleZ.Point)}
	g.Transparency = proto.Float32(grid.Transparency)
	gs.Grid = g

	// Lights
	for i := 0; i < len(lights); i++ {
		lo := lights[i]
		l := &LightObject{}
		l.Title = proto.String(lo.Title)
		l.Description = proto.String(lo.Description)
		l.Type = proto.Int32(int32(lo.LightType))
		l.ShowLampObject = proto.Bool(lo.ShowLampObject)
		l.ShowLampDirection = proto.Bool(lo.ShowLampDirection)
		l.ShowInWire = proto.Bool(lo.ShowInWire)

		l.PositionX = &ObjectCoordinate{Animate: proto.Bool(lo.PositionX.Animate), Point: proto.Float32(lo.PositionX.Point)}
		l.PositionY = &ObjectCoordinate{Animate: proto.Bool(lo.PositionY.Animate), Point: proto.Float32(lo.PositionY.Point)}
		l.PositionZ = &ObjectCoordinate{Animate: proto.Bool(lo.PositionZ.Animate), Point: proto.Float32(lo.PositionZ.Point)}

		l.DirectionX = &ObjectCoordinate{Animate: proto.Bool(lo.DirectionX.Animate), Point: proto.Float32(lo.DirectionX.Point)}
		l.DirectionY = &ObjectCoordinate{Animate: proto.Bool(lo.DirectionY.Animate), Point: proto.Float32(lo.DirectionY.Point)}
		l.DirectionZ = &ObjectCoordinate{Animate: proto.Bool(lo.DirectionZ.Animate), Point: proto.Float32(lo.DirectionZ.Point)}

		l.ScaleX = &ObjectCoordinate{Animate: proto.Bool(lo.ScaleX.Animate), Point: proto.Float32(lo.ScaleX.Point)}
		l.ScaleY = &ObjectCoordinate{Animate: proto.Bool(lo.ScaleY.Animate), Point: proto.Float32(lo.ScaleY.Point)}
		l.ScaleZ = &ObjectCoordinate{Animate: proto.Bool(lo.ScaleZ.Animate), Point: proto.Float32(lo.ScaleZ.Point)}

		l.RotateX = &ObjectCoordinate{Animate: proto.Bool(lo.RotateX.Animate), Point: proto.Float32(lo.RotateX.Point)}
		l.RotateY = &ObjectCoordinate{Animate: proto.Bool(lo.RotateY.Animate), Point: proto.Float32(lo.RotateY.Point)}
		l.RotateZ = &ObjectCoordinate{Animate: proto.Bool(lo.RotateZ.Animate), Point: proto.Float32(lo.RotateZ.Point)}

		l.RotateCenterX = &ObjectCoordinate{Animate: proto.Bool(lo.RotateCenterX.Animate), Point: proto.Float32(lo.RotateCenterX.Point)}
		l.RotateCenterY = &ObjectCoordinate{Animate: proto.Bool(lo.RotateCenterY.Animate), Point: proto.Float32(lo.RotateCenterY.Point)}
		l.RotateCenterZ = &ObjectCoordinate{Animate: proto.Bool(lo.RotateCenterZ.Animate), Point: proto.Float32(lo.RotateCenterZ.Point)}

		lva := &Vec3{X: proto.Float32(lo.Ambient.Color.X()), Y: proto.Float32(lo.Ambient.Color.Y()), Z: proto.Float32(lo.Ambient.Color.Z())}
		lvd := &Vec3{X: proto.Float32(lo.Diffuse.Color.X()), Y: proto.Float32(lo.Diffuse.Color.Y()), Z: proto.Float32(lo.Diffuse.Color.Z())}
		lvs := &Vec3{X: proto.Float32(lo.Specular.Color.X()), Y: proto.Float32(lo.Specular.Color.Y()), Z: proto.Float32(lo.Specular.Color.Z())}
		l.Ambient = &MaterialColor{Animate: proto.Bool(lo.Ambient.Animate), Color: lva, ColorPickerOpen: proto.Bool(lo.Ambient.ColorPickerOpen), Strength: proto.Float32(lo.Ambient.Strength)}
		l.Diffuse = &MaterialColor{Animate: proto.Bool(lo.Diffuse.Animate), Color: lvd, ColorPickerOpen: proto.Bool(lo.Diffuse.ColorPickerOpen), Strength: proto.Float32(lo.Diffuse.Strength)}
		l.Specular = &MaterialColor{Animate: proto.Bool(lo.Specular.Animate), Color: lvs, ColorPickerOpen: proto.Bool(lo.Specular.ColorPickerOpen), Strength: proto.Float32(lo.Specular.Strength)}

		l.LCutOff = &ObjectCoordinate{Animate: proto.Bool(lo.LCutOff.Animate), Point: proto.Float32(lo.LCutOff.Point)}
		l.LOuterCutOff = &ObjectCoordinate{Animate: proto.Bool(lo.LOuterCutOff.Animate), Point: proto.Float32(lo.LOuterCutOff.Point)}
		l.LConstant = &ObjectCoordinate{Animate: proto.Bool(lo.LConstant.Animate), Point: proto.Float32(lo.LConstant.Point)}
		l.LLinear = &ObjectCoordinate{Animate: proto.Bool(lo.LLinear.Animate), Point: proto.Float32(lo.LLinear.Point)}
		l.LQuadratic = &ObjectCoordinate{Animate: proto.Bool(lo.LQuadratic.Animate), Point: proto.Float32(lo.LQuadratic.Point)}

		gs.Lights = append(gs.Lights, l)
	}

	// Save
	data, err := proto.Marshal(gs)
	if err != nil {
		settings.LogWarn("[SaveOpen-ProtoBufs] [storeRenderingSettings] Marshalling error: %v", err)
		return
	}
	err = ioutil.WriteFile(pm.fileNameSettings, data, 0644)
	if err != nil {
		settings.LogWarn("[SaveOpen-ProtoBufs] [storeRenderingSettings] Can't save byte data error: %v", err)
		return
	}
}

func (pm *ProtoBufsSaveOpen) storeObjects(meshModelFaces []*meshes.ModelFace) {
	gs := &Scene{}

	for i := 0; i < len(meshModelFaces); i++ {
		m := meshModelFaces[i]
		mm := &MeshModel{}

		mm.ModelID = proto.Int32(m.ModelID)
		mm.Settings_DeferredRender = proto.Bool(m.DeferredRender)
		mm.Setting_CelShading = proto.Bool(m.CelShading)
		mm.Setting_Wireframe = proto.Bool(m.Wireframe)
		mm.Setting_UseTessellation = proto.Bool(m.UseTessellation)
		mm.Setting_UseCullFace = proto.Bool(m.UseCullFace)
		mm.Setting_Alpha = proto.Float32(m.Alpha)
		mm.Setting_TessellationSubdivision = proto.Int32(m.TessellationSubdivision)
		mm.PositionX = &ObjectCoordinate{Animate: proto.Bool(m.PositionX.Animate), Point: proto.Float32(m.PositionX.Point)}
		mm.PositionY = &ObjectCoordinate{Animate: proto.Bool(m.PositionY.Animate), Point: proto.Float32(m.PositionY.Point)}
		mm.PositionZ = &ObjectCoordinate{Animate: proto.Bool(m.PositionZ.Animate), Point: proto.Float32(m.PositionZ.Point)}
		mm.ScaleX = &ObjectCoordinate{Animate: proto.Bool(m.ScaleX.Animate), Point: proto.Float32(m.ScaleX.Point)}
		mm.ScaleY = &ObjectCoordinate{Animate: proto.Bool(m.ScaleY.Animate), Point: proto.Float32(m.ScaleY.Point)}
		mm.ScaleZ = &ObjectCoordinate{Animate: proto.Bool(m.ScaleZ.Animate), Point: proto.Float32(m.ScaleZ.Point)}
		mm.RotateX = &ObjectCoordinate{Animate: proto.Bool(m.RotateX.Animate), Point: proto.Float32(m.RotateX.Point)}
		mm.RotateY = &ObjectCoordinate{Animate: proto.Bool(m.RotateY.Animate), Point: proto.Float32(m.RotateY.Point)}
		mm.RotateZ = &ObjectCoordinate{Animate: proto.Bool(m.RotateZ.Animate), Point: proto.Float32(m.RotateZ.Point)}
		mm.DisplaceX = &ObjectCoordinate{Animate: proto.Bool(m.DisplaceX.Animate), Point: proto.Float32(m.DisplaceX.Point)}
		mm.DisplaceY = &ObjectCoordinate{Animate: proto.Bool(m.DisplaceY.Animate), Point: proto.Float32(m.DisplaceY.Point)}
		mm.DisplaceZ = &ObjectCoordinate{Animate: proto.Bool(m.DisplaceZ.Animate), Point: proto.Float32(m.DisplaceZ.Point)}
		mm.Setting_MaterialRefraction = &ObjectCoordinate{Animate: proto.Bool(m.MaterialRefraction.Animate), Point: proto.Float32(m.MaterialRefraction.Point)}
		mm.Setting_MaterialSpecularExp = &ObjectCoordinate{Animate: proto.Bool(m.MaterialSpecularExp.Animate), Point: proto.Float32(m.MaterialSpecularExp.Point)}

		mm.Setting_ModelViewSkin = proto.Int32(int32(m.ModelViewSkin))
		mm.SolidLightSkin_MaterialColor = &Vec3{X: proto.Float32(m.SolidLightSkinMaterialColor.X()), Y: proto.Float32(m.SolidLightSkinMaterialColor.Y()), Z: proto.Float32(m.SolidLightSkinMaterialColor.Z())}
		mm.SolidLightSkin_Ambient = &Vec3{X: proto.Float32(m.SolidLightSkinAmbient.X()), Y: proto.Float32(m.SolidLightSkinAmbient.Y()), Z: proto.Float32(m.SolidLightSkinAmbient.Z())}
		mm.SolidLightSkin_Diffuse = &Vec3{X: proto.Float32(m.SolidLightSkinDiffuse.X()), Y: proto.Float32(m.SolidLightSkinDiffuse.Y()), Z: proto.Float32(m.SolidLightSkinDiffuse.Z())}
		mm.SolidLightSkin_Specular = &Vec3{X: proto.Float32(m.SolidLightSkinSpecular.X()), Y: proto.Float32(m.SolidLightSkinSpecular.Y()), Z: proto.Float32(m.SolidLightSkinSpecular.Z())}
		mm.SolidLightSkin_Ambient_Strength = proto.Float32(m.SolidLightSkinAmbientStrength)
		mm.SolidLightSkin_Diffuse_Strength = proto.Float32(m.SolidLightSkinDiffuseStrength)
		mm.SolidLightSkin_Specular_Strength = proto.Float32(m.SlidLightSkinSpecularStrength)

		mm.Setting_LightPosition = &Vec3{X: proto.Float32(m.LightPosition.X()), Y: proto.Float32(m.LightPosition.Y()), Z: proto.Float32(m.LightPosition.Z())}
		mm.Setting_LightDirection = &Vec3{X: proto.Float32(m.LightDirection.X()), Y: proto.Float32(m.LightDirection.Y()), Z: proto.Float32(m.LightDirection.Z())}
		mm.Setting_LightAmbient = &Vec3{X: proto.Float32(m.LightAmbient.X()), Y: proto.Float32(m.LightAmbient.Y()), Z: proto.Float32(m.LightAmbient.Z())}
		mm.Setting_LightDiffuse = &Vec3{X: proto.Float32(m.LightDiffuse.X()), Y: proto.Float32(m.LightDiffuse.Y()), Z: proto.Float32(m.LightDiffuse.Z())}
		mm.Setting_LightSpecular = &Vec3{X: proto.Float32(m.LightSpecular.X()), Y: proto.Float32(m.LightSpecular.Y()), Z: proto.Float32(m.LightSpecular.Z())}
		mm.Setting_LightStrengthAmbient = proto.Float32(m.LightStrengthAmbient)
		mm.Setting_LightStrengthDiffuse = proto.Float32(m.LightStrengthDiffuse)
		mm.Setting_LightStrengthSpecular = proto.Float32(m.LightStrengthSpecular)

		mm.MaterialIlluminationModel = proto.Int32(int32(m.MaterialIlluminationModel))
		mm.DisplacementHeightScale = &ObjectCoordinate{Animate: proto.Bool(m.DisplacementHeightScale.Animate), Point: proto.Float32(m.DisplacementHeightScale.Point)}
		mm.ShowMaterialEditor = proto.Bool(m.ShowMaterialEditor)

		lvma := &Vec3{X: proto.Float32(m.MaterialAmbient.Color.X()), Y: proto.Float32(m.MaterialAmbient.Color.Y()), Z: proto.Float32(m.MaterialAmbient.Color.Z())}
		lvmd := &Vec3{X: proto.Float32(m.MaterialDiffuse.Color.X()), Y: proto.Float32(m.MaterialDiffuse.Color.Y()), Z: proto.Float32(m.MaterialDiffuse.Color.Z())}
		lvms := &Vec3{X: proto.Float32(m.MaterialSpecular.Color.X()), Y: proto.Float32(m.MaterialSpecular.Color.Y()), Z: proto.Float32(m.MaterialSpecular.Color.Z())}
		lvme := &Vec3{X: proto.Float32(m.MaterialEmission.Color.X()), Y: proto.Float32(m.MaterialEmission.Color.Y()), Z: proto.Float32(m.MaterialEmission.Color.Z())}
		mm.MaterialAmbient = &MaterialColor{Animate: proto.Bool(m.MaterialAmbient.Animate), Color: lvma, ColorPickerOpen: proto.Bool(m.MaterialAmbient.ColorPickerOpen), Strength: proto.Float32(m.MaterialAmbient.Strength)}
		mm.MaterialDiffuse = &MaterialColor{Animate: proto.Bool(m.MaterialDiffuse.Animate), Color: lvmd, ColorPickerOpen: proto.Bool(m.MaterialDiffuse.ColorPickerOpen), Strength: proto.Float32(m.MaterialDiffuse.Strength)}
		mm.MaterialSpecular = &MaterialColor{Animate: proto.Bool(m.MaterialSpecular.Animate), Color: lvms, ColorPickerOpen: proto.Bool(m.MaterialSpecular.ColorPickerOpen), Strength: proto.Float32(m.MaterialSpecular.Strength)}
		mm.MaterialEmission = &MaterialColor{Animate: proto.Bool(m.MaterialEmission.Animate), Color: lvme, ColorPickerOpen: proto.Bool(m.MaterialEmission.ColorPickerOpen), Strength: proto.Float32(m.MaterialEmission.Strength)}

		mm.Setting_ParallaxMapping = proto.Bool(m.ParallaxMapping)

		mm.Effect_GBlur_Mode = proto.Int32(m.EffectGBlurMode)
		mm.Effect_GBlur_Radius = &ObjectCoordinate{Animate: proto.Bool(m.EffectGBlurRadius.Animate), Point: proto.Float32(m.EffectGBlurRadius.Point)}
		mm.Effect_GBlur_Width = &ObjectCoordinate{Animate: proto.Bool(m.EffectGBlurWidth.Animate), Point: proto.Float32(m.EffectGBlurWidth.Point)}

		mm.Effect_Bloom_DoBloom = proto.Bool(m.EffectBloomDoBloom)
		mm.Effect_Bloom_WeightA = proto.Float32(m.EffectBloomWeightA)
		mm.Effect_Bloom_WeightB = proto.Float32(m.EffectBloomWeightB)
		mm.Effect_Bloom_WeightC = proto.Float32(m.EffectBloomWeightC)
		mm.Effect_Bloom_WeightD = proto.Float32(m.EffectBloomWeightD)
		mm.Effect_Bloom_Vignette = proto.Float32(m.EffectBloomVignette)
		mm.Effect_Bloom_VignetteAtt = proto.Float32(m.EffectBloomVignetteAtt)

		mm.Setting_LightingPass_DrawMode = proto.Int32(int32(m.LightingPassDrawMode))

		mm.EffectToneMappingACESFilmRec2020 = proto.Bool(m.EffectToneMappingACESFilmRec2020)
		mm.EffectHDRTonemapping = proto.Bool(m.EffectHDRTonemapping)

		mm.ShowShadows = proto.Bool(m.ShowShadows)

		mm.RenderingPBR = proto.Bool(m.RenderingPBR)
		mm.RenderingPBRMetallic = proto.Float32(m.RenderingPBRMetallic)
		mm.RenderingPBRRoughness = proto.Float32(m.RenderingPBRRoughness)
		mm.RenderingPBRAO = proto.Float32(m.RenderingPBRAO)

		mm.SolidLightSkinMaterialColor = &Vec3{X: proto.Float32(m.SolidLightSkinMaterialColor.X()), Y: proto.Float32(m.SolidLightSkinMaterialColor.Y()), Z: proto.Float32(m.SolidLightSkinMaterialColor.Z())}
		mm.SolidLightSkinAmbient = &Vec3{X: proto.Float32(m.SolidLightSkinAmbient.X()), Y: proto.Float32(m.SolidLightSkinAmbient.Y()), Z: proto.Float32(m.SolidLightSkinAmbient.Z())}
		mm.SolidLightSkinDiffuse = &Vec3{X: proto.Float32(m.SolidLightSkinDiffuse.X()), Y: proto.Float32(m.SolidLightSkinDiffuse.Y()), Z: proto.Float32(m.SolidLightSkinDiffuse.Z())}
		mm.SolidLightSkinSpecular = &Vec3{X: proto.Float32(m.SolidLightSkinSpecular.X()), Y: proto.Float32(m.SolidLightSkinSpecular.Y()), Z: proto.Float32(m.SolidLightSkinSpecular.Z())}

		mo := &Mesh{}
		mo.ID = proto.Int32(int32(m.MeshModel.ID))
		mo.ModelTitle = proto.String(m.MeshModel.ModelTitle)
		mo.MaterialTitle = proto.String(m.MeshModel.MaterialTitle)
		mo.CountVertices = proto.Int32(m.MeshModel.CountVertices)
		mo.CountTextureCoordinates = proto.Int32(m.MeshModel.CountTextureCoordinates)
		mo.CountNormals = proto.Int32(m.MeshModel.CountNormals)
		mo.CountIndices = proto.Int32(m.MeshModel.CountIndices)

		for i := 0; i < len(m.MeshModel.Vertices); i++ {
			mo.Vertices = append(mo.Vertices, &Vec3{X: proto.Float32(m.MeshModel.Vertices[i].X()), Y: proto.Float32(m.MeshModel.Vertices[i].Y()), Z: proto.Float32(m.MeshModel.Vertices[i].Z())})
		}
		for i := 0; i < len(m.MeshModel.TextureCoordinates); i++ {
			mo.TextureCoordinates = append(mo.TextureCoordinates, &Vec2{X: proto.Float32(m.MeshModel.TextureCoordinates[i].X()), Y: proto.Float32(m.MeshModel.TextureCoordinates[i].Y())})
		}
		for i := 0; i < len(m.MeshModel.Normals); i++ {
			mo.Normals = append(mo.Normals, &Vec3{X: proto.Float32(m.MeshModel.Normals[i].X()), Y: proto.Float32(m.MeshModel.Normals[i].Y()), Z: proto.Float32(m.MeshModel.Normals[i].Z())})
		}
		for i := 0; i < len(m.MeshModel.Indices); i++ {
			mo.Indices = append(mo.Indices, m.MeshModel.Indices[i])
		}

		mo.File = proto.String(m.MeshModel.File)
		mo.FilePath = proto.String(m.MeshModel.FilePath)

		mom := &MeshModelMaterial{}
		mom.MaterialID = proto.Int32(int32(m.MeshModel.ModelMaterial.MaterialID))
		mom.MaterialTitle = proto.String(m.MeshModel.ModelMaterial.MaterialTitle)

		mom.AmbientColor = &Vec3{X: proto.Float32(m.MeshModel.ModelMaterial.AmbientColor.X()), Y: proto.Float32(m.MeshModel.ModelMaterial.AmbientColor.Y()), Z: proto.Float32(m.MeshModel.ModelMaterial.AmbientColor.Z())}
		mom.DiffuseColor = &Vec3{X: proto.Float32(m.MeshModel.ModelMaterial.DiffuseColor.X()), Y: proto.Float32(m.MeshModel.ModelMaterial.DiffuseColor.Y()), Z: proto.Float32(m.MeshModel.ModelMaterial.DiffuseColor.Z())}
		mom.SpecularColor = &Vec3{X: proto.Float32(m.MeshModel.ModelMaterial.SpecularColor.X()), Y: proto.Float32(m.MeshModel.ModelMaterial.SpecularColor.Y()), Z: proto.Float32(m.MeshModel.ModelMaterial.SpecularColor.Z())}
		mom.EmissionColor = &Vec3{X: proto.Float32(m.MeshModel.ModelMaterial.EmissionColor.X()), Y: proto.Float32(m.MeshModel.ModelMaterial.EmissionColor.Y()), Z: proto.Float32(m.MeshModel.ModelMaterial.EmissionColor.Z())}
		mom.SpecularExp = proto.Float32(m.MeshModel.ModelMaterial.SpecularExp)
		mom.Transparency = proto.Float32(m.MeshModel.ModelMaterial.Transparency)
		mom.IlluminationMode = proto.Uint32(m.MeshModel.ModelMaterial.IlluminationMode)
		mom.OpticalDensity = proto.Float32(m.MeshModel.ModelMaterial.OpticalDensity)

		mmmtia := &MeshMaterialTextureImage{
			Filename:   proto.String(m.MeshModel.ModelMaterial.TextureAmbient.Filename),
			Image:      proto.String(m.MeshModel.ModelMaterial.TextureAmbient.Image),
			Width:      proto.Int32(m.MeshModel.ModelMaterial.TextureAmbient.Width),
			Height:     proto.Int32(m.MeshModel.ModelMaterial.TextureAmbient.Height),
			UseTexture: proto.Bool(m.MeshModel.ModelMaterial.TextureAmbient.UseTexture),
			Commands:   m.MeshModel.ModelMaterial.TextureAmbient.Commands}
		mom.TextureAmbient = mmmtia
		mmmtid := &MeshMaterialTextureImage{
			Filename:   proto.String(m.MeshModel.ModelMaterial.TextureDiffuse.Filename),
			Image:      proto.String(m.MeshModel.ModelMaterial.TextureDiffuse.Image),
			Width:      proto.Int32(m.MeshModel.ModelMaterial.TextureDiffuse.Width),
			Height:     proto.Int32(m.MeshModel.ModelMaterial.TextureDiffuse.Height),
			UseTexture: proto.Bool(m.MeshModel.ModelMaterial.TextureDiffuse.UseTexture),
			Commands:   m.MeshModel.ModelMaterial.TextureDiffuse.Commands}
		mom.TextureDiffuse = mmmtid
		mmmtis := &MeshMaterialTextureImage{
			Filename:   proto.String(m.MeshModel.ModelMaterial.TextureSpecular.Filename),
			Image:      proto.String(m.MeshModel.ModelMaterial.TextureSpecular.Image),
			Width:      proto.Int32(m.MeshModel.ModelMaterial.TextureSpecular.Width),
			Height:     proto.Int32(m.MeshModel.ModelMaterial.TextureSpecular.Height),
			UseTexture: proto.Bool(m.MeshModel.ModelMaterial.TextureSpecular.UseTexture),
			Commands:   m.MeshModel.ModelMaterial.TextureSpecular.Commands}
		mom.TextureSpecular = mmmtis
		mmmtise := &MeshMaterialTextureImage{
			Filename:   proto.String(m.MeshModel.ModelMaterial.TextureSpecularExp.Filename),
			Image:      proto.String(m.MeshModel.ModelMaterial.TextureSpecularExp.Image),
			Width:      proto.Int32(m.MeshModel.ModelMaterial.TextureSpecularExp.Width),
			Height:     proto.Int32(m.MeshModel.ModelMaterial.TextureSpecularExp.Height),
			UseTexture: proto.Bool(m.MeshModel.ModelMaterial.TextureSpecularExp.UseTexture),
			Commands:   m.MeshModel.ModelMaterial.TextureSpecularExp.Commands}
		mom.TextureSpecularExp = mmmtise
		mmmtids := &MeshMaterialTextureImage{
			Filename:   proto.String(m.MeshModel.ModelMaterial.TextureDissolve.Filename),
			Image:      proto.String(m.MeshModel.ModelMaterial.TextureDissolve.Image),
			Width:      proto.Int32(m.MeshModel.ModelMaterial.TextureDissolve.Width),
			Height:     proto.Int32(m.MeshModel.ModelMaterial.TextureDissolve.Height),
			UseTexture: proto.Bool(m.MeshModel.ModelMaterial.TextureDissolve.UseTexture),
			Commands:   m.MeshModel.ModelMaterial.TextureDissolve.Commands}
		mom.TextureDissolve = mmmtids
		mmmtib := &MeshMaterialTextureImage{
			Filename:   proto.String(m.MeshModel.ModelMaterial.TextureBump.Filename),
			Image:      proto.String(m.MeshModel.ModelMaterial.TextureBump.Image),
			Width:      proto.Int32(m.MeshModel.ModelMaterial.TextureBump.Width),
			Height:     proto.Int32(m.MeshModel.ModelMaterial.TextureBump.Height),
			UseTexture: proto.Bool(m.MeshModel.ModelMaterial.TextureBump.UseTexture),
			Commands:   m.MeshModel.ModelMaterial.TextureBump.Commands}
		mom.TextureBump = mmmtib
		mmmtidi := &MeshMaterialTextureImage{
			Filename:   proto.String(m.MeshModel.ModelMaterial.TextureDisplacement.Filename),
			Image:      proto.String(m.MeshModel.ModelMaterial.TextureDisplacement.Image),
			Width:      proto.Int32(m.MeshModel.ModelMaterial.TextureDisplacement.Width),
			Height:     proto.Int32(m.MeshModel.ModelMaterial.TextureDisplacement.Height),
			UseTexture: proto.Bool(m.MeshModel.ModelMaterial.TextureDisplacement.UseTexture),
			Commands:   m.MeshModel.ModelMaterial.TextureDisplacement.Commands}
		mom.TextureDisplacement = mmmtidi
		mo.ModelMaterial = mom

		mm.MeshObject = mo

		gs.Models = append(gs.Models, mm)
	}

	// Save
	data, err := proto.Marshal(gs)
	if err != nil {
		settings.LogWarn("[SaveOpen-ProtoBufs] [storeObjects] Marshalling error: %v", err)
		return
	}
	err = ioutil.WriteFile(pm.fileNameScene, data, 0644)
	if err != nil {
		settings.LogWarn("[SaveOpen-ProtoBufs] [storeObjects] Can't save byte data error: %v", err)
		return
	}
}
