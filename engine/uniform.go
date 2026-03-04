package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joaqu1m/gogl-playground/gmath"
)

func SetUniformMat4(program uint32, name string, m gmath.Mat4) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(loc, 1, false, &m[0])
}

func SetUniformVec3(program uint32, name string, v [3]float32) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.Uniform3f(loc, v[0], v[1], v[2])
}

func SetUniformFloat(program uint32, name string, v float32) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.Uniform1f(loc, v)
}
