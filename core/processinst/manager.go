package processinst

import (
	"github.com/TIBCOSoftware/flogo-lib/core/process"
)

// Manager is used to create or prepare process instance for start, restart or resume
type Manager struct {
	processProvider process.Provider
	idGenerator     IDGenerator
}

// NewManager creates a new Process Instance manager (todo: probably needs a better name)
func NewManager(processProvider process.Provider, idGenerator IDGenerator) *Manager {

	var manager Manager
	manager.processProvider = processProvider
	manager.idGenerator = idGenerator
	return &manager
}

// StartInstance creates a new ProcessInstance and prepares it to be executed
func (mgr *Manager) StartInstance(processURI string, processData map[string]string, replyHandler ReplyHandler, execOptions *ExecOptions) *Instance {

	process := mgr.processProvider.GetProcess(processURI)

	if process == nil {
		log.Errorf("Process [%s] not found", processURI)
		return nil
	}

	instanceID := mgr.idGenerator.NewProcessInstanceID()
	log.Info("Creating Instance: ", instanceID)

	instance := NewProcessInstance(instanceID, processURI, process)

	applyExecOptions(instance, execOptions)

	log.Info("Starting Instance: ", instanceID)

	instance.Start(processData)

	return instance
}

// RestartInstance creates a ProcessInstance from an initial state and prepares
// it to be executed
func (mgr *Manager) RestartInstance(initialState *Instance, processData map[string]string, replyHandler ReplyHandler, execOptions *ExecOptions) *Instance {

	//todo: handle process not found
	instance := initialState
	instanceID := mgr.idGenerator.NewProcessInstanceID()
	instance.Restart(instanceID, mgr.processProvider)

	log.Info("Restarting Instance: ", instanceID)

	applyExecOptions(instance, execOptions)

	instance.UpdateAttrs(processData)

	return instance
}

// ResumeInstance reconstitutes and prepares a ProcessInstance to be resumed
func (mgr *Manager) ResumeInstance(initialState *Instance, processData map[string]string, replyHandler ReplyHandler, execOptions *ExecOptions) *Instance {

	//todo: handle process not found
	instance := initialState
	applyExecOptions(instance, execOptions)
	//instance.Resume(data interface{})
	instance.UpdateAttrs(processData)

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
