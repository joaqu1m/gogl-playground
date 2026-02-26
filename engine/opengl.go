package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
}
