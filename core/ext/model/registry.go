package model

import (
	"sync"
)

var (
	modelsMu sync.RWMutex
	models   = make(map[string]*ProcessModel)
)

// Register registers the specified process model
func Register(model *ProcessModel) {
	modelsMu.Lock()
	defer modelsMu.Unlock()

	if model == nil {
		panic("model.Register: model cannot be nil")
	}

	id := model.Name()

	if _, dup := models[id]; dup {
		panic("model.Register: model " + id + " already registered")
	}

	log.Debugf("Registering ProcessModel: [%s]-%v\n", id, model)

	models[id] = model
}

// Registered gets all the registered process models
func Registered() []*ProcessModel {

	modelsMu.RLock()
	defer modelsMu.RUnlock()

	list := make([]*ProcessModel, 0, len(models))

	for _, value := range models {
		list = append(list, value)
	}

	return list
}

// Get gets specified ProcessModel
func Get(id string) *ProcessModel {
	return models[id]
}
