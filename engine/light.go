package engine

import (
	"math"

	"github.com/joaqu1m/gogl-playground/gmath"
)

type LightType int

const (
	LightDirectional LightType = iota
	LightPoint
	LightSpot
)

// Light is a unified light source that can represent directional, point,
// or spot lights depending on the Type field.
//
// Directional: uses Direction and Color. No position or attenuation.
// Point: uses Position, Color, and attenuation (Constant/Linear/Quadratic).
// Spot: uses Position, Direction, Color, attenuation, and cone angles (CutOff/OuterCutOff).
type Light struct {
	Type            LightType
	Position        gmath.Vec3
	Direction       gmath.Vec3
	Color           gmath.Vec3
	AmbientStrength float32
	Constant        float32
	Linear          float32
	Quadratic       float32
	CutOff          float32 // cosine of the inner cone angle
	OuterCutOff     float32 // cosine of the outer cone angle
}

func DefaultDirectionalLight() Light {
	return Light{
		Type:            LightDirectional,
		Direction:       gmath.Vec3{X: -0.3, Y: -0.8, Z: -0.5},
		Color:           gmath.Vec3{X: 0.4, Y: 0.4, Z: 0.5},
		AmbientStrength: 0.08,
	}
}

func DefaultPointLight() Light {
	return Light{
		Type:            LightPoint,
		Position:        gmath.Vec3{X: 0.5, Y: 1.5, Z: 1.5},
		Color:           gmath.Vec3{X: 1.0, Y: 0.5, Z: 0.1},
		AmbientStrength: 0.05,
		Constant:        1.0,
		Linear:          0.14,
		Quadratic:       0.07,
	}
}

func DefaultSpotLight() Light {
	return Light{
		Type:            LightSpot,
		Position:        gmath.Vec3{X: 0, Y: 3.0, Z: 1.0},
		Direction:       gmath.Vec3{X: 0, Y: -1.0, Z: -0.3},
		Color:           gmath.Vec3{X: 0.2, Y: 0.6, Z: 1.0},
		AmbientStrength: 0.02,
		CutOff:          float32(math.Cos(float64(math.Pi / 180.0 * 15.0))),
		OuterCutOff:     float32(math.Cos(float64(math.Pi / 180.0 * 25.0))),
		Constant:        1.0,
		Linear:          0.09,
		Quadratic:       0.032,
	}
}
