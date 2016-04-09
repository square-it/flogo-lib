package model

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("model")

// ProcessModel defines the execution Model for a Process.  It contains the
// execution behaviors for Processes and Tasks.
type ProcessModel struct {
	name             string
	processBehaviors map[int]ProcessBehavior
	taskBehaviors    map[int]TaskBehavior
	linkBehaviors    map[int]LinkBehavior
}

// New creates a new ProcessModel from the specified Behaviors
func New(name string) *ProcessModel {

	var processModel ProcessModel
	processModel.name = name
	processModel.processBehaviors = make(map[int]ProcessBehavior)
	processModel.taskBehaviors = make(map[int]TaskBehavior)
	processModel.linkBehaviors = make(map[int]LinkBehavior)

	return &processModel
}

// Name returns the name of the ProcessModel
func (pm *ProcessModel) Name() string {
	return pm.name
}

// RegisterProcessBehavior registers the specified ProcessBehavior with the Model
func (pm *ProcessModel) RegisterProcessBehavior(id int, processBehavior ProcessBehavior) {

	log.Debugf("Registering Process Behavior: [%d]-%v\n", id, processBehavior)
	pm.processBehaviors[id] = processBehavior
}

// GetProcessBehavior returns ProcessBehavior with the specified ID in the ProcessModel
func (pm *ProcessModel) GetProcessBehavior(id int) ProcessBehavior {
	return pm.processBehaviors[id]
}

// RegisterTaskBehavior registers the specified TaskBehavior with the Model
func (pm *ProcessModel) RegisterTaskBehavior(id int, taskBehavior TaskBehavior) {

	log.Debugf("Registering Task Behavior: [%d]-%v\n", id, taskBehavior)
	pm.taskBehaviors[id] = taskBehavior
}

// GetTaskBehavior returns TaskBehavior with the specified ID in he ProcessModel
func (pm *ProcessModel) GetTaskBehavior(id int) TaskBehavior {
	return pm.taskBehaviors[id]
}

// RegisterLinkBehavior registers the specified LinkBehavior with the Model
func (pm *ProcessModel) RegisterLinkBehavior(id int, linkBehavior LinkBehavior) {

	log.Debugf("Registering Link Behavior: [%d]-%v\n", id, linkBehavior)
	pm.linkBehaviors[id] = linkBehavior
}

// GetLinkBehavior returns LinkBehavior with the specified ID in the ProcessModel
func (pm *ProcessModel) GetLinkBehavior(id int) LinkBehavior {
	return pm.linkBehaviors[id]
}
