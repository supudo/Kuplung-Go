package objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/supudo/Kuplung-Go/interfaces"
	"github.com/supudo/Kuplung-Go/types"
)

// Camera ...
type Camera struct {
	window interfaces.Window

	MatrixCamera   mgl32.Mat4
	CameraPosition mgl32.Vec3

	EyeSettings                                 types.ObjectEye
	PositionX, PositionY, PositionZ             types.ObjectCoordinate
	RotateX, RotateY, RotateZ                   types.ObjectCoordinate
	RotateCenterX, RotateCenterY, RotateCenterZ types.ObjectCoordinate
}

// InitCamera ...
func InitCamera(window interfaces.Window) *Camera {
	camera := &Camera{}
	camera.window = window
	camera.InitProperties()
	return camera
}

// InitProperties ...
func (camera *Camera) InitProperties() {
	camera.EyeSettings = types.ObjectEye{}
	camera.EyeSettings.ViewEye = mgl32.Vec3{0, 0, 10}
	camera.EyeSettings.ViewCenter = mgl32.Vec3{0, 0, 0}
	camera.EyeSettings.ViewUp = mgl32.Vec3{0, -1, 0}

	camera.PositionX = types.ObjectCoordinate{Animate: false, Point: 0}
	camera.PositionY = types.ObjectCoordinate{Animate: false, Point: 0}
	camera.PositionZ = types.ObjectCoordinate{Animate: false, Point: -16}

	camera.RotateX = types.ObjectCoordinate{Animate: false, Point: 160}
	camera.RotateY = types.ObjectCoordinate{Animate: false, Point: 140}
	camera.RotateZ = types.ObjectCoordinate{Animate: false, Point: 0}

	camera.RotateCenterX = types.ObjectCoordinate{Animate: false, Point: 0}
	camera.RotateCenterY = types.ObjectCoordinate{Animate: false, Point: 0}
	camera.RotateCenterZ = types.ObjectCoordinate{Animate: false, Point: 0}

	camera.MatrixCamera = mgl32.Ident4()
}

// Dispose ...
func (camera *Camera) Dispose() {
}

// Render ...
func (camera *Camera) Render() {
	camera.MatrixCamera = mgl32.LookAtV(camera.EyeSettings.ViewEye, camera.EyeSettings.ViewCenter, camera.EyeSettings.ViewUp)

	camera.MatrixCamera = camera.MatrixCamera.Mul4(mgl32.Translate3D(camera.PositionX.Point, camera.PositionY.Point, camera.PositionZ.Point))

	camera.MatrixCamera = camera.MatrixCamera.Mul4(mgl32.Translate3D(0, 0, 0))
	camera.MatrixCamera = camera.MatrixCamera.Mul4(mgl32.HomogRotate3D(camera.RotateX.Point, mgl32.Vec3{1, 0, 0}))
	camera.MatrixCamera = camera.MatrixCamera.Mul4(mgl32.HomogRotate3D(camera.RotateY.Point, mgl32.Vec3{0, 1, 0}))
	camera.MatrixCamera = camera.MatrixCamera.Mul4(mgl32.HomogRotate3D(camera.RotateZ.Point, mgl32.Vec3{0, 0, 1}))
	camera.MatrixCamera = camera.MatrixCamera.Mul4(mgl32.Translate3D(0, 0, 0))

	camera.MatrixCamera = camera.MatrixCamera.Mul4(mgl32.HomogRotate3D(camera.RotateCenterX.Point, mgl32.Vec3{1, 0, 0}))
	camera.MatrixCamera = camera.MatrixCamera.Mul4(mgl32.HomogRotate3D(camera.RotateCenterY.Point, mgl32.Vec3{0, 1, 0}))
	camera.MatrixCamera = camera.MatrixCamera.Mul4(mgl32.HomogRotate3D(camera.RotateCenterZ.Point, mgl32.Vec3{0, 0, 1}))

	camera.CameraPosition = mgl32.Vec3{camera.MatrixCamera[4*3+0], camera.MatrixCamera[4*3+1], camera.MatrixCamera[4*3+2]}
}
