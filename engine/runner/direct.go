package runner

import (
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("runner")

// DirectRunner is a process runner that executes a process directly on the same
// thread
// todo: rename to SyncProcessRunner?
type DirectRunner struct {
	maxStepCount  int
	stateRecorder processinst.StateRecorder
	record        bool
}

// DirectConfig is the configuration object for a DirectRunner
type DirectConfig struct {
	MaxStepCount int `json:"maxStepCount"`
}

// NewDirectRunner create a new DirectRunner
func NewDirectRunner(stateRecorder processinst.StateRecorder, maxStepCount int) *DirectRunner {

	var directRunner DirectRunner
	directRunner.stateRecorder = stateRecorder

	if maxStepCount < 1 {
		directRunner.maxStepCount = int(^uint16(0))
	} else {
		directRunner.maxStepCount = maxStepCount
	}

	return &directRunner
}

// Start will start the engine, by starting all of its workers
func (runner *DirectRunner) Start() {
	//op-op
	log.Debug("Started Direct Process Instance Runner")
}

// Stop will stop the engine, by stopping all of its workers
func (runner *DirectRunner) Stop() {
	//no-op
	log.Debug("Stopped Direct Process Instance Runner")
}

// RunInstance runs the specified Process Instance until it is complete
// or it no longer has any tasks to execute
func (runner *DirectRunner) RunInstance(instance *processinst.Instance) bool {

	log.Debugf("Executing Instance: %s\n", instance.ID())

	stepCount := 0
	hasWork := true

	for hasWork && instance.Status() < processinst.StatusCompleted && stepCount < runner.maxStepCount {
		stepCount++
		log.Debugf("Step: %d\n", stepCount)
		hasWork = instance.DoStep()

		if runner.record {
			runner.stateRecorder.RecordSnapshot(instance)
			runner.stateRecorder.RecordStep(instance)
		}
	}

	log.Debugf("Done Executing Instance: %s\n", instance.ID())

	if instance.Status() == processinst.StatusCompleted {
		log.Infof("Process [%s] Completed", instance.ID())
		return true
	}

	return false
}
