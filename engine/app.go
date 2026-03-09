package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/joaqu1m/gogl-playground/domain/model"
	"github.com/joaqu1m/gogl-playground/libs/logger"
)

type App struct {
	Window         *Window
	DrawingContext DrawingContext
	ModelManager   *model.ModelManager
	Camera         Camera
	Lights         []Light
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
		ModelManager: model.NewModelManager(),
		Camera: NewFPSCamera(window.GLFWWindow),
		Lights: []Light{
			DefaultDirectionalLight(),
			DefaultPointLight(),
			DefaultSpotLight(),
		},
	}
}

func (a *App) Shutdown() {
	a.ModelManager.Destroy()
	gl.DeleteProgram(a.DrawingContext.ShaderProgram)
	a.DrawingContext.ShaderProgram = 0
}

func (a *App) Run() {
	previousTime := glfw.GetTime()

	for !a.Window.GLFWWindow.ShouldClose() {
		currentTime := glfw.GetTime()
		delta := currentTime - previousTime
		previousTime = currentTime

		a.Camera.Update(float32(delta))
		a.Draw()

		a.Window.GLFWWindow.SwapBuffers()
		glfw.PollEvents()
	}

	logger.Infof("Exiting game loop")
	a.Shutdown()
	glfw.Terminate()
	logger.Infof("Game closed")
}

func initOpenGL() {
	if err := gl.Init(); err != nil {
		logger.Fatalf("Failed to initialize OpenGL: %v", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
}
