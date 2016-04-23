package srlocal

import (
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("staterecorder")

// LocalStateRecorder is an implementation of StateRecorder service
// that stores the state locally
type LocalStateRecorder struct {
	host string
}

// NewLocalStateRecorder creates a new LocalStateRecorder
func NewLocalStateRecorder() *LocalStateRecorder {

	return &LocalStateRecorder{}
}

// Start implements util.Managed.Start()
func (sr *LocalStateRecorder) Start() error {
	// no-op
	return nil
}

// Stop implements util.Managed.Stop()
func (sr *LocalStateRecorder) Stop() {
	// no-op
}

// Init implements services.StateRecorderService.Init()
func (sr *LocalStateRecorder) Init(settings map[string]string) {
	//sr.host = settings["host"]
}

// RecordSnapshot implements flowinst.StateRecorder.RecordSnapshot
func (sr *LocalStateRecorder) RecordSnapshot(instance *flowinst.Instance) {

}

// RecordStep implements flowinst.StateRecorder.RecordStep
func (sr *LocalStateRecorder) RecordStep(instance *flowinst.Instance) {

}
