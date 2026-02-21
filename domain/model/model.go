package model

import (
	"github.com/joaqu1m/gogl-playground/libs/entities"
	"github.com/joaqu1m/gogl-playground/libs/gltfloader"
	"github.com/joaqu1m/gogl-playground/libs/logger"
)

type Model struct {
	Name        string
	FilePath    string
	Scale       entities.Dimensions
	Translation entities.Dimensions
	LoadedModel gltfloader.GLTFModel
}

func NewModel(name, filePath string, scale, translation entities.Dimensions) Model {

	logger.Debugf("Loading model %s from path %s", name, filePath)
	loaded, err := gltfloader.LoadGLB(filePath)
	if err != nil || loaded == nil {
		logger.Fatalf("Failed to load model %s from path %s: %v", name, filePath, err)
	}
	logger.Debugf("Loaded model %s from path %s", name, filePath)

	return Model{
		Name:        name,
		FilePath:    filePath,
		Scale:       scale,
		Translation: translation,
		LoadedModel: *loaded,
	}
}
