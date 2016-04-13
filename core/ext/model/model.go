package model

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("model")

// FlowModel defines the execution Model for a Flow.  It contains the
// execution behaviors for Flowes and Tasks.
type FlowModel struct {
	name          string
	flowBehaviors map[int]FlowBehavior
	taskBehaviors map[int]TaskBehavior
	linkBehaviors map[int]LinkBehavior
}

// New creates a new FlowModel from the specified Behaviors
func New(name string) *FlowModel {

	var flowModel FlowModel
	flowModel.name = name
	flowModel.flowBehaviors = make(map[int]FlowBehavior)
	flowModel.taskBehaviors = make(map[int]TaskBehavior)
	flowModel.linkBehaviors = make(map[int]LinkBehavior)

	return &flowModel
}

// Name returns the name of the FlowModel
func (pm *FlowModel) Name() string {
	return pm.name
}

// RegisterFlowBehavior registers the specified FlowBehavior with the Model
func (pm *FlowModel) RegisterFlowBehavior(id int, flowBehavior FlowBehavior) {

	log.Debugf("Registering Flow Behavior: [%d]-%v\n", id, flowBehavior)
	pm.flowBehaviors[id] = flowBehavior
}

// GetFlowBehavior returns FlowBehavior with the specified ID in the FlowModel
func (pm *FlowModel) GetFlowBehavior(id int) FlowBehavior {
	return pm.flowBehaviors[id]
}

// RegisterTaskBehavior registers the specified TaskBehavior with the Model
func (pm *FlowModel) RegisterTaskBehavior(id int, taskBehavior TaskBehavior) {

	log.Debugf("Registering Task Behavior: [%d]-%v\n", id, taskBehavior)
	pm.taskBehaviors[id] = taskBehavior
}

// GetTaskBehavior returns TaskBehavior with the specified ID in he FlowModel
func (pm *FlowModel) GetTaskBehavior(id int) TaskBehavior {
	return pm.taskBehaviors[id]
}

// RegisterLinkBehavior registers the specified LinkBehavior with the Model
func (pm *FlowModel) RegisterLinkBehavior(id int, linkBehavior LinkBehavior) {

	log.Debugf("Registering Link Behavior: [%d]-%v\n", id, linkBehavior)
	pm.linkBehaviors[id] = linkBehavior
}

// GetLinkBehavior returns LinkBehavior with the specified ID in the FlowModel
func (pm *FlowModel) GetLinkBehavior(id int) LinkBehavior {
	return pm.linkBehaviors[id]
}
