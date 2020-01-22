package renderers

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/objects"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// RendererDefered ...
type RendererDefered struct {
	window interfaces.Window

	shaderProgramGeometryPass     uint32
	shaderProgramLightingPass     uint32
	shaderProgramLightBox         uint32
	glGeometryPassTextureDiffuse  int32
	glGeometryPassTextureSpecular int32

	matrixProjection, matrixCamera mgl32.Mat4

	GLSLLightSourceNumberDirectional uint32
	GLSLLightSourceNumberPoint       uint32
	GLSLLightSourceNumberSpot        uint32

	mfLightsDirectional []*types.ModelFaceLightSourceDirectional
	mfLightsPoint       []*types.ModelFaceLightSourcePoint
	mfLightsSpot        []*types.ModelFaceLightSourceSpot

	gBuffer, gPosition, gNormal, gAlbedoSpec uint32

	lightPositions  []mgl32.Vec3
	lightColors     []mgl32.Vec3
	objectPositions []mgl32.Vec3

	NRLIGHTS uint16

	quadVAO uint32
	quadVBO uint32
	cubeVAO uint32
	cubeVBO uint32
}

// NewRendererDefered ...
func NewRendererDefered(window interfaces.Window) *RendererDefered {
	rend := &RendererDefered{}
	rend.window = window
	rend.Init()
	return rend
}

// Init ...
func (rend *RendererDefered) Init() {
	rend.NRLIGHTS = 32
	rend.GLSLLightSourceNumberDirectional = 0
	rend.GLSLLightSourceNumberPoint = 0
	rend.GLSLLightSourceNumberSpot = 0

	rend.initGeometryPass()
	rend.initLighingPass()
	rend.initLightObjects()
	rend.initProps()
	rend.initGBuffer()
	rend.initLights()
}

func (rend *RendererDefered) initGeometryPass() {
	sett := settings.GetSettings()
	gl := rend.window.OpenGL()

	sVertex := engine.GetShaderSource(sett.App.AppFolder + "shaders/deferred_g_buffer.vert")
	sFragment := engine.GetShaderSource(sett.App.AppFolder + "shaders/deferred_g_buffer.frag")

	var err error
	rend.shaderProgramGeometryPass, err = engine.LinkNewStandardProgram(gl, sVertex, sFragment)
	if err != nil {
		settings.LogWarn("[RendererDefered] Can't load the renderer defered shaders: %v", err)
	}

	rend.glGeometryPassTextureDiffuse = gl.GLGetUniformLocation(rend.shaderProgramGeometryPass, gl.Str("texture_diffuse\x00"))
	rend.glGeometryPassTextureSpecular = gl.GLGetUniformLocation(rend.shaderProgramGeometryPass, gl.Str("texture_specular\x00"))

	gl.CheckForOpenGLErrors("DeferedRenderer - initGeometryPass")
	settings.LogInfo("[Defered Renderer] Geometry Pass initialized.")
}

func (rend *RendererDefered) initLighingPass() {
	sett := settings.GetSettings()
	gl := rend.window.OpenGL()

	sVertex := engine.GetShaderSource(sett.App.AppFolder + "shaders/deferred_shading.vert")
	sFragment := engine.GetShaderSource(sett.App.AppFolder + "shaders/deferred_shading.frag")

	var err error
	rend.shaderProgramLightingPass, err = engine.LinkNewStandardProgram(gl, sVertex, sFragment)
	if err != nil {
		settings.LogWarn("[RendererDefered] Can't load the renderer defered shaders: %v", err)
	}

	gl.CheckForOpenGLErrors("DeferedRenderer - initLighingPass")
	settings.LogInfo("[Defered Renderer] Lighting Pass initialized.")
}

func (rend *RendererDefered) initLightObjects() {
	sett := settings.GetSettings()
	gl := rend.window.OpenGL()

	sVertex := engine.GetShaderSource(sett.App.AppFolder + "shaders/deferred_light_box.vert")
	sFragment := engine.GetShaderSource(sett.App.AppFolder + "shaders/deferred_light_box.frag")

	var err error
	rend.shaderProgramLightBox, err = engine.LinkNewStandardProgram(gl, sVertex, sFragment)
	if err != nil {
		settings.LogWarn("[RendererDefered] Can't load the renderer defered shaders: %v", err)
	}

	gl.CheckForOpenGLErrors("DeferedRenderer - initLightObjects")
	settings.LogInfo("[Defered Renderer] Light Objects initialized.")
}

func (rend *RendererDefered) initProps() {
	rsett := settings.GetRenderingSettings()
	gl := rend.window.OpenGL()

	gl.UseProgram(rend.shaderProgramLightingPass)

	gl.Uniform1i(gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("sampler_position\x00")), 0)
	gl.Uniform1i(gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("sampler_normal\x00")), 1)
	gl.Uniform1i(gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("sampler_albedospec\x00")), 2)

	rend.objectPositions = []mgl32.Vec3{}
	if rsett.Defered.DeferredTestMode {
		rend.objectPositions = append(rend.objectPositions, mgl32.Vec3{0.0, -3.0, 0.0})
		rend.objectPositions = append(rend.objectPositions, mgl32.Vec3{-3.0, -3.0, -3.0})
		rend.objectPositions = append(rend.objectPositions, mgl32.Vec3{0.0, -3.0, -3.0})
		rend.objectPositions = append(rend.objectPositions, mgl32.Vec3{3.0, -3.0, -3.0})
		rend.objectPositions = append(rend.objectPositions, mgl32.Vec3{-3.0, -3.0, 0.0})
		rend.objectPositions = append(rend.objectPositions, mgl32.Vec3{3.0, -3.0, 0.0})
		rend.objectPositions = append(rend.objectPositions, mgl32.Vec3{-3.0, -3.0, 3.0})
		rend.objectPositions = append(rend.objectPositions, mgl32.Vec3{0.0, -3.0, 3.0})
		rend.objectPositions = append(rend.objectPositions, mgl32.Vec3{3.0, -3.0, 3.0})
	} else {
		rend.objectPositions = append(rend.objectPositions, mgl32.Vec3{0.0, 0.0, 0.0})
	}

	for i := 0; i < len(rend.objectPositions); i++ {
		rend.objectPositions[i] = mgl32.Vec3{rend.objectPositions[i].X(), 0.0, rend.objectPositions[i].Y()}
	}

	rend.lightPositions = []mgl32.Vec3{}
	rend.lightColors = []mgl32.Vec3{}
	for i := uint16(0); i < rend.NRLIGHTS; i++ {
		xPos := float32(((math.Mod(rand.Float64(), 100))/100.0)*6.0 - 3.0)
		yPos := float32(((math.Mod(rand.Float64(), 100))/100.0)*6.0 - 0.0)
		zPos := float32(((math.Mod(rand.Float64(), 100))/100.0)*6.0 - 3.0)
		rend.lightPositions = append(rend.lightPositions, mgl32.Vec3{xPos, yPos, zPos})

		rColor := float32((math.Mod(rand.Float64(), 100) / 200.0) + 0.5)
		gColor := float32((math.Mod(rand.Float64(), 100) / 200.0) + 0.5)
		bColor := float32((math.Mod(rand.Float64(), 100) / 200.0) + 0.5)
		rend.lightColors = append(rend.lightColors, mgl32.Vec3{rColor, gColor, bColor})
	}

	gl.CheckForOpenGLErrors("DeferedRenderer - initProps")
	settings.LogInfo("[Defered Renderer] Properties initialized.")
}

func (rend *RendererDefered) initGBuffer() {
	sett := settings.GetSettings()
	gl := rend.window.OpenGL()

	rend.gBuffer = gl.GenFramebuffers(1)[0]
	gl.BindFramebuffer(oglconsts.FRAMEBUFFER, rend.gBuffer)

	// - Position color buffer
	rend.gPosition = gl.GenTextures(1)[0]
	gl.BindTexture(oglconsts.TEXTURE_2D, rend.gPosition)
	gl.TexImage2D(oglconsts.TEXTURE_2D, 0, oglconsts.RGB16F, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight), 0, oglconsts.RGB, oglconsts.FLOAT, nil)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MIN_FILTER, oglconsts.NEAREST)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MAG_FILTER, oglconsts.NEAREST)
	gl.FramebufferTexture2D(oglconsts.FRAMEBUFFER, oglconsts.COLOR_ATTACHMENT0, oglconsts.TEXTURE_2D, rend.gPosition, 0)

	// - Normal color buffer
	rend.gNormal = gl.GenTextures(1)[0]
	gl.BindTexture(oglconsts.TEXTURE_2D, rend.gNormal)
	gl.TexImage2D(oglconsts.TEXTURE_2D, 0, oglconsts.RGB16F, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight), 0, oglconsts.RGB, oglconsts.FLOAT, nil)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MIN_FILTER, oglconsts.NEAREST)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MAG_FILTER, oglconsts.NEAREST)
	gl.FramebufferTexture2D(oglconsts.FRAMEBUFFER, oglconsts.COLOR_ATTACHMENT1, oglconsts.TEXTURE_2D, rend.gNormal, 0)

	// - Color + Specular color buffer
	rend.gAlbedoSpec = gl.GenTextures(1)[0]
	gl.BindTexture(oglconsts.TEXTURE_2D, rend.gAlbedoSpec)
	gl.TexImage2D(oglconsts.TEXTURE_2D, 0, oglconsts.RGBA, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight), 0, oglconsts.RGBA, oglconsts.UNSIGNED_BYTE, nil)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MIN_FILTER, oglconsts.NEAREST)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MAG_FILTER, oglconsts.NEAREST)
	gl.FramebufferTexture2D(oglconsts.FRAMEBUFFER, oglconsts.COLOR_ATTACHMENT2, oglconsts.TEXTURE_2D, rend.gAlbedoSpec, 0)

	// - Tell OpenGL which color attachments we'll use (of this framebuffer) for rendering
	attachments := []uint32{oglconsts.COLOR_ATTACHMENT0, oglconsts.COLOR_ATTACHMENT1, oglconsts.COLOR_ATTACHMENT2}
	gl.DrawBuffers(3, &attachments[0])

	// - Create and attach depth buffer (renderbuffer)
	rboDepth := gl.GenRenderbuffers(1)[0]
	gl.BindRenderbuffer(oglconsts.RENDERBUFFER, rboDepth)
	gl.RenderbufferStorage(oglconsts.RENDERBUFFER, oglconsts.DEPTH_COMPONENT, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight))
	gl.FramebufferRenderbuffer(oglconsts.FRAMEBUFFER, oglconsts.DEPTH_ATTACHMENT, oglconsts.RENDERBUFFER, rboDepth)
	// - Finally check if framebuffer is complete
	if gl.CheckFramebufferStatus(oglconsts.FRAMEBUFFER) != oglconsts.FRAMEBUFFER_COMPLETE {
		settings.LogError("[Deferred Rendering T GBuffer] Framebuffer not complete!")
	}

	gl.BindFramebuffer(oglconsts.FRAMEBUFFER, 0)

	gl.CheckForOpenGLErrors("DeferedRenderer - initGBuffer")
	settings.LogInfo("[Defered Renderer] GBuffer initialized.")
}

func (rend *RendererDefered) initLights() {
	gl := rend.window.OpenGL()

	var i uint32

	// light - directional
	for i = 0; i < rend.GLSLLightSourceNumberDirectional; i++ {
		f := &types.ModelFaceLightSourceDirectional{}

		f.InUse = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("directionalLights["+fmt.Sprint(i)+"].inUse\x00"))

		f.Direction = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("directionalLights["+fmt.Sprint(i)+"].direction\x00"))

		f.Ambient = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("directionalLights["+fmt.Sprint(i)+"].ambient\x00"))
		f.Diffuse = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("directionalLights["+fmt.Sprint(i)+"].diffuse\x00"))
		f.Specular = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("directionalLights["+fmt.Sprint(i)+"].specular\x00"))

		f.StrengthAmbient = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("directionalLights["+fmt.Sprint(i)+"].strengthAmbient\x00"))
		f.StrengthDiffuse = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("directionalLights["+fmt.Sprint(i)+"].strengthDiffuse\x00"))
		f.StrengthSpecular = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("directionalLights["+fmt.Sprint(i)+"].strengthSpecular\x00"))
		rend.mfLightsDirectional = append(rend.mfLightsDirectional, f)
	}

	// light - point
	for i = 0; i < rend.GLSLLightSourceNumberPoint; i++ {
		f := &types.ModelFaceLightSourcePoint{}

		f.InUse = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("pointLights["+fmt.Sprint(i)+"].inUse\x00"))
		f.Position = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("pointLights["+fmt.Sprint(i)+"].position\x00"))

		f.Constant = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("pointLights["+fmt.Sprint(i)+"].constant\x00"))
		f.Linear = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("pointLights["+fmt.Sprint(i)+"].linear\x00"))
		f.Quadratic = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("pointLights["+fmt.Sprint(i)+"].quadratic\x00"))

		f.Ambient = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("pointLights["+fmt.Sprint(i)+"].ambient\x00"))
		f.Diffuse = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("pointLights["+fmt.Sprint(i)+"].diffuse\x00"))
		f.Specular = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("pointLights["+fmt.Sprint(i)+"].specular\x00"))

		f.StrengthAmbient = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("pointLights["+fmt.Sprint(i)+"].strengthAmbient\x00"))
		f.StrengthDiffuse = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("pointLights["+fmt.Sprint(i)+"].strengthDiffuse\x00"))
		f.StrengthSpecular = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("pointLights["+fmt.Sprint(i)+"].strengthSpecular\x00"))
		rend.mfLightsPoint = append(rend.mfLightsPoint, f)
	}

	// light - spot
	for i = 0; i < rend.GLSLLightSourceNumberSpot; i++ {
		f := &types.ModelFaceLightSourceSpot{}

		f.InUse = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].inUse\x00"))

		f.Position = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].position\x00"))
		f.Direction = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].direction\x00"))

		f.CutOff = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].cutOff\x00"))
		f.OuterCutOff = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].outerCutOff\x00"))

		f.Constant = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].constant\x00"))
		f.Linear = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].linear\x00"))
		f.Quadratic = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].quadratic\x00"))

		f.Ambient = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].ambient\x00"))
		f.Diffuse = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].diffuse\x00"))
		f.Specular = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].specular\x00"))

		f.StrengthAmbient = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].strengthAmbient\x00"))
		f.StrengthDiffuse = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].strengthDiffuse\x00"))
		f.StrengthSpecular = gl.GLGetUniformLocation(rend.shaderProgramLightingPass, gl.Str("spotLights["+fmt.Sprint(i)+"].strengthSpecular\x00"))
		rend.mfLightsSpot = append(rend.mfLightsSpot, f)
	}

	gl.CheckForOpenGLErrors("DeferedRenderer - initLights")
	settings.LogInfo("[Defered Renderer] Lights initialized.")
}

// Render ...
func (rend *RendererDefered) Render(rp types.RenderProperties, meshModelFaces []*meshes.ModelFace, matrixGrid mgl32.Mat4, camPos mgl32.Vec3, selectedModel int32, lightSources []*objects.Light) {
	sett := settings.GetSettings()
	rsett := settings.GetRenderingSettings()
	gl := rend.window.OpenGL()

	if rsett.Defered.DeferredRandomizeLightPositions {
		rend.Init()
		rsett.Defered.DeferredRandomizeLightPositions = false
	}

	rend.renderGBuffer(meshModelFaces, selectedModel)
	rend.renderLightingPass(lightSources)
	if rsett.Defered.DeferredTestLights {
		rend.renderLightObjects()
	} else {
		gl.BindFramebuffer(oglconsts.READ_FRAMEBUFFER, rend.gBuffer)
		gl.BindFramebuffer(oglconsts.DRAW_FRAMEBUFFER, 0)
		gl.BlitFramebuffer(0, 0, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight), 0, 0, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight), oglconsts.DEPTH_BUFFER_BIT, oglconsts.NEAREST)
		gl.BindFramebuffer(oglconsts.FRAMEBUFFER, 0)
	}

	gl.CheckForOpenGLErrors("DeferedRenderer - Render")
}

func (rend *RendererDefered) renderGBuffer(meshModelFaces []*meshes.ModelFace, selectedModel int32) {
	gl := rend.window.OpenGL()
	rsett := settings.GetRenderingSettings()

	gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.FILL)

	gl.BindFramebuffer(oglconsts.FRAMEBUFFER, rend.gBuffer)
	gl.Clear(oglconsts.COLOR_BUFFER_BIT | oglconsts.DEPTH_BUFFER_BIT)
	gl.UseProgram(rend.shaderProgramGeometryPass)
	gl.GLUniformMatrix4fv(gl.GLGetUniformLocation(rend.shaderProgramGeometryPass, gl.Str("projection\x00")), 1, false, &rsett.MatrixProjection[0])
	gl.GLUniformMatrix4fv(gl.GLGetUniformLocation(rend.shaderProgramGeometryPass, gl.Str("view\x00")), 1, false, &rsett.MatrixCamera[0])

	op := int32(0)
	if rsett.Defered.DeferredTestMode {
		op = int32(len(rend.objectPositions))
	} else {
		op = 1
	}

	for i := int32(0); i < op; i++ {
		matrixModel := mgl32.Ident4()

		matrixModel = matrixModel.Mul4(mgl32.Translate3D(rend.objectPositions[i].X(), rend.objectPositions[i].Y(), rend.objectPositions[i].Z()))
		if len(rend.objectPositions) > 1 {
			matrixModel = matrixModel.Mul4(mgl32.Scale3D(0.25, 0.25, 0.25))
		}

		matrixModel = matrixModel.Mul4(mgl32.Translate3D(0, 0, 0))
		matrixModel = matrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(-90.0), mgl32.Vec3{1, 0, 0}))
		matrixModel = matrixModel.Mul4(mgl32.Translate3D(0, 0, 0))

		for j := 0; j < len(meshModelFaces); j++ {
			mfd := meshModelFaces[j]

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

			gl.GLUniformMatrix4fv(gl.GLGetUniformLocation(rend.shaderProgramGeometryPass, gl.Str("model\x00")), 1, false, &matrixModel[0])

			if mfd.HasTextureDiffuse && mfd.MeshModel.ModelMaterial.TextureDiffuse.UseTexture {
				gl.Uniform1i(rend.glGeometryPassTextureDiffuse, 0)
				gl.ActiveTexture(oglconsts.TEXTURE0)
				gl.BindTexture(oglconsts.TEXTURE_2D, mfd.VboTextureDiffuse)
			}

			if mfd.HasTextureDiffuse && mfd.MeshModel.ModelMaterial.TextureSpecular.UseTexture {
				gl.Uniform1i(rend.glGeometryPassTextureSpecular, 1)
				gl.ActiveTexture(oglconsts.TEXTURE1)
				gl.BindTexture(oglconsts.TEXTURE_2D, mfd.VboTextureSpecular)
			}

			mfd.MatrixModel = matrixModel

			mfd.ModelViewSkin = rsett.General.SelectedViewModelSkin
			mfd.OutlineColor = rsett.General.OutlineColor
			mfd.OutlineThickness = rsett.General.OutlineThickness
			mfd.IsModelSelected = i == selectedModel
			mfd.VertexSphereVisible = rsett.General.VertexSphereVisible
			mfd.VertexSphereRadius = rsett.General.VertexSphereRadius
			mfd.VertexSphereSegments = rsett.General.VertexSphereSegments
			mfd.VertexSphereColor = rsett.General.VertexSphereColor
			mfd.VertexSphereIsSphere = rsett.General.VertexSphereIsSphere
			mfd.VertexSphereShowWireframes = rsett.General.VertexSphereShowWireframes

			mfd.Render(true)
		}
	}
	gl.BindFramebuffer(oglconsts.FRAMEBUFFER, 0)
	gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.FILL)

	gl.CheckForOpenGLErrors("DeferedRenderer - renderGBuffer")
}

func (rend *RendererDefered) renderLightingPass(lightSources []*objects.Light) {
	gl := rend.window.OpenGL()

	gl.Clear(oglconsts.COLOR_BUFFER_BIT | oglconsts.DEPTH_BUFFER_BIT)
	gl.UseProgram(rend.shaderProgramLightingPass)

	gl.ActiveTexture(oglconsts.TEXTURE0)
	gl.BindTexture(oglconsts.TEXTURE_2D, rend.gPosition)
	gl.ActiveTexture(oglconsts.TEXTURE1)
	gl.BindTexture(oglconsts.TEXTURE_2D, rend.gNormal)
	gl.ActiveTexture(oglconsts.TEXTURE2)
	gl.BindTexture(oglconsts.TEXTURE_2D, rend.gAlbedoSpec)

	// lights
	// lights
	lightsCountDirectional := uint32(0)
	lightsCountPoint := uint32(0)
	lightsCountSpot := uint32(0)
	for j := 0; j < len(lightSources); j++ {
		light := lightSources[j]
		switch light.LightType {
		case types.LightSourceTypeDirectional:
			if lightsCountDirectional < rend.GLSLLightSourceNumberDirectional {
				f := rend.mfLightsDirectional[lightsCountDirectional]

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
			if lightsCountPoint < rend.GLSLLightSourceNumberPoint {
				f := rend.mfLightsPoint[lightsCountPoint]

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
			if lightsCountSpot < rend.GLSLLightSourceNumberSpot {
				f := rend.mfLightsSpot[lightsCountSpot]

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

	gl.CheckForOpenGLErrors("DeferedRenderer - renderLightingPass")
}

func (rend *RendererDefered) renderLightObjects() {
	sett := settings.GetSettings()
	rsett := settings.GetRenderingSettings()
	gl := rend.window.OpenGL()

	gl.BindFramebuffer(oglconsts.READ_FRAMEBUFFER, rend.gBuffer)
	gl.BindFramebuffer(oglconsts.DRAW_FRAMEBUFFER, 0)
	gl.BlitFramebuffer(0, 0, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight), 0, 0, int32(sett.AppWindow.SDLWindowWidth), int32(sett.AppWindow.SDLWindowHeight), oglconsts.DEPTH_BUFFER_BIT, oglconsts.NEAREST)
	gl.BindFramebuffer(oglconsts.FRAMEBUFFER, 0)

	// 3. Render lights on top of scene, by blitting
	gl.UseProgram(rend.shaderProgramLightBox)
	gl.GLUniformMatrix4fv(gl.GLGetUniformLocation(rend.shaderProgramLightBox, gl.Str("projection\x00")), 1, false, &rsett.MatrixProjection[0])
	gl.GLUniformMatrix4fv(gl.GLGetUniformLocation(rend.shaderProgramLightBox, gl.Str("view\x00")), 1, false, &rsett.MatrixCamera[0])
	for i := int32(0); i < rsett.Defered.DeferredTestLightsNumber; i++ {
		matrixModel := mgl32.Ident4()
		matrixModel = matrixModel.Mul4(mgl32.Translate3D(rend.lightPositions[i].X(), rend.lightPositions[i].Y(), rend.lightPositions[i].Z()))
		matrixModel = matrixModel.Mul4(mgl32.Scale3D(0.25, 0.25, 0.25))
		gl.GLUniformMatrix4fv(gl.GLGetUniformLocation(rend.shaderProgramLightBox, gl.Str("model\x00")), 1, false, &matrixModel[0])
		gl.Uniform3fv(gl.GLGetUniformLocation(rend.shaderProgramLightBox, gl.Str("lightColor\x00")), 1, &rend.lightColors[i][0])
		rend.renderCube()
	}

	gl.CheckForOpenGLErrors("DeferedRenderer - renderLightObjects")
}

func (rend *RendererDefered) renderQuad() {
	gl := rend.window.OpenGL()

	if rend.quadVAO == 0 {
		quadVertices := []float32{
			// Positions
			-1.0, 1.0, 0.0,
			// Texture Coords
			0.0, 1.0, -1.0, -1.0, 0.0, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0, 1.0, 1.0, -1.0, 0.0, 1.0, 0.0}
		// Setup plane VAO

		rend.quadVAO = gl.GenVertexArrays(1)[0]
		rend.quadVBO = gl.GenBuffers(1)[0]
		gl.BindVertexArray(rend.quadVAO)
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, rend.quadVBO)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(quadVertices)*4, gl.Ptr(quadVertices), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 5*4, gl.PtrOffset(0))
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, oglconsts.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	}
	gl.BindVertexArray(rend.quadVAO)
	gl.DrawArrays(oglconsts.TRIANGLE_STRIP, 0, 4)
	gl.BindVertexArray(0)

	gl.CheckForOpenGLErrors("DeferedRenderer - renderQuad")
}

func (rend *RendererDefered) renderCube() {
	gl := rend.window.OpenGL()

	if rend.cubeVAO == 0 {
		vertices := []float32{
			// Back face
			-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0, // Bottom-left
			0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0, // top-right
			0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 0.0, // bottom-right
			0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 1.0, 1.0, // top-right
			-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0, // bottom-left
			-0.5, 0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 1.0, // top-left
			// Front face
			-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
			0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 0.0, // bottom-right
			0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
			0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 1.0, 1.0, // top-right
			-0.5, 0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0, // top-left
			-0.5, -0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom-left
			// Left face
			-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0, // top-right
			-0.5, 0.5, -0.5, -1.0, 0.0, 0.0, 1.0, 1.0, // top-left
			-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0, // bottom-left
			-0.5, -0.5, -0.5, -1.0, 0.0, 0.0, 0.0, 1.0, // bottom-left
			-0.5, -0.5, 0.5, -1.0, 0.0, 0.0, 0.0, 0.0, // bottom-right
			-0.5, 0.5, 0.5, -1.0, 0.0, 0.0, 1.0, 0.0, // top-right
			// Right face
			0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0, // top-left
			0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0, // bottom-right
			0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 1.0, // top-right
			0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 1.0, // bottom-right
			0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 1.0, 0.0, // top-left
			0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0, // bottom-left
			// Bottom face
			-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0, // top-right
			0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 1.0, 1.0, // top-left
			0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0, // bottom-left
			0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 1.0, 0.0, // bottom-left
			-0.5, -0.5, 0.5, 0.0, -1.0, 0.0, 0.0, 0.0, // bottom-right
			-0.5, -0.5, -0.5, 0.0, -1.0, 0.0, 0.0, 1.0, // top-right
			// Top face
			-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0, // top-left
			0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom-right
			0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0, // top-right
			0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom-right
			-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0, // top-left
			-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0} // bottom-left
		rend.cubeVAO = gl.GenVertexArrays(1)[0]
		rend.cubeVBO = gl.GenBuffers(1)[0]
		// Fill buffer
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, rend.cubeVBO)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), oglconsts.STATIC_DRAW)
		// Link vertex attributes
		gl.BindVertexArray(rend.cubeVAO)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 8*4, gl.PtrOffset(0))
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, oglconsts.FLOAT, false, 8*4, gl.PtrOffset(3*4))
		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 3, oglconsts.FLOAT, false, 8*4, gl.PtrOffset(6*4))
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}
	// Render Cube
	gl.BindVertexArray(rend.cubeVAO)
	gl.DrawArrays(oglconsts.TRIANGLES, 0, 36)
	gl.BindVertexArray(0)

	gl.CheckForOpenGLErrors("DeferedRenderer - renderCube")
}

// Dispose ...
func (rend *RendererDefered) Dispose() {
	gl := rend.window.OpenGL()

	gl.DeleteProgram(rend.shaderProgramGeometryPass)
	gl.DeleteProgram(rend.shaderProgramLightingPass)
	gl.DeleteProgram(rend.shaderProgramLightBox)
}
