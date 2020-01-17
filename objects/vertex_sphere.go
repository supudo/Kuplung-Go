package objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/engine"
	"github.com/supudo/Kuplung-Go/engine/oglconsts"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/settings"
	"github.com/supudo/Kuplung-Go/types"
	"math"
)

// VertexSphere ...
type VertexSphere struct {
	window interfaces.Window

	shaderProgram                uint32
	shaderVertex, shaderFragment uint32
	glVAO                        uint32

	glUniformMVPMatrix           int32
	glUniformColor               int32
	glUniformInnerLightDirection int32

	IsSphere       bool
	ShowWireframes bool
	CircleSegments int32
	dataVertices   []mgl32.Vec3
	dataNormals    []mgl32.Vec3
	dataIndices    []uint32
}

// InitVertexSphere ...
func InitVertexSphere(window interfaces.Window) *VertexSphere {
	vs := &VertexSphere{
		window:         window,
		IsSphere:       true,
		ShowWireframes: true,
		CircleSegments: 0,
	}
	return vs
}

// InitShaderProgram ...
func (vs *VertexSphere) InitShaderProgram() {
	sett := settings.GetSettings()
	gl := vs.window.OpenGL()

	vertexShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/vertex_sphere.vert")
	fragmentShader := engine.GetShaderSource(sett.App.AppFolder + "shaders/vertex_sphere.frag")

	var err error
	vs.shaderProgram, err = engine.LinkNewStandardProgram(gl, vertexShader, fragmentShader)
	if err != nil {
		settings.LogWarn("[VertexSphere] Can't load the vertex sphere shaders: %v", err)
	}

	vs.glUniformMVPMatrix = gl.GLGetUniformLocation(vs.shaderProgram, gl.Str("u_MVPMatrix\x00"))
	vs.glUniformColor = gl.GLGetUniformLocation(vs.shaderProgram, gl.Str("fs_color\x00"))
	vs.glUniformInnerLightDirection = gl.GLGetUniformLocation(vs.shaderProgram, gl.Str("fs_innerLightDirection\x00"))

	gl.CheckForOpenGLErrors("VertexSphere")
}

// InitBuffers ...
func (vs *VertexSphere) InitBuffers(meshModel types.MeshModel, segments int32, radius float32) {
	if segments == 0 {
		return
	}
	gl := vs.window.OpenGL()

	vs.CircleSegments = segments

	gl.Enable(oglconsts.DEPTH_TEST)
	gl.DepthFunc(oglconsts.LESS)
	gl.Disable(oglconsts.BLEND)
	gl.BlendFunc(oglconsts.SRC_ALPHA, oglconsts.ONE_MINUS_SRC_ALPHA)

	vs.glVAO = gl.GenVertexArrays(1)[0]
	gl.BindVertexArray(vs.glVAO)

	vs.dataVertices = []mgl32.Vec3{}
	vs.dataNormals = []mgl32.Vec3{}
	vs.dataIndices = []uint32{}

	if !vs.IsSphere {
		theta := float64(float64(2*math.Pi) / float64(vs.CircleSegments))
		c := float32(math.Cos(theta))
		s := float32(math.Sin(theta))

		r := radius
		x := r
		y := float32(0)
		var i, j uint32

		for i = 0; i < uint32(len(meshModel.Vertices)); i++ {
			vertex := meshModel.Vertices[i]
			for j = 0; j < uint32(vs.CircleSegments); j++ {
				vs.dataVertices = append(vs.dataVertices, mgl32.Vec3{x + vertex.X(), y + vertex.Y(), vertex.Z()})
				vs.dataIndices = append(vs.dataIndices, i)

				x = c*x - s*y
				y = s*x + c*y
			}
		}
		// vertices
		vboVertices := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(vs.dataVertices)*3*4, gl.Ptr(vs.dataVertices), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		// indices
		vboIndices := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, vboIndices)
		gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, len(vs.dataIndices)*4, gl.Ptr(vs.dataIndices), oglconsts.STATIC_DRAW)
	} else {
		rings := float32(1.0 / (vs.CircleSegments - 1))
		sectors := float32(1.0 / (vs.CircleSegments - 1))

		pi := float32(math.Pi)
		pi2 := float32(math.Phi)
		var position mgl32.Vec3
		var v1, v2, v3, normal mgl32.Vec3
		var x, y int32

		for i := 0; i < len(meshModel.Vertices); i++ {
			for y = 0; y < vs.CircleSegments; y++ {
				for x = 0; x < vs.CircleSegments; x++ {
					px := float32(math.Cos(float64(2*pi*float32(x)*sectors)) * math.Sin(float64(pi*float32(y)*rings)))
					py := float32(math.Sin(float64(-pi2 + pi*float32(y)*rings)))
					pz := float32(math.Sin(float64(2*pi*float32(x)*sectors)) * math.Sin(float64(pi*float32(y)*rings)))
					position = mgl32.Vec3{px * radius, py * radius, pz * radius}
					position = position.Add(meshModel.Vertices[i])
					vs.dataVertices = append(vs.dataVertices, position)
					vs.dataNormals = append(vs.dataNormals, position)
				}
			}

			for y = 0; y < vs.CircleSegments-1; y++ {
				for x = 0; x < vs.CircleSegments-1; x++ {
					start := (y*vs.CircleSegments + x) + (int32(i) * (vs.CircleSegments * vs.CircleSegments))

					vs.dataIndices = append(vs.dataIndices, uint32(start))
					vs.dataIndices = append(vs.dataIndices, uint32(start+1))
					vs.dataIndices = append(vs.dataIndices, uint32(start+vs.CircleSegments))
					v1 = vs.dataVertices[start]
					v2 = vs.dataVertices[start+1]
					v3 = vs.dataVertices[start+vs.CircleSegments]
					vn := v2.Sub(v1)
					vdest := v3.Sub(v1)
					normal = vn.Cross(vdest)
					vs.dataNormals[start] = normal
					vs.dataNormals[start+1] = normal
					vs.dataNormals[start+vs.CircleSegments] = normal

					vs.dataIndices = append(vs.dataIndices, uint32(start+1))
					vs.dataIndices = append(vs.dataIndices, uint32(start+1+vs.CircleSegments))
					vs.dataIndices = append(vs.dataIndices, uint32(start+vs.CircleSegments))
					v1 = vs.dataVertices[start]
					v2 = vs.dataVertices[start+vs.CircleSegments]
					v3 = vs.dataVertices[start+1+vs.CircleSegments]
					vn = v2.Sub(v1)
					vdest = v3.Sub(v1)
					normal = vn.Cross(vdest)
					vs.dataNormals[start] = normal
					vs.dataNormals[start+vs.CircleSegments] = normal
					vs.dataNormals[start+1+vs.CircleSegments] = normal
				}
			}
		}

		// vertices
		vboVertices := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboVertices)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(vs.dataVertices)*3*4, gl.Ptr(vs.dataVertices), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		// normals
		vboNormals := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ARRAY_BUFFER, vboNormals)
		gl.BufferData(oglconsts.ARRAY_BUFFER, len(vs.dataNormals)*3*4, gl.Ptr(vs.dataNormals), oglconsts.STATIC_DRAW)
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, oglconsts.FLOAT, false, 3*4, gl.PtrOffset(0))

		// indices
		vboIndices := gl.GenBuffers(1)[0]
		gl.BindBuffer(oglconsts.ELEMENT_ARRAY_BUFFER, vboIndices)
		gl.BufferData(oglconsts.ELEMENT_ARRAY_BUFFER, len(vs.dataIndices)*4, gl.Ptr(vs.dataIndices), oglconsts.STATIC_DRAW)

		gl.BindVertexArray(0)
		gl.DeleteBuffers([]uint32{vboVertices, vboIndices})
	}

	gl.CheckForOpenGLErrors("VertexSphere-InitBuffers")
}

// Render ...
func (vs *VertexSphere) Render(mtxModel mgl32.Mat4, color mgl32.Vec4) {
	rsett := settings.GetRenderingSettings()
	gl := vs.window.OpenGL()

	gl.UseProgram(vs.shaderProgram)
	gl.BindVertexArray(vs.glVAO)

	mvpMatrix := rsett.MatrixProjection.Mul4(rsett.MatrixCamera.Mul4(mtxModel))
	gl.GLUniformMatrix4fv(vs.glUniformMVPMatrix, 1, false, &mvpMatrix[0])
	gl.Uniform3f(vs.glUniformColor, color.X(), color.Y(), color.Z())
	gl.Uniform3f(vs.glUniformInnerLightDirection, 1.0, 0.55, 0.206)

	if !vs.IsSphere {
		var i int32
		for i = 0; i < vs.CircleSegments; i++ {
			gl.DrawArrays(oglconsts.LINE_LOOP, int32(vs.CircleSegments*i), int32(vs.CircleSegments))
		}
	} else {
		if vs.ShowWireframes {
			gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.LINE)
		}
		gl.DrawElements(oglconsts.TRIANGLES, int32(len(vs.dataIndices)), oglconsts.UNSIGNED_INT, 0)
		if vs.ShowWireframes {
			gl.PolygonMode(oglconsts.FRONT_AND_BACK, oglconsts.FILL)
		}
	}

	gl.BindVertexArray(0)
	gl.UseProgram(0)

	gl.CheckForOpenGLErrors("VertexSphere-Render")
}

// Dispose will cleanup everything
func (vs *VertexSphere) Dispose() {
	gl := vs.window.OpenGL()
	gl.DeleteVertexArrays([]uint32{vs.glVAO})
	gl.DeleteProgram(vs.shaderProgram)
}
