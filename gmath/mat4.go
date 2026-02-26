package gmath

import (
	"math"
	smath "math"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Mat4 [16]float32

type Quaternion struct {
	X float32
	Y float32
	Z float32
	W float32
}

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

func MatIdentity() Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func MatRotateX(angle float32) Mat4 {
	c := float32(smath.Cos(float64(angle)))
	s := float32(smath.Sin(float64(angle)))

	return Mat4{
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
		0, 0, 0, 1,
	}
}

func MatRotateY(angle float32) Mat4 {
	c := float32(smath.Cos(float64(angle)))
	s := float32(smath.Sin(float64(angle)))

	return Mat4{
		c, 0, -s, 0,
		0, 1, 0, 0,
		s, 0, c, 0,
		0, 0, 0, 1,
	}
}

func MatRotateZ(angle float32) Mat4 {
	c := float32(smath.Cos(float64(angle)))
	s := float32(smath.Sin(float64(angle)))

	return Mat4{
		c, s, 0, 0,
		-s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

func MatMul(a, b Mat4) Mat4 {
	var r Mat4

	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			var sum float32

			for k := 0; k < 4; k++ {
				sum += a[k*4+row] * b[col*4+k]
			}

			r[col*4+row] = sum
		}
	}

	return r
}

func MatScaleUniform(s float32) Mat4 {
	return MatScale(s, s, s)
}

func MatScale(sx, sy, sz float32) Mat4 {
	return Mat4{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
}

func MatPerspective(fovy, aspect, near, far float32) Mat4 {
	f := float32(1.0 / smath.Tan(float64(fovy/2.0)))
	nf := near - far

	return Mat4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (far + near) / nf, -1,
		0, 0, (2 * far * near) / nf, 0,
	}
}

func MatLookAt(eye, center, up [3]float32) Mat4 {
	f := vecNormalize(vecSub(center, eye))
	s := vecNormalize(vecCross(f, up))
	u := vecCross(s, f)

	return Mat4{
		s[0], u[0], -f[0], 0,
		s[1], u[1], -f[1], 0,
		s[2], u[2], -f[2], 0,
		-vecDot(s, eye), -vecDot(u, eye), vecDot(f, eye), 1,
	}
}

func vecSub(a, b [3]float32) [3]float32 {
	return [3]float32{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func vecCross(a, b [3]float32) [3]float32 {
	return [3]float32{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func vecDot(a, b [3]float32) float32 {
	return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func vecNormalize(v [3]float32) [3]float32 {
	l := float32(smath.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))

	if l == 0 {
		return v
	}

	return [3]float32{v[0] / l, v[1] / l, v[2] / l}
}

func SetUniformMat4(program uint32, name string, m Mat4) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(loc, 1, false, &m[0])
}

func SetUniformVec3(program uint32, name string, v [3]float32) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.Uniform3f(loc, v[0], v[1], v[2])
}

func MatTranslate(v Vec3) Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		v.X, v.Y, v.Z, 1,
	}
}

//Func para rotações com quaternions

func QuatFromAxisAngle(axis Vec3, angleRad float32) Quaternion {
	half := angleRad / 2
	s := float32(math.Sin(float64(half)))
	c := float32(math.Cos(float64(half)))

	return Quaternion{
		X: axis.X * s,
		Y: axis.Y * s,
		Z: axis.Z * s,
		W: c,
	}
}

func MultiplyQuat(a, b Quaternion) Quaternion {
	return Quaternion{
		W: a.W*b.W - a.X*b.X - a.Y*b.Y - a.Z*b.Z,
		X: a.W*b.X + a.X*b.W + a.Y*b.Z - a.Z*b.Y,
		Y: a.W*b.Y - a.X*b.Z + a.Y*b.W + a.Z*b.X,
		Z: a.W*b.Z + a.X*b.Y - a.Y*b.X + a.Z*b.W,
	}
}

func (q Quaternion) ToMat4() Mat4 {
	x := q.X
	y := q.Y
	z := q.Z
	w := q.W

	return Mat4{
		1 - 2*y*y - 2*z*z,
		2*x*y + 2*w*z,
		2*x*z - 2*w*y,
		0,

		2*x*y - 2*w*z,
		1 - 2*x*x - 2*z*z,
		2*y*z + 2*w*x,
		0,

		2*x*z + 2*w*y,
		2*y*z - 2*w*x,
		1 - 2*x*x - 2*y*y,
		0,

		0, 0, 0, 1,
	}
}

func (q Quaternion) Normalize() Quaternion {
	length := float32(math.Sqrt(float64(
		q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W,
	)))

	return Quaternion{
		X: q.X / length,
		Y: q.Y / length,
		Z: q.Z / length,
		W: q.W / length,
	}
}
