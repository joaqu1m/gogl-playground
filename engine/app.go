package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joaqu1m/gogl-playground/domain/model"
	"github.com/joaqu1m/gogl-playground/libs/logger"
)

type App struct {
	Window         *Window
	DrawingContext DrawingContext
	TimeAccum      float64
	Models         []model.Model
	Camera         Camera
	Light          Light
}

type DrawingContext struct {
	ShaderProgram uint32
}

func NewApp(width, height int, title string) *App {

	window := NewWindow(width, height, title)

	window.GLFWWindow.MakeContextCurrent()

	initOpenGL()

	return &App{
		Window: window,
		DrawingContext: DrawingContext{
			ShaderProgram: createShaderProgram(),
		},
		Models: []model.Model{},
		Camera: NewFPSCamera(window.GLFWWindow),
		Light:  DefaultLight(),
	}
}

func initOpenGL() {
	if err := gl.Init(); err != nil {
		logger.Fatalf("Failed to initialize OpenGL: %v", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
}
