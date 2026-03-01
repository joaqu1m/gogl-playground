package engine

import "github.com/joaqu1m/gogl-playground/gmath"

type Camera interface {
	ViewMatrix() gmath.Mat4
	ProjectionMatrix(aspect float32) gmath.Mat4
	Position() gmath.Vec3
	Update(dt float32)
}
