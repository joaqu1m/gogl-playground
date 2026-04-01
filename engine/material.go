package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joaqu1m/gogl-playground/libs/gltfloader"
)

// Material holds the visual properties of a mesh surface.
type Material struct {
	BaseColor  [4]float32
	TextureID  uint32
	HasTexture bool
}

// MaterialFromMesh extracts a Material from the loader's raw mesh data.
func MaterialFromMesh(m *gltfloader.GLTFMesh) Material {
	return Material{
		BaseColor:  m.BaseColor,
		TextureID:  m.TextureID,
		HasTexture: m.HasTexture,
	}
}

// Bind sends this material's uniforms to the shader and binds textures.
func (mat *Material) Bind(sp uint32) {
	SetUniformVec4(sp, "baseColor", mat.BaseColor)

	if mat.HasTexture {
		SetUniformInt(sp, "useTexture", 1)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, mat.TextureID)
		SetUniformInt(sp, "diffuseMap", 0)
	} else {
		SetUniformInt(sp, "useTexture", 0)
	}
}
