package objects

import (
	"image"
	"image/draw"
	"os"
	"time"

	"github.com/inkyblackness/imgui-go"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
)

// Shadertoy ...
type Shadertoy struct {
	window interfaces.Window

	IChannel0Image     string
	IChannel1Image     string
	IChannel2Image     string
	IChannel3Image     string
	IChannel0CubeImage string
	IChannel1CubeImage string
	IChannel2CubeImage string
	IChannel3CubeImage string
	TextureWidth       int32
	TextureHeight      int32

	iChannelResolution0 [2]float32
	iChannelResolution1 [2]float32
	iChannelResolution2 [2]float32
	iChannelResolution3 [2]float32

	shaderProgram      uint32
	glVAO              uint32
	vsInFBO            int32
	vsScreenResolution int32

	iResolution int32
	iGlobalTime int32
	iTimeDelta  int32
	iFrame      int32
	iFrameRate  int32

	iChannelTime       [4]int32
	iChannelResolution [4]int32

	iMouse int32
	iDate  int32

	iChannel0 uint32
	iChannel1 uint32
	iChannel2 uint32
	iChannel3 uint32

	tFBO uint32
	tRBO uint32

	version string
}

// InitShadertoy ...
func InitShadertoy(window interfaces.Window) *Shadertoy {
	return &Shadertoy{
		window:              window,
		iChannelResolution0: [2]float32{0.0, 0.0},
		iChannelResolution1: [2]float32{0.0, 0.0},
		iChannelResolution2: [2]float32{0.0, 0.0},
		iChannelResolution3: [2]float32{0.0, 0.0},
		TextureWidth:        0,
		TextureHeight:       0,
		version:             "#version 410 core",
	}
}

// InitShaderProgram ...
func (st *Shadertoy) InitShaderProgram(shaderSource string) {
	sett := settings.GetSettings()
	gl := st.window.OpenGL()

	vertexShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/shadertoy.vert")

	fragmentShader := st.version + `

out vec4 outFragmentColor;
uniform vec3 iResolution;
uniform float iGlobalTime;
uniform float iTimeDelta;
uniform int iFrame;
uniform int iFrameRate;
uniform float iChannelTime[4];
uniform vec3 iChannelResolution[4];
uniform vec4 iMouse;
uniform vec4 iDate;
uniform float iSampleRate;
` + "\x00"

	if len(st.IChannel0Image) > 0 {
		fragmentShader += `uniform sampler2D iChannel0;`
	} else if len(st.IChannel0CubeImage) > 0 {
		fragmentShader += `uniform samplerCube iChannel0;`
	}

	if len(st.IChannel1Image) > 0 {
		fragmentShader += `uniform sampler2D iChannel1;`
	} else if len(st.IChannel1CubeImage) > 0 {
		fragmentShader += `uniform samplerCube iChannel1;`
	}

	if len(st.IChannel2Image) > 0 {
		fragmentShader += `uniform sampler2D iChannel2;`
	} else if len(st.IChannel2CubeImage) > 0 {
		fragmentShader += `uniform samplerCube iChannel2;`
	}

	if len(st.IChannel3Image) > 0 {
		fragmentShader += `uniform sampler2D iChannel3;`
	} else if len(st.IChannel3CubeImage) > 0 {
		fragmentShader += `uniform samplerCube iChannel3;`
	}

	fragmentShader += `

#define texture2D texture
#define textureCube texture

`

	fragmentShader += shaderSource

	fragmentShader += `

void main() {
	vec4 color = vec4(0.0, 0.0, 0.0, 1.0);
	mainImage(color, gl_FragCoord.xy);
	outFragmentColor = color;
}
` + "\x00"

	var err error
	st.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[Shadertoy] Can't load the Shadertoy shaders: %v", err)
	}
	st.vsInFBO = gl.GLGetUniformLocation(st.shaderProgram, gl.Str("vs_inFBO\x00"))
	st.vsScreenResolution = gl.GLGetUniformLocation(st.shaderProgram, gl.Str("vs_screenResolution\x00"))

	st.iResolution = gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iResolution\x00"))
	st.iGlobalTime = gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iGlobalTime\x00"))
	st.iTimeDelta = gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iTimeDelta\x00"))
	st.iFrameRate = gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iFrameRate\x00"))
	st.iFrame = gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iFrame\x00"))
	st.iChannelTime = [4]int32{
		gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannelTime[0]\x00")),
		gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannelTime[1]\x00")),
		gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannelTime[2]\x00")),
		gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannelTime[3]\x00"))}
	st.iChannelResolution = [4]int32{
		gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannelResolution[0]\x00")),
		gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannelResolution[1]\x00")),
		gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannelResolution[2]\x00")),
		gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannelResolution[3]\x00"))}
	st.iMouse = gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iMouse\x00"))
	st.iDate = gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iDate\x00"))
	st.iChannel0 = uint32(gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannel0\x00")))
	st.iChannel1 = uint32(gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannel1\x00")))
	st.iChannel2 = uint32(gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannel2\x00")))
	st.iChannel3 = uint32(gl.GLGetUniformLocation(st.shaderProgram, gl.Str("iChannel3\x00")))
}

// InitBuffers ...
func (st *Shadertoy) InitBuffers() {
	sett := settings.GetSettings()
	gl := st.window.OpenGL()

	st.glVAO = gl.GenVertexArrays(1)[0]

	gl.BindVertexArray(st.glVAO)

	vertices := []float32{
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		1.0, 1.0, 0.0,
		1.0, 1.0, 0.0,
		-1.0, 1.0, 0.0,
		-1.0, -1.0, 0.0}

	vboVertices := gl.GenBuffers(1)[0]
	gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
	gl.BufferData(oglconsts.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), oglconsts.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 4, gl.PtrOffset(0))

	st.IChannel0Image = sett.App.AppFolder + "/shadertoy/noise16.png"
	st.initTextures()

	gl.BindVertexArray(0)
}

func (st *Shadertoy) initTextures() {
	tc := int32(0)
	if len(st.IChannel0Image) > 0 {
		st.addTexture(st.IChannel0Image, &st.iChannel0, tc)
		tc++
	}
	if len(st.IChannel1Image) > 0 {
		st.addTexture(st.IChannel1Image, &st.iChannel1, tc)
		tc++
	}
	if len(st.IChannel2Image) > 0 {
		st.addTexture(st.IChannel2Image, &st.iChannel2, tc)
		tc++
	}
	if len(st.IChannel3Image) > 0 {
		st.addTexture(st.IChannel3Image, &st.iChannel3, tc)
	}
}

func (st *Shadertoy) addTexture(textureImage string, vboTexture *uint32, textureID int32) {
	gl := st.window.OpenGL()

	imgFile, err := os.Open(textureImage)
	if err != nil {
		settings.LogError("[Shadertoy] Texture file not found: %v", err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		settings.LogError("[Shadertoy] Can't decode texture: %v", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		settings.LogError("[Shadertoy] Texture unsupported stride!")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	*vboTexture = gl.GenTextures(1)[0]
	gl.ActiveTexture(oglconsts.TEXTURE0)
	gl.BindTexture(oglconsts.TEXTURE_2D, *vboTexture)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MIN_FILTER, oglconsts.LINEAR)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MAG_FILTER, oglconsts.LINEAR)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_WRAP_S, oglconsts.CLAMP_TO_EDGE)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_WRAP_T, oglconsts.CLAMP_TO_EDGE)
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

	tWidth := float32(rgba.Bounds().Size().X)
	tHeight := float32(rgba.Bounds().Size().Y)
	switch textureID {
	case 0:
		st.iChannelResolution0[0] = tWidth
		st.iChannelResolution0[1] = tHeight
	case 1:
		st.iChannelResolution1[0] = tWidth
		st.iChannelResolution1[1] = tHeight
	case 2:
		st.iChannelResolution2[0] = tWidth
		st.iChannelResolution2[1] = tHeight
	case 3:
		st.iChannelResolution3[0] = tWidth
		st.iChannelResolution3[1] = tHeight
	}
}

// Render ...
func (st *Shadertoy) Render(deltaTime float32, mouseX, mouseY int32, seconds float32) {
	gl := st.window.OpenGL()
	if st.glVAO > 0 {
		gl.UseProgram(st.shaderProgram)

		gl.Enable(oglconsts.TEXTURE_2D)

		gl.Uniform1i(st.vsInFBO, 1)

		tc := int32(0)
		if len(st.IChannel0Image) > 0 {
			gl.ActiveTexture(uint32(oglconsts.TEXTURE0 + tc))
			gl.BindTexture(oglconsts.TEXTURE_2D, st.iChannel0)
			gl.Uniform1i(int32(st.iChannel0), tc)
			gl.Uniform3f(st.iChannelResolution[0], st.iChannelResolution0[0], st.iChannelResolution0[1], 0.0)
			tc++
		}
		if len(st.IChannel1Image) > 0 {
			gl.ActiveTexture(uint32(oglconsts.TEXTURE1 + tc))
			gl.BindTexture(oglconsts.TEXTURE_2D, st.iChannel1)
			gl.Uniform1i(int32(st.iChannel0), tc)
			gl.Uniform3f(st.iChannelResolution[0], st.iChannelResolution1[0], st.iChannelResolution1[1], 0.0)
			tc++
		}
		if len(st.IChannel2Image) > 0 {
			gl.ActiveTexture(uint32(oglconsts.TEXTURE2 + tc))
			gl.BindTexture(oglconsts.TEXTURE_2D, st.iChannel2)
			gl.Uniform1i(int32(st.iChannel0), tc)
			gl.Uniform3f(st.iChannelResolution[0], st.iChannelResolution2[0], st.iChannelResolution2[1], 0.0)
			tc++
		}
		if len(st.IChannel3Image) > 0 {
			gl.ActiveTexture(uint32(oglconsts.TEXTURE3 + tc))
			gl.BindTexture(oglconsts.TEXTURE_2D, st.iChannel3)
			gl.Uniform1i(int32(st.iChannel0), tc)
			gl.Uniform3f(st.iChannelResolution[0], st.iChannelResolution3[0], st.iChannelResolution3[1], 0.0)
		}

		gl.Uniform2i(st.vsScreenResolution, st.TextureWidth, st.TextureHeight)
		gl.Uniform3i(st.iResolution, st.TextureWidth, st.TextureHeight, 0)
		gl.Uniform1f(st.iGlobalTime, seconds)
		gl.Uniform4f(st.iMouse, float32(mouseX), float32(mouseY), 0.0, 0.0)
		gl.Uniform1f(st.iChannelTime[0], seconds)
		gl.Uniform1f(st.iChannelTime[1], seconds)
		gl.Uniform1f(st.iChannelTime[2], seconds)
		gl.Uniform1f(st.iChannelTime[3], seconds)
		gl.Uniform1f(st.iTimeDelta, deltaTime)

		timet := time.Now()
		gl.Uniform4f(st.iDate, float32(timet.Year()), float32(timet.Month()), float32(timet.Day()), float32(timet.Second()))

		gl.Uniform1f(st.iFrameRate, imgui.CurrentIO().Framerate())
		gl.Uniform1f(st.iFrame, 0.0)

		// draw
		gl.BindVertexArray(st.glVAO)
		gl.DrawArrays(oglconsts.TRIANGLES, 0, 6)
		gl.BindVertexArray(0)

		gl.UseProgram(0)
	}
}

// InitFBO ...
func (st *Shadertoy) InitFBO(windowWidth, windowHeight int32, vboTexture *uint32) {
	gl := st.window.OpenGL()

	st.TextureWidth = windowWidth
	st.TextureHeight = windowHeight

	*vboTexture = gl.GenTextures(1)[0]
	gl.BindTexture(oglconsts.TEXTURE_2D, *vboTexture)
	gl.TexParameterf(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MAG_FILTER, oglconsts.LINEAR)
	gl.TexParameterf(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_MIN_FILTER, oglconsts.LINEAR_MIPMAP_LINEAR)
	gl.TexParameterf(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_WRAP_S, oglconsts.CLAMP_TO_EDGE)
	gl.TexParameterf(oglconsts.TEXTURE_2D, oglconsts.TEXTURE_WRAP_T, oglconsts.CLAMP_TO_EDGE)
	gl.TexParameteri(oglconsts.TEXTURE_2D, oglconsts.GENERATE_MIPMAP, oglconsts.TRUE)
	gl.TexImage2D(oglconsts.TEXTURE_2D, 0, oglconsts.RGBA8, windowWidth, windowHeight, 0, oglconsts.RGBA, oglconsts.UNSIGNED_BYTE, nil)
	gl.BindTexture(oglconsts.TEXTURE_2D, 0)

	st.tRBO = gl.GenRenderbuffers(1)[0]
	gl.BindRenderbuffer(oglconsts.RENDERBUFFER, st.tRBO)
	gl.RenderbufferStorage(oglconsts.RENDERBUFFER, oglconsts.DEPTH_COMPONENT, windowWidth, windowHeight)
	gl.BindRenderbuffer(oglconsts.RENDERBUFFER, 0)

	st.tFBO = gl.GenFramebuffers(1)[0]
	gl.BindFramebuffer(oglconsts.FRAMEBUFFER, st.tFBO)

	gl.FramebufferTexture2D(oglconsts.FRAMEBUFFER, oglconsts.COLOR_ATTACHMENT0, oglconsts.TEXTURE_2D, *vboTexture, 0)
	gl.FramebufferRenderbuffer(oglconsts.FRAMEBUFFER, oglconsts.DEPTH_ATTACHMENT, oglconsts.RENDERBUFFER, st.tRBO)

	if gl.CheckFramebufferStatus(oglconsts.FRAMEBUFFER) != oglconsts.FRAMEBUFFER_COMPLETE {
		settings.LogError("[Shadertoy] Error creating FBO!")
	}

	gl.BindFramebuffer(oglconsts.FRAMEBUFFER, 0)
}

func (st *Shadertoy) bindFBO() {
	gl := st.window.OpenGL()
	gl.BindFramebuffer(oglconsts.FRAMEBUFFER, st.tFBO)
	gl.Clear(oglconsts.COLOR_BUFFER_BIT | oglconsts.DEPTH_BUFFER_BIT)
}

func (st *Shadertoy) unbindFBO(vboTexture *uint32) {
	gl := st.window.OpenGL()
	gl.BindFramebuffer(oglconsts.FRAMEBUFFER, 0)
	gl.BindTexture(oglconsts.TEXTURE_2D, *vboTexture)
	gl.GenerateMipmap(oglconsts.TEXTURE_2D)
	gl.BindTexture(oglconsts.TEXTURE_2D, 0)
}

// RenderToTexture ...
func (st *Shadertoy) RenderToTexture(deltaTime float32, mouseX, mouseY int32, seconds float32, vboTexture *uint32) {
	st.bindFBO()
	st.Render(deltaTime, mouseX, mouseY, seconds)
	st.unbindFBO(vboTexture)
}
