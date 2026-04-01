package engine

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joaqu1m/gogl-playground/gmath"
)

func (a *App) Draw() {

	sp := a.DrawingContext.ShaderProgram

	gl.ClearColor(0.1, 0.1, 0.15, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(sp)

	// ---- Global uniforms ----

	SetUniformMat4(sp, "view", a.Camera.ViewMatrix())
	SetUniformMat4(sp, "projection", a.Camera.ProjectionMatrix(
		float32(a.Window.Width)/float32(a.Window.Height),
	))

	pos := a.Camera.Position()
	SetUniformVec3(sp, "viewPos", [3]float32{pos.X, pos.Y, pos.Z})

	// Material specular properties (will move to per-mesh material later)
	SetUniformFloat(sp, "specularStrength", 0.5)
	SetUniformFloat(sp, "shininess", 32.0)

	a.sendLights(sp)

	// ---- Render models ----

	for _, entry := range a.ModelManager.GetModels() {

		baseMat := entry.Transform.ToMat4()

		for _, m := range entry.LoadedModel.Meshes {

			modelMat := gmath.MatMul(baseMat, gmath.Mat4(m.Transform))
			SetUniformMat4(sp, "model", modelMat)

			mat := MaterialFromMesh(m)
			mat.Bind(sp)

			gl.BindVertexArray(m.VAO)

			if m.HasIndices {
				gl.DrawElements(gl.TRIANGLES, m.IndexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
			} else {
				gl.DrawArrays(gl.TRIANGLES, 0, m.VertexCount)
			}
		}
	}
}

// sendLights uploads the light array uniforms to the shader.
func (a *App) sendLights(sp uint32) {
	numLights := len(a.Lights)
	if numLights > 8 {
		numLights = 8
	}
	SetUniformInt(sp, "numLights", int32(numLights))

	for i := 0; i < numLights; i++ {
		l := a.Lights[i]
		prefix := fmt.Sprintf("lights[%d].", i)

		SetUniformInt(sp, prefix+"type", int32(l.Type))
		SetUniformVec3(sp, prefix+"position", [3]float32{l.Position.X, l.Position.Y, l.Position.Z})
		SetUniformVec3(sp, prefix+"direction", [3]float32{l.Direction.X, l.Direction.Y, l.Direction.Z})
		SetUniformVec3(sp, prefix+"color", [3]float32{l.Color.X, l.Color.Y, l.Color.Z})
		SetUniformFloat(sp, prefix+"ambient", l.AmbientStrength)
		SetUniformFloat(sp, prefix+"constant", l.Constant)
		SetUniformFloat(sp, prefix+"linear", l.Linear)
		SetUniformFloat(sp, prefix+"quadratic", l.Quadratic)
		SetUniformFloat(sp, prefix+"cutOff", l.CutOff)
		SetUniformFloat(sp, prefix+"outerCutOff", l.OuterCutOff)
	}
}
