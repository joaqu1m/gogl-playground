package model

type ModelManager struct {
	Models []Model
}

func NewModelManager() *ModelManager {
	return &ModelManager{
		Models: []Model{},
	}
}

func (mm *ModelManager) AddModel(model Model) {
	mm.Models = append(mm.Models, model)
}

func (mm *ModelManager) GetModels() []Model {
	return mm.Models
}
