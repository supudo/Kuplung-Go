package saveopen

import (
	"io/ioutil"
	"strings"

	proto "github.com/golang/protobuf/proto"
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
	utilities.ZipFiles(fileName, zfiles, true)
}

// Open ...
func (pm *ProtoBufsSaveOpen) Open(file *types.FBEntity) []*meshes.ModelFace {
	meshes := []*meshes.ModelFace{}
	return meshes
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
		settings.LogError("[SaveOpen-ProtoBufs] [storeRenderingSettings] Marshalling error: %v", err)
	}
	err = ioutil.WriteFile(pm.fileNameSettings, data, 0644)
	if err != nil {
		settings.LogError("[SaveOpen-ProtoBufs] [storeRenderingSettings] Can't save byte data error: %v", err)
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
		settings.LogError("[SaveOpen-ProtoBufs] [storeObjects] Marshalling error: %v", err)
	}
	err = ioutil.WriteFile(pm.fileNameScene, data, 0644)
	if err != nil {
		settings.LogError("[SaveOpen-ProtoBufs] [storeObjects] Can't save byte data error: %v", err)
	}
}
