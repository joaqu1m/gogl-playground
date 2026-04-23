package engine

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joaqu1m/gogl-playground/libs/gltfloader"
)

// BeginFrame clears the screen and activates the shader program.
func (a *App) BeginFrame() {
	gl.ClearColor(0.1, 0.1, 0.15, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(a.DrawingContext.ShaderProgram)
}

// SetGlobalUniforms sends camera and lighting data to the shader.
func (a *App) SetGlobalUniforms() {
	uc := a.DrawingContext.Uniforms

	SetUniformMat4(uc, "view", a.Camera.ViewMatrix())
	SetUniformMat4(uc, "projection", a.Camera.ProjectionMatrix(
		float32(a.Window.Width)/float32(a.Window.Height),
	))

	pos := a.Camera.Position()
	SetUniformVec3(uc, "viewPos", [3]float32{pos.X, pos.Y, pos.Z})

	a.sendLights(uc)
}

// DrawMesh issues the draw call for a single mesh (VAO must already be bound).
func (a *App) DrawMesh(m *gltfloader.GLTFMesh) {
	if m.HasIndices {
		gl.DrawElements(gl.TRIANGLES, m.IndexCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
	} else {
		gl.DrawArrays(gl.TRIANGLES, 0, m.VertexCount)
	}
}

func (a *App) Draw() {
	passes := []RenderPass{
		{
			Name:        "geometry",
			Framebuffer: 0, // default framebuffer (HDR framebuffer later)
			DrawFunc:    a.geometryPass,
		},
	}

	a.BeginFrame()

	for i := range passes {
		passes[i].Execute()
	}
}

func (a *App) geometryPass() {
	uc := a.DrawingContext.Uniforms

	a.SetGlobalUniforms()

	for _, entry := range a.ModelManager.GetModels() {
		for _, m := range entry.LoadedModel.Meshes {
			modelMat := entry.MeshWorldMatrix(m.Transform)
			SetUniformMat4(uc, "model", modelMat)

			mat := MaterialFromMesh(m)
			mat.Bind(uc)

			gl.BindVertexArray(m.VAO)
			a.DrawMesh(m)
		}
	}
}

// sendLights uploads the light array uniforms to the shader.
func (a *App) sendLights(uc *UniformCache) {
	numLights := len(a.Lights)
	if numLights > 8 {
		numLights = 8
	}
	SetUniformInt(uc, "numLights", int32(numLights))

	for i := 0; i < numLights; i++ {
		l := a.Lights[i]
		prefix := fmt.Sprintf("lights[%d].", i)

		SetUniformInt(uc, prefix+"type", int32(l.Type))
		SetUniformVec3(uc, prefix+"position", [3]float32{l.Position.X, l.Position.Y, l.Position.Z})
		SetUniformVec3(uc, prefix+"direction", [3]float32{l.Direction.X, l.Direction.Y, l.Direction.Z})
		SetUniformVec3(uc, prefix+"color", [3]float32{l.Color.X, l.Color.Y, l.Color.Z})
		SetUniformFloat(uc, prefix+"ambient", l.AmbientStrength)
		SetUniformFloat(uc, prefix+"constant", l.Constant)
		SetUniformFloat(uc, prefix+"linear", l.Linear)
		SetUniformFloat(uc, prefix+"quadratic", l.Quadratic)
		SetUniformFloat(uc, prefix+"cutOff", l.CutOff)
		SetUniformFloat(uc, prefix+"outerCutOff", l.OuterCutOff)
	}
}
