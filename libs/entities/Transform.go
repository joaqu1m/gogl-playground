package entities

import "github.com/joaqu1m/gogl-playground/gmath"

type Transform struct {
	Position gmath.Vec3
	Rotation gmath.Quaternion
	Scale    gmath.Vec3
}

// ToMat4 computes the world matrix from Position, Rotation and Scale (T * R * S).
func (t Transform) ToMat4() gmath.Mat4 {
	rotMat := t.Rotation.Normalize().ToMat4()
	transMat := gmath.MatTranslate(t.Position)
	scaleMat := gmath.MatScale(t.Scale.X, t.Scale.Y, t.Scale.Z)
	return gmath.MatMul(transMat, gmath.MatMul(rotMat, scaleMat))
}
