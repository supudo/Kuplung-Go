package objects

import (
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
)

// Light ...
type Light struct {
	window interfaces.Window

	shaderProgram      uint32
	glVAO              uint32
	glUniformMVPMatrix int32
	glUniformSampler   int32
	glUniformUseColor  int32
	glUniformColor     int32

	MatrixModel mgl32.Mat4

	model types.MeshModel

	lightDirectionRay *RayLine

	hasTexture        bool
	vboTextureDiffuse uint32

	Title, Description                string
	LightType                         types.LightSourceType
	ShowLampObject, ShowLampDirection bool
	TurnOffPosition, ShowInWire       bool

	PositionX, PositionY, PositionZ             types.ObjectCoordinate
	DirectionX, DirectionY, DirectionZ          types.ObjectCoordinate
	ScaleX, ScaleY, ScaleZ                      types.ObjectCoordinate
	RotateX, RotateY, RotateZ                   types.ObjectCoordinate
	RotateCenterX, RotateCenterY, RotateCenterZ types.ObjectCoordinate
	Ambient, Diffuse, Specular                  types.MaterialColor
	LCutOff, LOuterCutOff                       types.ObjectCoordinate
	LConstant, LLinear, LQuadratic              types.ObjectCoordinate
}

// InitLight ...
func InitLight(window interfaces.Window) *Light {
	sett := settings.GetSettings()
	gl := window.OpenGL()

	lightModel := &Light{}
	lightModel.window = window

	vertexShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/light.vert")
	fragmentShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/light.frag")

	var err error
	lightModel.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[Light] Can't load the light shaders: %v", err)
	}

	lightModel.glUniformMVPMatrix = gl.GLGetUniformLocation(lightModel.shaderProgram, gl.Str("u_MVPMatrix\x00"))
	lightModel.glUniformSampler = gl.GLGetUniformLocation(lightModel.shaderProgram, gl.Str("u_sampler\x00"))
	lightModel.glUniformUseColor = gl.GLGetUniformLocation(lightModel.shaderProgram, gl.Str("fs_useColor\x00"))
	lightModel.glUniformColor = gl.GLGetUniformLocation(lightModel.shaderProgram, gl.Str("fs_color\x00"))

	return lightModel
}

// InitProperties ...
func (l *Light) InitProperties(shape types.LightSourceType) {
	l.ShowLampObject = true
	l.ShowLampDirection = true
	l.ShowInWire = false
	l.LightType = shape

	l.PositionX = types.ObjectCoordinate{Animate: false, Point: 0.0}
	l.PositionY = types.ObjectCoordinate{Animate: false, Point: 5.0}
	l.PositionZ = types.ObjectCoordinate{Animate: false, Point: 0.0}

	l.DirectionX = types.ObjectCoordinate{Animate: false, Point: 0.0}
	l.DirectionY = types.ObjectCoordinate{Animate: false, Point: 1.0}
	l.DirectionZ = types.ObjectCoordinate{Animate: false, Point: 0.0}

	l.ScaleX = types.ObjectCoordinate{Animate: false, Point: 1.0}
	l.ScaleY = types.ObjectCoordinate{Animate: false, Point: 1.0}
	l.ScaleZ = types.ObjectCoordinate{Animate: false, Point: 1.0}

	l.RotateX = types.ObjectCoordinate{Animate: false, Point: 0.0} // -71.0f
	l.RotateY = types.ObjectCoordinate{Animate: false, Point: 0.0} // -36.0f
	l.RotateZ = types.ObjectCoordinate{Animate: false, Point: 0.0}

	l.RotateCenterX = types.ObjectCoordinate{Animate: false, Point: 0.0}
	l.RotateCenterY = types.ObjectCoordinate{Animate: false, Point: 0.0}
	l.RotateCenterZ = types.ObjectCoordinate{Animate: false, Point: 0.0}

	l.Ambient = types.MaterialColor{ColorPickerOpen: false, Animate: false, Strength: 0.3, Color: mgl32.Vec3{1.0, 1.0, 1.0}}
	l.Diffuse = types.MaterialColor{ColorPickerOpen: false, Animate: false, Strength: 1.0, Color: mgl32.Vec3{1.0, 1.0, 1.0}}

	switch shape {
	case types.LightSourceTypeDirectional:
		{
			l.LConstant = types.ObjectCoordinate{Animate: false, Point: 0.0}
			l.LLinear = types.ObjectCoordinate{Animate: false, Point: 0.0}
			l.LQuadratic = types.ObjectCoordinate{Animate: false, Point: 0.0}
			l.Specular = types.MaterialColor{ColorPickerOpen: false, Animate: false, Strength: 0.0, Color: mgl32.Vec3{1.0, 1.0, 1.0}}
			l.LCutOff = types.ObjectCoordinate{Animate: false, Point: -180.0}
			l.LOuterCutOff = types.ObjectCoordinate{Animate: false, Point: 160.0}
		}
	case types.LightSourceTypePoint:
		{
			l.LConstant = types.ObjectCoordinate{Animate: false, Point: 0.0}
			l.LLinear = types.ObjectCoordinate{Animate: false, Point: 0.2}
			l.LQuadratic = types.ObjectCoordinate{Animate: false, Point: 0.05}
			l.Specular = types.MaterialColor{ColorPickerOpen: false, Animate: false, Strength: 0.0, Color: mgl32.Vec3{1.0, 1.0, 1.0}}
			l.LCutOff = types.ObjectCoordinate{Animate: false, Point: -180.0}
			l.LOuterCutOff = types.ObjectCoordinate{Animate: false, Point: 160.0}
		}
	case types.LightSourceTypeSpot:
		{
			l.LConstant = types.ObjectCoordinate{Animate: false, Point: 1.0}
			l.LLinear = types.ObjectCoordinate{Animate: false, Point: 0.09}
			l.LQuadratic = types.ObjectCoordinate{Animate: false, Point: 0.032}
			l.Specular = types.MaterialColor{ColorPickerOpen: false, Animate: false, Strength: 1.0, Color: mgl32.Vec3{1.0, 1.0, 1.0}}
			l.LCutOff = types.ObjectCoordinate{Animate: false, Point: 12.5}
			l.LOuterCutOff = types.ObjectCoordinate{Animate: false, Point: 15.0}
		}
	}

	l.MatrixModel = mgl32.Ident4()

	lrFrom := mgl32.Vec3{0, 0, 0}
	lrTo := mgl32.Vec3{l.PositionX.Point, l.PositionY.Point * -100.0, l.PositionZ.Point}
	l.lightDirectionRay = NewLightRay(l.window)
	l.lightDirectionRay.InitBuffers(lrFrom, lrTo)

	l.TurnOffPosition = false
}

// SetModel ...
func (l *Light) SetModel(model types.MeshModel) {
	l.model = model
}

// InitBuffers ...
func (l *Light) InitBuffers() {
	gl := l.window.OpenGL()

	l.glVAO = gl.GenVertexArrays(1)[0]

	gl.BindVertexArray(l.glVAO)

	vboVertices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(l.model.Vertices)*3*4, gl.Ptr(l.model.Vertices), oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

	vboNormals := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboNormals)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(l.model.Normals)*3*4, gl.Ptr(l.model.Normals), oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

	// textures and colors
	l.hasTexture = false
	if len(l.model.TextureCoordinates) > 0 {
		l.hasTexture = true
		vboTextureCoordinates := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboTextureCoordinates)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(l.model.TextureCoordinates)*3*4, gl.Ptr(l.model.TextureCoordinates), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 2, oglconsts.FLOAT, false, 2*4, gl.PtrOffset(0))

		if len(l.model.ModelMaterial.TextureDiffuse.Image) > 0 {
			sett := settings.GetSettings()
			file := sett.App.AppFolder + "gui/" + l.model.ModelMaterial.TextureDiffuse.Image
			imgFile, err := os.Open(file)
			if err != nil {
				settings.LogError("[Light] Texture file not found: %v", err)
			}
			img, _, err := image.Decode(imgFile)
			if err != nil {
				settings.LogError("[Light] Can't decode texture: %v", err)
			}

			rgba := image.NewRGBA(img.Bounds())
			if rgba.Stride != rgba.Rect.Size().X*4 {
				settings.LogError("[Light] Texture unsupported stride!")
			}
			draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

			l.vboTextureDiffuse = gl.GenTextures(1)[0]
			gl.ActiveTexture(oglconsts.TEXTURE0)
			gl.BindTexture(oglconsts.TEXTURE_2D, l.vboTextureDiffuse)
			gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MIN_FILTER, oglconsts.LINEAR)
			gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MAG_FILTER, oglconsts.LINEAR)
			gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_WRAP_S, oglconsts.REPEAT)
			gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_WRAP_T, oglconsts.REPEAT)
			gl.TexImage2D(
				oglconsts.TEXTURE_2D,
				0,
				oglconsts.RGBA,
				int32(rgba.Rect.Size().X),
				int32(rgba.Rect.Size().Y),
				0,
				oglconsts.RGBA,
				oglconsts.UNSIGNED_BYTE,
				gl.Ptr(rgba.Pix))
		}
	}

	vboIndices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, vboIndices)
	gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, int(l.model.CountIndices)*4, gl.Ptr(l.model.Indices), oglconsts.STATIC_DRAW)

	gl.BindVertexArray(0)

	gl.DeleteBuffers([]uint32{vboVertices, vboNormals, vboIndices})
}

// Render ...
func (l *Light) Render() {
	if l.ShowLampObject {
		gl := l.window.OpenGL()
		rsett := settings.GetRenderingSettings()

		gl.UseProgram(l.shaderProgram)

		l.MatrixModel = mgl32.Ident4()

		l.MatrixModel = l.MatrixModel.Mul4(mgl32.Scale3D(l.ScaleX.Point, l.ScaleY.Point, l.ScaleZ.Point))

		l.MatrixModel = l.MatrixModel.Mul4(mgl32.Translate3D(0, 0, 0))
		l.MatrixModel = l.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(l.RotateX.Point), mgl32.Vec3{1, 0, 0}))
		l.MatrixModel = l.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(l.RotateY.Point), mgl32.Vec3{0, 1, 0}))
		l.MatrixModel = l.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(l.RotateZ.Point), mgl32.Vec3{0, 0, 1}))
		l.MatrixModel = l.MatrixModel.Mul4(mgl32.Translate3D(0, 0, 0))

		if !l.TurnOffPosition {
			l.MatrixModel = l.MatrixModel.Mul4(mgl32.Translate3D(l.PositionX.Point, l.PositionY.Point, l.PositionZ.Point))
		}

		l.MatrixModel = l.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(l.RotateCenterX.Point), mgl32.Vec3{1, 0, 0}))
		l.MatrixModel = l.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(l.RotateCenterY.Point), mgl32.Vec3{0, 1, 0}))
		l.MatrixModel = l.MatrixModel.Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(l.RotateCenterZ.Point), mgl32.Vec3{0, 0, 1}))

		if l.hasTexture {
			gl.BindTexture(oglconsts.TEXTURE_2D, l.vboTextureDiffuse)
		}

		mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(l.MatrixModel))
		gl.GLUniformMatrix4fv(l.glUniformMVPMatrix, 1, false, &mvpMatrix[0])

		if l.vboTextureDiffuse == 0 {
			gl.Uniform1i(l.glUniformUseColor, 1)
			gl.Uniform3f(l.glUniformColor, l.model.ModelMaterial.DiffuseColor.X(), l.model.ModelMaterial.DiffuseColor.Y(), l.model.ModelMaterial.DiffuseColor.Z())
		}

		gl.BindVertexArray(l.glVAO)
		if l.ShowInWire {
			gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.LINE)
		}
		gl.DrawElements(oglconsts.TRIANGLES, l.model.CountIndices, oglconsts.UNSIGNED_INT, 0)
		if l.ShowInWire {
			gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.FILL)
		}
		gl.BindVertexArray(0)

		if l.hasTexture {
			gl.BindTexture(oglconsts.TEXTURE_2D, 0)
		}

		gl.UseProgram(0)

		if l.lightDirectionRay != nil {
			l.lightDirectionRay.MatrixModel = l.MatrixModel
			if l.ShowLampDirection {
				l.lightDirectionRay.Render(l.MatrixModel)
			}
		}
	}
}

// Dispose will cleanup everything
func (l *Light) Dispose() {
	gl := l.window.OpenGL()

	gl.DeleteVertexArrays([]uint32{l.glVAO})
	gl.DeleteProgram(l.shaderProgram)
}
