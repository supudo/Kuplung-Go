package renderers

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// RendererSimple ...
type RendererSimple struct {
	window interfaces.Window

	shaderProgram uint32
	glVAO         uint32

	glMVPMatrix, glWorldMatrix            int32
	glSamplerTexture, glHasSamplerTexture int32
	glCameraPosition, glUIAmbient         int32

	solidLight *types.ModelFaceLightSourceDirectional
}

// NewRendererSimple ...
func NewRendererSimple(window interfaces.Window) *RendererSimple {
	sett := settings.GetSettings()
	gl := window.OpenGL()

	rend := &RendererSimple{}
	rend.window = window

	sVertex := engine.GetShaderSource(sett.App.AppFolder + "shaders/rendering_simple.vert")
	sTcs := engine.GetShaderSource(sett.App.AppFolder + "shaders/rendering_simple.tcs")
	sTes := engine.GetShaderSource(sett.App.AppFolder + "shaders/rendering_simple.tes")
	sGeom := engine.GetShaderSource(sett.App.AppFolder + "shaders/rendering_simple.geom")
	sFragment := engine.GetShaderSource(sett.App.AppFolder + "shaders/rendering_simple.frag")

	var err error
	rend.shaderProgram, err = engine.LinkMultiProgram(gl, sVertex, sTcs, sTes, sGeom, sFragment)
	if err != nil {
		settings.LogWarn("[RendererSimple] Can't load the renderer simple shaders: %v", err)
	}

	rend.glMVPMatrix = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("vs_MVPMatrix\x00"))

	rend.glMVPMatrix = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("vs_MVPMatrix\x00"))
	rend.glWorldMatrix = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("vs_WorldMatrix\x00"))

	rend.glHasSamplerTexture = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("has_texture\x00"))
	rend.glSamplerTexture = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("sampler_texture\x00"))

	rend.glCameraPosition = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_cameraPosition\x00"))
	rend.glUIAmbient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("fs_UIAmbient\x00"))

	rend.solidLight = &types.ModelFaceLightSourceDirectional{}
	rend.solidLight.InUse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.inUse\x00"))
	rend.solidLight.Direction = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.direction\x00"))
	rend.solidLight.Ambient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.ambient\x00"))
	rend.solidLight.Diffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.diffuse\x00"))
	rend.solidLight.Specular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.specular\x00"))
	rend.solidLight.StrengthAmbient = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.strengthAmbient\x00"))
	rend.solidLight.StrengthDiffuse = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.strengthDiffuse\x00"))
	rend.solidLight.StrengthSpecular = gl.GLGetUniformLocation(rend.shaderProgram, gl.Str("solidSkin_Light.strengthSpecular\x00"))

	gl.CheckForOpenGLErrors("ForwardRenderer")

	return rend
}

// Render ...
func (rend *RendererSimple) Render(rp types.RenderProperties, meshModelFaces []*meshes.ModelFace, matrixGrid mgl32.Mat4, camPos mgl32.Vec3) {
	gl := rend.window.OpenGL()
	rsett := settings.GetRenderingSettings()

	gl.UseProgram(rend.shaderProgram)

	for i := 0; i < len(meshModelFaces); i++ {
		mfd := meshModelFaces[i]

		matrixModel := mgl32.Ident4()
		matrixModel = matrixModel.Mul4(matrixGrid)
		// scale
		matrixModel = matrixModel.Mul4(mgl32.Scale3D(mfd.ScaleX.Point, mfd.ScaleY.Point, mfd.ScaleZ.Point))
		// rotate
		matrixModel = matrixModel.Mul4(mgl32.Translate3D(0, 0, 0))
		matrixModel = matrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(mfd.RotateX.Point), mgl32.Vec3{1, 0, 0}))
		matrixModel = matrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(mfd.RotateY.Point), mgl32.Vec3{0, 1, 0}))
		matrixModel = matrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(mfd.RotateZ.Point), mgl32.Vec3{0, 0, 1}))
		matrixModel = matrixModel.Mul4(mgl32.Translate3D(0, 0, 0))
		// translate
		matrixModel = matrixModel.Mul4(mgl32.Translate3D(mfd.PositionX.Point, mfd.PositionY.Point, mfd.PositionZ.Point))

		mfd.MatrixModel = matrixModel
		// mfd->Setting_ModelViewSkin = rp.viewModelSkin;
		// mfd->lightSources = rp.lightSources;
		// mfd->setOptionsFOV(rp.Setting_FOV);
		// mfd->setOptionsOutlineColor(rp.Setting_OutlineColor);
		// mfd->setOptionsOutlineThickness(rp.Setting_OutlineThickness);

		mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(matrixModel))
		gl.GLUniformMatrix4fv(rend.glMVPMatrix, 1, false, &mvpMatrix[0])

		gl.GLUniformMatrix4fv(rend.glWorldMatrix, 1, false, &matrixModel[0])

		gl.Uniform3f(rend.glCameraPosition, camPos.X(), camPos.Y(), camPos.Z())
		gl.Uniform3f(rend.glUIAmbient, rp.UIAmbientLightX, rp.UIAmbientLightY, rp.UIAmbientLightZ)

		gl.Uniform1i(rend.solidLight.InUse, 1)
		gl.Uniform3f(rend.solidLight.Direction, rp.SolidLightDirectionX, rp.SolidLightDirectionY, rp.SolidLightDirectionZ)
		gl.Uniform3f(rend.solidLight.Ambient, rp.SolidLightAmbient.X(), rp.SolidLightAmbient.Y(), rp.SolidLightAmbient.Z())
		gl.Uniform3f(rend.solidLight.Diffuse, rp.SolidLightDiffuse.X(), rp.SolidLightDiffuse.Y(), rp.SolidLightDiffuse.Z())
		gl.Uniform3f(rend.solidLight.Specular, rp.SolidLightSpecular.X(), rp.SolidLightSpecular.Y(), rp.SolidLightSpecular.Z())
		gl.Uniform1f(rend.solidLight.StrengthAmbient, rp.SolidLightAmbientStrength)
		gl.Uniform1f(rend.solidLight.StrengthDiffuse, rp.SolidLightDiffuseStrength)
		gl.Uniform1f(rend.solidLight.StrengthSpecular, rp.SolidLightSpecularStrength)

		if mfd.HasTextureDiffuse && mfd.MeshModel.ModelMaterial.TextureDiffuse.UseTexture {
			gl.Uniform1i(rend.glHasSamplerTexture, 1)
			gl.Uniform1i(rend.glSamplerTexture, 0)
			gl.ActiveTexture(oglconsts.TEXTURE0)
			gl.BindTexture(oglconsts.TEXTURE_2D, mfd.VboTextureDiffuse)
		} else {
			gl.Uniform1i(rend.glHasSamplerTexture, 0)
		}

		mfd.Render(true)
	}

	gl.UseProgram(0)
}

// Dispose ...
func (rend *RendererSimple) Dispose() {
	gl := rend.window.OpenGL()

	gl.DeleteVertexArrays([]uint32{rend.glVAO})
	gl.DeleteProgram(rend.shaderProgram)
}
