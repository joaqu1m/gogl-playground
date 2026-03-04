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

func (mm *ModelManager) GetModel(name string) *Model {
	for i := range mm.Models {
		if mm.Models[i].Name == name {
			return &mm.Models[i]
		}
	}
	return nil
}

func (mm *ModelManager) RemoveModel(name string) {
	for i, m := range mm.Models {
		if m.Name == name {
			mm.Models = append(mm.Models[:i], mm.Models[i+1:]...)
			return
		}
	}
}
