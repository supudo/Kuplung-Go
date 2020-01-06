package renderers

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/objects"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
	"math"
)

// RendererForward ...
type RendererForward struct {
	window interfaces.Window

	shaderProgram uint32
	glVAO         uint32

	matrixProjection, matrixCamera    mgl32.Mat4
	vecCameraPosition, uiAmbientLight mgl32.Vec3
	lightingPass_DrawMode             int32

	solidLight *types.ModelFaceLightSourceDirectional

	glFS_solidSkin_materialColor int32

	GLSL_LightSourceNumber_Directional, GLSL_LightSourceNumber_Point, GLSL_LightSourceNumber_Spot uint32

	mfLights_Directional []*types.ModelFaceLightSourceDirectional
	mfLights_Point       []*types.ModelFaceLightSourcePoint
	mfLights_Spot        []*types.ModelFaceLightSourceSpot

	// variables
	glVS_MVPMatrix, glFS_MMatrix, glVS_WorldMatrix, glVS_NormalMatrix, glFS_MVMatrix int32

	// general
	glGS_GeomDisplacementLocation, glFS_AlphaBlending, glFS_CameraPosition, glFS_CelShading   int32
	glFS_OutlineColor, glVS_IsBorder, glFS_ScreenResX, glFS_ScreenResY, glFS_UIAmbient        int32
	glTCS_UseCullFace, glTCS_UseTessellation, glTCS_TessellationSubdivision, gl_ModelViewSkin int32
	glFS_GammaCoeficient, glFS_showShadows, glFS_ShadowPass                                   int32

	// depth color
	glFS_planeClose, glFS_planeFar, glFS_showDepthColor int32

	// material
	glMaterial_Ambient, glMaterial_Diffuse, glMaterial_Specular, glMaterial_SpecularExp                                           int32
	glMaterial_Emission, glMaterial_Refraction, glMaterial_IlluminationModel, glMaterial_HeightScale                              int32
	glMaterial_SamplerAmbient, glMaterial_SamplerDiffuse, glMaterial_SamplerSpecular                                              int32
	glMaterial_SamplerSpecularExp, glMaterial_SamplerDissolve, glMaterial_SamplerBump, glMaterial_SamplerDisplacement             int32
	glMaterial_HasTextureAmbient, glMaterial_HasTextureDiffuse, glMaterial_HasTextureSpecular                                     int32
	glMaterial_HasTextureSpecularExp, glMaterial_HasTextureDissolve, glMaterial_HasTextureBump, glMaterial_HasTextureDisplacement int32
	glMaterial_ParallaxMapping                                                                                                    int32

	// effects - gaussian blur
	glEffect_GB_W, glEffect_GB_Radius, glEffect_GB_Mode int32

	// effects - bloom
	glEffect_Bloom_doBloom, glEffect_Bloom_WeightA, glEffect_Bloom_WeightB, glEffect_Bloom_WeightC int32
	glEffect_Bloom_WeightD, glEffect_Bloom_Vignette, glEffect_Bloom_VignetteAtt                    int32

	// effects - tone mapping
	glEffect_ToneMapping_ACESFilmRec2020, glEffect_HDR_Tonemapping int32

	// PBR
	glPBR_UsePBR, glPBR_Metallic, glPBR_Rougness, glPBR_AO int32
}

// NewRendererForward ...
func NewRendererForward(window interfaces.Window) *RendererForward {
	sett := settings.GetSettings()
	gl := window.OpenGL()

	rend := &RendererForward{}
	rend.window = window

	rend.GLSL_LightSourceNumber_Directional = 8
	rend.GLSL_LightSourceNumber_Point = 4
	rend.GLSL_LightSourceNumber_Spot = 4

	sVertex := engine.GetShaderSource(sett.App.CurrentPath + "shaders/model_face.vert")
	sTcs := engine.GetShaderSource(sett.App.CurrentPath + "shaders/model_face.tcs")
	sTes := engine.GetShaderSource(sett.App.CurrentPath + "shaders/model_face.tes")
	sGeom := engine.GetShaderSource(sett.App.CurrentPath + "shaders/model_face.geom")
	sFragment := engine.GetShaderSourcePartial(sett.App.CurrentPath + "shaders/model_face_vars.frag")
	sFragment += engine.GetShaderSourcePartial(sett.App.CurrentPath + "shaders/model_face_effects.frag")
	sFragment += engine.GetShaderSourcePartial(sett.App.CurrentPath + "shaders/model_face_lights.frag")
	sFragment += engine.GetShaderSourcePartial(sett.App.CurrentPath + "shaders/model_face_mapping.frag")
	sFragment += engine.GetShaderSourcePartial(sett.App.CurrentPath + "shaders/model_face_shadow_mapping.frag")
	sFragment += engine.GetShaderSourcePartial(sett.App.CurrentPath + "shaders/model_face_misc.frag")
	sFragment += engine.GetShaderSourcePartial(sett.App.CurrentPath + "shaders/model_face_pbr.frag")
	sFragment += engine.GetShaderSource(sett.App.CurrentPath + "shaders/model_face.frag")

	var err error
	rend.shaderProgram, err = engine.LinkMultiProgram(gl, sVertex, sTcs, sTes, sGeom, sFragment)
	if err != nil {
		settings.LogWarn("[RendererForward] Can't load the forward renderer shaders: %v", err)
	}

	gl.PatchParameteri(oglconsts.PATCH_VERTICES, 3)

	rend.glFS_showShadows = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_showShadows\x00"))
	rend.glFS_ShadowPass = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_shadowPass\x00"))

	rend.glFS_planeClose = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_planeClose\x00"))
	rend.glFS_planeFar = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_planeFar\x00"))
	rend.glFS_showDepthColor = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_showDepthColor\x00"))

	rend.glGS_GeomDisplacementLocation = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("vs_displacementLocation\x00"))
	rend.glTCS_UseCullFace = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("tcs_UseCullFace\x00"))
	rend.glTCS_UseTessellation = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("tcs_UseTessellation\x00"))
	rend.glTCS_TessellationSubdivision = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("tcs_TessellationSubdivision\x00"))

	rend.glFS_AlphaBlending = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_alpha\x00"))
	rend.glFS_CelShading = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_celShading\x00"))
	rend.glFS_CameraPosition = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_cameraPosition\x00"))
	rend.glVS_IsBorder = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("vs_isBorder\x00"))
	rend.glFS_OutlineColor = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_outlineColor\x00"))
	rend.glFS_UIAmbient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_UIAmbient\x00"))
	rend.glFS_GammaCoeficient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_gammaCoeficient\x00"))

	rend.glVS_MVPMatrix = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("vs_MVPMatrix\x00"))
	rend.glFS_MMatrix = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_ModelMatrix\x00"))
	rend.glVS_WorldMatrix = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("vs_WorldMatrix\x00"))
	rend.glFS_MVMatrix = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("vs_MVMatrix\x00"))
	rend.glVS_NormalMatrix = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("vs_normalMatrix\x00"))

	rend.glFS_ScreenResX = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_screenResX\x00"))
	rend.glFS_ScreenResY = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_screenResY\x00"))

	rend.glMaterial_ParallaxMapping = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_userParallaxMapping\x00"))

	rend.gl_ModelViewSkin = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_modelViewSkin\x00"))
	rend.glFS_solidSkin_materialColor = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_materialColor\x00"))
	rend.solidLight = &types.ModelFaceLightSourceDirectional{}
	rend.solidLight.InUse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.inUse\x00"))
	rend.solidLight.Direction = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.direction\x00"))
	rend.solidLight.Ambient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.ambient\x00"))
	rend.solidLight.Diffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.diffuse\x00"))
	rend.solidLight.Specular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.specular\x00"))
	rend.solidLight.StrengthAmbient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.strengthAmbient\x00"))
	rend.solidLight.StrengthDiffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.strengthDiffuse\x00"))
	rend.solidLight.StrengthSpecular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.strengthSpecular\x00"))

	// light - directional
	var i uint32
	for i = 0; i < rend.GLSL_LightSourceNumber_Directional; i++ {
		f := &types.ModelFaceLightSourceDirectional{}

		f.InUse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("directionalLights["+string(i)+"].inUse\x00"))

		f.Direction = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("directionalLights["+string(i)+"].direction\x00"))

		f.Ambient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("directionalLights["+string(i)+"].ambient\x00"))
		f.Diffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("directionalLights["+string(i)+"].diffuse\x00"))
		f.Specular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("directionalLights["+string(i)+"].specular\x00"))

		f.StrengthAmbient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("directionalLights["+string(i)+"].strengthAmbient\x00"))
		f.StrengthDiffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("directionalLights["+string(i)+"].strengthDiffuse\x00"))
		f.StrengthSpecular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("directionalLights["+string(i)+"].strengthSpecular\x00"))
		rend.mfLights_Directional = append(rend.mfLights_Directional, f)
	}

	// light - point
	for i = 0; i < rend.GLSL_LightSourceNumber_Point; i++ {
		f := &types.ModelFaceLightSourcePoint{}

		f.InUse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("pointLights["+string(i)+"].inUse\x00"))
		f.Position = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("pointLights["+string(i)+"].position\x00"))

		f.Constant = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("pointLights["+string(i)+"].constant\x00"))
		f.Linear = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("pointLights["+string(i)+"].linear\x00"))
		f.Quadratic = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("pointLights["+string(i)+"].quadratic\x00"))

		f.Ambient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("pointLights["+string(i)+"].ambient\x00"))
		f.Diffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("pointLights["+string(i)+"].diffuse\x00"))
		f.Specular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("pointLights["+string(i)+"].specular\x00"))

		f.StrengthAmbient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("pointLights["+string(i)+"].strengthAmbient\x00"))
		f.StrengthDiffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("pointLights["+string(i)+"].strengthDiffuse\x00"))
		f.StrengthSpecular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("pointLights["+string(i)+"].strengthSpecular\x00"))
		rend.mfLights_Point = append(rend.mfLights_Point, f)
	}

	// light - spot
	for i = 0; i < rend.GLSL_LightSourceNumber_Spot; i++ {
		f := &types.ModelFaceLightSourceSpot{}

		f.InUse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].inUse\x00"))

		f.Position = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].position\x00"))
		f.Direction = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].direction\x00"))

		f.CutOff = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].cutOff\x00"))
		f.OuterCutOff = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].outerCutOff\x00"))

		f.Constant = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].constant\x00"))
		f.Linear = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].linear\x00"))
		f.Quadratic = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].quadratic\x00"))

		f.Ambient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].ambient\x00"))
		f.Diffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].diffuse\x00"))
		f.Specular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].specular\x00"))

		f.StrengthAmbient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].strengthAmbient\x00"))
		f.StrengthDiffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].strengthDiffuse\x00"))
		f.StrengthSpecular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("spotLights["+string(i)+"].strengthSpecular\x00"))
		rend.mfLights_Spot = append(rend.mfLights_Spot, f)
	}

	// material
	rend.glMaterial_Refraction = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.refraction\x00"))
	rend.glMaterial_SpecularExp = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.specularExp\x00"))
	rend.glMaterial_IlluminationModel = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.illumination_model\x00"))
	rend.glMaterial_HeightScale = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.heightScale\x00"))

	rend.glMaterial_Ambient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.ambient\x00"))
	rend.glMaterial_Diffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.diffuse\x00"))
	rend.glMaterial_Specular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.specular\x00"))
	rend.glMaterial_Emission = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.emission\x00"))

	rend.glMaterial_SamplerAmbient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.sampler_ambient\x00"))
	rend.glMaterial_SamplerDiffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.sampler_diffuse\x00"))
	rend.glMaterial_SamplerSpecular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.sampler_specular\x00"))
	rend.glMaterial_SamplerSpecularExp = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.sampler_specularExp\x00"))
	rend.glMaterial_SamplerDissolve = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.sampler_dissolve\x00"))
	rend.glMaterial_SamplerBump = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.sampler_bump\x00"))
	rend.glMaterial_SamplerDisplacement = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.sampler_displacement\x00"))

	rend.glMaterial_HasTextureAmbient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.has_texture_ambient\x00"))
	rend.glMaterial_HasTextureDiffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.has_texture_diffuse\x00"))
	rend.glMaterial_HasTextureSpecular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.has_texture_specular\x00"))
	rend.glMaterial_HasTextureSpecularExp = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.has_texture_specularExp\x00"))
	rend.glMaterial_HasTextureDissolve = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.has_texture_dissolve\x00"))
	rend.glMaterial_HasTextureBump = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.has_texture_bump\x00"))
	rend.glMaterial_HasTextureDisplacement = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("material.has_texture_displacement\x00"))

	// effects - gaussian blur
	rend.glEffect_GB_W = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("effect_GBlur.gauss_w\x00"))
	rend.glEffect_GB_Radius = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("effect_GBlur.gauss_radius\x00"))
	rend.glEffect_GB_Mode = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("effect_GBlur.gauss_mode\x00"))

	// effects - bloom
	rend.glEffect_Bloom_doBloom = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("effect_Bloom.doBloom\x00"))
	rend.glEffect_Bloom_WeightA = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("effect_Bloom.bloom_WeightA\x00"))
	rend.glEffect_Bloom_WeightB = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("effect_Bloom.bloom_WeightB\x00"))
	rend.glEffect_Bloom_WeightC = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("effect_Bloom.bloom_WeightC\x00"))
	rend.glEffect_Bloom_WeightD = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("effect_Bloom.bloom_WeightD\x00"))
	rend.glEffect_Bloom_Vignette = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("effect_Bloom.bloom_Vignette\x00"))
	rend.glEffect_Bloom_VignetteAtt = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("effect_Bloom.bloom_VignetteAtt\x00"))

	// effects - tone mapping
	rend.glEffect_ToneMapping_ACESFilmRec2020 = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_ACESFilmRec2020\x00"))
	rend.glEffect_HDR_Tonemapping = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_HDRTonemapping\x00"))

	// PBR
	rend.glPBR_UsePBR = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_renderPBR\x00"))
	rend.glPBR_Metallic = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_PBR_Metallic\x00"))
	rend.glPBR_Rougness = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_PBR_Roughness\x00"))
	rend.glPBR_AO = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_PBR_AO\x00"))

	gl.CheckForOpenGLErrors("ForwardRenderer")

	return rend
}

// Render ...
func (rend *RendererForward) Render(rp types.RenderProperties, meshModelFaces []*meshes.ModelFace, matrixGrid mgl32.Mat4, camPos mgl32.Vec3, selectedModel int32, lightSources []*objects.Light) {
	gl := rend.window.OpenGL()
	rsett := settings.GetRenderingSettings()

	gl.UseProgram(rend.shaderProgram)

	querycount := int32(5)
	queries := make([]uint32, querycount)
	currentQuery := int32(0)
	queries = gl.GenQueries(querycount)

	gl.BeginQuery(oglconsts.TIME_ELAPSED, queries[currentQuery])

	//selectedModelID := int32(-1)
	i := int32(-1)
	for i = 0; i < int32(len(meshModelFaces)); i++ {
		mfd := meshModelFaces[i]

		j := int32(0)
		if rsett.General.OcclusionCulling {
			gl.Disable(oglconsts.CULL_FACE)
			gl.DepthMask(false)
			gl.ColorMask(false, false, false, false)
			for ; j < int32(len(meshModelFaces)); j++ {
				gl.BeginQuery(oglconsts.ANY_SAMPLES_PASSED, mfd.OccQuery)
				gl.BindVertexArray(mfd.GLVAO)
				gl.DrawElements(oglconsts.TRIANGLES, 6*6, oglconsts.UNSIGNED_INT, 0)
				gl.EndQuery(oglconsts.ANY_SAMPLES_PASSED)
			}
			j = i
		}
		gl.Enable(oglconsts.CULL_FACE)

		//selectedModelID = i

		matrixModel := mgl32.Ident4()
		matrixModel = matrixModel.Mul4(matrixGrid)
		// scale
		matrixModel = matrixModel.Mul4(mgl32.Scale3D(mfd.ScaleX.Point, mfd.ScaleY.Point, mfd.ScaleZ.Point))
		// translate
		matrixModel = matrixModel.Mul4(mgl32.Translate3D(mfd.PositionX.Point, mfd.PositionY.Point, mfd.PositionZ.Point))
		// rotate
		matrixModel = matrixModel.Mul4(mgl32.Translate3D(0, 0, 0))
		matrixModel = matrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(mfd.RotateX.Point), mgl32.Vec3{1, 0, 0}))
		matrixModel = matrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(mfd.RotateY.Point), mgl32.Vec3{0, 1, 0}))
		matrixModel = matrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(mfd.RotateZ.Point), mgl32.Vec3{0, 0, 1}))
		matrixModel = matrixModel.Mul4(mgl32.Translate3D(0, 0, 0))

		mfd.MatrixModel = matrixModel

		mfd.ModelViewSkin = rsett.General.SelectedViewModelSkin
		mfd.OutlineColor = rsett.General.OutlineColor
		mfd.OutlineThickness = rsett.General.OutlineThickness
		mfd.IsModelSelected = i == selectedModel

		mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(matrixModel))
		gl.GLUniformMatrix4fv(rend.glVS_MVPMatrix, 1, false, &mvpMatrix[0])
		gl.GLUniformMatrix4fv(rend.glFS_MMatrix, 1, false, &matrixModel[0])

		matrixModelView := rsett.MatrixCamera.Mul4(matrixModel)
		gl.GLUniformMatrix4fv(rend.glFS_MVMatrix, 1, false, &matrixModelView[0])

		matrixNormal := matrixModelView.Inv().Transpose()
		gl.UniformMatrix3fv(rend.glVS_NormalMatrix, 1, false, &matrixNormal[0])

		gl.GLUniformMatrix4fv(rend.glVS_WorldMatrix, 1, false, &matrixModel[0])

		// blending
		if mfd.MeshModel.ModelMaterial.Transparency < 1.0 || mfd.SettingAlpha < 1.0 {
			gl.Disable(oglconsts.DEPTH_TEST)
			gl.BlendFunc(oglconsts.SRC_ALPHA, oglconsts.ONE_MINUS_SRC_ALPHA)
			gl.Enable(oglconsts.BLEND)
			if mfd.MeshModel.ModelMaterial.Transparency < 1.0 {
				gl.Uniform1f(rend.glFS_AlphaBlending, mfd.MeshModel.ModelMaterial.Transparency)
			} else {
				gl.Uniform1f(rend.glFS_AlphaBlending, mfd.SettingAlpha)
			}
		} else {
			gl.Enable(oglconsts.DEPTH_TEST)
			gl.DepthFunc(oglconsts.LESS)
			gl.Disable(oglconsts.BLEND)
			gl.BlendFunc(oglconsts.SRC_ALPHA, oglconsts.ONE_MINUS_SRC_ALPHA)
			gl.Uniform1f(rend.glFS_AlphaBlending, 1.0)
		}

		// depth color
		pc := float32(1.0)
		if rsett.General.PlaneClose >= 1.0 {
			pc = rsett.General.PlaneClose
		}
		gl.Uniform1f(rend.glFS_planeClose, pc)
		gl.Uniform1f(rend.glFS_planeFar, rsett.General.PlaneFar/100.0)
		if rsett.General.RenderingDepth {
			gl.Uniform1i(rend.glFS_showDepthColor, 1)
		} else {
			gl.Uniform1i(rend.glFS_showDepthColor, 0)
		}
		gl.Uniform1i(rend.glFS_ShadowPass, 0)

		// tessellation
		if mfd.SettingUseCullFace {
			gl.Uniform1i(rend.glTCS_UseCullFace, 1)
		} else {
			gl.Uniform1i(rend.glTCS_UseCullFace, 0)
		}
		if mfd.SettingUseTessellation {
			gl.Uniform1i(rend.glTCS_UseTessellation, 1)
		} else {
			gl.Uniform1i(rend.glTCS_UseTessellation, 0)
		}
		gl.Uniform1i(rend.glTCS_TessellationSubdivision, int32(mfd.SettingTessellationSubdivision))

		// cel-shading
		if mfd.SettingCelShading {
			gl.Uniform1i(rend.glFS_CelShading, 1)
		} else {
			gl.Uniform1i(rend.glFS_CelShading, 0)
		}

		// camera position
		gl.Uniform3f(rend.glFS_CameraPosition, camPos.X(), camPos.Y(), camPos.Z())

		// screen size
		w, h := rend.window.Size()
		gl.Uniform1f(rend.glFS_ScreenResX, float32(w))
		gl.Uniform1f(rend.glFS_ScreenResY, float32(h))

		// Outline color
		gl.Uniform3f(rend.glFS_OutlineColor, mfd.OutlineColor.X(), mfd.OutlineColor.Y(), mfd.OutlineColor.Z())

		// ambient color for editor
		gl.Uniform3f(rend.glFS_UIAmbient, rend.uiAmbientLight.X(), rend.uiAmbientLight.Y(), rend.uiAmbientLight.Z())

		// geometry shader displacement
		gl.Uniform3f(rend.glGS_GeomDisplacementLocation, mfd.DisplaceX.Point, mfd.DisplaceY.Point, mfd.DisplaceZ.Point)

		// mapping
		if mfd.SettingParallaxMapping {
			gl.Uniform1i(rend.glMaterial_ParallaxMapping, 1)
		} else {
			gl.Uniform1i(rend.glMaterial_ParallaxMapping, 0)
		}

		// gamma correction
		gl.Uniform1f(rend.glFS_GammaCoeficient, rsett.General.GammaCoeficient)

		// render skin
		gl.Uniform1i(rend.gl_ModelViewSkin, int32(mfd.SettingModelViewSkin))
		gl.Uniform3f(rend.glFS_solidSkin_materialColor, mfd.SolidLightSkinMaterialColor.X(), mfd.SolidLightSkinMaterialColor.Y(), mfd.SolidLightSkinMaterialColor.Z())

		// shadows
		gl.Uniform1i(rend.glFS_showShadows, 0)

		gl.Uniform1i(rend.solidLight.InUse, 1)
		gl.Uniform3f(rend.solidLight.Direction, rp.SolidLightDirectionX, rp.SolidLightDirectionY, rp.SolidLightDirectionZ)
		gl.Uniform3f(rend.solidLight.Ambient, rp.SolidLightAmbient.X(), rp.SolidLightAmbient.Y(), rp.SolidLightAmbient.Z())
		gl.Uniform3f(rend.solidLight.Diffuse, rp.SolidLightDiffuse.X(), rp.SolidLightDiffuse.Y(), rp.SolidLightDiffuse.Z())
		gl.Uniform3f(rend.solidLight.Specular, rp.SolidLightSpecular.X(), rp.SolidLightSpecular.Y(), rp.SolidLightSpecular.Z())
		gl.Uniform1f(rend.solidLight.StrengthAmbient, rp.SolidLightAmbientStrength)
		gl.Uniform1f(rend.solidLight.StrengthDiffuse, rp.SolidLightDiffuseStrength)
		gl.Uniform1f(rend.solidLight.StrengthSpecular, rp.SolidLightSpecularStrength)

		// lights
		lightsCountDirectional := uint32(0)
		lightsCountPoint := uint32(0)
		lightsCountSpot := uint32(0)
		for j := 0; j < len(lightSources); j++ {
			light := lightSources[j]
			switch light.LightType {
			case types.LightSourceTypeDirectional:
				if lightsCountDirectional < rend.GLSL_LightSourceNumber_Directional {
					f := rend.mfLights_Directional[lightsCountDirectional]

					gl.Uniform1i(f.InUse, 1)

					// light
					gl.Uniform3f(f.Direction, light.PositionX.Point, light.PositionY.Point, light.PositionZ.Point)

					// color
					gl.Uniform3f(f.Ambient, light.Ambient.Color.X(), light.Ambient.Color.Y(), light.Ambient.Color.Z())
					gl.Uniform3f(f.Diffuse, light.Diffuse.Color.X(), light.Diffuse.Color.Y(), light.Diffuse.Color.Z())
					gl.Uniform3f(f.Specular, light.Specular.Color.X(), light.Specular.Color.Y(), light.Specular.Color.Z())

					// light factors
					gl.Uniform1f(f.StrengthAmbient, light.Ambient.Strength)
					gl.Uniform1f(f.StrengthDiffuse, light.Diffuse.Strength)
					gl.Uniform1f(f.StrengthSpecular, light.Specular.Strength)

					lightsCountDirectional++
				}
			case types.LightSourceTypePoint:
				if lightsCountPoint < rend.GLSL_LightSourceNumber_Point {
					f := rend.mfLights_Point[lightsCountPoint]

					gl.Uniform1i(f.InUse, 1)

					// light
					gl.Uniform3f(f.Position, light.MatrixModel[4*3+0], light.MatrixModel[4*3+1], light.MatrixModel[4*3+2])

					// factors
					gl.Uniform1f(f.Constant, light.LConstant.Point)
					gl.Uniform1f(f.Linear, light.LLinear.Point)
					gl.Uniform1f(f.Quadratic, light.LQuadratic.Point)

					// color
					gl.Uniform3f(f.Ambient, light.Ambient.Color.X(), light.Ambient.Color.Y(), light.Ambient.Color.Z())
					gl.Uniform3f(f.Diffuse, light.Diffuse.Color.X(), light.Diffuse.Color.Y(), light.Diffuse.Color.Z())
					gl.Uniform3f(f.Specular, light.Specular.Color.X(), light.Specular.Color.Y(), light.Specular.Color.Z())

					// light factors
					gl.Uniform1f(f.StrengthAmbient, light.Ambient.Strength)
					gl.Uniform1f(f.StrengthDiffuse, light.Diffuse.Strength)
					gl.Uniform1f(f.StrengthSpecular, light.Specular.Strength)

					lightsCountPoint++
				}
			case types.LightSourceTypeSpot:
				if lightsCountSpot < rend.GLSL_LightSourceNumber_Spot {
					f := rend.mfLights_Spot[lightsCountSpot]

					gl.Uniform1i(f.InUse, 1)

					// light
					gl.Uniform3f(f.Direction, light.PositionX.Point, light.PositionY.Point, light.PositionZ.Point)
					gl.Uniform3f(f.Position, light.MatrixModel[4*3+0], light.MatrixModel[4*3+1], light.MatrixModel[4*3+2])

					// cutoff
					gl.Uniform1f(f.CutOff, float32(math.Cos(float64(mgl32.DegToRad(light.LCutOff.Point)))))
					gl.Uniform1f(f.OuterCutOff, float32(math.Cos(float64(mgl32.DegToRad(light.LOuterCutOff.Point)))))

					// factors
					gl.Uniform1f(f.Constant, light.LConstant.Point)
					gl.Uniform1f(f.Linear, light.LLinear.Point)
					gl.Uniform1f(f.Quadratic, light.LQuadratic.Point)

					// color
					gl.Uniform3f(f.Ambient, light.Ambient.Color.X(), light.Ambient.Color.Y(), light.Ambient.Color.Z())
					gl.Uniform3f(f.Diffuse, light.Diffuse.Color.X(), light.Diffuse.Color.Y(), light.Diffuse.Color.Z())
					gl.Uniform3f(f.Specular, light.Specular.Color.X(), light.Specular.Color.Y(), light.Specular.Color.Z())

					// light factors
					gl.Uniform1f(f.StrengthAmbient, light.Ambient.Strength)
					gl.Uniform1f(f.StrengthDiffuse, light.Diffuse.Strength)
					gl.Uniform1f(f.StrengthSpecular, light.Specular.Strength)

					lightsCountSpot++
				}
			}
		}

		for j := lightsCountDirectional; j < rend.GLSL_LightSourceNumber_Directional; j++ {
			gl.Uniform1i(rend.mfLights_Directional[j].InUse, 0)
		}

		for j := lightsCountPoint; j < rend.GLSL_LightSourceNumber_Point; j++ {
			gl.Uniform1i(rend.mfLights_Point[j].InUse, 0)
		}

		for j := lightsCountSpot; j < rend.GLSL_LightSourceNumber_Spot; j++ {
			gl.Uniform1i(rend.mfLights_Spot[j].InUse, 0)
		}

		// material
		gl.Uniform1f(rend.glMaterial_Refraction, mfd.SettingMaterialRefraction.Point)
		gl.Uniform1f(rend.glMaterial_SpecularExp, mfd.SettingMaterialSpecularExp.Point)
		gl.Uniform1i(rend.glMaterial_IlluminationModel, int32(mfd.MaterialIlluminationModel))
		gl.Uniform1f(rend.glMaterial_HeightScale, mfd.DisplacementHeightScale.Point)
		gl.Uniform3f(rend.glMaterial_Ambient, mfd.MaterialAmbient.Color.X(), mfd.MaterialAmbient.Color.Y(), mfd.MaterialAmbient.Color.Z())
		gl.Uniform3f(rend.glMaterial_Diffuse, mfd.MaterialDiffuse.Color.X(), mfd.MaterialDiffuse.Color.Y(), mfd.MaterialDiffuse.Color.Z())
		gl.Uniform3f(rend.glMaterial_Specular, mfd.MaterialSpecular.Color.X(), mfd.MaterialSpecular.Color.Y(), mfd.MaterialSpecular.Color.Z())
		gl.Uniform3f(rend.glMaterial_Emission, mfd.MaterialEmission.Color.X(), mfd.MaterialEmission.Color.Y(), mfd.MaterialEmission.Color.Z())

		if mfd.HasTextureAmbient && mfd.MeshModel.ModelMaterial.TextureAmbient.UseTexture {
			gl.Uniform1i(rend.glMaterial_HasTextureAmbient, 1)
			gl.Uniform1i(rend.glMaterial_SamplerAmbient, 0)
			gl.ActiveTexture(oglconsts.TEXTURE0)
			gl.BindTexture(oglconsts.TEXTURE_2D, mfd.VboTextureAmbient)
		} else {
			gl.Uniform1i(rend.glMaterial_HasTextureAmbient, 0)
		}

		if mfd.HasTextureDiffuse && mfd.MeshModel.ModelMaterial.TextureDiffuse.UseTexture {
			gl.Uniform1i(rend.glMaterial_HasTextureDiffuse, 1)
			gl.Uniform1i(rend.glMaterial_SamplerDiffuse, 1)
			gl.ActiveTexture(oglconsts.TEXTURE1)
			gl.BindTexture(oglconsts.TEXTURE_2D, mfd.VboTextureDiffuse)
		} else {
			gl.Uniform1i(rend.glMaterial_HasTextureDiffuse, 0)
		}

		if mfd.HasTextureSpecular && mfd.MeshModel.ModelMaterial.TextureSpecular.UseTexture {
			gl.Uniform1i(rend.glMaterial_HasTextureSpecular, 1)
			gl.Uniform1i(rend.glMaterial_SamplerSpecular, 2)
			gl.ActiveTexture(oglconsts.TEXTURE2)
			gl.BindTexture(oglconsts.TEXTURE_2D, mfd.VboTextureSpecular)
		} else {
			gl.Uniform1i(rend.glMaterial_HasTextureSpecular, 0)
		}

		if mfd.HasTextureSpecularExp && mfd.MeshModel.ModelMaterial.TextureSpecularExp.UseTexture {
			gl.Uniform1i(rend.glMaterial_HasTextureSpecularExp, 1)
			gl.Uniform1i(rend.glMaterial_SamplerSpecularExp, 3)
			gl.ActiveTexture(oglconsts.TEXTURE3)
			gl.BindTexture(oglconsts.TEXTURE_2D, mfd.VboTextureSpecularExp)
		} else {
			gl.Uniform1i(rend.glMaterial_HasTextureSpecularExp, 0)
		}

		if mfd.HasTextureDissolve && mfd.MeshModel.ModelMaterial.TextureDissolve.UseTexture {
			gl.Uniform1i(rend.glMaterial_HasTextureDissolve, 1)
			gl.Uniform1i(rend.glMaterial_SamplerDissolve, 4)
			gl.ActiveTexture(oglconsts.TEXTURE4)
			gl.BindTexture(oglconsts.TEXTURE_2D, mfd.VboTextureDissolve)
		} else {
			gl.Uniform1i(rend.glMaterial_HasTextureDissolve, 0)
		}

		if mfd.HasTextureBump && mfd.MeshModel.ModelMaterial.TextureBump.UseTexture {
			gl.Uniform1i(rend.glMaterial_HasTextureBump, 1)
			gl.Uniform1i(rend.glMaterial_SamplerBump, 5)
			gl.ActiveTexture(oglconsts.TEXTURE5)
			gl.BindTexture(oglconsts.TEXTURE_2D, mfd.VboTextureBump)
		} else {
			gl.Uniform1i(rend.glMaterial_HasTextureBump, 0)
		}

		if mfd.HasTextureDisplacement && mfd.MeshModel.ModelMaterial.TextureDisplacement.UseTexture {
			gl.Uniform1i(rend.glMaterial_HasTextureDisplacement, 1)
			gl.Uniform1i(rend.glMaterial_SamplerDisplacement, 6)
			gl.ActiveTexture(oglconsts.TEXTURE6)
			gl.BindTexture(oglconsts.TEXTURE_2D, mfd.VboTextureDisplacement)
		} else {
			gl.Uniform1i(rend.glMaterial_HasTextureDisplacement, 0)
		}

		// effects - gaussian blur
		gl.Uniform1i(rend.glEffect_GB_Mode, mfd.EffectGBlurMode-1)
		gl.Uniform1f(rend.glEffect_GB_W, mfd.EffectGBlurWidth.Point)
		gl.Uniform1f(rend.glEffect_GB_Radius, mfd.EffectGBlurRadius.Point)

		// effects - bloom
		// TODO: Bloom effect
		if mfd.EffectBloomDoBloom {
			gl.Uniform1i(rend.glEffect_Bloom_doBloom, 1)
		} else {
			gl.Uniform1i(rend.glEffect_Bloom_doBloom, 0)
		}
		gl.Uniform1f(rend.glEffect_Bloom_WeightA, mfd.EffectBloomWeightA)
		gl.Uniform1f(rend.glEffect_Bloom_WeightB, mfd.EffectBloomWeightB)
		gl.Uniform1f(rend.glEffect_Bloom_WeightC, mfd.EffectBloomWeightC)
		gl.Uniform1f(rend.glEffect_Bloom_WeightD, mfd.EffectBloomWeightD)
		gl.Uniform1f(rend.glEffect_Bloom_Vignette, mfd.EffectBloomVignette)
		gl.Uniform1f(rend.glEffect_Bloom_VignetteAtt, mfd.EffectBloomVignetteAtt)

		// effects - tone mapping
		if mfd.EffectToneMappingACESFilmRec2020 {
			gl.Uniform1i(rend.glEffect_ToneMapping_ACESFilmRec2020, 1)
		} else {
			gl.Uniform1i(rend.glEffect_ToneMapping_ACESFilmRec2020, 0)
		}
		if mfd.EffectHDRTonemapping {
			gl.Uniform1i(rend.glEffect_HDR_Tonemapping, 1)
		} else {
			gl.Uniform1i(rend.glEffect_HDR_Tonemapping, 0)
		}

		// PBR
		if mfd.SettingRenderingPBR {
			gl.Uniform1i(rend.glPBR_UsePBR, 1)
		} else {
			gl.Uniform1i(rend.glPBR_UsePBR, 0)
		}
		gl.Uniform1f(rend.glPBR_Metallic, mfd.SettingRenderingPBRMetallic)
		gl.Uniform1f(rend.glPBR_Rougness, mfd.SettingRenderingPBRRoughness)
		gl.Uniform1f(rend.glPBR_AO, mfd.SettingRenderingPBRAO)

		gl.Uniform1f(rend.glVS_IsBorder, 0.0)

		mtxModel := mgl32.Ident4()

		// model draw
		gl.Uniform1f(rend.glVS_IsBorder, 0.0)
		mtxModel = matrixModel.Mul4(mgl32.Scale3D(1.0, 1.0, 1.0))
		gl.GLUniformMatrix4fv(rend.glFS_MMatrix, 1, false, &mtxModel[0])

		mfd.VertexSphereVisible = rsett.General.VertexSphereVisible
		mfd.VertexSphereRadius = rsett.General.VertexSphereRadius
		mfd.VertexSphereSegments = rsett.General.VertexSphereSegments
		mfd.VertexSphereColor = rsett.General.VertexSphereColor
		mfd.VertexSphereIsSphere = rsett.General.VertexSphereIsSphere
		mfd.VertexSphereShowWireframes = rsett.General.VertexSphereShowWireframes

		mfd.Render(true)
	}

	gl.EndQuery(oglconsts.TIME_ELAPSED)

	if gl.IsQuery(queries[(currentQuery+1)%querycount]) {
		var result uint64
		gl.GetQueryObjectui64v(queries[(currentQuery+1)%querycount], oglconsts.QUERY_RESULT, &result)
		settings.LogWarn("[RenderingForward - renderModels] OccQuery = %v ms/frame", (result * 1.e6))
	}
	currentQuery = (currentQuery + 1) % querycount

	gl.UseProgram(0)
}

// Dispose ...
func (rend *RendererForward) Dispose() {
	gl := rend.window.OpenGL()

	gl.DeleteVertexArrays([]uint32{rend.glVAO})
	gl.DeleteProgram(rend.shaderProgram)
}
