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
	ModelManager   *model.ModelManager
	Camera         Camera
	Lights         []Light
}

type DrawingContext struct {
	ShaderProgram uint32
	Uniforms      *UniformCache
}

func NewApp(width, height int, title string) *App {

	window := NewWindow(width, height, title)

	window.GLFWWindow.MakeContextCurrent()

	initOpenGL()

	return &App{
		Window: window,
		DrawingContext: func() DrawingContext {
			sp := createShaderProgram()
			return DrawingContext{
				ShaderProgram: sp,
				Uniforms:      NewUniformCache(sp),
			}
		}(),
		ModelManager: model.NewModelManager(),
		Camera: NewFPSCamera(window.GLFWWindow),
		Lights: []Light{
			DefaultDirectionalLight(),
			DefaultPointLight(),
			DefaultSpotLight(),
		},
	}
}

func initOpenGL() {
	if err := gl.Init(); err != nil {
		logger.Fatalf("Failed to initialize OpenGL: %v", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
}
