package runner

import (
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("runner")

// DirectRunner is a flow runner that executes a flow directly on the same
// thread
// todo: rename to SyncFlowRunner?
type DirectRunner struct {
	maxStepCount  int
	stateRecorder flowinst.StateRecorder
	record        bool
}

// DirectConfig is the configuration object for a DirectRunner
type DirectConfig struct {
	MaxStepCount int `json:"maxStepCount"`
}

// NewDirectRunner create a new DirectRunner
func NewDirectRunner(stateRecorder flowinst.StateRecorder, maxStepCount int) *DirectRunner {

	var directRunner DirectRunner
	directRunner.stateRecorder = stateRecorder
	directRunner.record = stateRecorder != nil

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
	log.Debug("Started Direct Flow Instance Runner")
}

// Stop will stop the engine, by stopping all of its workers
func (runner *DirectRunner) Stop() {
	//no-op
	log.Debug("Stopped Direct Flow Instance Runner")
}

// RunInstance runs the specified Flow Instance until it is complete
// or it no longer has any tasks to execute
func (runner *DirectRunner) RunInstance(instance *flowinst.Instance) bool {

	log.Debugf("Executing Instance: %s\n", instance.ID())

	stepCount := 0
	hasWork := true

	for hasWork && instance.Status() < flowinst.StatusCompleted && stepCount < runner.maxStepCount {
		stepCount++
		log.Debugf("Step: %d\n", stepCount)
		hasWork = instance.DoStep()

		if runner.record {
			runner.stateRecorder.RecordSnapshot(instance)
			runner.stateRecorder.RecordStep(instance)
		}
	}

	log.Debugf("Done Executing Instance: %s\n", instance.ID())

	if instance.Status() == flowinst.StatusCompleted {
		log.Infof("Flow [%s] Completed", instance.ID())
		return true
	}

	return false
}
