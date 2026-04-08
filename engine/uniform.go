package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joaqu1m/gogl-playground/gmath"
)

// UniformCache resolves and caches uniform locations for a shader program.
// First call to Loc() queries the GL driver; subsequent calls return the cached value.
type UniformCache struct {
	program   uint32
	locations map[string]int32
}

func NewUniformCache(program uint32) *UniformCache {
	return &UniformCache{
		program:   program,
		locations: make(map[string]int32),
	}
}

// Loc returns the cached uniform location, resolving it on first access.
func (c *UniformCache) Loc(name string) int32 {
	if loc, ok := c.locations[name]; ok {
		return loc
	}
	loc := gl.GetUniformLocation(c.program, gl.Str(name+"\x00"))
	c.locations[name] = loc
	return loc
}

func SetUniformMat4(c *UniformCache, name string, m gmath.Mat4) {
	gl.UniformMatrix4fv(c.Loc(name), 1, false, &m[0])
}

func SetUniformVec3(c *UniformCache, name string, v [3]float32) {
	gl.Uniform3f(c.Loc(name), v[0], v[1], v[2])
}

func SetUniformFloat(c *UniformCache, name string, v float32) {
	gl.Uniform1f(c.Loc(name), v)
}

func SetUniformInt(c *UniformCache, name string, v int32) {
	gl.Uniform1i(c.Loc(name), v)
}

func SetUniformVec4(c *UniformCache, name string, v [4]float32) {
	gl.Uniform4f(c.Loc(name), v[0], v[1], v[2], v[3])
}
