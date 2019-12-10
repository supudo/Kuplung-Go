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
	cameraPosition mgl32.Vec3

	eyeSettings                                 types.ObjectEye
	positionX, positionY, positionZ             types.ObjectCoordinate
	rotateX, rotateY, rotateZ                   types.ObjectCoordinate
	rotateCenterX, rotateCenterY, rotateCenterZ types.ObjectCoordinate
}

// InitCamera ...
func InitCamera(window interfaces.Window) *Camera {
	camera := &Camera{}
	camera.window = window

	camera.eyeSettings = types.ObjectEye{}
	camera.eyeSettings.ViewEye = mgl32.Vec3{0, 3, 10}
	camera.eyeSettings.ViewCenter = mgl32.Vec3{0, 0, 0}
	camera.eyeSettings.ViewUp = mgl32.Vec3{0, 1, 0}

	camera.positionX = types.ObjectCoordinate{Animate: false, Point: 0}
	camera.positionY = types.ObjectCoordinate{Animate: false, Point: -1}
	camera.positionZ = types.ObjectCoordinate{Animate: false, Point: -10}

	camera.rotateX = types.ObjectCoordinate{Animate: false, Point: 160}
	camera.rotateY = types.ObjectCoordinate{Animate: false, Point: 140}
	camera.rotateZ = types.ObjectCoordinate{Animate: false, Point: 0}

	camera.rotateCenterX = types.ObjectCoordinate{Animate: false, Point: 0}
	camera.rotateCenterY = types.ObjectCoordinate{Animate: false, Point: 0}
	camera.rotateCenterZ = types.ObjectCoordinate{Animate: false, Point: 0}

	camera.MatrixCamera = mgl32.Ident4()

	return camera
}

// Dispose ...
func (camera *Camera) Dispose() {
}

// Render ...
func (camera *Camera) Render() {
	camera.MatrixCamera = mgl32.LookAtV(camera.eyeSettings.ViewEye, camera.eyeSettings.ViewCenter, camera.eyeSettings.ViewUp)

	camera.MatrixCamera = mgl32.Translate3D(camera.positionX.Point, camera.positionY.Point, camera.positionZ.Point)

	// camera.MatrixCamera = mgl32.Translate3D(0, 0, 0)
	// camera.MatrixCamera = mgl32.HomogRotate3D(camera.rotateX.Point, mgl32.Vec3{1, 0, 0})
	// camera.MatrixCamera = mgl32.HomogRotate3D(camera.rotateY.Point, mgl32.Vec3{0, 1, 0})
	// camera.MatrixCamera = mgl32.HomogRotate3D(camera.rotateZ.Point, mgl32.Vec3{0, 0, 1})
	// camera.MatrixCamera = mgl32.Translate3D(0, 0, 0)

	// camera.MatrixCamera = mgl32.HomogRotate3D(camera.rotateCenterX.Point, mgl32.Vec3{1, 0, 0})
	// camera.MatrixCamera = mgl32.HomogRotate3D(camera.rotateCenterY.Point, mgl32.Vec3{0, 1, 0})
	// camera.MatrixCamera = mgl32.HomogRotate3D(camera.rotateCenterZ.Point, mgl32.Vec3{0, 0, 1})

	// camera.cameraPosition = mgl32.Vec3{camera.MatrixCamera[4*3+0]}
	// camera.cameraPosition = mgl32.Vec3{camera.MatrixCamera[4*3+1]}
	// camera.cameraPosition = mgl32.Vec3{camera.MatrixCamera[4*3+2]}
}
