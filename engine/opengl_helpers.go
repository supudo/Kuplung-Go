package engine

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"io/ioutil"
	"os"

	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
)

// LinkNewProgram creates a new shader program based on the provided shaders.
func LinkNewProgram(gl interfaces.OpenGL, shaders ...uint32) (program uint32, err error) {
	program = gl.CreateProgram()

	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}
	gl.LinkProgram(program)

	if gl.GetProgramParameter(program, oglconsts.LINK_STATUS) == 0 {
		err = fmt.Errorf("%v", gl.GetProgramInfoLog(program))
		gl.DeleteProgram(program)
		program = 0
	}

	return
}

// LinkMultiProgram creates a new shader program based on the provided shaders.
func LinkMultiProgram(gl interfaces.OpenGL, vertexShaderSource, tcsShaderSource, tesShaderSource, geomShaderSource, fragmentShaderSource string) (program uint32, err error) {
	vertexShader, vertexErr := CompileNewShader(gl, oglconsts.VERTEX_SHADER, vertexShaderSource)
	defer gl.DeleteShader(vertexShader)
	tcsShader, tcsErr := CompileNewShader(gl, oglconsts.TESS_CONTROL_SHADER, tcsShaderSource)
	defer gl.DeleteShader(tcsShader)
	tesShader, tesErr := CompileNewShader(gl, oglconsts.TESS_EVALUATION_SHADER, tesShaderSource)
	defer gl.DeleteShader(tesShader)
	geomShader, geomErr := CompileNewShader(gl, oglconsts.GEOMETRY_SHADER, geomShaderSource)
	defer gl.DeleteShader(geomShader)
	fragmentShader, fragmentErr := CompileNewShader(gl, oglconsts.FRAGMENT_SHADER, fragmentShaderSource)
	defer gl.DeleteShader(fragmentShader)

	if (vertexErr == nil) && (tcsErr == nil) && (tesErr == nil) && (geomErr == nil) && (fragmentErr == nil) {
		program, err = LinkNewProgram(gl, vertexShader, tcsShader, tesShader, geomShader, fragmentShader)
	} else {
		err = fmt.Errorf("[OpenGL Utils] Error compiling shaders:\nVertex = %v\nTSC = %v\nTES = %v\nGEOM = %v\nFragment = %v", vertexErr, tcsErr, tesErr, geomErr, fragmentErr)
	}

	return
}

// LinkNewStandardProgram creates a new shader based on two shader sources.
func LinkNewStandardProgram(gl interfaces.OpenGL, vertexShaderSource, fragmentShaderSource string) (program uint32, err error) {
	vertexShader, vertexErr := CompileNewShader(gl, oglconsts.VERTEX_SHADER, vertexShaderSource)
	defer gl.DeleteShader(vertexShader)
	fragmentShader, fragmentErr := CompileNewShader(gl, oglconsts.FRAGMENT_SHADER, fragmentShaderSource)
	defer gl.DeleteShader(fragmentShader)

	if (vertexErr == nil) && (fragmentErr == nil) {
		program, err = LinkNewProgram(gl, vertexShader, fragmentShader)
	} else {
		err = fmt.Errorf("[OpenGL Utils] Error compiling shaders:\n%v\n%v", vertexErr, fragmentErr)
	}

	return
}

// CompileNewShader creates a shader of given type and compiles the provided source.
func CompileNewShader(gl interfaces.OpenGL, shaderType uint32, source string) (shader uint32, err error) {
	shader = gl.CreateShader(shaderType)

	gl.ShaderSource(shader, source)
	gl.CompileShader(shader)

	compileStatus := gl.GetShaderParameter(shader, oglconsts.COMPILE_STATUS)
	if compileStatus == 0 {
		err = fmt.Errorf("%s", gl.GetShaderInfoLog(shader))
		gl.DeleteShader(shader)
		shader = 0
	}

	return
}

// GetShaderSource ...
func GetShaderSource(filename string) string {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		settings.LogWarn("[OpenGL Utils] Can't load shader source for %v", filename)
		return ""
	}
	return string(source) + "\x00"
}

// LoadTexture ...
func LoadTexture(gl interfaces.OpenGL, file string) uint32 {
	imgFile, err := os.Open(file)
	if err != nil {
		settings.LogError("[OpenGL Utils] Texture file not found: %v", err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		settings.LogError("[OpenGL Utils] Can't decode texture: %v", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		settings.LogError("[OpenGL Utils] Texture unsupported stride!")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	texture := gl.GenTextures(1)[0]
	gl.ActiveTexture(oglconsts.TEXTURE0)
	gl.BindTexture(oglconsts.TEXTURE_2D, texture)
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

	return texture
}

// LoadTextureRepeat ...
func LoadTextureRepeat(gl interfaces.OpenGL, file string) uint32 {
	imgFile, err := os.Open(file)
	if err != nil {
		settings.LogError("[OpenGL Utils] Texture file (%v) not found: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		settings.LogError("[OpenGL Utils] Can't decode texture (%v): %v", file, err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		settings.LogError("[OpenGL Utils] Texture unsupported stride! (%v)", file)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	texture := gl.GenTextures(1)[0]
	gl.ActiveTexture(oglconsts.TEXTURE0)
	gl.BindTexture(oglconsts.TEXTURE_2D, texture)
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

	return texture
}

// LoadCubemapTexture ...
func LoadCubemapTexture(gl interfaces.OpenGL, images []string) uint32 {
	tcms := []uint32{
		oglconsts.TEXTURE_CUBE_MAP_POSITIVE_X, // Right
		oglconsts.TEXTURE_CUBE_MAP_NEGATIVE_X, // Left
		oglconsts.TEXTURE_CUBE_MAP_POSITIVE_Y, // Top
		oglconsts.TEXTURE_CUBE_MAP_NEGATIVE_Y, // Bottom
		oglconsts.TEXTURE_CUBE_MAP_POSITIVE_Z, // Back
		oglconsts.TEXTURE_CUBE_MAP_NEGATIVE_Z, // Front
	}

	texture := gl.GenTextures(1)[0]
	gl.ActiveTexture(oglconsts.TEXTURE0)
	gl.BindTexture(oglconsts.TEXTURE_CUBE_MAP, texture)

	sett := settings.GetSettings()
	for i := 0; i < len(images); i++ {
		file := sett.App.CurrentPath + "skybox/" + images[i]
		imgFile, err := os.Open(file)
		if err != nil {
			settings.LogError("[OpenGL Utils] Cubemap texture (%v) - file not found: %v", file, err)
		}
		img, _, err := image.Decode(imgFile)
		if err != nil {
			settings.LogError("[OpenGL Utils] Cubemap texture (%v) - can't decode texture: %v", file, err)
		}

		rgba := image.NewRGBA(img.Bounds())
		if rgba.Stride != rgba.Rect.Size().X*4 {
			settings.LogError("[OpenGL Utils] Cubemap texture (%v) - unsupported stride!", file)
		}
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

		gl.TexImage2D(tcms[i], 0, oglconsts.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, oglconsts.RGBA, oglconsts.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	}

	gl.TexParameteri(oglconsts.TEXTURE_CUBE_MAP, oglconsts.TEXTURE_MAG_FILTER, oglconsts.LINEAR)
	gl.TexParameteri(oglconsts.TEXTURE_CUBE_MAP, oglconsts.TEXTURE_MIN_FILTER, oglconsts.LINEAR)
	gl.TexParameteri(oglconsts.TEXTURE_CUBE_MAP, oglconsts.TEXTURE_WRAP_S, oglconsts.CLAMP_TO_EDGE)
	gl.TexParameteri(oglconsts.TEXTURE_CUBE_MAP, oglconsts.TEXTURE_WRAP_T, oglconsts.CLAMP_TO_EDGE)
	gl.TexParameteri(oglconsts.TEXTURE_CUBE_MAP, oglconsts.TEXTURE_WRAP_R, oglconsts.CLAMP_TO_EDGE)
	gl.BindTexture(oglconsts.TEXTURE_CUBE_MAP, 0)

	return texture
}
