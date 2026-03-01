package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joaqu1m/gogl-playground/gmath"
)

func (a *App) Draw() {

	sp := a.DrawingContext.ShaderProgram

	gl.ClearColor(0.1, 0.1, 0.15, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(sp)

	viewMat := a.Camera.ViewMatrix()

	projMat := a.Camera.ProjectionMatrix(
		float32(a.Window.Width) / float32(a.Window.Height),
	)

	gmath.SetUniformMat4(sp, "view", viewMat)
	gmath.SetUniformMat4(sp, "projection", projMat)

	pos := a.Camera.Position()
	gmath.SetUniformVec3(sp, "viewPos",
		[3]float32{pos.X, pos.Y, pos.Z},
	)

	gmath.SetUniformMat4(sp, "view", viewMat)
	gmath.SetUniformMat4(sp, "projection", projMat)
	gmath.SetUniformVec3(sp, "lightDir", [3]float32{-0.3, -0.8, -0.5})

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
			gmath.SetUniformMat4(sp, "model", modelMat)

			// Material
			bcLoc := gl.GetUniformLocation(sp, gl.Str("baseColor\x00"))
			gl.Uniform4f(
				bcLoc,
				m.BaseColor[0],
				m.BaseColor[1],
				m.BaseColor[2],
				m.BaseColor[3],
			)

			utLoc := gl.GetUniformLocation(sp, gl.Str("useTexture\x00"))

			if m.HasTexture {
				gl.Uniform1i(utLoc, 1)
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, m.TextureID)

				dmLoc := gl.GetUniformLocation(sp, gl.Str("diffuseMap\x00"))
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
