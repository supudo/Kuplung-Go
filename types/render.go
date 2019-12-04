package types

// RenderSettings holds all calculation variables
type RenderSettings struct {
	GLSLVersion string

	Fov         float32
	RatioWidth  float32
	RatioHeight float32
	PlaneClose  float32
	PlaneFar    float32
}
