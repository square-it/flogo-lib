package engine

import (
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/TIBCOSoftware/flogo-lib/engine/starter"
	"github.com/TIBCOSoftware/flogo-lib/core/process"
)

// Runner manages the creation, reconstitution and execution of
// ProcessInstances
type Runner struct {
	maxStepCount    int
	processProvider  process.Provider
	stateRecorder    processinst.StateRecorder
}

// NewRunner create a new Runner associated with the specified System
func NewRunner(system *System, maxStepCount int) *Runner {

	var runner Runner
	runner.processProvider = system.ProcessProvider()
	runner.stateRecorder = system.StateRecorder()

	if maxStepCount < 1 {
		runner.maxStepCount = int(^uint16(0))
	} else {
		runner.maxStepCount = maxStepCount
	}

	log.Debugf("Max Step Count: %d\n", runner.maxStepCount)

	return &runner
}

// StartInstance creates a new ProcessInstance and prepares it to be executed
func (runner *Runner) StartInstance(instanceID string, startRequest *starter.StartRequest) *processinst.Instance {

	process := runner.processProvider.GetProcess(startRequest.ProcessURI)

	if process == nil {
		log.Errorf("Process [%s] not found", startRequest.ProcessURI)
		return nil
	}

	log.Info("Starting Instance: ", instanceID)

	instance := processinst.NewProcessInstance(instanceID, startRequest.ProcessURI, process)

	if startRequest.Patch != nil {
		log.Infof("Instance [%s] has patch", instanceID)
		instance.Patch = startRequest.Patch
		instance.Patch.Init()
	}

	if startRequest.Interceptor != nil {
		log.Infof("Instance [%s] has interceptor", instanceID)
		instance.Interceptor = startRequest.Interceptor
		instance.Interceptor.Init()
	}

	instance.Start(startRequest.Data)

	return instance
}

// RestartInstance creates a ProcessInstance from an initial state and prepares
// it to be executed
func (runner *Runner) RestartInstance(instanceID string, restartRequest *starter.RestartRequest) *processinst.Instance {

	//todo: handle process not found
	instance := restartRequest.IntialState
	instance.Restart(instanceID, runner.processProvider)

	log.Info("Restarting Instance: ", instanceID)

	if restartRequest.Patch != nil {
		log.Infof("Instance [%s] has patch", instanceID)
		instance.Patch = restartRequest.Patch
		instance.Patch.Init()
	}

	if restartRequest.Interceptor != nil {
		log.Infof("Instance [%s] has interceptor", instanceID)
		instance.Interceptor = restartRequest.Interceptor
		instance.Interceptor.Init()
	}

	instance.UpdateAttrs(restartRequest.Data)

	runner.stateRecorder.RecordSnapshot(instance)
	runner.stateRecorder.RecordStep(instance)

	return instance
}

// ResumeInstance reconstitutes and prepares a ProcessInstance to be resumed
func (runner *Runner) ResumeInstance(resumeRequest *starter.ResumeRequest) *processinst.Instance {

	//todo: handle process not found
	instance := resumeRequest.State

	if resumeRequest.Patch != nil {
		instance.Patch = resumeRequest.Patch
		instance.Patch.Init()
	}

	if resumeRequest.Interceptor != nil {
		instance.Interceptor = resumeRequest.Interceptor
		instance.Interceptor.Init()
	}

	instance.UpdateAttrs(resumeRequest.Data)

	return instance
}

// ExecuteInstance executes the specified ProcessInstance until it is complete
// or it no longer has any tasks to execute
func (runner *Runner) ExecuteInstance(instance *processinst.Instance) bool {

	log.Debugf("Executing Instance: %s\n", instance.ID())

	stepCount := 0
	hasWork := true

	for hasWork && instance.Status() < processinst.StatusCompleted && stepCount < runner.maxStepCount {
		stepCount++
		log.Debugf("Step: %d\n", stepCount)
		hasWork = instance.DoStep()

		runner.stateRecorder.RecordSnapshot(instance)
		runner.stateRecorder.RecordStep(instance)
	}

	log.Debugf("Done Executing Instance: %s\n", instance.ID())

	if instance.Status() == processinst.StatusCompleted {
		log.Infof("Process [%s] Completed", instance.ID())
		return true
	}

	return false
}

