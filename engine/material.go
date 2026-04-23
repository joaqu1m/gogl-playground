package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/joaqu1m/gogl-playground/libs/gltfloader"
)

// Material holds the visual properties of a mesh surface.
type Material struct {
	BaseColor        [4]float32
	TextureID        uint32
	HasTexture       bool
	SpecularStrength float32
	Shininess        float32
}

// MaterialFromMesh extracts a Material from the loader's raw mesh data.
func MaterialFromMesh(m *gltfloader.GLTFMesh) Material {
	return Material{
		BaseColor:        m.BaseColor,
		TextureID:        m.TextureID,
		HasTexture:       m.HasTexture,
		SpecularStrength: 0.5,
		Shininess:        32.0,
	}
}

// Bind sends this material's uniforms to the shader and binds textures.
func (mat *Material) Bind(uc *UniformCache) {
	SetUniformVec4(uc, "baseColor", mat.BaseColor)

	if mat.HasTexture {
		SetUniformInt(uc, "useTexture", 1)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, mat.TextureID)
		SetUniformInt(uc, "diffuseMap", 0)
	} else {
		SetUniformInt(uc, "useTexture", 0)
	}

	SetUniformFloat(uc, "specularStrength", mat.SpecularStrength)
	SetUniformFloat(uc, "shininess", mat.Shininess)
}
