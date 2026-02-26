package engine

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/joaqu1m/gogl-playground/domain/model"
	"github.com/joaqu1m/gogl-playground/libs/logger"
	gmath "github.com/joaqu1m/gogl-playground/math"
)

type App struct {
	window        *glfw.Window
	width         int
	height        int
	shaderProgram uint32
	timeAccum     float64
	angle         float64
	models        []model.Model
}

func NewApp(width, height int, title string) *App {
	initGLFW()

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	initOpenGL()

	return &App{
		window:        window,
		width:         width,
		height:        height,
		shaderProgram: createShaderProgram(),
		models:        []model.Model{},
	}
}

func (a *App) SetModels(models []model.Model) {
	a.models = models
}

func (a *App) Run() {
	previousTime := glfw.GetTime()

	for !a.window.ShouldClose() {
		currentTime := glfw.GetTime()
		a.timeAccum = currentTime - previousTime
		previousTime = currentTime

		a.draw()

		a.window.SwapBuffers()
		glfw.PollEvents()
	}

	logger.Infof("Exiting game loop")
	glfw.Terminate()
}

func (a *App) draw() {
	a.angle += 1.0 * a.timeAccum

	gl.ClearColor(0.1, 0.1, 0.15, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(a.shaderProgram)

	// Rotação global compartilhada por todos os modelos
	rotMat := gmath.MatRotateY(float32(a.angle))

	viewMat := gmath.MatLookAt(
		[3]float32{0, 0.8, 3.0},
		[3]float32{0, 0, 0},
		[3]float32{0, 1, 0},
	)
	projMat := gmath.MatPerspective(
		float32(45.0*math.Pi/180.0),
		float32(a.width)/float32(a.height),
		0.1, 100.0,
	)

	gmath.SetUniformMat4(a.shaderProgram, "view", viewMat)
	gmath.SetUniformMat4(a.shaderProgram, "projection", projMat)
	gmath.SetUniformVec3(a.shaderProgram, "lightDir", [3]float32{-0.3, -0.8, -0.5})

	// Renderiza cada modelo da cena
	for _, entry := range a.models {
		// Transform por modelo: translate * rotate * scale
		scaleMat := gmath.MatScale(entry.Scale.Width, entry.Scale.Height, entry.Scale.Depth)
		transMat := gmath.MatTranslate(entry.Translation.ToArray())
		baseMat := gmath.MatMul(transMat, gmath.MatMul(rotMat, scaleMat))

		logger.Debugf(entry.Name)
		for _, m := range entry.LoadedModel.Meshes {
			modelMat := gmath.MatMul(baseMat, gmath.Mat4(m.Transform))
			gmath.SetUniformMat4(a.shaderProgram, "model", modelMat)

			// Material
			bcLoc := gl.GetUniformLocation(a.shaderProgram, gl.Str("baseColor\x00"))
			gl.Uniform4f(bcLoc, m.BaseColor[0], m.BaseColor[1], m.BaseColor[2], m.BaseColor[3])

			utLoc := gl.GetUniformLocation(a.shaderProgram, gl.Str("useTexture\x00"))
			if m.HasTexture {
				gl.Uniform1i(utLoc, 1)
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, m.TextureID)
				dmLoc := gl.GetUniformLocation(a.shaderProgram, gl.Str("diffuseMap\x00"))
				gl.Uniform1i(dmLoc, 0)
			} else {
				gl.Uniform1i(utLoc, 0)
			}

			gl.BindVertexArray(m.VAO)
			if m.HasIndices {
				gl.DrawElements(gl.TRIANGLES, m.IndexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
			} else {
				gl.DrawArrays(gl.TRIANGLES, 0, m.VertexCount)
			}
		}
	}
}
