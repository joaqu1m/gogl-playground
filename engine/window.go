package engine

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/joaqu1m/gogl-playground/libs/logger"
)

type Window struct {
	GLFWWindow *glfw.Window
	Width      int
	Height     int
	Title      string
}

func NewWindow(height, width int, title string) *Window {

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

	return &Window{
		GLFWWindow: glfwWindow,
		Width:      width,
		Height:     height,
		Title:      title,
	}

}
