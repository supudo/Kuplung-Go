package utilities

import "github.com/go-gl/mathgl/mgl32"

// ComputeTangentBasis ...
func ComputeTangentBasis(uvs []mgl32.Vec2, vertices, normals []mgl32.Vec3) (tangents, bitangents []mgl32.Vec3) {
	for i := 0; i < len(vertices); i += 3 {
		// Shortcuts for vertices
		v0 := vertices[i+0]
		v1 := vertices[i+1]
		v2 := vertices[i+2]

		// Shortcuts for UVs
		uv0 := uvs[i+0]
		uv1 := uvs[i+1]
		uv2 := uvs[i+2]

		// Edges of the triangle : postion delta
		deltaPos1 := v1.Sub(v0)
		deltaPos2 := v2.Sub(v0)

		// UV delta
		deltaUV1 := uv1.Sub(uv0)
		deltaUV2 := uv2.Sub(uv0)

		r := 1.0 / (deltaUV1.X()*deltaUV2.Y() - deltaUV1.Y()*deltaUV2.X())
		tangent := deltaPos1.Mul(deltaUV2.Y()).Sub(deltaPos2.Mul(deltaUV1.Y())).Mul(r)
		bitangent := deltaPos2.Mul(deltaUV1.X()).Sub(deltaPos1.Mul(deltaUV2.X())).Mul(r)

		// Set the same tangent for all three vertices of the triangle. They will be merged later.
		tangents = append(tangents, tangent)
		tangents = append(tangents, tangent)
		tangents = append(tangents, tangent)

		// Same thing for binormals
		bitangents = append(bitangents, bitangent)
		bitangents = append(bitangents, bitangent)
		bitangents = append(bitangents, bitangent)
	}

	// for i:=0; i<len(vertices); i+=1 {
	//   n := normals[i]
	//   t := tangents[i]
	//   b := bitangents[i]

	//   // Gram-Schmidt orthogonalize
	//   t = glm::normalize(t - n * glm::dot(n, t));

	//   // Calculate handedness
	//   if (glm::dot(glm::cross(n, t), b) < 0.0f)
	//     t = t * -1.0f;
	// }
	return tangents, bitangents
}
