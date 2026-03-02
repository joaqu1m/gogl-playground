package engine

import "github.com/joaqu1m/gogl-playground/gmath"

// Light represents a directional light source (like the sun).
// A directional light has no position — it illuminates everything
// from the same direction with the same intensity.
type Light struct {
	Direction       gmath.Vec3 // where the light is coming FROM (will be negated in shader)
	Color           gmath.Vec3 // RGB color of the light, e.g. {1,1,1} = white
	AmbientStrength float32    // how bright the base/ambient light is (0.0 - 1.0)
	SpecularStrength float32   // how strong the shiny highlight is   (0.0 - 1.0)
	Shininess        float32   // how tight/small the specular spot is (8, 32, 128...)
}

func DefaultLight() Light {
	return Light{
		Direction:        gmath.Vec3{X: -0.3, Y: -0.8, Z: -0.5},
		Color:            gmath.Vec3{X: 1.0, Y: 1.0, Z: 1.0},
		AmbientStrength:  0.2,
		SpecularStrength: 0.5,
		Shininess:        32.0,
	}
}
