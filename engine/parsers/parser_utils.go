package parsers

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// FixVectorAxis ...
func FixVectorAxis(v mgl32.Vec3, indexForward, indexUp int32) mgl32.Vec3 {
	v2 := v

	//
	//                 +Z
	//                  |
	//                  |           +Y
	//                  |         /
	//                  |       /
	//                  |     /
	//                  |   /
	//                  | /
	//  -X--------------|-----------------+X
	//                 /|
	//               /  |
	//             /    |
	//           /      |
	//         /        |
	//       /          |
	//     -Y           |
	//                 -Z
	//

	switch indexForward {
	case 0: // -X Forward
		v2 = Rotate3DZ(v, mgl32.DegToRad(-90.0))
	case 1: // -Y Forward
		v2 = Rotate3DZ(v, mgl32.DegToRad(180.0))
	case 2: // -Z Forward
		v2 = Rotate3DX(v, mgl32.DegToRad(90.0))
	case 3: // X Forward
		v2 = Rotate3DZ(v, mgl32.DegToRad(90.0))
	// case 4: // Y Forward
	case 5: // Z Forward
		v2 = Rotate3DX(v, mgl32.DegToRad(-90.0))
	}

	switch indexUp {
	case 0: // -X Up
		v2 = Rotate3DY(v, mgl32.DegToRad(-90.0))
	//case 1: // -Y Up
	case 2: // -Z Up
		v2 = Rotate3DY(v, mgl32.DegToRad(180.0))
	case 3: // X Up
		v2 = Rotate3DY(v, mgl32.DegToRad(90.0))
	case 4: // Y Up
		v2 = Rotate3DY(v, mgl32.DegToRad(180.0))
		// case 5: // Z Up
	}

	return v2
}

// Rotate3DX ...
func Rotate3DX(v mgl32.Vec3, rads float32) mgl32.Vec3 {
	ca := float32(math.Cos(float64(rads)))
	sa := float32(math.Sin(float64(rads)))
	return mgl32.Vec3{v.X(), v.Y()*ca + v.Z()*sa, v.Y()*sa + v.Z()*ca}
}

// Rotate3DY ...
func Rotate3DY(v mgl32.Vec3, rads float32) mgl32.Vec3 {
	ca := float32(math.Cos(float64(rads)))
	sa := float32(math.Sin(float64(rads)))
	return mgl32.Vec3{v.X()*ca + v.Z()*sa, v.Y(), -v.X()*sa + v.Z()*ca}
}

// Rotate3DZ ...
func Rotate3DZ(v mgl32.Vec3, rads float32) mgl32.Vec3 {
	ca := float32(math.Cos(float64(rads)))
	sa := float32(math.Sin(float64(rads)))
	return mgl32.Vec3{v.X()*ca - v.Y()*sa, v.X()*sa + v.Y()*ca, v.Z()}
}
