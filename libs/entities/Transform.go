package entities

import "github.com/joaqu1m/gogl-playground/gmath"

type Transform struct {
	Position gmath.Vec3
	Rotation gmath.Quaternion
	Scale    gmath.Vec3
}
