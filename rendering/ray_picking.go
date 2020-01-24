package rendering

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/meshes"
	"github.com/supudo/Kuplung-Go/objects"
	"github.com/supudo/Kuplung-Go/settings"
)

// RayPicking ...
type RayPicking struct {
	window interfaces.Window

	rayLines []*objects.RayLine

	sceneSelectedModelObject int32
}

// InitRayPicking ...
func InitRayPicking(window interfaces.Window) *RayPicking {
	return &RayPicking{
		window: window,
	}
}

// SelectModel ...
func (rp *RayPicking) SelectModel(meshModelFaces []*meshes.ModelFace, sceneSelectedModelObject *int32) []*objects.RayLine {
	rp.sceneSelectedModelObject = *sceneSelectedModelObject
	rp.rayLines = []*objects.RayLine{}
	rp.pickModel(meshModelFaces)
	*sceneSelectedModelObject = rp.sceneSelectedModelObject
	return rp.rayLines
}

// SelectVertex ...
func (rp *RayPicking) SelectVertex(meshModelFaces []*meshes.ModelFace, rayLines []*objects.RayLine, sceneSelectedModelObject *int32) {
}

func (rp *RayPicking) pickModel(meshModelFaces []*meshes.ModelFace) {
	rsett := settings.GetRenderingSettings()
	mousex := rsett.Controls.MouseX
	mousey := rsett.Controls.MouseY

	w, h := rp.window.Size()

	vFrom, _ := rp.getRay(int32(w/2), int32(h/2), int32(w), int32(h), rsett.MatrixCamera, rsett.MatrixProjection)

	normalizedCoordinates := rp.getNormalizeDeviceCordinates(mousex, mousey, int32(w), int32(h))
	clipCoordinates := mgl32.Vec4{normalizedCoordinates.X(), normalizedCoordinates.Y(), -1.0, 1.0}
	eyeCoordinates := rp.getEyeCoordinates(clipCoordinates, rsett.MatrixProjection)

	invertedViewMatrix := rsett.MatrixCamera.Inv()
	rayWorld := rp.mat4MulVec4(invertedViewMatrix, eyeCoordinates)
	vTo := mgl32.Vec3{rayWorld.X(), rayWorld.Y(), rayWorld.Z()}
	vTo = vTo.Mul(rsett.General.PlaneFar)

	if rsett.General.ShowPickRays {
		rl := objects.NewLightRay(rp.window)
		rl.InitBuffers(vFrom, vTo)
		if rsett.General.ShowPickRaysSingle {
			for i := 0; i < len(rp.rayLines); i++ {
				rp.rayLines[i].Dispose()
			}
			rp.rayLines = nil
		}
		rp.rayLines = append(rp.rayLines, rl)
	}

	rp.sceneSelectedModelObject = -1
	for i := 0; i < len(meshModelFaces); i++ {
		mmf := meshModelFaces[i]
		for j := 0; j < len(mmf.MeshModel.Vertices); j++ {
			if (j+1)%3 == 0 {
				id := float32(-1)
				aabbMin := mgl32.Vec3{mmf.BoundingBox.MinX, mmf.BoundingBox.MinY, mmf.BoundingBox.MinZ}
				aabbMax := mgl32.Vec3{mmf.BoundingBox.MaxX, mmf.BoundingBox.MaxY, mmf.BoundingBox.MaxZ}
				if rp.testRayOBBIntersection(vFrom, vTo, aabbMin, aabbMax, mmf.MatrixModel, id) {
					rp.sceneSelectedModelObject = int32(i)
				}
			}
		}
	}
}

func (rp *RayPicking) getRay(mouseX, mouseY, screenWidth, screenHeight int32, ViewMatrix, ProjectionMatrix mgl32.Mat4) (outОrigin, outDirection mgl32.Vec3) {
	lRayStartNDC := mgl32.Vec4{
		(float32(mouseX/screenWidth) - 0.5) * 2.0,
		(float32(mouseY/screenHeight) - 0.5) * 2.0,
		-1.0, 1.0}
	lRayEndNDC := mgl32.Vec4{
		(float32(mouseX/screenWidth) - 0.5) * 2.0,
		(float32(mouseY/screenHeight) - 0.5) * 2.0,
		0.0, 1.0}

	mtx := ProjectionMatrix.Mul4(ViewMatrix)
	M := mtx.Inv()
	lRayStartWorld := rp.mat4MulVec4(M, lRayStartNDC)
	lRayStartWorld = rp.vec4DivFloat(lRayStartWorld, lRayStartWorld.W())
	lRayEndWorld := rp.mat4MulVec4(M, lRayEndNDC)
	lRayEndWorld = rp.vec4DivFloat(lRayEndWorld, lRayEndWorld.W())

	lRayDirWorld := lRayEndWorld.Sub(lRayStartWorld)
	lRayDirWorld = lRayDirWorld.Normalize()

	outОrigin = mgl32.Vec3{lRayStartWorld.X(), lRayStartWorld.Y(), lRayStartWorld.Z()}
	od := lRayDirWorld.Normalize()
	outDirection = mgl32.Vec3{od.X(), od.Y(), od.Z()}
	return outОrigin, outDirection
}

func (rp *RayPicking) getNormalizeDeviceCordinates(X, Y int32, screenWidth, screenHeight int32) mgl32.Vec2 {
	return mgl32.Vec2{
		float32((2*X)/screenWidth - 1),
		float32((2*Y)/screenHeight - 1)}
}

func (rp *RayPicking) getEyeCoordinates(coordinates mgl32.Vec4, mtxProjection mgl32.Mat4) mgl32.Vec4 {
	invertedProjectionMatrix := mtxProjection.Inv()
	eyeCoordinates := rp.mat4MulVec4(invertedProjectionMatrix, coordinates)
	return mgl32.Vec4{eyeCoordinates.X(), eyeCoordinates.Y(), -1.0, 0.0}
}

func (rp *RayPicking) testRayOBBIntersection(rayOrigin, rayDirection, aabbMin, aabbMax mgl32.Vec3, ModelMatrix mgl32.Mat4, intersectionDistance float32) bool {
	tMin := float32(0.0)
	tMax := float32(100000.0)
	OBBpositionWorldspace := mgl32.Vec3{ModelMatrix[3*4+0] + 0, ModelMatrix[3*4+1], ModelMatrix[3*4+2]}
	delta := OBBpositionWorldspace.Sub(rayOrigin)
	xaxis := mgl32.Vec3{ModelMatrix[0], ModelMatrix[0], ModelMatrix[0]}
	e := xaxis.Dot(delta)
	f := rayDirection.Dot(xaxis)
	if math.Abs(float64(f)) > 0.001 {
		t1 := (e + aabbMin.X()) / f
		t2 := (e + aabbMax.Y()) / f

		if t1 > t2 {
			t2, t1 = t1, t2
		}
		if t2 < tMax {
			tMax = t2
		}
		if t1 > tMin {
			tMin = t1
		}
		if tMax < tMin {
			return false
		}
	} else {
		if -e+aabbMin.X() > 0.0 || -e+aabbMax.X() < 0.0 {
			return false
		}
	}

	yaxis := mgl32.Vec3{ModelMatrix[1*4+0], ModelMatrix[1*4+1], ModelMatrix[1*4+2]}
	e = yaxis.Dot(delta)
	f = rayDirection.Dot(yaxis)
	if math.Abs(float64(f)) > 0.001 {
		t1 := (e + aabbMin.Y()) / f
		t2 := (e + aabbMax.Y()) / f

		if t1 > t2 {
			t1, t2 = t2, t1
		}

		if t2 < tMax {
			tMax = t2
		}
		if t1 > tMin {
			tMin = t1
		}
		if tMin > tMax {
			return false
		}
	} else {
		if -e+aabbMin.Y() > 0.0 || -e+aabbMax.Y() < 0.0 {
			return false
		}
	}

	zaxis := mgl32.Vec3{ModelMatrix[2*4+0], ModelMatrix[2*4+1], ModelMatrix[2*4+2]}
	e = zaxis.Dot(delta)
	f = rayDirection.Dot(zaxis)
	if math.Abs(float64(f)) > 0.001 {
		t1 := (e + aabbMin.Z()) / f
		t2 := (e + aabbMax.Z()) / f

		if t1 > t2 {
			t2, t1 = t1, t2
		}
		if t2 < tMax {
			tMax = t2
		}
		if t1 > tMin {
			tMin = t1
		}
		if tMin > tMax {
			return false
		}
	} else {
		if -e+aabbMin.Z() > 0.0 || -e+aabbMax.Z() < 0.0 {
			return false
		}
	}

	intersectionDistance = tMin
	return true
}

func (rp *RayPicking) mat4MulVec4(mtx mgl32.Mat4, vec mgl32.Vec4) mgl32.Vec4 {
	x := vec.X()*mtx[0] + vec.Y()*mtx[1] + vec.Z()*mtx[2] + vec.W()*mtx[3]
	y := vec.X()*mtx[4] + vec.Y()*mtx[5] + vec.Z()*mtx[6] + vec.W()*mtx[7]
	z := vec.X()*mtx[8] + vec.Y()*mtx[9] + vec.Z()*mtx[10] + vec.W()*mtx[11]
	w := vec.X()*mtx[12] + vec.Y()*mtx[13] + vec.Z()*mtx[14] + vec.W()*mtx[15]
	return mgl32.Vec4{x, y, z, w}
}

func (rp *RayPicking) vec4DivFloat(vec mgl32.Vec4, div float32) mgl32.Vec4 {
	x := vec.X() / div
	y := vec.Y() / div
	z := vec.Z() / div
	w := vec.W() / div
	return mgl32.Vec4{x, y, z, w}
}
