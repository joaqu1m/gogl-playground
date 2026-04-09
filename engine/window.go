package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/joaqu1m/gogl-playground/libs/logger"
)

type Window struct {
	GLFWWindow *glfw.Window
	Width      int
	Height     int
	Title      string
}

func NewWindow(width, height int, title string) *Window {

	if err := glfw.Init(); err != nil {
		logger.Fatalf(err.Error())
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	glfwWindow, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	w := &Window{
		GLFWWindow: glfwWindow,
		Width:      width,
		Height:     height,
		Title:      title,
	}

	glfwWindow.SetFramebufferSizeCallback(func(_ *glfw.Window, fbWidth, fbHeight int) {
		w.Width = fbWidth
		w.Height = fbHeight
		gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
	})

	return w
}
