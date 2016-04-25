package flowinst

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/flow"
)

// Manager is used to create or prepare flow instance for start, restart or resume
type Manager struct {
	flowProvider flow.Provider
	idGenerator  IDGenerator
}

// NewManager creates a new Flow Instance manager (todo: probably needs a better name)
func NewManager(flowProvider flow.Provider, idGenerator IDGenerator) *Manager {

	var manager Manager
	manager.flowProvider = flowProvider
	manager.idGenerator = idGenerator
	return &manager
}

// StartInstance creates a new FlowInstance and prepares it to be executed
func (mgr *Manager) StartInstance(flowURI string, startAttrs []*data.Attribute, replyHandler ReplyHandler, execOptions *ExecOptions) (*Instance, error) {

	flow := mgr.flowProvider.GetFlow(flowURI)

	if flow == nil {
		err := fmt.Errorf("Flow [%s] not found", flowURI)
		return nil, err
	}

	instanceID := mgr.idGenerator.NewFlowInstanceID()
	log.Info("Creating Instance: ", instanceID)

	instance := NewFlowInstance(instanceID, flowURI, flow)

	applyExecOptions(instance, execOptions)

	log.Info("Starting Instance: ", instanceID)

	instance.Start(startAttrs)

	return instance, nil
}

// RestartInstance creates a FlowInstance from an initial state and prepares
// it to be executed
func (mgr *Manager) RestartInstance(initialState *Instance, flowData map[string]interface{}, replyHandler ReplyHandler, execOptions *ExecOptions) *Instance {

	//todo: handle flow not found
	instance := initialState
	instanceID := mgr.idGenerator.NewFlowInstanceID()
	instance.Restart(instanceID, mgr.flowProvider)

	log.Info("Restarting Instance: ", instanceID)

	applyExecOptions(instance, execOptions)

	instance.UpdateAttrs(flowData)

	return instance
}

// ResumeInstance reconstitutes and prepares a FlowInstance to be resumed
func (mgr *Manager) ResumeInstance(initialState *Instance, flowData map[string]interface{}, replyHandler ReplyHandler, execOptions *ExecOptions) *Instance {

	//todo: handle flow not found
	instance := initialState
	applyExecOptions(instance, execOptions)
	//instance.Resume(data interface{})
	instance.UpdateAttrs(flowData)

	return instance
}

func applyExecOptions(instance *Instance, execOptions *ExecOptions) {

	if execOptions != nil {

		if execOptions.Patch != nil {
			log.Infof("Instance [%s] has patch", instance.ID())
			instance.Patch = execOptions.Patch
			instance.Patch.Init()
		}

		if execOptions.Interceptor != nil {
			log.Infof("Instance [%s] has interceptor", instance.ID)
			instance.Interceptor = execOptions.Interceptor
			instance.Interceptor.Init()
		}
	}
}
