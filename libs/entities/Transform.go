package entities

import "github.com/joaqu1m/gogl-playground/gmath"

func IdentityQuat() gmath.Quaternion {
	return gmath.Quaternion{X: 0, Y: 0, Z: 0, W: 1}
}

type Transform struct {
	Position gmath.Vec3
	Rotation gmath.Quaternion
	Scale    gmath.Vec3
}
