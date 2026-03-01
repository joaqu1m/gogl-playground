package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/joaqu1m/gogl-playground/domain/model"
	"github.com/joaqu1m/gogl-playground/gmath"
)

type App struct {
	Window        *glfw.Window
	Width         int
	Height        int
	ShaderProgram uint32
	TimeAccum     float64
	Angle         float64
	Models        []model.Model

	Camera Camera
}

func (a *App) Update() {
	a.Camera.Update(float32(a.TimeAccum))
}

func NewApp(width, height int, title string) *App {
	initGLFW()

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	initOpenGL()

	// construct the App first so that any future helper can access
	// the app context (window, models, etc.).  the camera lives
	// inside the app, so we create it after the struct exists and
	// then assign it.  this also makes it easy to replace or swap
	// cameras later if needed.
	app := &App{
		Window:        window,
		Width:         width,
		Height:        height,
		ShaderProgram: createShaderProgram(),
		Models:        []model.Model{},
	}

	app.Camera = NewFPSCamera(app.Window)

	return app
}

func (a *App) Draw() {

	gl.ClearColor(0.1, 0.1, 0.15, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(a.ShaderProgram)

	viewMat := a.Camera.ViewMatrix()

	projMat := a.Camera.ProjectionMatrix(
		float32(a.Width) / float32(a.Height),
	)

	gmath.SetUniformMat4(a.ShaderProgram, "view", viewMat)
	gmath.SetUniformMat4(a.ShaderProgram, "projection", projMat)

	pos := a.Camera.Position()
	gmath.SetUniformVec3(a.ShaderProgram, "viewPos",
		[3]float32{pos.X, pos.Y, pos.Z},
	)

	gmath.SetUniformMat4(a.ShaderProgram, "view", viewMat)
	gmath.SetUniformMat4(a.ShaderProgram, "projection", projMat)
	gmath.SetUniformVec3(a.ShaderProgram, "lightDir", [3]float32{-0.3, -0.8, -0.5})

	// ----------- Render por modelo -----------

	for _, entry := range a.Models {

		t := entry.Transform

		// Usa apenas a rotação definida no Transform
		rotMat := t.Rotation.Normalize().ToMat4()
		transMat := gmath.MatTranslate(t.Position)
		scaleMat := gmath.MatScale(t.Scale.X, t.Scale.Y, t.Scale.Z)

		// Ordem correta: T * R * S
		baseMat := gmath.MatMul(
			transMat,
			gmath.MatMul(rotMat, scaleMat),
		)

		for _, m := range entry.LoadedModel.Meshes {

			modelMat := gmath.MatMul(baseMat, gmath.Mat4(m.Transform))
			gmath.SetUniformMat4(a.ShaderProgram, "model", modelMat)

			// Material
			bcLoc := gl.GetUniformLocation(a.ShaderProgram, gl.Str("baseColor\x00"))
			gl.Uniform4f(
				bcLoc,
				m.BaseColor[0],
				m.BaseColor[1],
				m.BaseColor[2],
				m.BaseColor[3],
			)

			utLoc := gl.GetUniformLocation(a.ShaderProgram, gl.Str("useTexture\x00"))

			if m.HasTexture {
				gl.Uniform1i(utLoc, 1)
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, m.TextureID)

				dmLoc := gl.GetUniformLocation(a.ShaderProgram, gl.Str("diffuseMap\x00"))
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
