package model

import (
	"github.com/joaqu1m/gogl-playground/gmath"
	"github.com/joaqu1m/gogl-playground/libs/entities"
	"github.com/joaqu1m/gogl-playground/libs/gltfloader"
	"github.com/joaqu1m/gogl-playground/libs/logger"
)

type Model struct {
	Name        string
	FilePath    string
	Transform   entities.Transform
	LoadedModel gltfloader.GLTFModel
}

// MeshWorldMatrix computes the final world matrix for a mesh within this model.
// Combines the entity transform with the mesh's node transform from glTF.
func (m *Model) MeshWorldMatrix(meshTransform [16]float32) gmath.Mat4 {
	baseMat := m.Transform.ToMat4()
	return gmath.MatMul(baseMat, gmath.Mat4(meshTransform))
}

func (m *Model) Destroy() {
	m.LoadedModel.Destroy()
}

func NewModel(name, filePath string, transform entities.Transform) Model {

	logger.Debugf("Loading model %s from path %s", name, filePath)

	loaded, err := gltfloader.LoadGLB(filePath)
	if err != nil || loaded == nil {
		logger.Fatalf("Failed to load model %s from path %s: %v", name, filePath, err)
	}

	logger.Debugf("Loaded model %s from path %s", name, filePath)

	return Model{
		Name:        name,
		FilePath:    filePath,
		Transform:   transform,
		LoadedModel: *loaded,
	}
}
